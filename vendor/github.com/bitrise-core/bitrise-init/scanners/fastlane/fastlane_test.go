package fastlane

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterFastFiles(t *testing.T) {

	t.Log(`Contains "Fastfile" files`)
	{
		fileList := []string{
			"/Users/bitrise/Develop/bitrise/sample-apps/sample-apps-android/fastlane/Fastfile",
			"/Users/bitrise/Develop/bitrise/sample-apps/sample-apps-android/Fastfile",
			"path/to/my/gradlew/file",
			"path/to/my",
		}

		files, err := FilterFastfiles(fileList)
		require.NoError(t, err)
		require.Equal(t, 2, len(files))

		// Also sorts "Fastfile" files by path components length
		require.Equal(t, "/Users/bitrise/Develop/bitrise/sample-apps/sample-apps-android/Fastfile", files[0])
		require.Equal(t, "/Users/bitrise/Develop/bitrise/sample-apps/sample-apps-android/fastlane/Fastfile", files[1])
	}

	t.Log(`Do not contains "Fastfile" file`)
	{
		fileList := []string{
			"path/to/my/gradlew/build.",
			"path/to/my/gradle",
		}

		files, err := FilterFastfiles(fileList)
		require.NoError(t, err)
		require.Equal(t, 0, len(files))
	}
}

func TestInspectFastFileContent(t *testing.T) {
	lines := []string{
		" test ",
		" lane ",
		":xcode",

		"  lane :xcode do",
		"lane :deploy do",
		"  lane :unit_tests do |params|",

		"  private_lane :post_to_slack do |options|",
		"  private_lane :verify_xcode_version do",
	}
	content := strings.Join(lines, "\n")

	expectedLanes := []string{
		"xcode",
		"deploy",
		"unit_tests",
	}

	lanes, err := inspectFastfileContent(content)
	require.NoError(t, err)
	require.Equal(t, expectedLanes, lanes)

	t.Log("ios test")
	{
		lanes, err := inspectFastfileContent(iosTesFastfileContent)
		require.NoError(t, err)

		expectedLanes := []string{
			"ios test",
		}

		require.Equal(t, expectedLanes, lanes)
	}

	t.Log("experimental ios test")
	{
		lanes, err := inspectFastfileContent(complexIosTestFastFileContent)
		require.NoError(t, err)

		expectedLanes := []string{
			"ios analyze",
			"ios testAndPushBeta",
			"ios submitAndPushToMaster",
			"ios verifyTestPlatforms",
			"ios verify",
			"ios bumpPatch",
			"ios bumpMinor",
			"ios bumpMajor",
			"ios bump",
			"ios bumpAndTagBeta",
			"ios bumpAndTagRelease",
			"ios default_changelog",
			"ios beta",
			"ios store",
			"ios dev",
		}

		require.Equal(t, expectedLanes, lanes)
	}
}

func TestFastlaneWorkDir(t *testing.T) {
	t.Log("Fastfile's dir, if Fastfile is NOT in fastlane dir")
	{
		expected := "."
		actual := WorkDir("Fastfile")
		require.Equal(t, expected, actual)
	}

	t.Log("fastlane dir's parent, if Fastfile is in fastlane dir")
	{
		expected := "."
		actual := WorkDir("fastlane/Fastfile")
		require.Equal(t, expected, actual)
	}

	t.Log("Fastfile's dir, if Fastfile is NOT in fastlane dir")
	{
		expected := "test"
		actual := WorkDir("test/Fastfile")
		require.Equal(t, expected, actual)
	}

	t.Log("fastlane dir's parent, if Fastfile is in fastlane dir")
	{
		expected := "test"
		actual := WorkDir("test/fastlane/Fastfile")
		require.Equal(t, expected, actual)
	}

	t.Log("Fastfile's dir, if Fastfile is NOT in fastlane dir")
	{
		expected := "my/app/test"
		actual := WorkDir("my/app/test/Fastfile")
		require.Equal(t, expected, actual)
	}

	t.Log("fastlane dir's parent, if Fastfile is in fastlane dir")
	{
		expected := "my/app/test"
		actual := WorkDir("my/app/test/fastlane/Fastfile")
		require.Equal(t, expected, actual)
	}
}
