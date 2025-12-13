package mcagent

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	DefaultHTTPPort = 8080
)

type PluginConfig struct {
	HTTPPort   int    `mapstructure:"httpPort"`
	PlayAPIURL string `mapstructure:"playAPIURL"` // https://play.mineviet.com/api
}

func (c *PluginConfig) Sanitize() error {
	if c.HTTPPort == 0 {
		c.HTTPPort = DefaultHTTPPort
	}
	return nil
}

func loadPluginConfig(filename string) (*PluginConfig, error) {
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config PluginConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := config.Sanitize(); err != nil {
		return nil, err
	}
	return &config, nil
}
