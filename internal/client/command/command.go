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
	Delete(context.Context, string) error
	Show(context.Context, string, string) (string, error)
}

type MainCommand struct {
	passKeeperService passKeeperService
	buildVersion      string
	buildDate         string
}

func New(passKeeperService passKeeperService, buildVersion, buildDate string) *MainCommand {
	return &MainCommand{
		passKeeperService: passKeeperService,
		buildDate:         buildDate,
		buildVersion:      buildVersion,
	}
}

func (mc *MainCommand) GetRootCommandList() []*cobra.Command {
	return []*cobra.Command{
		mc.getLoginCommand(),
		mc.getAllCommand(),
		mc.getRegisterCommand(),
		mc.getSaveCommand(),
		mc.getDeleteCommand(),
		mc.getShowCommand(),
		mc.getVersionCommand(mc.buildVersion, mc.buildDate),
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
