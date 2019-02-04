package xamarin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterSolutionFiles(t *testing.T) {
	t.Log(`Contains solution files`)
	{
		fileList := []string{
			"/Users/bitrise/Develop/bitrise/sample-apps/sample-apps-xamarin-ios/CreditCardValidator.iOS.sln",
			"/Users/bitrise/Develop/bitrise/sample-apps/sample-apps-android/sln",
			"path/to/my/gradlew/file",
			"path/to/my",
		}

		files, err := FilterSolutionFiles(fileList)
		require.NoError(t, err)
		require.Equal(t, 1, len(files))

		// Also sorts solution files by path components length
		require.Equal(t, "/Users/bitrise/Develop/bitrise/sample-apps/sample-apps-xamarin-ios/CreditCardValidator.iOS.sln", files[0])
	}

	t.Log(`Do not contains solution file`)
	{
		fileList := []string{
			"path/to/my/gradlew/build.",
			"path/to/my/gradle",
		}

		files, err := FilterSolutionFiles(fileList)
		require.NoError(t, err)
		require.Equal(t, 0, len(files))
	}
}
