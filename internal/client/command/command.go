package command

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
	"github.com/spf13/cobra"
)

type passKeeperService interface {
	Login(context.Context, string, string) error
	All(ctx context.Context) ([]domain.ShortData, error)
}

type MainCommand struct {
	passKeeperService passKeeperService
}

func New(passKeeperService passKeeperService) *MainCommand {
	return &MainCommand{
		passKeeperService: passKeeperService,
	}
}

func (mc *MainCommand) GetRootCommandList() []*cobra.Command {
	return []*cobra.Command{
		mc.getLoginCommand(),
		mc.getAllCommand(),
	}
}

func (mc *MainCommand) Execute() error {
	var RootCommand = cobra.Command{Use: "passkeep"}
	RootCommand.AddCommand(mc.GetRootCommandList()...)
	if err := RootCommand.Execute(); err != nil {
		return fmt.Errorf("failed execute root command: %w", err)
	}

	return nil
}
