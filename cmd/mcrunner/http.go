package main

import (
	"github.com/gofiber/fiber/v2"
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

func serveHttp(listenAddr string, mcserver *core.MCServerCmd) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	apiGroup := app.Group("/api")
	apiGroup.Get("/mc/status", getServerStatusHandler(mcserver))
	apiGroup.Post("/mc/command", postCommandHandler(mcserver))
	apiGroup.Post("/mc/stop", postStopServerHandler(mcserver))
	apiGroup.Post("/mc/kill", postKillServerHandler(mcserver))
	app.Listen(listenAddr)
}
