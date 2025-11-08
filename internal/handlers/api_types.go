package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const apiVersion = "1.0"

// APIResponse represents a standard API response
type APIResponse struct {
	APIVersion string       `json:"apiVersion,omitempty"`
	Data       interface{}  `json:"data,omitempty"`
	Error      *fiber.Error `json:"error,omitempty"`
}

// ServerStatus represents the current server status
type ServerStatus string

const (
	StatusRunning  ServerStatus = "running"
	StatusStopping ServerStatus = "stopping"
	StatusStopped  ServerStatus = "stopped"
)

// StatusResponse represents the server status response
type StatusResponse struct {
	Status    ServerStatus   `json:"status"`
	PID       int            `json:"pid,omitempty"`
	Uptime    *time.Duration `json:"uptime,omitempty"`
	StartTime *time.Time     `json:"startTime,omitempty"`
}

// CommandRequest represents a command request
type CommandRequest struct {
	Command string `json:"command"`
}

func BadRequestError(msg string) error {
	return fiber.NewError(fiber.StatusBadRequest, msg)
}

func InternalServerError(err error) error {
	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
}

var (
	ErrServerNotRunning     = fiber.NewError(fiber.StatusConflict, "server is not running")
	ErrServerAlreadyRunning = fiber.NewError(fiber.StatusConflict, "server is already running")
)
