// Package client provides a Go client for the MCRunner API
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ServerStatus string

const (
	StatusRunning  ServerStatus = "running"
	StatusStopping ServerStatus = "stopping"
	StatusStopped  ServerStatus = "stopped"
)

// MCRunnerAPI represents an MCRunner API client
type MCRunnerAPI struct {
	baseURL    string
	httpClient *http.Client
}

// APIResponse represents a standard API response
type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

// APIError represents an API error response
type APIError struct {
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// StatusResponse represents the server status response
type StatusResponse struct {
	Status    ServerStatus `json:"status"`
	Pid       int          `json:"pid,omitempty"`
	Uptime    *Duration    `json:"uptime,omitempty"`
	StartTime *time.Time   `json:"startTime,omitempty"`
}

// CommandRequest represents a command request
type CommandRequest struct {
	Command string `json:"command"`
}

// Duration is a custom duration type for JSON marshaling
type Duration time.Duration

// MarshalJSON implements json.Marshaler for Duration
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// UnmarshalJSON implements json.Unmarshaler for Duration
func (d *Duration) UnmarshalJSON(data []byte) error {
	// Accept either a string duration (e.g., "1h2m3s") or a numeric nanoseconds value
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		duration, err := time.ParseDuration(s)
		if err != nil {
			return err
		}
		*d = Duration(duration)
		return nil
	}
	// Fallback to numeric (assumed nanoseconds)
	var ns int64
	if err := json.Unmarshal(data, &ns); err == nil {
		*d = Duration(time.Duration(ns))
		return nil
	}
	return fmt.Errorf("invalid duration format: %s", string(data))
}

// NewMCRunnerAPI creates a new MCRunner API client
func NewMCRunnerAPI(baseURL string) *MCRunnerAPI {
	return &MCRunnerAPI{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetHTTPClient allows setting a custom HTTP client
func (m *MCRunnerAPI) SetHTTPClient(client *http.Client) {
	m.httpClient = client
}

// GetStatus retrieves the current server status
func (m *MCRunnerAPI) GetStatus(ctx context.Context) (*StatusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", m.baseURL+"/status", nil)
	if err != nil {
		return nil, err
	}
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode into a lightweight struct to avoid re-marshal cycle
	tmp := struct {
		Data  *StatusResponse `json:"data"`
		Error *APIError       `json:"error"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	if tmp.Error != nil {
		return nil, fmt.Errorf("API error (%d): %s", tmp.Error.Status, tmp.Error.Message)
	}
	if tmp.Data == nil {
		return nil, fmt.Errorf("missing status data in response")
	}
	return tmp.Data, nil
}

// StartServer starts the Minecraft server
func (m *MCRunnerAPI) StartServer(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/start", nil)
	if err != nil {
		return err
	}
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if apiResp.Error != nil {
		return fmt.Errorf("API error (%d): %s", apiResp.Error.Status, apiResp.Error.Message)
	}

	return nil
}

// StopServer stops the Minecraft server gracefully
func (m *MCRunnerAPI) StopServer(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/stop", nil)
	if err != nil {
		return err
	}
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if apiResp.Error != nil {
		return fmt.Errorf("API error (%d): %s", apiResp.Error.Status, apiResp.Error.Message)
	}

	return nil
}

// KillServer forcefully kills the Minecraft server
func (m *MCRunnerAPI) KillServer(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/kill", nil)
	if err != nil {
		return err
	}
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if apiResp.Error != nil {
		return fmt.Errorf("API error (%d): %s", apiResp.Error.Status, apiResp.Error.Message)
	}

	return nil
}

// KillServer forcefully kills the Minecraft server
func (m *MCRunnerAPI) Restart(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/restart", nil)
	if err != nil {
		return err
	}
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if apiResp.Error != nil {
		return fmt.Errorf("API error (%d): %s", apiResp.Error.Status, apiResp.Error.Message)
	}

	return nil
}

// SendCommand sends a command to the Minecraft server
func (m *MCRunnerAPI) SendCommand(ctx context.Context, command string) error {
	reqBody := CommandRequest{Command: command}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/command", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if apiResp.Error != nil {
		return fmt.Errorf("API error (%d): %s", apiResp.Error.Status, apiResp.Error.Message)
	}

	return nil
}
