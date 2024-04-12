package app

import (
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/config"
	"github.com/ilya372317/pass-keeper/internal/server/logger"
)

type App struct {
	conf config.Config
}

func New(configPath string) (App, error) {
	conf, err := config.New(configPath)
	logger.InitMust()
	if err != nil {
		return App{}, fmt.Errorf("failed create new config: %w", err)
	}

	return App{
		conf: conf,
	}, nil
}
