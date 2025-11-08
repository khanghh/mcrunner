package core

import (
	"io"
	"sync"
)

type OutputStream struct {
	mu      sync.Mutex
	pr      *io.PipeReader
	pw      *io.PipeWriter
	done    chan struct{}
	writeCh chan []byte
}

func NewOutputStream(buffer int) *OutputStream {
	pr, pw := io.Pipe()
	s := &OutputStream{
		pr:      pr,
		pw:      pw,
		done:    make(chan struct{}),
		writeCh: make(chan []byte, buffer),
	}

	go func() {
		defer pw.Close()
		for {
			select {
			case <-s.done:
				return
			case p := <-s.writeCh:
				pw.Write(p) // Only one goroutine writes → order preserved
			}
		}
	}()

	return s
}

func (s *OutputStream) Write(p []byte) (int, error) {
	data := make([]byte, len(p))
	copy(data, p)

	select {
	case s.writeCh <- data:
		return len(p), nil
	case <-s.done:
		return 0, io.ErrClosedPipe
	default:
		// Non-blocking: drop if buffer full
		return len(p), nil
	}
}

func (s *OutputStream) Read(p []byte) (int, error) {
	return s.pr.Read(p)
}

func (s *OutputStream) Close() error {
	close(s.done)
	return s.pw.Close()
}
