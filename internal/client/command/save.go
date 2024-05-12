package command

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

const (
	saveArgCount      = 0
	saveLoginArgCount = 3
)

func (mc *MainCommand) getSaveCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		fmt.Println("for save any information, try specific command for do this:")
		fmt.Println("passkeep save login [login] [password] [url] - for save login pass pair info")
		fmt.Println("passkeep save card [card number] [expiration date in format 02/24] [CVV] [bank name]")
		fmt.Println("passkeep save text [info about saving text] [text data]")
		fmt.Println("passkeep save binary [info about saving file] [file_path]")
	}

	cmd := &cobra.Command{
		Use:   "save [type]",
		Short: "base command for save different types of data. not do anything on it own.",
		Long: `base command for save differen types of data. should be used with specified type.
on it own not do anything. for actual save data you need use subcommand like login or card`,
		Example: "passkeep register",
		Args:    cobra.MinimumNArgs(saveArgCount),
		Run:     run,
	}
	cmd.AddCommand(mc.getSaveLoginCommand())

	return cmd
}

type saveLoginValidate struct {
	URL      string `validate:"required,url"`
	Login    string `validate:"required"`
	Password string `validate:"required"`
}

func (mc *MainCommand) getSaveLoginCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		login, password, url := args[0], args[1], args[2]
		validateStruct := saveLoginValidate{
			URL:      url,
			Login:    login,
			Password: password,
		}
		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(&validateStruct); err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		if err := mc.passKeeperService.SaveLogin(cmd.Context(), login, password, url); err != nil {
			fmt.Printf("%v", err)
			return
		}

		fmt.Println("login info successfully saved!")
	}

	return &cobra.Command{
		Use:     "login [login] [password] [url]",
		Short:   "command for save login pass pair type data",
		Long:    `command for save login pass pair type data`,
		Example: "passkeep save login test 123 'https://localhost'",
		Args:    cobra.MinimumNArgs(saveLoginArgCount),
		Run:     run,
	}
}
