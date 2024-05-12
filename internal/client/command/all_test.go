package command

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	command_mock "github.com/ilya372317/pass-keeper/internal/client/command/mocks"
	"github.com/ilya372317/pass-keeper/internal/client/domain"
	"github.com/stretchr/testify/require"
)

func TestMainCommand_getAllCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := command_mock.NewMockpassKeeperService(ctrl)
	cmd := MainCommand{passKeeperService: serv}
	allCmd := cmd.getAllCommand()
	allCmd.SetArgs([]string{})

	t.Run("success get all case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().All(gomock.Any()).Times(1).Return([]domain.ShortData{{}}, nil)

		// Execute.
		err := allCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	t.Run("unauthenticated", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().All(gomock.Any()).Times(1).Return(nil, domain.ErrUnauthenticated)

		// Execute.
		err := allCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	t.Run("internal error", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().All(gomock.Any()).Times(1).Return(nil, fmt.Errorf("internal"))

		// Execute.
		err := allCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})
}
