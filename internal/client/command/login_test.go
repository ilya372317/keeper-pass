package command

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	command_mock "github.com/ilya372317/pass-keeper/internal/client/command/mocks"
	"github.com/stretchr/testify/require"
)

func TestMainCommand_getLoginCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := command_mock.NewMockpassKeeperService(ctrl)
	cmd := MainCommand{passKeeperService: serv}
	loginCmd := cmd.getLoginCommand()
	validEmail := "1@gmail.com"
	validPass := "123"
	validEmailPassPair := []string{validEmail, validPass}

	t.Run("success login case", func(t *testing.T) {
		// Prepare.
		loginCmd.SetArgs(validEmailPassPair)
		serv.EXPECT().Login(gomock.Any(), validEmail, validPass).Times(1).Return(nil)

		// Execute.
		err := loginCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	t.Run("failed perform login in service", func(t *testing.T) {
		// Prepare.
		loginCmd.SetArgs(validEmailPassPair)
		serv.EXPECT().Login(gomock.Any(), validEmail, validPass).Times(1).Return(fmt.Errorf("internal"))

		// Execute.
		err := loginCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	validationTests := []struct {
		name     string
		email    string
		password string
	}{
		{
			name:     "invalid email case",
			email:    "invalid-email",
			password: validPass,
		},
		{
			name:     "invalid password case",
			email:    validEmail,
			password: "",
		},
	}
	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare.
			loginCmd.SetArgs([]string{tt.email, tt.password})
			serv.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Times(0).Return(nil)

			// Execute.
			err := loginCmd.Execute()

			// Assert.
			require.NoError(t, err)
		})
	}
}
