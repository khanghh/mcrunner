package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/khanghh/mcrunner/internal/core"
)

type MCRunnerHandler struct {
	mcserver *core.MCServerCmd
	*mcrunnerWSHandler
}

func (h *MCRunnerHandler) PostCommand(ctx *fiber.Ctx) error {
	var req CommandRequest
	if err := ctx.BodyParser(&req); err != nil {
		return BadRequestError("invalid request payload")
	}
	if req.Command == "" {
		return BadRequestError("missing command")
	}
	if err := h.mcserver.SendCommand(req.Command); err != nil {
		return InternalServerError(err)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) PostStartServer(ctx *fiber.Ctx) error {
	if h.mcserver.GetStatus() == core.StateRunning {
		return ErrServerAlreadyRunning
	}
	h.buffer.Reset()
	if err := h.mcserver.Start(); err != nil {
		return InternalServerError(err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) PostStopServer(ctx *fiber.Ctx) error {
	if h.mcserver.GetStatus() != core.StateRunning {
		return ErrServerNotRunning
	}
	if err := h.mcserver.Stop(); err != nil {
		return InternalServerError(err)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) PostRestartServer(ctx *fiber.Ctx) error {
	if h.mcserver.GetStatus() != core.StateRunning {
		return ErrServerNotRunning
	}
	if err := h.mcserver.Stop(); err != nil {
		return InternalServerError(err)
	}
	if err := h.mcserver.Start(); err != nil {
		return InternalServerError(err)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) PostKillServer(ctx *fiber.Ctx) error {
	if h.mcserver.GetStatus() != core.StateRunning {
		return ErrServerNotRunning
	}
	if err := h.mcserver.Kill(); err != nil {
		return InternalServerError(err)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) GetStatus(ctx *fiber.Ctx) error {
	status := h.mcserver.GetStatus()
	startTime := h.mcserver.GetStartTime()

	var pid int
	var uptime *time.Duration

	if status == core.StateRunning && h.mcserver.GetProcess() != nil {
		pid = h.mcserver.GetProcess().Pid
		if startTime != nil {
			uptimeDuration := time.Since(*startTime)
			uptime = &uptimeDuration
		}
	}

	response := StatusResponse{
		Status:    ServerStatus(status),
		PID:       pid,
		Uptime:    uptime,
		StartTime: startTime,
	}

	return ctx.JSON(APIResponse{
		Data: response,
	})
}

func NewMCRunnerHandler(mcserver *core.MCServerCmd) *MCRunnerHandler {

	return &MCRunnerHandler{
		mcserver: mcserver,
		mcrunnerWSHandler: &mcrunnerWSHandler{
			mcserver: mcserver,
			buffer:   core.NewRingBuffer(1 << 20),
		},
	}
}
