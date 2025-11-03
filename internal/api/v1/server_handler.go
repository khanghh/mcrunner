package api

import (
	"github.com/gofiber/fiber/v2"
)

// ServerHandler handles server lifecycle management requests
type ServerHandler struct {
	*BaseHandler
	mcserver MCServer
	Status   string
}

// NewServerHandler creates a new ServerHandler instance
func NewServerHandler(mcserver MCServer) *ServerHandler {
	return &ServerHandler{
		BaseHandler: &BaseHandler{},
		mcserver:    mcserver,
		Status:      "stopped",
	}
}

// GetStatus returns the current server status
func (h *ServerHandler) GetStatus(c *fiber.Ctx) error {
	return h.SendJSON(c, fiber.Map{
		"status": h.Status,
	})
}

// PostStart starts the server
func (h *ServerHandler) PostStart(c *fiber.Ctx) error {
	h.Status = "running"
	return c.SendStatus(fiber.StatusOK)
}

// PostStop stops the server
func (h *ServerHandler) PostStop(c *fiber.Ctx) error {
	h.Status = "stopped"
	return c.SendStatus(fiber.StatusOK)
}

// PostRestart restarts the server
func (h *ServerHandler) PostRestart(c *fiber.Ctx) error {
	h.Status = "restarting"
	// Add restart logic here
	h.Status = "running"
	return c.SendStatus(fiber.StatusOK)
}

// PostCommand sends a command to the server
func (h *ServerHandler) PostCommand(c *fiber.Ctx) error {
	var body struct {
		Command string `json:"command"`
	}
	if err := c.BodyParser(&body); err != nil {
		return h.SendError(c, 400, "invalid body")
	}
	if err := h.mcserver.SendCommand(body.Command); err != nil {
		return h.SendError(c, 500, "failed to send command")
	}
	return h.SendJSON(c, fiber.Map{"status": "sent"})
}
