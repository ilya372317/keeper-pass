package command

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
	"github.com/spf13/cobra"
)

type passKeeperService interface {
	Login(context.Context, string, string) error
	All(context.Context) ([]domain.ShortData, error)
	Register(context.Context, string, string) error
	SaveLogin(context.Context, string, string, string) error
	SaveCard(context.Context, string, string, string, string) error
	SaveText(context.Context, string, string) error
	SaveBinary(context.Context, string, string) error
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
		mc.getRegisterCommand(),
		mc.getSaveCommand(),
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
