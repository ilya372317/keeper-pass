package app

import (
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/config"
)

type App struct {
	conf *config.Config
}

func New(configPath string) (*App, error) {
	conf, err := config.New(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed create new config: %w", err)
	}

	return &App{
		conf: conf,
	}, nil
}
