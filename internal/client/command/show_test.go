package command

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	command_mock "github.com/ilya372317/pass-keeper/internal/client/command/mocks"
	"github.com/stretchr/testify/require"
)

func TestMainCommand_getShowCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	serv := command_mock.NewMockpassKeeperService(ctrl)
	validID := "1"
	validType := "login-pass"
	validArgs := []string{validID, validType}
	mainCmd := MainCommand{passKeeperService: serv}
	showCmd := mainCmd.getShowCommand()

	t.Run("success case", func(t *testing.T) {
		// Prepare.
		serv.EXPECT().Show(gomock.Any(), validID, validType).Times(1).Return("result", nil)
		showCmd.SetArgs(validArgs)

		// Execute.
		err := showCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	t.Run("fail in service", func(t *testing.T) {
		// Prepare.
		serv.
			EXPECT().
			Show(gomock.Any(), validID, validType).
			Times(1).
			Return("", fmt.Errorf("internal"))
		showCmd.SetArgs(validArgs)

		// Execute.
		err := showCmd.Execute()

		// Assert.
		require.NoError(t, err)
	})

	validateTests := []struct {
		name string
		id   string
		kind string
	}{
		{
			name: "invalid id",
			id:   "invalid id",
			kind: validType,
		},
		{
			name: "invalid type",
			id:   validID,
			kind: "",
		},
	}
	for _, tt := range validateTests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare.
			showCmd.SetArgs([]string{tt.id, tt.kind})

			// Execute.
			err := showCmd.Execute()

			// Assert.
			require.NoError(t, err)
		})
	}
}
