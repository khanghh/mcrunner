package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/khanghh/mcrunner/internal/mcagent"
	"github.com/khanghh/mcrunner/internal/mccmd"
	"github.com/khanghh/mcrunner/internal/sysmetrics"
	"github.com/khanghh/mcrunner/pkg/api"
)

var (
	ErrServerNotRunning     = fiber.NewError(fiber.StatusConflict, "server is not running")
	ErrServerAlreadyRunning = fiber.NewError(fiber.StatusConflict, "server is already running")
)

type MCRunnerHandler struct {
	mcserver *mccmd.MCServerCmd     // server process
	mcagent  *mcagent.MCAgentBridge // plugin bridge
}

func (h *MCRunnerHandler) getServerState() api.ServerState {
	status := h.mcserver.GetStatus()

	var serverIPAddr string
	if ipAddr, err := sysmetrics.GetOutboundIP(); err == nil {
		serverIPAddr = ipAddr.String()
	}

	serverState := api.ServerState{
		Status:    api.ServerStatus(status),
		IPAddress: serverIPAddr,
	}

	usage := sysmetrics.GetResourceUsage()
	serverState.MemoryUsage = &usage.MemoryUsage
	serverState.MemoryLimit = &usage.MemoryLimit
	serverState.CPUUsage = &usage.CPUUsage
	serverState.CPULimit = &usage.CPULimit
	serverState.DiskUsage = &usage.DiskUsage
	serverState.DiskSize = &usage.DiskSize

	process := h.mcserver.GetProcess()
	if process == nil {
		return serverState
	}
	serverState.PID = process.Pid
	if startTime := h.mcserver.GetStartTime(); startTime != nil {
		serverState.UptimeSec = uint64(time.Since(*startTime).Seconds())
	}

	if serverInfo, err := h.mcagent.GetServerInfo(); err == nil {
		serverState.Server = &api.ServerInfo{
			Name:          serverInfo.Name,
			Version:       serverInfo.Version,
			TPS:           serverInfo.TPS,
			PlayersOnline: serverInfo.PlayersOnline,
			PlayersMax:    serverInfo.PlayersMax,
		}
	}

	return serverState
}

func (h *MCRunnerHandler) runWithTimeout(ctx *fiber.Ctx, fn func() error) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- fn()
	}()
	execCtx, cancel := context.WithTimeout(ctx.Context(), apiRequestTimeout)
	defer cancel()

	select {
	case err := <-errCh:
		return err
	case <-execCtx.Done():
		return fiber.NewError(fiber.StatusRequestTimeout, "request timed out")
	}
}

func (h *MCRunnerHandler) PostCommand(ctx *fiber.Ctx) error {
	var req api.CommandRequest
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
	if h.mcserver.GetStatus() == mccmd.StatusRunning {
		return ErrServerAlreadyRunning
	}

	err := h.runWithTimeout(ctx, func() error {
		if err := h.mcserver.Start(); err != nil {
			return InternalServerError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) PostStopServer(ctx *fiber.Ctx) error {
	if h.mcserver.GetStatus() != mccmd.StatusRunning {
		return ErrServerNotRunning
	}

	err := h.runWithTimeout(ctx, func() error {
		if err := h.mcserver.Stop(); err != nil {
			return InternalServerError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) PostRestartServer(ctx *fiber.Ctx) error {
	if h.mcserver.GetStatus() != mccmd.StatusRunning {
		return ErrServerNotRunning
	}

	err := h.runWithTimeout(ctx, func() error {
		if err := h.mcserver.Stop(); err != nil {
			return InternalServerError(err)
		}
		if err := h.mcserver.Start(); err != nil {
			return InternalServerError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) PostKillServer(ctx *fiber.Ctx) error {
	if h.mcserver.GetStatus() == mccmd.StatusStopped {
		return ErrServerNotRunning
	}

	err := h.runWithTimeout(ctx, func() error {
		if err := h.mcserver.Kill(); err != nil {
			return InternalServerError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCRunnerHandler) GetState(ctx *fiber.Ctx) error {
	return ctx.JSON(APIResponse{
		Data: h.getServerState(),
	})
}

func NewMCRunnerHandler(mcserver *mccmd.MCServerCmd) *MCRunnerHandler {
	return &MCRunnerHandler{
		mcserver: mcserver,
	}
}
