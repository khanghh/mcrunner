package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/khanghh/mcrunner/internal/mcagent"
	"github.com/khanghh/mcrunner/pkg/logger"
)

type MCAgentHandler struct {
	mcagent *mcagent.MCAgentBridge
}

func (h *MCAgentHandler) getCallbackURL(ctx *fiber.Ctx) string {
	callbackURL := ctx.Request().URI()
	callbackURL.QueryArgs().Del("ticket")
	return callbackURL.String()
}

func (h *MCAgentHandler) PostLoginCallback(ctx *fiber.Ctx) error {
	if !h.mcagent.IsAuthServer() {
		return fiber.ErrNotFound
	}

	uuid := ctx.FormValue("uuid")
	token := ctx.FormValue("token")
	ticket := ctx.FormValue("ticket")
	callbackURL := h.getCallbackURL(ctx)

	userInfo, err := h.mcagent.ValidateLoginTicket(ctx.Context(), callbackURL, ticket)
	if err != nil {
		if mcagent.IsValidateTicketErr(err) {
			return fiber.ErrUnauthorized
		}
		return fiber.ErrServiceUnavailable
	}

	err = h.mcagent.LoginPlayer(ctx.Context(), userInfo, uuid, token, ticket)
	if err != nil {
		logger.Warn(fmt.Sprintf("Could not login player %s", userInfo.Username), "error", err)
		return fiber.ErrUnauthorized
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *MCAgentHandler) PostLogoutCallback(ctx *fiber.Ctx) error {
	if !h.mcagent.IsAuthServer() {
		return fiber.ErrNotFound
	}
	ticket := ctx.FormValue("ticket")
	username := ctx.FormValue("username")
	return h.mcagent.LogoutPlayer(ctx.Context(), ticket, username)
}

func NewMCAgentPluginHandler(mcagent *mcagent.MCAgentBridge) *MCAgentHandler {
	return &MCAgentHandler{
		mcagent: mcagent,
	}
}
