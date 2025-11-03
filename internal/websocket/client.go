package websocket

import (
	"net"
	"sync"

	"github.com/khanghh/mcrunner/internal/params"
)

type Client struct {
	conn   net.Conn
	out    chan []byte
	closed chan struct{}
	mu     sync.Mutex // protects the closed state
}

// NewClient creates a new websocket client with initialized channels.
func NewClient(conn net.Conn) *Client {
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
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.IsAlive() {
		return ErrClientDisconnected
	}

	select {
	case c.out <- frame:
		return nil
	default:
		close(c.closed)
		_ = c.conn.Close() // slow client protection
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
			if _, err := c.conn.Write(frame); err != nil {
				c.Close()
				return
			}
		}
	}
}

// Close closes the client connection and cleans up resources.
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if already closed
	select {
	case <-c.closed:
		return // already closed
	default:
	}

	close(c.closed)
	_ = c.conn.Close()
}
