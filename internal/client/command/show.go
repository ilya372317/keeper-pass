package command

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

const showArgCount = 2

type showValidator struct {
	ID   string `validate:"required,number"`
	Kind string `validate:"required"`
}

func (mc *MainCommand) getShowCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		id, kind := args[0], args[1]
		validateStruct := showValidator{ID: id, Kind: kind}
		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(&validateStruct); err != nil {
			fmt.Println(err)
			return
		}

		info, err := mc.passKeeperService.Show(cmd.Context(), id, kind)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(info)
	}

	return &cobra.Command{
		Use:     "show [id] [type]",
		Short:   "command for show data from server",
		Long:    `command for show data from server`,
		Example: "passkeep show 1",
		Args:    cobra.MinimumNArgs(showArgCount),
		Run:     run,
	}
}
