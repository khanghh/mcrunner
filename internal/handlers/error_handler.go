package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if apiErr, ok := err.(*APIError); ok {
		return ctx.Status(apiErr.Code).JSON(APIResponse{
			APIVersion: apiVersion,
			Error:      apiErr,
		})
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return ctx.Status(fiberErr.Code).JSON(APIResponse{
			APIVersion: apiVersion,
			Error: &APIError{
				Code:    fiberErr.Code,
				Message: fiberErr.Message,
			},
		})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(APIResponse{
		APIVersion: apiVersion,
		Error: &APIError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		},
	})
}
