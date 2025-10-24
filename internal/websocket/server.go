package websocket

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"sync"

	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/khanghh/mcrunner/internal/ptyproc"
)

// Options configures the multiplex websocket server behavior.
type Options struct {
	// Route to register. Defaults to "/ws" when empty.
	Route string
	// BufferSize for per-session ring buffer in bytes. Defaults to 1 MiB when <= 0.
	BufferSize int
	// ClientQueue is the per-client outbound queue length. Defaults to 64 when <= 0.
	ClientQueue int
}

// Server multiplexes multiple PTY sessions over a single websocket connection.
// One reader goroutine per PTY session, fanout to subscribed clients.
type Server struct {
	m   *ptyproc.PTYManager
	opt Options

	mu       sync.RWMutex
	sessions map[string]*muxSession
}

type muxSession struct {
	id   string
	rw   io.ReadWriter
	ring *ringBuffer

	mu       sync.Mutex
	subs     map[*client]struct{}
	readOnce sync.Once
}

type client struct {
	conn   *fiberws.Conn
	out    chan []byte
	subsMu sync.Mutex
	subs   map[string]struct{}
	closed chan struct{}
}

// control frames
type msgBase struct {
	Type string `json:"type"`
}

type msgSubscribe struct {
	Type     string   `json:"type"`
	Sessions []string `json:"sessions"`
	Replay   bool     `json:"replay"`
}

type msgUnsubscribe struct {
	Type     string   `json:"type"`
	Sessions []string `json:"sessions"`
}

type msgInput struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId"`
	DataB64   string `json:"data"`
}

type msgResize struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId"`
	Cols      int    `json:"cols"`
	Rows      int    `json:"rows"`
}

type msgBufferReq struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId"`
}

// server -> client frames
type msgOutput struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId"`
	DataB64   string `json:"data"`
}

type msgError struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId,omitempty"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

// NewServer constructs a multiplex Server bound to a PTYManager.
func NewServer(m *ptyproc.PTYManager, opt Options) *Server {
	if opt.BufferSize <= 0 {
		opt.BufferSize = 1 << 20
	}
	if opt.ClientQueue <= 0 {
		opt.ClientQueue = 64
	}
	if opt.Route == "" {
		opt.Route = "/ws"
	}
	return &Server{m: m, opt: opt, sessions: make(map[string]*muxSession)}
}

// RegisterRoutes registers the multiplexing websocket endpoint at s.opt.Route (default "/ws").
func (s *Server) RegisterRoutes(router fiber.Router) {
	route := s.opt.Route
	router.Use(route, func(c *fiber.Ctx) error {
		if fiberws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	router.Get(route, fiberws.New(func(conn *fiberws.Conn) {
		cl := &client{
			conn:   conn,
			out:    make(chan []byte, s.opt.ClientQueue),
			subs:   make(map[string]struct{}),
			closed: make(chan struct{}),
		}

		// writer goroutine
		go func() {
			for {
				select {
				case frame, ok := <-cl.out:
					if !ok {
						return
					}
					if err := cl.conn.WriteMessage(fiberws.TextMessage, frame); err != nil {
						_ = cl.conn.Close()
						return
					}
				case <-cl.closed:
					return
				}
			}
		}()

		// cleanup on exit
		defer func() {
			// remove from each subscribed session
			s.mu.RLock()
			var list []*muxSession
			for sid := range cl.subs {
				if ms, ok := s.sessions[sid]; ok {
					list = append(list, ms)
				}
			}
			s.mu.RUnlock()
			for _, ms := range list {
				ms.mu.Lock()
				delete(ms.subs, cl)
				ms.mu.Unlock()
			}
			close(cl.closed)
			close(cl.out)
			_ = conn.Close()
		}()

		// read loop
		for {
			mt, data, err := conn.ReadMessage()
			if err != nil {
				break
			}
			if mt != fiberws.TextMessage && mt != fiberws.BinaryMessage {
				continue
			}
			s.handleCommand(cl, data)
		}
	}))
}

// handleCommand parses a single client frame and routes it to appropriate handlers.
func (s *Server) handleCommand(cl *client, data []byte) {
	var base msgBase
	if err := json.Unmarshal(data, &base); err != nil {
		s.sendErr(cl, "", "bad_json", err.Error())
		return
	}
	switch base.Type {
	case "subscribe":
		var msub msgSubscribe
		if err := json.Unmarshal(data, &msub); err != nil {
			s.sendErr(cl, "", "bad_subscribe", err.Error())
			return
		}
		s.handleSubscribe(cl, msub.Sessions, msub.Replay)
	case "unsubscribe":
		var mu msgUnsubscribe
		if err := json.Unmarshal(data, &mu); err != nil {
			s.sendErr(cl, "", "bad_unsubscribe", err.Error())
			return
		}
		s.handleUnsubscribe(cl, mu.Sessions)
	case "input":
		var mi msgInput
		if err := json.Unmarshal(data, &mi); err != nil {
			s.sendErr(cl, "", "bad_input", err.Error())
			return
		}
		s.handleInput(mi.SessionID, mi.DataB64)
	case "resize":
		var mr msgResize
		if err := json.Unmarshal(data, &mr); err != nil {
			s.sendErr(cl, "", "bad_resize", err.Error())
			return
		}
		if err := s.m.Resize(mr.SessionID, mr.Cols, mr.Rows); err != nil {
			s.sendErr(cl, mr.SessionID, "resize_failed", err.Error())
		}
	case "buffer":
		var mb msgBufferReq
		if err := json.Unmarshal(data, &mb); err != nil {
			s.sendErr(cl, "", "bad_buffer", err.Error())
			return
		}
		s.sendBuffer(cl, mb.SessionID)
	default:
		s.sendErr(cl, "", "unknown_type", base.Type)
	}
}

func (s *Server) sendErr(cl *client, sid, code, msg string) {
	b, _ := json.Marshal(msgError{Type: "error", SessionID: sid, Code: code, Message: msg})
	s.enqueue(cl, b)
}

func (s *Server) enqueue(cl *client, frame []byte) {
	select {
	case cl.out <- frame:
	default:
		_ = cl.conn.Close() // slow client protection
	}
}

func (s *Server) handleSubscribe(cl *client, sessionIDs []string, replay bool) {
	for _, sid := range sessionIDs {
		ms, err := s.getOrInitSession(sid)
		if err != nil {
			s.sendErr(cl, sid, "not_found", err.Error())
			continue
		}
		ms.mu.Lock()
		if ms.subs == nil {
			ms.subs = make(map[*client]struct{})
		}
		ms.subs[cl] = struct{}{}
		ms.mu.Unlock()

		cl.subsMu.Lock()
		cl.subs[sid] = struct{}{}
		cl.subsMu.Unlock()

		if replay {
			s.sendBuffer(cl, sid)
		}
	}
}

func (s *Server) handleUnsubscribe(cl *client, sessionIDs []string) {
	for _, sid := range sessionIDs {
		s.mu.RLock()
		ms, ok := s.sessions[sid]
		s.mu.RUnlock()
		if !ok {
			continue
		}
		ms.mu.Lock()
		delete(ms.subs, cl)
		ms.mu.Unlock()
		cl.subsMu.Lock()
		delete(cl.subs, sid)
		cl.subsMu.Unlock()
	}
}

func (s *Server) handleInput(sid, dataB64 string) {
	if sid == "" {
		return
	}
	data, err := base64.StdEncoding.DecodeString(dataB64)
	if err != nil || len(data) == 0 {
		return
	}
	rw, err := s.m.PTY(sid)
	if err != nil {
		return
	}
	_, _ = rw.Write(data)
}

func (s *Server) sendBuffer(cl *client, sid string) {
	s.mu.RLock()
	ms, ok := s.sessions[sid]
	s.mu.RUnlock()
	if !ok {
		s.sendErr(cl, sid, "not_found", "session not found")
		return
	}
	snap := ms.ring.Snapshot()
	if len(snap) == 0 {
		return
	}
	out := msgOutput{Type: "buffer", SessionID: sid, DataB64: base64.StdEncoding.EncodeToString(snap)}
	b, _ := json.Marshal(out)
	s.enqueue(cl, b)
}

func (s *Server) getOrInitSession(sid string) (*muxSession, error) {
	s.mu.RLock()
	if ms, ok := s.sessions[sid]; ok {
		s.mu.RUnlock()
		return ms, nil
	}
	s.mu.RUnlock()

	rw, err := s.m.PTY(sid)
	if err != nil {
		return nil, err
	}
	ms := &muxSession{
		id:   sid,
		rw:   rw,
		ring: newRingBuffer(s.opt.BufferSize),
		subs: make(map[*client]struct{}),
	}

	s.mu.Lock()
	if exist, ok := s.sessions[sid]; ok {
		s.mu.Unlock()
		return exist, nil
	}
	s.sessions[sid] = ms
	s.mu.Unlock()

	ms.readOnce.Do(func() { go s.readLoop(ms) })
	return ms, nil
}

func (s *Server) readLoop(ms *muxSession) {
	buf := make([]byte, 4096)
	for {
		n, err := ms.rw.Read(buf)
		if n > 0 {
			chunk := append([]byte(nil), buf[:n]...)
			ms.ring.Write(chunk)
			out := msgOutput{Type: "output", SessionID: ms.id, DataB64: base64.StdEncoding.EncodeToString(chunk)}
			frame, _ := json.Marshal(out)
			ms.mu.Lock()
			for cl := range ms.subs {
				select {
				case cl.out <- frame:
				default:
					_ = cl.conn.Close() // slow client
				}
			}
			ms.mu.Unlock()
		}
		if err != nil {
			if !errors.Is(err, io.EOF) {
				// log or handle error if desired
			}
			return
		}
	}
}
