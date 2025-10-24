package core

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	DefaultListenAddr = ":3000"
)

type Config struct {
	Debug      bool   `mapstructure:"debug"`
	RootDir    string `mapstructure:"rootDir"`
	JarFile    string `mapstructure:"jarFile"`
	ListenAddr string `mapstructure:"listenAddr"`
}

func (c *Config) Sanitize() error {
	if c.ListenAddr == "" {
		c.ListenAddr = DefaultListenAddr
	}
	return nil
}

func LoadConfig(filename string) (*Config, error) {
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := config.Sanitize(); err != nil {
		return nil, err
	}
	return &config, nil
}
