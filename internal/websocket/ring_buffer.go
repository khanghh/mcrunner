package websocket

import "sync"

// ringBuffer is a fixed-size ring buffer to keep recent output.
type ringBuffer struct {
	mu    sync.Mutex
	buf   []byte
	cap   int
	start int
	size  int
}

func newRingBuffer(capacity int) *ringBuffer {
	if capacity <= 0 {
		capacity = 1 << 20 // 1 MiB default
	}
	return &ringBuffer{buf: make([]byte, capacity), cap: capacity}
}

func (r *ringBuffer) Write(p []byte) {
	if len(p) == 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(p) >= r.cap {
		copy(r.buf, p[len(p)-r.cap:])
		r.start = 0
		r.size = r.cap
		return
	}
	needDrop := r.size + len(p) - r.cap
	if needDrop > 0 {
		r.start = (r.start + needDrop) % r.cap
		r.size -= needDrop
		if r.size < 0 {
			r.size = 0
		}
	}
	widx := (r.start + r.size) % r.cap
	tail := r.cap - widx
	if len(p) <= tail {
		copy(r.buf[widx:], p)
	} else {
		copy(r.buf[widx:], p[:tail])
		copy(r.buf[0:], p[tail:])
	}
	r.size += len(p)
	if r.size > r.cap {
		r.size = r.cap
	}
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
