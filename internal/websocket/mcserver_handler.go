package websocket

import "github.com/khanghh/mcrunner/internal/ptyproc"

type MCServerHandler struct {
	BaseHandler
	mcserver *ptyproc.PTYSession
}

func (h *MCServerHandler) OnConnect(ctx *Ctx) {
}

func (h *MCServerHandler) Handle(ctx *Ctx, message *Message) error {
	return nil
}

func NewMCServerHandler(mcserver *ptyproc.PTYSession) *MCServerHandler {
	return &MCServerHandler{
		mcserver: mcserver,
	}
}
