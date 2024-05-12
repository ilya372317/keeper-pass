package command

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

const registerArgCount = 2

type registerValidator struct {
	Email    string `validate:"email,required"`
	Password string `validate:"required,min=3,max=255"`
}

func (mc *MainCommand) getRegisterCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		email, password := args[0], args[1]
		validateStruct := registerValidator{
			Email:    email,
			Password: password,
		}

		validate := validator.New(validator.WithRequiredStructEnabled())

		if err := validate.Struct(&validateStruct); err != nil {
			fmt.Printf("invalid argument given: %v\n", err)
			return
		}

		if err := mc.passKeeperService.Register(cmd.Context(), email, password); err != nil {
			fmt.Printf("failed register: %v\n", err)
			return
		}

		fmt.Println("your successfully registered. Your creds for login:")
		fmt.Printf("email: %-10s password: %-10s\n", email, password)
	}

	return &cobra.Command{
		Use:     "register [email] [password]",
		Short:   "register in server by given email and password",
		Long:    `register in server by given email and password. remember this data, it need to user in login command`,
		Example: "passkeep register 1@gmail.com 123",
		Args:    cobra.MinimumNArgs(registerArgCount),
		Run:     run,
	}
}
