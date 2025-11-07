package main

import (
	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/khanghh/mcrunner/internal/core"
)

func postCommandHandler(mcserver *core.MCServerCmd) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		type CommandRequest struct {
			Command string `json:"command"`
		}
		var req CommandRequest
		if err := ctx.BodyParser(&req); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if req.Command == "" {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		if err := mcserver.SendCommand(req.Command); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return nil
	}
}

func postStopServerHandler(mcserver *core.MCServerCmd) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if err := mcserver.Stop(); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func postKillServerHandler(mcserver *core.MCServerCmd) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if err := mcserver.Kill(); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func getServerStatusHandler(mcserver *core.MCServerCmd) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status": "running",
		})
	}
}

func streamLogsHandler(mcserver *core.MCServerCmd) fiber.Handler {
	return fiberws.New(func(c *fiberws.Conn) {
		outCh := make(chan []byte, 128)
		go mcserver.StreamOutput(func(dataCh <-chan []byte) {
			defer close(outCh)
			for data := range dataCh {
				select {
				case outCh <- data:
				default:
					return
				}
			}
		})
		for data := range outCh {
			err := c.WriteMessage(fiberws.TextMessage, data)
			if err != nil {
				return
			}
		}
	})
}

func serveHttp(listenAddr string, mcserver *core.MCServerCmd) {
	wsUpgradeRequired := func(ctx *fiber.Ctx) error {
		if !fiberws.IsWebSocketUpgrade(ctx) {
			return fiber.ErrUpgradeRequired
		}
		return ctx.Next()
	}
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Get("/status", getServerStatusHandler(mcserver))
	app.Post("/command", postCommandHandler(mcserver))
	app.Post("/stop", postStopServerHandler(mcserver))
	app.Post("/kill", postKillServerHandler(mcserver))
	app.Get("/logs/stream", wsUpgradeRequired, streamLogsHandler(mcserver))
	app.Listen(listenAddr)
}
