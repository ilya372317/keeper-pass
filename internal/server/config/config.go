package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type GRPCConfig struct {
	Host string
}

type Config struct {
	GRPC GRPCConfig
}

func New(configPath string) (*Config, error) {
	var cnfg *Config
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed load config file. file path [%s] is invalid: %w", configPath, err)
	}

	viper.SetDefault("grpc.host", ":3200")

	cnfg = &Config{
		GRPC: GRPCConfig{Host: viper.GetString("grpc.host")},
	}

	return cnfg, nil
}
