package pathfilters

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/xcodeproj"
)

const (
	embeddedWorkspacePathPattern = `.+\.xcodeproj/.+\.xcworkspace`

	gitDirName         = ".git"
	podsDirName        = "Pods"
	carthageDirName    = "Carthage"
	cordovaLibDirName  = "CordovaLib"
	nodeModulesDirName = "node_modules"
	frameworkExt       = ".framework"
	podfileBase        = "Podfile"
)

// XcodeProjectType ...
type XcodeProjectType string

const (
	// XcodeProjectTypeIOS ...
	XcodeProjectTypeIOS XcodeProjectType = "ios"
	// XcodeProjectTypeMacOS ...
	XcodeProjectTypeMacOS XcodeProjectType = "macos"
)

// AllowXcodeProjExtFilter ...
var AllowXcodeProjExtFilter = pathutil.ExtensionFilter(xcodeproj.XCodeProjExt, true)

// AllowXCWorkspaceExtFilter ...
var AllowXCWorkspaceExtFilter = pathutil.ExtensionFilter(xcodeproj.XCWorkspaceExt, true)

// AllowIsDirectoryFilter ...
var AllowIsDirectoryFilter = pathutil.IsDirectoryFilter(true)

// AllowWorkspaceWithContentsFile ...
var AllowWorkspaceWithContentsFile = pathutil.DirectoryContainsFileFilter("contents.xcworkspacedata")

// ForbidEmbeddedWorkspaceRegexpFilter ...
var ForbidEmbeddedWorkspaceRegexpFilter = pathutil.RegexpFilter(embeddedWorkspacePathPattern, false)

// ForbidGitDirComponentFilter ...
var ForbidGitDirComponentFilter = pathutil.ComponentFilter(gitDirName, false)

// AllowPodfileBaseFilter ...
var AllowPodfileBaseFilter = pathutil.BaseFilter(podfileBase, true)

// ForbidPodsDirComponentFilter ...
var ForbidPodsDirComponentFilter = pathutil.ComponentFilter(podsDirName, false)

// ForbidCarthageDirComponentFilter ...
var ForbidCarthageDirComponentFilter = pathutil.ComponentFilter(carthageDirName, false)

// ForbidCordovaLibDirComponentFilter ...
var ForbidCordovaLibDirComponentFilter = pathutil.ComponentFilter(cordovaLibDirName, false)

// ForbidFramworkComponentWithExtensionFilter ...
var ForbidFramworkComponentWithExtensionFilter = pathutil.ComponentWithExtensionFilter(frameworkExt, false)

// ForbidNodeModulesComponentFilter ...
var ForbidNodeModulesComponentFilter = pathutil.ComponentFilter(nodeModulesDirName, false)

// AllowIphoneosSDKFilter ...
var AllowIphoneosSDKFilter = SDKFilter("iphoneos", true)

// AllowMacosxSDKFilter ...
var AllowMacosxSDKFilter = SDKFilter("macosx", true)

// SDKFilter ...
func SDKFilter(sdk string, allowed bool) pathutil.FilterFunc {
	return func(pth string) (bool, error) {
		found := false

		projectFiles := []string{}

		if xcodeproj.IsXCodeProj(pth) {
			projectFiles = append(projectFiles, pth)
		} else if xcodeproj.IsXCWorkspace(pth) {
			projects, err := xcodeproj.WorkspaceProjectReferences(pth)
			if err != nil {
				return false, err
			}

			for _, project := range projects {
				exist, err := pathutil.IsPathExists(project)
				if err != nil {
					return false, err
				}
				if !exist {
					continue
				}
				projectFiles = append(projectFiles, project)

			}
		} else {
			return false, fmt.Errorf("Not Xcode project nor workspace file: %s", pth)
		}

		for _, projectFile := range projectFiles {
			pbxprojPth := filepath.Join(projectFile, "project.pbxproj")
			projectSDKs, err := xcodeproj.GetBuildConfigSDKs(pbxprojPth)
			if err != nil {
				return false, err
			}

			for _, projectSDK := range projectSDKs {
				if projectSDK == sdk {
					found = true
					break
				}
			}
		}

		return (allowed == found), nil
	}
}
