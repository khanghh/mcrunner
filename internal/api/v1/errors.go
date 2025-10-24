package api

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/khanghh/mcrunner/internal/core"
)

var (
	ErrNotFound = errors.New("not found")
)

// error helpers
func errorMsg(msg string) fiber.Map { return fiber.Map{"error": msg} }

func badRequest(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(errorMsg(msg))
}

func mapLocalFileServiceError(c *fiber.Ctx, err error) error {
	switch {
	case errorsIs(err, core.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(errorMsg("path not found"))
	case errorsIs(err, core.ErrIsDirectory):
		return c.Status(fiber.StatusBadRequest).JSON(errorMsg("path is a directory"))
	case errorsIs(err, core.ErrNotDirectory):
		return c.Status(fiber.StatusBadRequest).JSON(errorMsg("path is not a directory"))
	case errorsIs(err, core.ErrAlreadyExists):
		return c.Status(fiber.StatusConflict).JSON(errorMsg("conflict: already exists"))
	case errorsIs(err, core.ErrPathTraversal):
		return c.Status(fiber.StatusBadRequest).JSON(errorMsg("invalid path"))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(errorMsg(err.Error()))
	}
}

// errorsIs wraps errors.Is without importing std errors repeatedly
func errorsIs(err, target error) bool {
	return err != nil && target != nil && (err == target || strings.Contains(err.Error(), target.Error()))
}
