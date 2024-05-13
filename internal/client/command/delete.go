package command

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

const deleteArgCount = 1

type deleteValidator struct {
	ID string `validate:"number"`
}

func (mc *MainCommand) getDeleteCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		id := args[0]
		structValidate := deleteValidator{ID: id}
		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(&structValidate); err != nil {
			fmt.Println(err)
			return
		}

		if err := mc.passKeeperService.Delete(cmd.Context(), id); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("data with id [%s] successfully delete\n", id)
	}

	return &cobra.Command{
		Use:     "delete [record ID]",
		Short:   "delete data from server using given ID",
		Long:    `delete data from server using given ID`,
		Example: "passkeep delete 1",
		Args:    cobra.MinimumNArgs(deleteArgCount),
		Run:     run,
	}
}
