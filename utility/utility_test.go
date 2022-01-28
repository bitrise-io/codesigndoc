package utility

import (
	"testing"

	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
	"github.com/stretchr/testify/require"
)

func TestPlatformsMatching_iOS(t *testing.T) {
	buildSettings := serialized.Object{}
	buildSettings["SDKROOT"] = "iphoneos"

	platform, err := getPlatform(buildSettings)

	require.Equal(t, "iOS", string(platform))
	require.Nil(t, err)
}

func TestPlatformsMatching_macOS(t *testing.T) {
	buildSettings := serialized.Object{}
	buildSettings["SDKROOT"] = "macosx"

	platform, err := getPlatform(buildSettings)

	require.Equal(t, "OS X", string(platform))
	require.Nil(t, err)
}

func TestPlatformsMatching_fails(t *testing.T) {
	buildSettings := serialized.Object{}
	platform, err := getPlatform(buildSettings)

	require.Empty(t, platform)
	require.NotNil(t, err)
}
