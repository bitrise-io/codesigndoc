package xcodebuild

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBuildSettings(t *testing.T) {
	settings, err := parseBuildSettings(testBuildSettingsOut)
	require.NoError(t, err)

	desired := map[string]string{
		"VERSION_INFO_STRING":                           "@(#)PROGRAM:sample-apps-osx-10-12  PROJECT:sample-apps-osx-10-12-",
		"AVAILABLE_PLATFORMS":                           "appletvos appletvsimulator iphoneos iphonesimulator macosx watchos watchsimulator",
		"EXCLUDED_RECURSIVE_SEARCH_PATH_SUBDIRECTORIES": "*.nib *.lproj *.framework *.gch *.xcode* *.xcassets (*) .DS_Store CVS .svn .git .hg *.pbproj *.pbxproj",
		"BUILD_STYLE":                                   "",
		"ACTION":                                        "build",
		"SDK_VERSION_MINOR":                             "1200",
	}

	require.Equal(t, desired, settings)
}

const testBuildSettingsOut = `Build settings for action build and target sample-apps-osx-10-12:
    VERSION_INFO_STRING = "@(#)PROGRAM:sample-apps-osx-10-12  PROJECT:sample-apps-osx-10-12-"
    AVAILABLE_PLATFORMS = appletvos appletvsimulator iphoneos iphonesimulator macosx watchos watchsimulator
    EXCLUDED_RECURSIVE_SEARCH_PATH_SUBDIRECTORIES = *.nib *.lproj *.framework *.gch *.xcode* *.xcassets (*) .DS_Store CVS .svn .git .hg *.pbproj *.pbxproj
    BUILD_STYLE =
    ACTION = build
    SDK_VERSION_MINOR = 1200`
