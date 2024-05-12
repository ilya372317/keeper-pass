package command

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

const loginArgCount = 2

type loginValidator struct {
	Email    string `validate:"email,required"`
	Password string `validate:"required"`
}

func (mc *MainCommand) getLoginCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		email, password := args[0], args[1]
		validateStruct := loginValidator{
			Email:    email,
			Password: password,
		}
		validate := validator.New(validator.WithRequiredStructEnabled())

		if validationErr := validate.Struct(&validateStruct); validationErr != nil {
			fmt.Printf("invalid argument given: %v\n", validationErr)
			return
		}

		if err := mc.passKeeperService.Login(cmd.Context(), email, password); err != nil {
			fmt.Printf("failed login: %v\n", err)
			return
		}

		fmt.Println("you successfully login.")
	}

	cmd := &cobra.Command{
		Use:   "login [login] [password]",
		Short: "login on server",
		Long: `by given arguments login and password perform login on server. creating file 
~/.passtoken with auth jwt token`,
		Example: "passkeep login 1@gmail.com 123",
		Args:    cobra.MinimumNArgs(loginArgCount),
		Run:     run,
	}

	return cmd
}
