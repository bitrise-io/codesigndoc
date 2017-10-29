package utility

import (
	"testing"

	"github.com/bitrise-tools/go-xcode/models"
	"github.com/stretchr/testify/require"
)

const testXcodebuildVersionOutput = `Xcode 8.2.1
Build version 8C1002`

func TestGetXcodeVersionFromXcodebuildOutput(t *testing.T) {
	t.Log("GetXcodeVersionFromXcodebuildOutput")
	{
		expectedVersion := models.XcodebuildVersionModel{
			Version:      "Xcode 8.2.1",
			BuildVersion: "Build version 8C1002",
			MajorVersion: 8,
		}

		currentVersion, err := getXcodeVersionFromXcodebuildOutput(testXcodebuildVersionOutput)
		require.NoError(t, err)
		require.Equal(t, expectedVersion, currentVersion)
	}
}
