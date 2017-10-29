package integration

import (
	"testing"

	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xcode/ipa"
	"github.com/stretchr/testify/require"
)

const (
	sampleArtifactsGitURI = "https://github.com/bitrise-samples/sample-artifacts.git"
)

func testIpasDirPth(t *testing.T) string {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__ipa__")
	require.NoError(t, err)

	cmd := command.New("git", "clone", sampleArtifactsGitURI, tmpDir)
	require.NoError(t, cmd.Run())

	ipasDir := filepath.Join(tmpDir, "ipas")
	exist, err := pathutil.IsDirExists(ipasDir)
	require.NoError(t, err, tmpDir)
	require.Equal(t, true, exist, tmpDir)

	return ipasDir
}

func TestIPAPackage(t *testing.T) {
	ipasDirPth := testIpasDirPth(t)

	t.Log("multiple files zipped in ipa")
	{
		ipaPth := filepath.Join(ipasDirPth, "watch-test.ipa")

		embeddedMobileProvisionPth, err := ipa.UnwrapEmbeddedMobileProvision(ipaPth)
		require.NoError(t, err)
		require.NotEqual(t, "", embeddedMobileProvisionPth)

		infoPlistPth, err := ipa.UnwrapEmbeddedInfoPlist(ipaPth)
		require.NoError(t, err)
		require.NotEqual(t, "", infoPlistPth)
	}

	t.Log("ipa file name != embedded app file name")
	{
		ipaPth := filepath.Join(ipasDirPth, "app-store-watch-test.ipa")

		embeddedMobileProvisionPth, err := ipa.UnwrapEmbeddedMobileProvision(ipaPth)
		require.NoError(t, err)
		require.NotEqual(t, "", embeddedMobileProvisionPth)

		infoPlistPth, err := ipa.UnwrapEmbeddedInfoPlist(ipaPth)
		require.NoError(t, err)
		require.NotEqual(t, "", infoPlistPth)
	}
}
