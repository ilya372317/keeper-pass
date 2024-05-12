package command

import (
	"errors"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
	"github.com/spf13/cobra"
)

func (mc *MainCommand) getAllCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		data, err := mc.passKeeperService.All(cmd.Context())
		if err != nil {
			if errors.Is(err, domain.ErrUnauthenticated) {
				fmt.Println("you need login before use this command. try passkeep login [email] [password]")
				return
			}
			fmt.Printf("failed get data: %v\n", err)
			return
		}
		fmt.Printf("%-10s %-40s %-10s\n", "ID", "INFO", "TYPE")
		for _, d := range data {
			fmt.Printf("%-10d %-40s %-10s\n", d.ID, d.Info, d.StringKind())
		}
	}

	return &cobra.Command{
		Use:     "all",
		Short:   "list your data. required be auth",
		Long:    `list your data in tabel view. with column ID, INFO, TYPE. required be auth`,
		Example: "passkeep all",
		Args:    cobra.MaximumNArgs(0),
		Run:     run,
	}
}
