package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/khanghh/mcrunner/internal/mcagent"
)

type MCAgentHandler struct {
	mcagent *mcagent.MCAgentBridge
}

func (h *MCAgentHandler) PostAuthLogin(ctx *fiber.Ctx) error {
	proxyURL := fmt.Sprintf("http://localhost:%d/auth/login", h.mcagent.HTTPPort())
	return proxy.Do(ctx, proxyURL)
}

func (h *MCAgentHandler) PostAuthLogout(ctx *fiber.Ctx) error {
	proxyURL := fmt.Sprintf("http://localhost:%d/auth/logout", h.mcagent.HTTPPort())
	return proxy.Do(ctx, proxyURL)
}

func NewMCAgentPluginHandler(mcagent *mcagent.MCAgentBridge) *MCAgentHandler {
	return &MCAgentHandler{
		mcagent: mcagent,
	}
}
