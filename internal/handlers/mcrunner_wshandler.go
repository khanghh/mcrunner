package handlers

import (
	"time"

	"github.com/khanghh/mcrunner/internal/core"
	"github.com/khanghh/mcrunner/internal/websocket"
	"github.com/khanghh/mcrunner/pkg/gen"
)

type mcrunnerWSHandler struct {
	mcserver *core.MCServerCmd
	buffer   *core.RingBuffer
}

func (h *mcrunnerWSHandler) WSOnClientConnect(cl *websocket.Client) error {
	msg := NewPTYOutputMessage(h.buffer.Snapshot())
	return cl.SendMessage(msg)
}

func (h *mcrunnerWSHandler) getServerStateMessage() *gen.Message {
	status := h.mcserver.GetStatus()
	pid := 0
	if process := h.mcserver.GetProcess(); process != nil {
		pid = process.Pid
	}
	uptimeSec := int64(0)
	if startTime := h.mcserver.GetStartTime(); startTime != nil {
		uptimeSec = int64(time.Since(*startTime).Seconds())
	}
	usage, err := core.GetServerUsage()
	if err != nil {
		usage = &core.ServerUsage{}
	}
	return NewServerStateMessage(status, pid, 0, uptimeSec, *usage)
}

func (h *mcrunnerWSHandler) broadcastServerState(broadcastCh chan *gen.Message, done chan struct{}) {
	statusCh := make(chan core.ServerStatus)
	h.mcserver.OnStatusChanged(func(status core.ServerStatus) {
		statusCh <- status
	})

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				broadcastCh <- h.getServerStateMessage()
			case <-statusCh:
				broadcastCh <- h.getServerStateMessage()
			case <-done:
				return
			}
		}
	}()
}

func (h *mcrunnerWSHandler) broadcastPTYOutput(broadcastCh chan *gen.Message, done chan struct{}) {
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
		msg := NewPTYOutputMessage(data)
		select {
		case broadcastCh <- msg:
		case <-done:
			return
		}
	}
}

func (h *mcrunnerWSHandler) WSBroadcast(broadcastCh chan *gen.Message, done chan struct{}) {
	go h.broadcastServerState(broadcastCh, done)
	h.broadcastPTYOutput(broadcastCh, done)
}

func (h *mcrunnerWSHandler) WSHandlePTYInput(cl *websocket.Client, msg *gen.Message) error {
	ptyBuffer := msg.GetPtyBuffer()
	if ptyBuffer != nil {
		_, err := h.mcserver.Write(ptyBuffer.Data)
		return err
	}
	return nil
}

func (h *mcrunnerWSHandler) WSHandlePTYResize(cl *websocket.Client, msg *gen.Message) error {
	ptyResize := msg.GetPtyResize()
	if ptyResize != nil {
		rows, cols := int(ptyResize.Rows), int(ptyResize.Cols)
		return h.mcserver.ResizeWindow(rows, cols)
	}
	return nil
}
