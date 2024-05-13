package command

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	command_mock "github.com/ilya372317/pass-keeper/internal/client/command/mocks"
	"github.com/stretchr/testify/require"
)

func TestMainCommand_getRegisterCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := command_mock.NewMockpassKeeperService(ctrl)
	mainCommand := MainCommand{passKeeperService: serv}
	registerCommand := mainCommand.getRegisterCommand()
	validEmail := "1@gmail.com"
	validPassword := "123"
	validArgs := []string{validEmail, validPassword}

	t.Run("success register case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Register(gomock.Any(), validEmail, validPassword).Times(1).Return(nil)
		registerCommand.SetArgs(validArgs)

		// Execute.
		err := registerCommand.Execute()

		// Assert.
		require.NoError(t, err)
	})

	t.Run("failed register in service", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Register(gomock.Any(), validEmail, validPassword).Times(1).Return(fmt.Errorf("internal"))
		registerCommand.SetArgs(validArgs)

		// Execute.
		err := registerCommand.Execute()

		// Assert.
		require.NoError(t, err)
	})

	validateTests := []struct {
		name     string
		email    string
		password string
	}{
		{
			name:     "invalid password",
			email:    validEmail,
			password: "1",
		},
		{
			name:     "invalid email",
			email:    "ivalid-email",
			password: validPassword,
		},
	}
	for _, tt := range validateTests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare.
			serv.EXPECT().Register(gomock.Any(), tt.email, tt.password).Times(0).Return(nil)
			registerCommand.SetArgs([]string{tt.email, tt.password})

			// Execute.
			err := registerCommand.Execute()

			// Assert.
			require.NoError(t, err)
		})
	}
}
