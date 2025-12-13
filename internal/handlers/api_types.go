package handlers

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	apiVersion        = "1.0"
	apiRequestTimeout = 30 * time.Second
)

// APIResponse represents a standard API response
type APIResponse struct {
	APIVersion string      `json:"apiVersion,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(code int, message string, reason string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Reason:  reason,
	}
}

// ServerStatus represents the current server status
type ServerStatus string

const (
	StatusRunning  ServerStatus = "running"
	StatusStopping ServerStatus = "stopping"
	StatusStopped  ServerStatus = "stopped"
)

func BadRequestError(msg string) error {
	return fiber.NewError(fiber.StatusBadRequest, msg)
}

func InternalServerError(err error) error {
	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
}

func ParseAPIError(body []byte) (*APIError, error) {
	apiResp := APIResponse{}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}
	return apiResp.Error, nil
}
