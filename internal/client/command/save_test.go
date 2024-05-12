package command

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	command_mock "github.com/ilya372317/pass-keeper/internal/client/command/mocks"
	"github.com/stretchr/testify/require"
)

func TestMainCommand_getSaveCommand(t *testing.T) {
	mainCmd := MainCommand{}
	saveCmd := mainCmd.getSaveCommand()

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		saveCmd.SetArgs([]string{})

		// Execute.
		err := saveCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})
}

func TestMainCommand_getSaveLoginCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := command_mock.NewMockpassKeeperService(ctrl)
	mainCmd := MainCommand{passKeeperService: serv}
	saveLoginCmd := mainCmd.getSaveLoginCommand()
	validLogin := "login"
	validPass := "pass"
	validURL := "https://localhost"
	validArgs := []string{validLogin, validPass, validURL}

	t.Run("success save case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().SaveLogin(gomock.Any(), validLogin, validPass, validURL).Times(1).Return(nil)
		saveLoginCmd.SetArgs(validArgs)

		// Execute.
		err := saveLoginCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	t.Run("error in service", func(t *testing.T) {
		// Prepare.
		serv.
			EXPECT().
			SaveLogin(gomock.Any(), validLogin, validPass, validURL).
			Times(1).
			Return(fmt.Errorf("internal"))
		saveLoginCmd.SetArgs(validArgs)

		// Execute.
		err := saveLoginCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	validateTests := []struct {
		name     string
		login    string
		password string
		url      string
	}{
		{
			name:     "invalid login",
			login:    "",
			password: validPass,
			url:      validURL,
		},
		{
			name:     "invalid pass",
			login:    validLogin,
			password: "",
			url:      validURL,
		},
		{
			name:     "invalid URL",
			login:    validLogin,
			password: validPass,
			url:      "invalid-url",
		},
	}
	for _, tt := range validateTests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare.
			saveLoginCmd.SetArgs([]string{tt.login, tt.password, tt.url})

			// Execute.
			err := saveLoginCmd.Execute()

			// Assert.
			require.NoError(t, err)
		})
	}
}

func TestMainCommand_getSaveCardCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := command_mock.NewMockpassKeeperService(ctrl)
	mainCmd := MainCommand{passKeeperService: serv}
	saveCardCmd := mainCmd.getSaveCardCommand()
	validNumber := "374245455400126"
	validExp := "02/24"
	validCode := "123"
	validBankName := "bank name"
	validaArgs := []string{validNumber, validExp, validCode, validBankName}

	t.Run("success save case", func(t *testing.T) {
		// Prepare.
		serv.
			EXPECT().
			SaveCard(gomock.Any(), validNumber, validExp, validCode, validBankName).
			Times(1).
			Return(nil)
		saveCardCmd.SetArgs(validaArgs)

		// Execute.
		err := saveCardCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	t.Run("error in service", func(t *testing.T) {
		// Prepare.
		serv.
			EXPECT().
			SaveCard(gomock.Any(), validNumber, validExp, validCode, validBankName).
			Times(1).
			Return(fmt.Errorf("internal"))
		saveCardCmd.SetArgs(validaArgs)

		// Execute.
		err := saveCardCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	validateTests := []struct {
		name       string
		cardNumber string
		exp        string
		code       string
		bankName   string
	}{
		{
			name:       "invalid card number",
			cardNumber: "123",
			exp:        validExp,
			code:       validCode,
			bankName:   validBankName,
		},
		{
			name:       "invalid exp",
			cardNumber: validNumber,
			exp:        "invalid-exp",
			code:       validCode,
			bankName:   validBankName,
		},
		{
			name:       "invalid code",
			cardNumber: validNumber,
			exp:        validExp,
			code:       "invalid code",
			bankName:   validBankName,
		},
		{
			name:       "invalid bank name",
			cardNumber: validNumber,
			exp:        validExp,
			code:       validCode,
			bankName:   "",
		},
	}
	for _, tt := range validateTests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare.
			saveCardCmd.SetArgs([]string{tt.cardNumber, tt.exp, tt.code, tt.bankName})

			// Execute.
			err := saveCardCmd.Execute()

			// Assert.
			require.NoError(t, err)
		})
	}
}

func TestMainCommand_getSaveTextCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := command_mock.NewMockpassKeeperService(ctrl)
	mainCmd := MainCommand{passKeeperService: serv}
	saveTextCmd := mainCmd.getSaveTextCommand()
	validInfo := "info"
	validData := "data"
	validArgs := []string{validInfo, validData}

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().SaveText(gomock.Any(), validInfo, validData).Times(1).Return(nil)
		saveTextCmd.SetArgs(validArgs)

		// Execute.
		err := saveTextCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	t.Run("fail in service", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().SaveText(gomock.Any(), validInfo, validData).Times(1).Return(fmt.Errorf("internal"))
		saveTextCmd.SetArgs(validArgs)

		// Execute.
		err := saveTextCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	validateTests := []struct {
		name string
		info string
		data string
	}{
		{
			name: "invalid info",
			info: "",
			data: validData,
		},
		{
			name: "invalid data",
			info: "",
			data: validData,
		},
	}
	for _, tt := range validateTests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare.
			saveTextCmd.SetArgs([]string{tt.info, tt.data})

			// Execute.
			err := saveTextCmd.Execute()

			// Assert.
			require.NoError(t, err)
		})
	}
}
