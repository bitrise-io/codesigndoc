package xcarchive

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMacosArchive(t *testing.T) {
	macosArchivePth := filepath.Join(sampleRepoPath(t), "archives/macos.xcarchive")
	archive, err := NewMacosArchive(macosArchivePth)
	require.NoError(t, err)
	require.Equal(t, 5, len(archive.InfoPlist))

	app := archive.Application
	require.Equal(t, 21, len(app.InfoPlist))
	require.Equal(t, 2, len(app.Entitlements))
	require.Nil(t, app.ProvisioningProfile)

	require.Equal(t, 1, len(app.Extensions))
	extension := app.Extensions[0]
	require.Equal(t, 22, len(extension.InfoPlist))
	require.Equal(t, 2, len(extension.Entitlements))
	require.Nil(t, extension.ProvisioningProfile)
}

func TestMacosIsXcodeManaged(t *testing.T) {
	macosArchivePth := filepath.Join(sampleRepoPath(t), "archives/macos.xcarchive")
	archive, err := NewMacosArchive(macosArchivePth)
	require.NoError(t, err)

	require.Equal(t, false, archive.IsXcodeManaged())
}

func TestMacosSigningIdentity(t *testing.T) {
	macosArchivePth := filepath.Join(sampleRepoPath(t), "archives/macos.xcarchive")
	archive, err := NewMacosArchive(macosArchivePth)
	require.NoError(t, err)

	require.Equal(t, "Mac Developer: Gödrei Krisztian (T3694PR6UJ)", archive.SigningIdentity())
}

func TestMacosBundleIDEntitlementsMap(t *testing.T) {
	macosArchivePth := filepath.Join(sampleRepoPath(t), "archives/macos.xcarchive")
	archive, err := NewMacosArchive(macosArchivePth)
	require.NoError(t, err)

	bundleIDEntitlementsMap := archive.BundleIDEntitlementsMap()
	require.Equal(t, 2, len(bundleIDEntitlementsMap))

	bundleIDs := []string{"io.bitrise.archive.Test", "io.bitrise.archive.Test.ActionExtension"}
	for _, bundleID := range bundleIDs {
		_, ok := bundleIDEntitlementsMap[bundleID]
		require.True(t, ok, fmt.Sprintf("%v", bundleIDEntitlementsMap))
	}
}

func TestMacosBundleIDProfileInfoMap(t *testing.T) {
	macosArchivePth := filepath.Join(sampleRepoPath(t), "archives/macos.xcarchive")
	archive, err := NewMacosArchive(macosArchivePth)
	require.NoError(t, err)

	bundleIDProfileInfoMap := archive.BundleIDProfileInfoMap()
	require.Equal(t, 0, len(bundleIDProfileInfoMap))
}

func TestMacosFindDSYMs(t *testing.T) {
	macosArchivePth := filepath.Join(sampleRepoPath(t), "archives/macos.xcarchive")
	archive, err := NewMacosArchive(macosArchivePth)
	require.NoError(t, err)

	appDsym, otherDsyms, err := archive.FindDSYMs()
	require.NoError(t, err)
	require.NotEmpty(t, appDsym)
	require.Equal(t, 1, len(otherDsyms))
}
