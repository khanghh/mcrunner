package main

import (
	"io"

	"github.com/khanghh/mcrunner/internal/ptyproc"
	"github.com/khanghh/mcrunner/internal/websocket"
)

type MCServerHandler struct {
	websocket.BaseHandler
	session *ptyproc.PTYSession
}

func (h *MCServerHandler) attachPTY(ctx *websocket.Ctx, session *ptyproc.PTYSession) {
	pr, pw := io.Pipe()
	go func() {
		session.Attach(nil, pw)
	}()
	defer pw.Close()
	for {
		buf := make([]byte, 4096)
		n, err := pr.Read(buf)
		if n > 0 {
			msg := websocket.NewPTYBufferMessage(session.Name(), buf[:n])
			if err := ctx.SendMessage(msg); err != nil {
				return
			}
		}
		if err != nil {
			break
		}
	}
}

func (h *MCServerHandler) OnConnect(ctx *websocket.Ctx) {
	h.attachPTY(ctx, h.session)
}

func (h *MCServerHandler) handlePTYInput(ctx *websocket.Ctx, msg *websocket.Message) error {
	ptmx, err := h.session.PTY()
	if err != nil {
		ctx.SendError("Minecraft server is not running")
		return nil
	}

	ptyInput := msg.GetPtyInput()
	inputData := ptyInput.GetData()
	if _, err = ptmx.Write(inputData); err != nil {
		ctx.SendError("Failed to send input to Minecraft server")
	}
	return nil
}

func (h *MCServerHandler) handlePTYResize(ctx *websocket.Ctx, msg *websocket.Message) error {
	ptyResize := msg.GetPtyResize()
	rows, cols := int(ptyResize.Rows), int(ptyResize.Cols)
	return h.session.Resize(rows, cols)
}

func (h *MCServerHandler) Handle(ctx *websocket.Ctx, msg *websocket.Message) error {
	switch msg.GetType() {
	case websocket.MessageType_PTY_INPUT:
		return h.handlePTYInput(ctx, msg)
	case websocket.MessageType_PTY_RESIZE:
		return h.handlePTYResize(ctx, msg)
	default:
		ctx.SendError("Unknown message type")
	}
	return nil
}

func NewMCServerHandler(mcserver *ptyproc.PTYSession) *MCServerHandler {
	return &MCServerHandler{
		session: mcserver,
	}
}
