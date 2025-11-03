package websocket

import (
	"sync"

	fiberws "github.com/gofiber/websocket/v2"
	"github.com/khanghh/mcrunner/internal/params"
)

type Client struct {
	conn   *fiberws.Conn
	out    chan []byte
	closed chan struct{}
	mu     sync.Mutex // protects the closed state
}

// NewClient creates a new websocket client with initialized channels.
func NewClient(conn *fiberws.Conn) *Client {
	cl := &Client{
		conn:   conn,
		out:    make(chan []byte, params.WSClientQueueSize),
		closed: make(chan struct{}),
	}
	go cl.writeLoop()
	return cl
}

func (c *Client) IsAlive() bool {
	select {
	case <-c.closed:
		return false
	default:
		return true
	}
}

// Send enqueues a frame to be sent to the client.
// If the client is slow and the queue is full, the connection will be closed.
func (c *Client) Send(frame []byte) error {
	if !c.IsAlive() {
		return ErrClientDisconnected
	}

	select {
	case c.out <- frame:
		return nil
	default:
		c.Close()
		return ErrClientDisconnected
	}
}

// writeLoop starts the writer goroutine that sends queued messages to the client.
func (c *Client) writeLoop() {
	defer close(c.out) // close out channel when exiting

	for {
		select {
		case <-c.closed:
			return
		case frame, ok := <-c.out:
			if !ok {
				return
			}
			if err := c.conn.WriteMessage(fiberws.BinaryMessage, frame); err != nil {
				c.Close()
				return
			}
		}
	}
}

// Close closes the client connection and cleans up s.
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
}
