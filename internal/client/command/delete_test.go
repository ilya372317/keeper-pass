package command

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	command_mock "github.com/ilya372317/pass-keeper/internal/client/command/mocks"
	"github.com/stretchr/testify/require"
)

func TestMainCommand_getDeleteCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := command_mock.NewMockpassKeeperService(ctrl)
	mainCmd := MainCommand{passKeeperService: serv}
	deleteCmd := mainCmd.getDeleteCommand()
	validID := "1"
	validArgs := []string{validID}

	t.Run("success delete case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Delete(gomock.Any(), validID).Times(1).Return(nil)
		deleteCmd.SetArgs(validArgs)

		// Execute.
		err := deleteCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	t.Run("fail in service", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Delete(gomock.Any(), validID).Times(1).Return(fmt.Errorf("internal"))
		deleteCmd.SetArgs(validArgs)

		// Execute.
		err := deleteCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	t.Run("invalid id given", func(t *testing.T) {
		// Prepare.
		deleteCmd.SetArgs([]string{"invalid-id"})

		// Execute.
		err := deleteCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})
}
