package handlers

import (
	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/khanghh/mcrunner/internal/core"
)

type mcrunnerWS struct {
	mcserver *core.MCServerCmd
}

func (h *mcrunnerWS) streamLoop(conn *fiberws.Conn) {
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
}

func (h *mcrunnerWS) WebsocketHandler() fiber.Handler {
	return fiberws.New(func(conn *fiberws.Conn) {
		go h.streamLoop(conn)
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			_, err = h.mcserver.Write(data)
			if err != nil {
				return
			}
		}
	})
}
