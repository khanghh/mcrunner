package websocket

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/khanghh/mcrunner/pkg/gen"
	"google.golang.org/protobuf/proto"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingInterval   = (pongWait * 9) / 10 // 90% of pongWait
	maxMessageSize = 1 << 20             // 1MB
)

type Client struct {
	conn   *fiberws.Conn
	out    chan []byte
	server *Server
	mu     sync.Mutex
	closed chan struct{}
}

func (c *Client) readPump() {
	defer c.Close()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		msgType, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v", err)
			}
			break
		}
		if msgType != websocket.BinaryMessage {
			// we expect binary protobuf messages
			continue
		}
		var msg gen.Message
		if err := proto.Unmarshal(data, &msg); err != nil {
			log.Printf("proto unmarshal: %v", err)
			continue
		}

		handler, ok := c.server.handlers[msg.Type]
		if ok {
			if err := handler(c, &msg); err != nil {
				log.Printf("handler message %d: %v", msg.Type, err)
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case msg, ok := <-c.out:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// write as binary frame
			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(msg); err != nil {
				_ = w.Close()
				return
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-c.closed:
			return
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) SendMessage(msg *gen.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return c.send(data)
}

func (c *Client) send(data []byte) error {
	select {
	case c.out <- data:
		return nil
	case <-c.closed:
		return errors.New("disconnected")
	default:
		return errors.New("send buffer full")
	}
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if already closed
	select {
	case <-c.closed:
		return // already closed
	default:
		close(c.closed)
	}
	c.conn.Close()
}
