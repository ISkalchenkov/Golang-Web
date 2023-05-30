package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort       = 8080
	defaultSecretKey      = "secret"
	defaultAccessTokenTTL = time.Hour * 24 * 7
	defaultLoggerLevel    = "info"
)

type Config struct {
	HTTP    HTTPConfig
	Session SessionConfig
	Logger  LoggerConfig
}

type HTTPConfig struct {
	Port string `mapstructure:"port"`
}

type SessionConfig struct {
	SecretKey      string        `mapstructure:"secretKey"`
	AccessTokenTTL time.Duration `mapstructure:"accessTokenTTL"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

func NewConfig(path string) (*Config, error) {
	populateDefaults()
	if err := parseConfigFile(path); err != nil {
		return nil, fmt.Errorf("parseConfigFile failed: %w", err)
	}
	cfg := Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config Unmarshal failed: %w", err)
	}
	return &cfg, nil
}

func parseConfigFile(path string) error {
	dir, file := filepath.Split(path)
	file, _ = strings.CutSuffix(file, ".yml")
	viper.AddConfigPath(dir)
	viper.SetConfigName(file)

	return viper.ReadInConfig()
}
func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("session.secretKey", defaultSecretKey)
	viper.SetDefault("session.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("logger.level", defaultLoggerLevel)
}
