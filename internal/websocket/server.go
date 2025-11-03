package websocket

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"google.golang.org/protobuf/proto"
)

// Server multiplexes multiple PTY sessions over a single websocket connection.
// One reader goroutine per PTY session, fanout to subscribed clients.
type Server struct {
	handlers []Handler
	topics   map[string]*Topic
	mu       sync.RWMutex
	connWg   sync.WaitGroup
	stopCh   chan struct{}
}

// NewServer constructs a multiplex Server bound to a PTYManager.
func NewServer() *Server {
	return &Server{
		topics: make(map[string]*Topic),
		stopCh: make(chan struct{}),
	}
}

func (s *Server) onClientConnect(ctx *Ctx) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, handler := range s.handlers {
		ctx.wg.Add(1)
		go func() {
			handler.OnConnect(ctx)
			ctx.wg.Done()
		}()
	}
}

func (s *Server) FiberHandler() fiber.Handler {
	return fiberws.New(func(conn *fiberws.Conn) {
		s.connWg.Add(1)
		cl := NewClient(conn)
		ctx, cancel := context.WithCancel(context.Background())
		connCtx := &Ctx{
			context: ctx,
			server:  s,
			client:  cl,
		}

		s.onClientConnect(connCtx)

		// cleanup on exit
		defer func() {
			s.mu.Lock()
			topics := make([]*Topic, 0, len(s.topics))
			for _, topic := range s.topics {
				topics = append(topics, topic)
			}
			s.mu.Unlock()

			for _, topic := range topics {
				topic.subsMu.Lock()
				delete(topic.subs, cl)
				topic.subsMu.Unlock()
			}

			cancel()
			connCtx.wg.Wait()
			cl.Close()
			s.connWg.Done()
		}()

		// server shutdown handling
		go func() {
			<-s.stopCh
			cl.Close()
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

			var msg Message
			if err := proto.Unmarshal(data, &msg); err != nil {
				// client sent bad message
				slog.Error("client sent bad message", "blob", string(data), "error", err)
				continue
			}

			if err := s.handleMessage(connCtx, &msg); err != nil {
				slog.Error("Could not handle message", "type", msg.Type, "error", err)
			}
		}
	})
}

func (s *Server) handleMessage(ctx *Ctx, msg *Message) error {
	for _, handler := range s.handlers {
		if err := handler.Handle(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) sendError(cl *Client, msg string) error {
	if !cl.IsAlive() {
		return errors.New("client disconnected")
	}
	data, _ := proto.Marshal(&Message{
		Type:  MessageType_ERROR,
		Error: msg,
	})
	cl.Send(data)
	return nil
}

func (s *Server) clientSubscribe(cl *Client, topicName string) error {
	s.mu.RLock()
	topic, ok := s.topics[topicName]
	s.mu.RUnlock()
	if !ok {
		return ErrTopicNotFound
	}
	topic.AddSubscriber(cl)
	return nil
}

func (s *Server) clientUnsubscribe(cl *Client, topicName string) error {
	s.mu.RLock()
	topic, ok := s.topics[topicName]
	s.mu.RUnlock()
	if !ok {
		return ErrTopicNotFound
	}
	topic.RemoveSubscriber(cl)
	return nil
}

func (s *Server) getOrCreateTopic(topicName string) *Topic {
	s.mu.Lock()
	defer s.mu.Unlock()

	topic, ok := s.topics[topicName]
	if ok {
		return topic
	}

	s.topics[topicName] = &Topic{
		id:   topicName,
		subs: make(map[*Client]struct{}),
	}
	return s.topics[topicName]
}

func (s *Server) Broadcast(topicName string, message *Message) {
	data, err := proto.Marshal(message)
	if err != nil {
		slog.Error("Could not marshal broadcast message", "error", err)
		return
	}
	s.getOrCreateTopic(topicName).Broadcast(data)
}

// RegisterHandler registers a single handler for the given message type.
func (s *Server) RegisterHandler(handler Handler) {
	s.handlers = append(s.handlers, handler)
}

func (s *Server) Shutdown() error {
	close(s.stopCh)
	s.connWg.Wait()
	return nil
}
