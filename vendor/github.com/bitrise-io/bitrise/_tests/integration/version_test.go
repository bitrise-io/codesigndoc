package integration

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/stretchr/testify/require"
)

func Test_VersionOutput(t *testing.T) {
	t.Log("Version")
	{
		out, err := command.RunCommandAndReturnCombinedStdoutAndStderr(binPath(), "version")
		require.NoError(t, err)
		require.Equal(t, "1.8.0", out)
	}

	t.Log("Version --full")
	{
		out, err := command.RunCommandAndReturnCombinedStdoutAndStderr(binPath(), "version", "--full")
		require.NoError(t, err)

		expectedOSVersion := fmt.Sprintf("%s (%s)", runtime.GOOS, runtime.GOARCH)
		expectedVersionOut := fmt.Sprintf(`version: 1.8.0
format version: 4
os: %s
go: %s
build number: 
commit:`, expectedOSVersion, runtime.Version())

		require.Equal(t, expectedVersionOut, out)
	}
}
