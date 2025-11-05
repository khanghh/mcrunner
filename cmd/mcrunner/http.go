package main

import (
	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/khanghh/mcrunner/internal/ptyproc"
	"github.com/khanghh/mcrunner/internal/websocket"
)

func postCommandHandler(mcserver *ptyproc.PTYSession) fiber.Handler {
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
		ptmx, err := mcserver.PTY()
		if err != nil {
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "minecraft server is not running",
			})
		}
		cmd := []byte(req.Command + "\n")
		if _, err := ptmx.Write(cmd); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return nil
	}
}

func postKillServerHandler(mcserver *ptyproc.PTYSession) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if err := mcserver.Kill(); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func getServerStatusHandler(mcserver *ptyproc.PTYSession) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status": "running",
		})
	}
}

func serveHttp(listenAddr string, mcserver *ptyproc.PTYSession) {
	wsUpgradeRequired := func(ctx *fiber.Ctx) error {
		if !fiberws.IsWebSocketUpgrade(ctx) {
			return fiber.ErrUpgradeRequired
		}
		return ctx.Next()
	}
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	wsServer := websocket.NewServer()
	wsServer.RegisterHandler(NewMCServerHandler(mcserver))
	apiGroup := app.Group("/api")
	apiGroup.Get("/mc/status", getServerStatusHandler(mcserver))
	apiGroup.Post("/mc/command", postCommandHandler(mcserver))
	apiGroup.Post("/mc/kill", postKillServerHandler(mcserver))
	app.Get("/ws", wsUpgradeRequired, wsServer.FiberHandler())
	app.Listen(listenAddr)
}
