package api

import (
	"github.com/gofiber/fiber/v2"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct{}

// SendJSON sends a successful JSON response
func (h *BaseHandler) SendJSON(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"apiVersion": "1.0",
		"data":       data,
	})
}

// SendError sends an error response
func (h *BaseHandler) SendError(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{
		"apiVersion": "1.0",
		"error": fiber.Map{
			"code":    code,
			"message": msg,
			"status":  fiber.ErrInternalServerError.Error(),
		},
	})
}
