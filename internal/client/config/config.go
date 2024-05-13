package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type GRPCConfig struct {
	OpenMethods     map[string]any
	Host            string
	TLSCertPath     string
	RefreshDuration time.Duration
}

type Config struct {
	GRPC GRPCConfig
}

func New(configPath string) Config {
	var cnfg Config
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("missing config file or it is invalid: %v\n", err)
	}

	viper.SetDefault("grpc.host", ":3200")
	viper.SetDefault("grpc.open_methods", map[string]any{
		"/PassService/Auth":    struct{}{},
		"/PassKeeper/Register": struct{}{},
	})
	viper.SetDefault("grpc.refresh_token_duration", time.Hour)

	cnfg = Config{
		GRPC: GRPCConfig{
			Host:            viper.GetString("grpc.host"),
			OpenMethods:     viper.GetStringMap("grpc.open_methods"),
			RefreshDuration: viper.GetDuration("grpc.refresh_token_duration"),
			TLSCertPath:     viper.GetString("grpc.tls_cert"),
		},
	}

	return cnfg
}
