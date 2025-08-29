package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	PORT             string `yaml:"PORT"`
	Type             string `yaml:"type"`
	RequirePassword  bool   `yaml:"RequirePassword"`
	Password         string `yaml:"Password"`
	MaxMessageLength int    `yaml:"MaxMessageLength"`
	MaxClients       int    `yaml:"MaxClients"`
	RateLimit        int    `yaml:"RateLimit"`
	EnableLogging    bool   `yaml:"EnableLogging"`
	LogFile          string `yaml:"LogFile"`
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
