package handlers

import (
	"github.com/khanghh/mcrunner/internal/websocket"
	"github.com/khanghh/mcrunner/pkg/gen"
)

func (h *MCRunnerHandler) WSOnClientConnect(cl *websocket.Client) error {
	return nil
}

func (h *MCRunnerHandler) WSOnClientDisconnect(cl *websocket.Client) error {
	return nil
}

func (h *MCRunnerHandler) WSOnServerShutdown(s *websocket.Server) error {
	return nil
}

func (h *MCRunnerHandler) WSBroadcast(broadcastCh chan *gen.Message, done chan struct{}) {
	stream := h.mcserver.OutputStream()
	buf := make([]byte, 4096)
	for {
		n, err := stream.Read(buf)
		if err != nil {
			return
		}
		data := make([]byte, n)
		copy(data, buf[:n])
		msg := gen.NewPTYBufferMessage(data)
		select {
		case broadcastCh <- msg:
		case <-done:
			return
		}
	}
}

func (h *MCRunnerHandler) WSHandlePTYInput(cl *websocket.Client, msg *gen.Message) error {
	return nil
}

func (h *MCRunnerHandler) WSHandlePTYResize(cl *websocket.Client, msg *gen.Message) error {
	return nil
}
