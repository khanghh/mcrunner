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

// ServerState represents the server status response
type ServerState struct {
	Status      ServerStatus `json:"status"`                // current server status
	TPS         float64      `json:"tps"`                   // ticks per second
	PID         int          `json:"pid,omitempty"`         // process ID
	MemoryUsage *uint64      `json:"memoryUsage,omitempty"` // current memory usage in bytes
	MemoryLimit *uint64      `json:"memoryLimit,omitempty"` // max allowed memory in bytes (0 = unlimited)
	CPUUsage    *float64     `json:"cpuUsage,omitempty"`    // current CPU usage percent
	CPULimit    *float64     `json:"cpuLimit,omitempty"`    // max CPUs allowed
	DiskUsage   *uint64      `json:"diskUsage,omitempty"`   // current disk usage in bytes
	DiskSize    *uint64      `json:"diskSize,omitempty"`    // disk size in bytes
	UptimeSec   uint64       `json:"uptimeSec,omitempty"`   // server uptime in seconds
}

type CommandRequest struct {
	Command string `json:"command"`
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

// GetServerState retrieves the current server status
func (m *MCRunnerAPI) GetServerState(ctx context.Context) (*ServerState, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", m.baseURL+"/state", nil)
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
		Data  *ServerState `json:"data"`
		Error *APIError    `json:"error"`
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
