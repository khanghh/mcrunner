package core

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"syscall"

	"golang.org/x/sys/unix"
)

func ensureFifoExists(fifoPath string) error {
	if _, err := os.Stat(fifoPath); errors.Is(err, os.ErrNotExist) {
		if err := unix.Mkfifo(fifoPath, 0666); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkfifo: %v", err)
		}
	}
	return nil
}

type FifoWriter struct {
	path string
	fd   atomic.Int32 // stores file descriptor (int32)
}

// NewFifoWriter creates a fire-and-forget FIFO writer.
func NewFifoWriter(path string) (*FifoWriter, error) {
	if err := ensureFifoExists(path); err != nil {
		return nil, err
	}
	return &FifoWriter{path: path}, nil
}

// openIfNeeded tries to open the FIFO if not already open.
// Returns true if fd is valid.
func (w *FifoWriter) openIfNeeded() bool {
	if w.fd.Load() != 0 {
		return true // already open
	}

	fd, err := syscall.Open(w.path, syscall.O_WRONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return false
	}

	oldFd := w.fd.Swap(int32(fd))
	if oldFd != 0 {
		syscall.Close(int(oldFd))
	}
	return true
}

// Write writes data non-blockingly. Always returns len(p), nil.
func (w *FifoWriter) Write(p []byte) (int, error) {
	// Fast path: try write without lock
	if w.openIfNeeded() {
		fd := int(w.fd.Load())
		n, err := syscall.Write(fd, p)
		if err == nil {
			return n, nil
		}
		// Reader gone or pipe full → close and discard
		if isBrokenPipe(err) {
			w.fd.Store(0)
			return len(p), nil
		}
	}

	// Slow path: no fd or error → discard
	return len(p), nil
}

func isBrokenPipe(err error) bool {
	var errno syscall.Errno
	return errors.As(err, &errno) &&
		(errno == syscall.EPIPE || errno == syscall.ENXIO || errno == syscall.EAGAIN)
}

// Close closes the FIFO if open.
func (w *FifoWriter) Close() error {
	fd := w.fd.Swap(0)
	if fd != 0 {
		syscall.Close(int(fd))
	}
	return nil
}
