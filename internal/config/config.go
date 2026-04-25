package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the global configuration structure.
type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
}

type ServerConfig struct {
	HTTP struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"http"`
	GRPC struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"grpc"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load loads the configuration from the given file path, falling back to env vars.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		// Default search paths
		v.AddConfigPath("configs")
		v.SetConfigName("server")
		v.SetConfigType("yaml")
	}

	// Read environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values
	v.SetDefault("app.name", "my-app")
	v.SetDefault("server.http.addr", ":8080")
	v.SetDefault("log.level", "info")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// It's ok if config file is not found, we fallback to defaults/env
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
