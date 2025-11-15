package handlers

import (
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

// ServerState represents the server status response
type ServerState struct {
	Status      ServerStatus `json:"status"`                // current server status
	TPS         float64      `json:"tps"`                   // ticks per second
	PID         int          `json:"pid,omitempty"`         // process ID
	MemoryUsage *uint64      `json:"memoryUsage,omitempty"` // current memory usage
	MemoryLimit *uint64      `json:"memoryLimit,omitempty"` // max allowed memory (0 = unlimited)
	CPUUsage    *float64     `json:"cpuUsage,omitempty"`    // current CPU usage %
	CPULimit    *float64     `json:"cpuLimit,omitempty"`    // max CPUs allowed
	UptimeSec   *int64       `json:"uptimeSec,omitempty"`   // server uptime in seconds
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
