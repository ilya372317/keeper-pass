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
