package handlers

import (
	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/khanghh/mcrunner/internal/core"
)

type mcrunnerWSHandler struct {
	mcserver *core.MCServerCmd
}

func (h *mcrunnerWSHandler) WebsocketHandler() fiber.Handler {
	return fiberws.New(func(conn *fiberws.Conn) {
		stream := h.mcserver.OutputStream()
		buf := make([]byte, 1024)
		for {
			n, err := stream.Read(buf)
			if n > 0 {
				err := conn.WriteMessage(fiberws.BinaryMessage, buf[:n])
				if err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		}
	})
}
