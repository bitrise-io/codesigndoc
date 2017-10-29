package xcarchive

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

var tmpDir = ""

func sampleRepoPath(t *testing.T) string {
	dir := ""
	if tmpDir != "" {
		dir = tmpDir
	} else {
		var err error
		dir, err = pathutil.NormalizedOSTempDirPath("__artifacts__")
		require.NoError(t, err)
		sampleArtifactsGitURI := "https://github.com/bitrise-samples/sample-artifacts.git"
		cmd := command.New("git", "clone", sampleArtifactsGitURI, dir)
		require.NoError(t, cmd.Run())
		tmpDir = dir
	}
	t.Logf("sample artifcats dir: %s\n", dir)
	return dir
}

func TestNewIosArchive(t *testing.T) {
	iosArchivePth := filepath.Join(sampleRepoPath(t), "archives/ios.xcarchive")
	archive, err := NewIosArchive(iosArchivePth)
	require.NoError(t, err)
	require.Equal(t, 5, len(archive.InfoPlist))

	app := archive.Application
	require.Equal(t, 26, len(app.InfoPlist))
	require.Equal(t, 2, len(app.Entitlements))
	require.Equal(t, "*", app.ProvisioningProfile.BundleID)

	require.Equal(t, 1, len(app.Extensions))
	extension := app.Extensions[0]
	require.Equal(t, 23, len(extension.InfoPlist))
	require.Equal(t, 2, len(extension.Entitlements))
	require.Equal(t, "*", extension.ProvisioningProfile.BundleID)

	require.NotNil(t, app.WatchApplication)
	watchApp := *app.WatchApplication
	require.Equal(t, 24, len(watchApp.InfoPlist))
	require.Equal(t, 2, len(watchApp.Entitlements))
	require.Equal(t, "*", watchApp.ProvisioningProfile.BundleID)

	require.Equal(t, 1, len(watchApp.Extensions))
	watchExtension := watchApp.Extensions[0]
	require.Equal(t, 23, len(watchExtension.InfoPlist))
	require.Equal(t, 2, len(watchExtension.Entitlements))
	require.Equal(t, "*", watchExtension.ProvisioningProfile.BundleID)
}

func TestIsXcodeManaged(t *testing.T) {
	iosArchivePth := filepath.Join(sampleRepoPath(t), "archives/ios.xcarchive")
	archive, err := NewIosArchive(iosArchivePth)
	require.NoError(t, err)

	require.Equal(t, false, archive.IsXcodeManaged())
}

func TestSigningIdentity(t *testing.T) {
	iosArchivePth := filepath.Join(sampleRepoPath(t), "archives/ios.xcarchive")
	archive, err := NewIosArchive(iosArchivePth)
	require.NoError(t, err)

	require.Equal(t, "iPhone Developer: Bitrise Bot (VV2J4SV8V4)", archive.SigningIdentity())
}

func TestBundleIDEntitlementsMap(t *testing.T) {
	iosArchivePth := filepath.Join(sampleRepoPath(t), "archives/ios.xcarchive")
	archive, err := NewIosArchive(iosArchivePth)
	require.NoError(t, err)

	bundleIDEntitlementsMap := archive.BundleIDEntitlementsMap()
	require.Equal(t, 4, len(bundleIDEntitlementsMap))

	bundleIDs := []string{"com.bitrise.code-sign-test.share-extension", "com.bitrise.code-sign-test.watchkitapp", "com.bitrise.code-sign-test.watchkitapp.watchkitextension", "com.bitrise.code-sign-test"}
	for _, bundleID := range bundleIDs {
		_, ok := bundleIDEntitlementsMap[bundleID]
		require.True(t, ok, fmt.Sprintf("%v", bundleIDEntitlementsMap))
	}
}

func TestBundleIDProfileInfoMap(t *testing.T) {
	iosArchivePth := filepath.Join(sampleRepoPath(t), "archives/ios.xcarchive")
	archive, err := NewIosArchive(iosArchivePth)
	require.NoError(t, err)

	bundleIDProfileInfoMap := archive.BundleIDProfileInfoMap()
	require.Equal(t, 4, len(bundleIDProfileInfoMap))

	bundleIDs := []string{"com.bitrise.code-sign-test.share-extension", "com.bitrise.code-sign-test.watchkitapp", "com.bitrise.code-sign-test.watchkitapp.watchkitextension", "com.bitrise.code-sign-test"}
	for _, bundleID := range bundleIDs {
		_, ok := bundleIDProfileInfoMap[bundleID]
		require.True(t, ok, fmt.Sprintf("%v", bundleIDProfileInfoMap))
	}
}

func TestFindDSYMs(t *testing.T) {
	iosArchivePth := filepath.Join(sampleRepoPath(t), "archives/ios.xcarchive")
	archive, err := NewIosArchive(iosArchivePth)
	require.NoError(t, err)

	appDsym, otherDsyms, err := archive.FindDSYMs()
	require.NoError(t, err)
	require.NotEmpty(t, appDsym)
	require.Equal(t, 2, len(otherDsyms))
}
