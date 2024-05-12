package command

import (
	"testing"

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
