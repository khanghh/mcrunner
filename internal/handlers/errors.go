package handlers

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/khanghh/mcrunner/internal/file"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

// Helper functions
func mapLocalFileServiceError(c *fiber.Ctx, err error) error {
	if os.IsNotExist(err) || errors.Is(err, file.ErrFileNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(JSONErrFileNotFound)
	}
	if os.IsPermission(err) {
		return c.Status(fiber.StatusForbidden).JSON(JSONErrNoPermissions)
	}
	if errors.Is(err, file.ErrDirNotEmpty) {
		return c.Status(fiber.StatusBadRequest).JSON(JSONErrDirectoryNotEmpty)
	}
	return c.Status(fiber.StatusInternalServerError).JSON(errorMsg(err.Error()))
}

func badRequest(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
}

func errorMsg(msg string) fiber.Map {
	return fiber.Map{"error": msg}
}
