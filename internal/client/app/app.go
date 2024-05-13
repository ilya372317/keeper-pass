package app

import (
	"fmt"
	"os"

	"github.com/ilya372317/pass-keeper/internal/client/config"
)

const tokenFilePerm = 0600

var tokenFile *os.File

type App struct {
	c            *Container
	buildDate    string
	buildVersion string
	conf         config.Config
}

func New(buildDate, buildVersion string) (*App, error) {
	a := &App{
		buildDate:    buildDate,
		buildVersion: buildVersion,
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed get user home dir: %w", err)
	}
	a.conf = config.New(homeDir + "/.passkeeper.yaml")
	tokenFile, err = os.OpenFile(homeDir+"/.passtoken", os.O_RDWR|os.O_CREATE, tokenFilePerm)
	if err != nil {
		return nil, fmt.Errorf("failed open token ")
	}

	a.c = NewContainer(a.conf, tokenFile)
	return a, nil
}

func (a *App) Stop() error {
	if tokenFile != nil {
		if err := tokenFile.Close(); err != nil {
			return fmt.Errorf("failed close token file: %w", err)
		}
	}
	return nil
}
