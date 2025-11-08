package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return ctx.Status(fiberErr.Code).JSON(APIResponse{
			APIVersion: apiVersion,
			Error:      fiberErr,
		})
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(APIResponse{
		APIVersion: apiVersion,
		Error: &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		},
	})
}
