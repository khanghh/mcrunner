package handlers

import (
	"github.com/khanghh/mcrunner/internal/core"
	"github.com/khanghh/mcrunner/internal/websocket"
	"github.com/khanghh/mcrunner/pkg/gen"
)

type mcrunnerWSHandler struct {
	mcserver *core.MCServerCmd
	buffer   *core.RingBuffer
}

func (h *mcrunnerWSHandler) WSOnClientConnect(cl *websocket.Client) error {
	msg := gen.NewPTYOutputMessage(h.buffer.Snapshot())
	return cl.SendMessage(msg)
}

func (h *mcrunnerWSHandler) WSBroadcast(broadcastCh chan *gen.Message, done chan struct{}) {
	stream := h.mcserver.OutputStream()
	buf := make([]byte, 4096)
	for {
		n, err := stream.Read(buf)
		if err != nil {
			return
		}
		data := make([]byte, n)
		copy(data, buf[:n])

		h.buffer.Write(data)
		msg := gen.NewPTYOutputMessage(data)
		select {
		case broadcastCh <- msg:
		case <-done:
			return
		}
	}
}

func (h *mcrunnerWSHandler) WSHandlePTYInput(cl *websocket.Client, msg *gen.Message) error {
	input := msg.GetPtyBuffer().Data
	_, err := h.mcserver.Write(input)
	return err
}
