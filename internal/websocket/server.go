package websocket

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/khanghh/mcrunner/pkg/gen"
	"google.golang.org/protobuf/proto"
)

type Broadcaster func(broadcastCh chan *gen.Message, done chan struct{})

type HandleFunc func(cl *Client, msg *gen.Message) error

type Server struct {
	clients      map[*Client]struct{}
	handlers     map[gen.MessageType]HandleFunc
	onConnect    []func(cl *Client) error
	onDisconnect []func(cl *Client) error
	broadcast    chan *gen.Message
	register     chan *Client
	unregister   chan *Client
	shutdown     chan struct{}
}

func (s *Server) OnConnect(handler func(cl *Client) error) {
	s.onConnect = append(s.onConnect, handler)
}

func (s *Server) OnDisconnect(handler func(cl *Client) error) {
	s.onDisconnect = append(s.onDisconnect, handler)
}

func (s *Server) OnMessage(msgtype gen.MessageType, handler HandleFunc) {
	s.handlers[msgtype] = handler
}

func (s *Server) StartBroadcast(broadcaster Broadcaster) {
	go broadcaster(s.broadcast, s.shutdown)
}

func (s *Server) Broadcast(msg *gen.Message) {
	s.broadcast <- msg
}

func (s *Server) ServeFiberWS() fiber.Handler {
	return fiberws.New(func(conn *fiberws.Conn) {
		cl := &Client{
			conn:   conn,
			out:    make(chan []byte, 256),
			server: s,
			closed: make(chan struct{}),
		}
		s.register <- cl
		go cl.readPump()
		cl.writePump()
		s.unregister <- cl
		fmt.Println("client disconnected")
	})
}

func (s *Server) loop() {
	defer func() {
		for client := range s.clients {
			for _, handler := range s.onDisconnect {
				if err := handler(client); err != nil {
					fmt.Println("disconnect:", err)
				}
			}
			client.Close()
		}
	}()

	for {
		select {
		case <-s.shutdown:
			return
		case msg := <-s.broadcast:
			data, err := proto.Marshal(msg)
			if err != nil {
				fmt.Println("proto marshal:", err)
				continue
			}
			for c := range s.clients {
				if err := c.send(data); err != nil {
					fmt.Println("send:", err)
				}
			}
		case cl := <-s.register:
			s.clients[cl] = struct{}{}
			for _, handler := range s.onConnect {
				if err := handler(cl); err != nil {
					fmt.Println("connect:", err)
				}
			}
		case cl := <-s.unregister:
			for _, handler := range s.onDisconnect {
				if err := handler(cl); err != nil {
					fmt.Println("disconnect:", err)
				}
			}
			cl.Close()
			delete(s.clients, cl)
		}
	}

}

func (s *Server) Shutdown() error {
	select {
	case <-s.shutdown:
	default:
		close(s.shutdown)
	}
	return nil
}

func (s *Server) Done() <-chan struct{} {
	return s.shutdown
}

func NewServer() *Server {
	s := &Server{
		clients:    map[*Client]struct{}{},
		handlers:   make(map[gen.MessageType]HandleFunc),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *gen.Message),
		shutdown:   make(chan struct{}),
	}

	go s.loop()
	return s
}
