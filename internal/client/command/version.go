package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

const versionArgCount = 0

func (mc *MainCommand) getVersionCommand(version, buildDate string) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nBuild date: %s\n", version, buildDate)
	}

	return &cobra.Command{
		Use:     "version",
		Short:   "print information about version and build date",
		Long:    `print information about version and build date`,
		Example: "passkeep show",
		Args:    cobra.MaximumNArgs(versionArgCount),
		Run:     run,
	}
}
