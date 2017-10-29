package xcodeproj

import (
	"testing"

	"github.com/bitrise-io/go-utils/testutil"
	"github.com/stretchr/testify/require"
)

func TestGetBuildConfigSDKRoot(t *testing.T) {
	t.Log("ios")
	{
		pbxprojPth, err := testIOSPbxprojPth()
		require.NoError(t, err)

		sdks, err := GetBuildConfigSDKs(pbxprojPth)
		require.NoError(t, err)
		testutil.EqualSlicesWithoutOrder(t, []string{"iphoneos"}, sdks)
	}

	t.Log("macos")
	{
		pbxprojPth, err := testMacOSPbxprojPth()
		require.NoError(t, err)

		sdks, err := GetBuildConfigSDKs(pbxprojPth)
		require.NoError(t, err)
		testutil.EqualSlicesWithoutOrder(t, []string{"macosx"}, sdks)
	}
}
