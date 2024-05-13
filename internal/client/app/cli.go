package app

import (
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/client/command"
)

func (a *App) ExecuteCommandCLI() error {
	mainCommand := command.New(a.c.GetDefaultPassKeeperService(), a.buildVersion, a.buildDate)

	if err := mainCommand.Execute(); err != nil {
		return fmt.Errorf("failed execute command CLI: %w", err)
	}

	return nil
}
