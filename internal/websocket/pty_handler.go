package websocket

import (
	"io"

	"github.com/khanghh/mcrunner/internal/ptyproc"
)

type PTYHandler struct {
	BaseHandler
	ptyManager *ptyproc.PTYManager
}

func (h *PTYHandler) attachPTY(ctx *Ctx, session *ptyproc.PTYSession) {
	pr, pw := io.Pipe()
	go func() {
		session.Attach(nil, pw)
	}()
	defer pw.Close()
	for {
		buf := make([]byte, 4096)
		n, err := pr.Read(buf)
		if n > 0 {
			msg := NewPTYBufferMessage(session.Name(), buf[:n])
			if err := ctx.SendMessage(msg); err != nil {
				return
			}
		}
		if err != nil {
			break
		}
	}
}

func (h *PTYHandler) OnConnect(ctx *Ctx) {
	mcserver, exsist := h.ptyManager.Get("mcserver")
	if !exsist {
		ctx.Disconnect("Minecraft server is not running")
		return
	}
	h.attachPTY(ctx, mcserver)
}

func (h *PTYHandler) Handle(ctx *Ctx, msg *Message) error {
	return nil
}

func NewPTYHandler(ptyManager *ptyproc.PTYManager) *PTYHandler {
	return &PTYHandler{
		ptyManager: ptyManager,
	}
}
