package core

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

// UnixLogWriter writes logs to a Unix socket.
// If no client is connected, writes are silently ignored (non-blocking).
type UnixLogWriter struct {
	sockPath  string
	listener  net.Listener
	output    chan []byte // buffered channel for output
	stop      chan struct{}
	connected atomic.Bool
}

func NewUnixLogWriter(sockPath string) (*UnixLogWriter, error) {
	if err := os.MkdirAll(filepath.Dir(sockPath), 0755); err != nil {
		return nil, err
	}

	// Reuse existing socket file if stale
	_ = os.Remove(sockPath)
	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		return nil, err
	}
	_ = os.Chmod(sockPath, 0600)

	w := &UnixLogWriter{
		sockPath: sockPath,
		listener: listener,
		stop:     make(chan struct{}),
	}

	// Accept clients in background
	go w.acceptLoop()

	return w, nil
}

func (w *UnixLogWriter) acceptLoop() {
	for {
		select {
		case <-w.stop:
			return
		default:
		}

		conn, err := w.listener.Accept()
		if err != nil {
			log.Println("UnixLogWriter accept error:", err)
			continue
		}

		w.output = make(chan []byte, 128)
		w.connected.Store(true)
		for data := range w.output {
			conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
			_, err := conn.Write(data)
			conn.SetWriteDeadline(time.Time{})
			if err != nil {
				conn.Close()
				break
			}
		}
		w.connected.Store(false)
	}
}

// Write implements io.Writer.
// If no client → silently drop.
func (w *UnixLogWriter) Write(p []byte) (int, error) {
	if !w.connected.Load() {
		return len(p), nil
	}
	select {
	case w.output <- p:
	default:
		close(w.output)
		w.connected.Store(false)
	}
	return len(p), nil
}

// Close shuts down the listener and stops the accept loop.
func (w *UnixLogWriter) Close() error {
	close(w.stop)
	close(w.output)
	return w.listener.Close()
}

func isAddrInUse(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Err.Error() == "address already in use"
	}
	return false
}
