package pathfilters

import (
	"fmt"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcworkspace"
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
var AllowXcodeProjExtFilter = pathutil.ExtensionFilter(xcodeproj.XcodeProjExtension, true)

// AllowXCWorkspaceExtFilter ...
var AllowXCWorkspaceExtFilter = pathutil.ExtensionFilter(xcworkspace.XCWorkspaceExtension, true)

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

		if xcodeproj.IsXcodeProj(pth) {
			projectFiles = append(projectFiles, pth)
		} else if xcworkspace.IsWorkspace(pth) {
			workspace, err := xcworkspace.Open(pth)
			if err != nil {
				return false, err
			}
			projects, err := workspace.ProjectFileLocations()
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
			return false, fmt.Errorf("not Xcode project nor workspace file: %s", pth)
		}

		sdkMap := map[string]bool{}
		for _, projectFile := range projectFiles {
			project, err := xcodeproj.Open(projectFile)
			if err != nil {
				return false, err
			}

			var buildConfigurations []xcodeproj.BuildConfiguration
			buildConfigurations = append(buildConfigurations, project.Proj.BuildConfigurationList.BuildConfigurations...)
			for _, target := range project.Proj.Targets {
				buildConfigurations = append(buildConfigurations, target.BuildConfigurationList.BuildConfigurations...)
			}

			for _, buildConfiguratioon := range buildConfigurations {
				sdk, err := buildConfiguratioon.BuildSettings.String("SDKROOT")
				if err == nil {
					sdkMap[sdk] = true
				}
			}

			for projectSDK := range sdkMap {
				if projectSDK == sdk {
					found = true
					break
				}
			}
		}

		return (allowed == found), nil
	}
}
