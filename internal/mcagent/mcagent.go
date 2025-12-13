package mcagent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/khanghh/mcrunner/pkg/logger"
)

type MCAgentBridge struct {
	configFile string
	config     *PluginConfig
}

func (m *MCAgentBridge) HTTPPort() int {
	return m.config.HTTPPort
}

func IsValidateTicketErr(err error) bool {
	return err == ErrTicketNotFound || err == ErrTicketExpired || err == ErrServiceMismatch
}

// LoginPlayer sends a login request to the MCAgent plugin
// to log in the player with the given user information and ticket.
func (m *MCAgentBridge) LoginPlayer(ctx context.Context, userInfo *UserInfo, playerUuid string, token string, ticket string) error {
	loginURL := fmt.Sprintf("http://localhost:%d/auth/login", m.config.HTTPPort)
	resp, err := http.PostForm(loginURL, url.Values{
		"userId":   {userInfo.UserID},
		"username": {userInfo.Username},
		"fullName": {userInfo.FullName},
		"email":    {userInfo.Email},
		"uuid":     {playerUuid},
		"token":    {token},
		"ticket":   {ticket},
	})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return fmt.Errorf("login failed with status code %d", resp.StatusCode)
	}
	return nil
}

// LogoutPlayer sends a logout request to the MCAgent plugin
// to log out the player with the given username or login ticket.
func (m *MCAgentBridge) LogoutPlayer(ctx context.Context, ticket string, username string) error {
	logoutURL := fmt.Sprintf("http://localhost:%d/auth/logout", m.config.HTTPPort)
	resp, err := http.PostForm(logoutURL, url.Values{
		"ticket":   {ticket},
		"username": {username},
	})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return fmt.Errorf("logout failed with status code %d", resp.StatusCode)
	}
	return nil
}

func (m *MCAgentBridge) GetServerInfo() (*ServerInfo, error) {
	statsURL := fmt.Sprintf("http://localhost:%d/stats", m.config.HTTPPort)
	resp, err := http.Get(statsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get server stats failed with status code %d", resp.StatusCode)
	}
	var stats ServerInfo
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func (m *MCAgentBridge) Reload() error {
	config, err := loadPluginConfig(m.configFile)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load config file %s", m.configFile), "error", err)
		return err
	}
	m.config = config
	return nil
}

func NewMCAgentBridge(configFile string) *MCAgentBridge {
	return &MCAgentBridge{
		configFile: configFile,
	}
}
