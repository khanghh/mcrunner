package ptyproc

import "sync"

// ringBuffer is a fixed-size ring buffer to keep recent PTY output per session.
type ringBuffer struct {
	mu    sync.Mutex
	buf   []byte
	cap   int
	start int
	size  int
}

func newRingBuffer(capacity int) *ringBuffer {
	if capacity <= 0 {
		capacity = 1 << 20 // default 1 MiB
	}
	return &ringBuffer{buf: make([]byte, capacity), cap: capacity}
}

func (r *ringBuffer) Write(p []byte) (int, error) {
	n := len(p)
	if n == 0 {
		return 0, nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if n >= r.cap {
		copy(r.buf, p[n-r.cap:])
		r.start = 0
		r.size = r.cap
		return n, nil
	}

	if r.size+n > r.cap {
		over := r.size + n - r.cap
		r.start = (r.start + over) % r.cap
		r.size -= over
	}

	widx := (r.start + r.size) % r.cap
	tail := r.cap - widx
	if n <= tail {
		copy(r.buf[widx:], p)
	} else {
		copy(r.buf[widx:], p[:tail])
		copy(r.buf[0:], p[tail:])
	}
	r.size += n
	return n, nil
}

func (r *ringBuffer) Snapshot() []byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]byte, r.size)
	if r.size == 0 {
		return out
	}
	tail := r.cap - r.start
	if r.size <= tail {
		copy(out, r.buf[r.start:r.start+r.size])
	} else {
		copy(out, r.buf[r.start:])
		copy(out[tail:], r.buf[:r.size-tail])
	}
	return out
}
