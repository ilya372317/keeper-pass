package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

const tuiArgumentCount = 0

func (mc *MainCommand) getTUICommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		if err := mc.terminalInterfaceApp.Run(cmd.Context()); err != nil {
			fmt.Println(err.Error())
		}
	}

	return &cobra.Command{
		Use:   "tui",
		Short: "run terminal user interface",
		Long: `run terminal user interface. in this interface user can perform 
operations available in other commands`,
		Example: "passkeep show 1",
		Args:    cobra.MaximumNArgs(tuiArgumentCount),
		Run:     run,
	}
}
