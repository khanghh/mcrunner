package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khanghh/mcrunner/internal/core"
)

type MCRunnerHandler struct {
	mcserver *core.MCServerCmd
	*mcrunnerWSHandler
}

func (h *MCRunnerHandler) PostCommand(ctx *fiber.Ctx) error {
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

	if err := h.mcserver.SendCommand(req.Command); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return nil
}

func (h *MCRunnerHandler) PostStartServer(ctx *fiber.Ctx) error {
	if err := h.mcserver.Start(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) PostStopServer(ctx *fiber.Ctx) error {
	if err := h.mcserver.Stop(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) PostKillServer(ctx *fiber.Ctx) error {
	if err := h.mcserver.Kill(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) GetStatus(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status": "running",
	})
}

func NewMCRunnerHandler(mcserver *core.MCServerCmd) *MCRunnerHandler {
	return &MCRunnerHandler{
		mcserver:          mcserver,
		mcrunnerWSHandler: &mcrunnerWSHandler{mcserver: mcserver},
	}
}
