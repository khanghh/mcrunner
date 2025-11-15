package mccmd

import (
	"io"
)

type outputStream struct {
	pr      *io.PipeReader
	pw      *io.PipeWriter
	done    chan struct{}
	writeCh chan []byte
}

func newOutputStream(buffer int) *outputStream {
	pr, pw := io.Pipe()
	s := &outputStream{
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
			case p, ok := <-s.writeCh:
				if !ok {
					return
				}
				pw.Write(p)
			}
		}
	}()

	return s
}

func (s *outputStream) Write(p []byte) (int, error) {
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

func (s *outputStream) Read(p []byte) (int, error) {
	return s.pr.Read(p)
}

func (s *outputStream) Close() error {
	close(s.done)
	close(s.writeCh)
	return s.pw.Close()
}
