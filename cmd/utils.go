package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xcode/xcodeproj"
)

// IDEType ...
type IDEType string

const (
	xcodeIDE IDEType = "iOS"

	//
	// iOS
	embeddedWorkspacePathPattern = `.+\.xcodeproj/.+\.xcworkspace`

	gitDirName        = ".git"
	podsDirName       = "Pods"
	carthageDirName   = "Carthage"
	cordovaLibDirName = "CordovaLib"

	frameworkExt = ".framework"

	//
	// Xamarin
	solutionExtension = ".sln"
	componentsDirName = "Components"

	// NodeModulesDirName ...
	NodeModulesDirName = "node_modules"

	solutionConfigurationStart = "GlobalSection(SolutionConfigurationPlatforms) = preSolution"
	solutionConfigurationEnd   = "EndGlobalSection"
)

// Scans the root dir for the provided project files
func scanForProjectFiles(ideType IDEType) ([]string, error) {
	searchDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fileList, err := utility.ListPathInDirSortedByComponents(searchDir, false)
	if err != nil {
		return nil, fmt.Errorf("failed to search for files in (%s), error: %s", searchDir, err)
	}

	var paths []string
	{
		if ideType == xcodeIDE {
			paths, err = FilterRelevantWorkspaceFiles(fileList)
			if err != nil {
				return nil, fmt.Errorf("failed to search for solution files, error: %s", err)
			}

			if len(paths) == 0 {
				paths, err = FilterRelevantProjectFiles(fileList)
				if err != nil {
					return nil, fmt.Errorf("failed to search for solution files, error: %s", err)
				}
			}
		} else {
			paths, err = FilterSolutionFiles(fileList)
			if err != nil {
				return nil, fmt.Errorf("failed to search for solution files, error: %s", err)
			}
		}

	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("no project file found: %s", searchDir)

	}
	return paths, nil
}

//
// iOS

// XcodeProjectType ...
type XcodeProjectType string

// AllowXcodeProjExtFilter ...
var AllowXcodeProjExtFilter = utility.ExtensionFilter(xcodeproj.XCodeProjExt, true)

// AllowXCWorkspaceExtFilter ...
var AllowXCWorkspaceExtFilter = utility.ExtensionFilter(xcodeproj.XCWorkspaceExt, true)

// AllowIsDirectoryFilter ...
var AllowIsDirectoryFilter = utility.IsDirectoryFilter(true)

// ForbidEmbeddedWorkspaceRegexpFilter ...
var ForbidEmbeddedWorkspaceRegexpFilter = utility.RegexpFilter(embeddedWorkspacePathPattern, false)

// ForbidGitDirComponentFilter ...
var ForbidGitDirComponentFilter = utility.ComponentFilter(gitDirName, false)

// ForbidPodsDirComponentFilter ...
var ForbidPodsDirComponentFilter = utility.ComponentFilter(podsDirName, false)

// ForbidCarthageDirComponentFilter ...
var ForbidCarthageDirComponentFilter = utility.ComponentFilter(carthageDirName, false)

// ForbidCordovaLibDirComponentFilter ...
var ForbidCordovaLibDirComponentFilter = utility.ComponentFilter(cordovaLibDirName, false)

// ForbidFramworkComponentWithExtensionFilter ...
var ForbidFramworkComponentWithExtensionFilter = utility.ComponentWithExtensionFilter(frameworkExt, false)

// ForbidNodeModulesComponentFilter ...
var ForbidNodeModulesComponentFilter = utility.ComponentFilter(NodeModulesDirName, false)

// AllowIphoneosSDKFilter ...
var AllowIphoneosSDKFilter = SDKFilter("iphoneos", true)

// AllowMacosxSDKFilter ...
var AllowMacosxSDKFilter = SDKFilter("macosx", true)

// SDKFilter ...
func SDKFilter(sdk string, allowed bool) utility.FilterFunc {
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

// FilterRelevantProjectFiles ...
func FilterRelevantProjectFiles(fileList []string, projectTypes ...XcodeProjectType) ([]string, error) {
	filters := []utility.FilterFunc{
		AllowXcodeProjExtFilter,
		AllowIsDirectoryFilter,
		ForbidEmbeddedWorkspaceRegexpFilter,
		ForbidGitDirComponentFilter,
		ForbidPodsDirComponentFilter,
		ForbidCarthageDirComponentFilter,
		ForbidFramworkComponentWithExtensionFilter,
		ForbidCordovaLibDirComponentFilter,
		ForbidNodeModulesComponentFilter,
	}

	return utility.FilterPaths(fileList, filters...)
}

// FilterRelevantWorkspaceFiles ...
func FilterRelevantWorkspaceFiles(fileList []string, projectTypes ...XcodeProjectType) ([]string, error) {
	filters := []utility.FilterFunc{
		AllowXCWorkspaceExtFilter,
		AllowIsDirectoryFilter,
		ForbidEmbeddedWorkspaceRegexpFilter,
		ForbidGitDirComponentFilter,
		ForbidPodsDirComponentFilter,
		ForbidCarthageDirComponentFilter,
		ForbidFramworkComponentWithExtensionFilter,
		ForbidCordovaLibDirComponentFilter,
		ForbidNodeModulesComponentFilter,
	}

	return utility.FilterPaths(fileList, filters...)
}

//
// Xamarin

var allowSolutionExtensionFilter = utility.ExtensionFilter(solutionExtension, true)
var forbidComponentsSolutionFilter = utility.ComponentFilter(componentsDirName, false)
var forbidNodeModulesDirComponentFilter = utility.ComponentFilter(NodeModulesDirName, false)

// FilterSolutionFiles ...
func FilterSolutionFiles(fileList []string) ([]string, error) {
	files, err := utility.FilterPaths(fileList,
		allowSolutionExtensionFilter,
		forbidComponentsSolutionFilter,
		forbidNodeModulesDirComponentFilter)
	if err != nil {
		return []string{}, err
	}

	return files, nil
}
