package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/khanghh/mcrunner/internal/core"
	"github.com/khanghh/mcrunner/pkg/api"
)

type MCRunnerHandler struct {
	mcserver *core.MCServerCmd
	*mcrunnerWSHandler
}

func (h *mcrunnerWSHandler) getServerState() api.ServerState {
	status := h.mcserver.GetStatus()
	serverState := api.ServerState{
		Status: api.ServerStatus(status),
	}
	usage, err := core.GetServerUsage()
	if err != nil {
		return serverState
	}
	serverState.MemoryUsage = &usage.MemoryUsage
	serverState.MemoryLimit = &usage.MemoryLimit
	serverState.CPUUsage = &usage.CPUUsage
	serverState.CPULimit = &usage.CPULimit

	process := h.mcserver.GetProcess()
	if process == nil {
		return serverState
	}
	if startTime := h.mcserver.GetStartTime(); startTime != nil {
		serverState.UptimeSec = uint64(time.Now().Sub(*startTime).Seconds())
	}
	serverState.PID = process.Pid
	return serverState
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
	if h.mcserver.GetStatus() == core.StatusRunning {
		return ErrServerAlreadyRunning
	}
	h.buffer.Reset()
	if err := h.mcserver.Start(); err != nil {
		return InternalServerError(err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) PostStopServer(ctx *fiber.Ctx) error {
	if h.mcserver.GetStatus() != core.StatusRunning {
		return ErrServerNotRunning
	}
	if err := h.mcserver.Stop(); err != nil {
		return InternalServerError(err)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) PostRestartServer(ctx *fiber.Ctx) error {
	if h.mcserver.GetStatus() != core.StatusRunning {
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
	if h.mcserver.GetStatus() != core.StatusRunning {
		return ErrServerNotRunning
	}
	if err := h.mcserver.Kill(); err != nil {
		return InternalServerError(err)
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) GetState(ctx *fiber.Ctx) error {
	return ctx.JSON(APIResponse{
		Data: h.getServerState(),
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
