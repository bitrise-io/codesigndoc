package xcodeproj

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func cloneSampleProject(t *testing.T, url, projectPth string) string {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__codesignproperties__")
	require.NoError(t, err)

	cmd := command.New("git", "clone", url, tmpDir)
	require.NoError(t, cmd.Run())

	return filepath.Join(tmpDir, projectPth)
}

func TestResolveCodeSignInfo(t *testing.T) {
	user := os.Getenv("USER")

	{
		projectPth := cloneSampleProject(t, "https://github.com/bitrise-samples/sample-apps-ios-multi-target.git", "code-sign-test.xcodeproj")

		targetCodeSignInfoMap, err := ResolveCodeSignInfo(projectPth, "code-sign-test", user)
		require.NoError(t, err)
		require.Equal(t, 4, len(targetCodeSignInfoMap), fmt.Sprintf("%s", targetCodeSignInfoMap))

		{
			properties, ok := targetCodeSignInfoMap["watchkit-app Extension"]
			require.True(t, ok)

			require.Equal(t, "com.bitrise.code-sign-test.watchkitapp.watchkitextension", properties.BundleIdentifier)
			require.Equal(t, "iPhone Developer", properties.CodeSignIdentity)
			require.Equal(t, "548cd560-c511-4540-8b6b-cbec4a22f49d", properties.ProvisioningProfile)
			require.Equal(t, "BitriseBot-Wildcard", properties.ProvisioningProfileSpecifier)
			require.Equal(t, "72SA8V3WYL", properties.DevelopmentTeam)
		}

		{
			properties, ok := targetCodeSignInfoMap["code-sign-test"]
			require.True(t, ok)

			require.Equal(t, "com.bitrise.code-sign-test", properties.BundleIdentifier)
			require.Equal(t, "iPhone Developer", properties.CodeSignIdentity)
			require.Equal(t, "548cd560-c511-4540-8b6b-cbec4a22f49d", properties.ProvisioningProfile)
			require.Equal(t, "BitriseBot-Wildcard", properties.ProvisioningProfileSpecifier)
			require.Equal(t, "72SA8V3WYL", properties.DevelopmentTeam)
		}

		{
			properties, ok := targetCodeSignInfoMap["share-extension"]
			require.True(t, ok)

			require.Equal(t, "com.bitrise.code-sign-test.share-extension", properties.BundleIdentifier)
			require.Equal(t, "iPhone Developer", properties.CodeSignIdentity)
			require.Equal(t, "548cd560-c511-4540-8b6b-cbec4a22f49d", properties.ProvisioningProfile)
			require.Equal(t, "BitriseBot-Wildcard", properties.ProvisioningProfileSpecifier)
			require.Equal(t, "72SA8V3WYL", properties.DevelopmentTeam)
		}

		{
			properties, ok := targetCodeSignInfoMap["watchkit-app"]
			require.True(t, ok)

			require.Equal(t, "com.bitrise.code-sign-test.watchkitapp", properties.BundleIdentifier)
			require.Equal(t, "iPhone Developer", properties.CodeSignIdentity)
			require.Equal(t, "548cd560-c511-4540-8b6b-cbec4a22f49d", properties.ProvisioningProfile)
			require.Equal(t, "BitriseBot-Wildcard", properties.ProvisioningProfileSpecifier)
			require.Equal(t, "72SA8V3WYL", properties.DevelopmentTeam)
		}
	}

	{
		projectPth := cloneSampleProject(t, "https://github.com/bitrise-samples/sample-apps-ios-simple-objc.git", "ios-simple-objc/ios-simple-objc.xcodeproj")

		targetCodeSignInfoMap, err := ResolveCodeSignInfo(projectPth, "ios-simple-objc", user)
		require.NoError(t, err)
		require.Equal(t, 1, len(targetCodeSignInfoMap), fmt.Sprintf("%s", targetCodeSignInfoMap))

		{
			properties, ok := targetCodeSignInfoMap["ios-simple-objc"]
			require.True(t, ok)

			require.Equal(t, "Bitrise.ios-simple-objc", properties.BundleIdentifier)
			require.Equal(t, "iPhone Developer", properties.CodeSignIdentity)
			require.Equal(t, "", properties.ProvisioningProfile)
			require.Equal(t, "BitriseBot-Wildcard", properties.ProvisioningProfileSpecifier)
			require.Equal(t, "72SA8V3WYL", properties.DevelopmentTeam)
		}
	}
}
