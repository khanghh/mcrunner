package websocket

import "sync"

type Topic struct {
	id     string
	subs   map[*Client]struct{}
	subsMu sync.RWMutex
}

func (t *Topic) AddSubscriber(c *Client) {
	t.subsMu.Lock()
	defer t.subsMu.Unlock()
	t.subs[c] = struct{}{}
}

func (t *Topic) RemoveSubscriber(c *Client) {
	t.subsMu.Lock()
	defer t.subsMu.Unlock()
	delete(t.subs, c)
}

func (t *Topic) Broadcast(data []byte) {
	t.subsMu.Lock()
	defer t.subsMu.Unlock()
	for cl := range t.subs {
		cl.Send(data)
	}
}
