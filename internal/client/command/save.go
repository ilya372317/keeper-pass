package command

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ilya372317/pass-keeper/pkg/validation"
	"github.com/spf13/cobra"
)

const (
	saveArgCount      = 0
	saveLoginArgCount = 3
	saveCardArgCount  = 4
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
	cmd.AddCommand(mc.getSaveCardCommand())

	return cmd
}

type saveLoginValidator struct {
	URL      string `validate:"required,url"`
	Login    string `validate:"required"`
	Password string `validate:"required"`
}

func (mc *MainCommand) getSaveLoginCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		login, password, url := args[0], args[1], args[2]
		validateStruct := saveLoginValidator{
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

type cardValidator struct {
	CardNumber string `validate:"required,credit_card"`
	Expiration string `validate:"required,cardexp"`
	CVV        string `validate:"required,number"`
	BankName   string `validate:"required"`
}

func (mc *MainCommand) getSaveCardCommand() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		number, exp, cvv, bankName := args[0], args[1], args[2], args[3]
		validateStruct := cardValidator{
			CardNumber: number,
			Expiration: exp,
			CVV:        cvv,
			BankName:   bankName,
		}
		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.RegisterValidation("cardexp", validation.CardExpValidation); err != nil {
			fmt.Println("failed register validator for check card expiration." +
				" this internal error can fix only developer of the app. sorry about that...")
			return
		}
		if err := validate.Struct(&validateStruct); err != nil {
			fmt.Println(err)
			return
		}

		if err := mc.passKeeperService.SaveCard(cmd.Context(), number, exp, cvv, bankName); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("successfully save card information!")
	}

	return &cobra.Command{
		Use:     "card [number] [exp] [cvv] [bank name]",
		Short:   "command for save credit card information",
		Long:    `command for save credit card information`,
		Example: "passkeep save credit 374245455400126 02/24 123",
		Args:    cobra.MinimumNArgs(saveCardArgCount),
		Run:     run,
	}
}
