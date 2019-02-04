package ios

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestAllowXcodeProjExtFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := utility.FilterPaths(paths, AllowXcodeProjExtFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestAllowXCWorkspaceExtFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
	}
	actualFiltered, err := utility.FilterPaths(paths, AllowXCWorkspaceExtFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestForbidEmbeddedWorkspaceRegexpFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := utility.FilterPaths(paths, ForbidEmbeddedWorkspaceRegexpFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestForbidGitDirComponentFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := utility.FilterPaths(paths, ForbidGitDirComponentFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestForbidPodsDirComponentFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := utility.FilterPaths(paths, ForbidPodsDirComponentFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestForbidCarthageDirComponentFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := utility.FilterPaths(paths, ForbidCarthageDirComponentFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestForbidFramworkComponentWithExtensionFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := utility.FilterPaths(paths, ForbidFramworkComponentWithExtensionFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestAllowIphoneosSDKFilter(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__xcodeproj_test__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	iphoneosProject := filepath.Join(tmpDir, "iphoneos.xcodeproj")
	require.NoError(t, os.MkdirAll(iphoneosProject, 0777))

	iphoneosPbxprojPth := filepath.Join(iphoneosProject, "project.pbxproj")
	require.NoError(t, fileutil.WriteStringToFile(iphoneosPbxprojPth, testIOSPbxprojContent))

	macosxProject := filepath.Join(tmpDir, "macosx.xcodeproj")
	require.NoError(t, os.MkdirAll(macosxProject, 0777))

	macosxPbxprojPth := filepath.Join(macosxProject, "project.pbxproj")
	require.NoError(t, fileutil.WriteStringToFile(macosxPbxprojPth, testMacOSPbxprojContent))

	t.Log("iphoneos sdk")
	{
		paths := []string{
			iphoneosProject,
			macosxProject,
		}
		expectedFiltered := []string{
			iphoneosProject,
		}
		actualFiltered, err := utility.FilterPaths(paths, AllowIphoneosSDKFilter)
		require.NoError(t, err)
		require.Equal(t, expectedFiltered, actualFiltered)
	}

	t.Log("macosx sdk")
	{
		paths := []string{
			iphoneosProject,
			macosxProject,
		}
		expectedFiltered := []string{
			macosxProject,
		}
		actualFiltered, err := utility.FilterPaths(paths, AllowMacosxSDKFilter)
		require.NoError(t, err)
		require.Equal(t, expectedFiltered, actualFiltered)
	}
}
