package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Security  SecurityConfig
	Message   MessageConfig
	RateLimit RateLimitConfig
	Log       LogConfig
}

type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Type         string `yaml:"type"`
	MaxClients   int    `yaml:"maxClients"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
}

type SecurityConfig struct {
	RequirePassword bool   `yaml:"requirePassword"`
	Password        string `yaml:"password"`
}

type MessageConfig struct {
	MaxLength int `yaml:"maxLength"`
}

type RateLimitConfig struct {
	MessagePerSecond int `yaml:"messagePerSecond"`
	Burst            int `yaml:"burst"`
}

type LogConfig struct {
	EnableLogging bool   `yaml:"enableLogging"`
	File          string `yaml:"file"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	return &config, nil
}
