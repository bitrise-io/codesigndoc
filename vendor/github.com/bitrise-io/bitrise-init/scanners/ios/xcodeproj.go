package ios

import (
	"path/filepath"

	"fmt"

	"github.com/bitrise-io/bitrise-init/scanners/xamarin"
	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/xcodeproj"
)

const (
	embeddedWorkspacePathPattern = `.+\.xcodeproj/.+\.xcworkspace`

	gitDirName        = ".git"
	podsDirName       = "Pods"
	carthageDirName   = "Carthage"
	cordovaLibDirName = "CordovaLib"

	frameworkExt = ".framework"
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
var AllowXcodeProjExtFilter = utility.ExtensionFilter(xcodeproj.XCodeProjExt, true)

// AllowXCWorkspaceExtFilter ...
var AllowXCWorkspaceExtFilter = utility.ExtensionFilter(xcodeproj.XCWorkspaceExt, true)

// AllowIsDirectoryFilter ...
var AllowIsDirectoryFilter = utility.IsDirectoryFilter(true)

var containsContentsXcworkspacedata = utility.DirectoryContainsFile("contents.xcworkspacedata")

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
var ForbidNodeModulesComponentFilter = utility.ComponentFilter(xamarin.NodeModulesDirName, false)

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

// FindWorkspaceInList ...
func FindWorkspaceInList(workspacePth string, workspaces []xcodeproj.WorkspaceModel) (xcodeproj.WorkspaceModel, bool) {
	for _, workspace := range workspaces {
		if workspace.Pth == workspacePth {
			return workspace, true
		}
	}
	return xcodeproj.WorkspaceModel{}, false
}

// FindProjectInList ...
func FindProjectInList(projectPth string, projects []xcodeproj.ProjectModel) (xcodeproj.ProjectModel, bool) {
	for _, project := range projects {
		if project.Pth == projectPth {
			return project, true
		}
	}
	return xcodeproj.ProjectModel{}, false
}

// RemoveProjectFromList ...
func RemoveProjectFromList(projectPth string, projects []xcodeproj.ProjectModel) []xcodeproj.ProjectModel {
	newProjects := []xcodeproj.ProjectModel{}
	for _, project := range projects {
		if project.Pth != projectPth {
			newProjects = append(newProjects, project)
		}
	}
	return newProjects
}

// ReplaceWorkspaceInList ...
func ReplaceWorkspaceInList(workspaces []xcodeproj.WorkspaceModel, workspace xcodeproj.WorkspaceModel) []xcodeproj.WorkspaceModel {
	updatedWorkspaces := []xcodeproj.WorkspaceModel{}
	for _, w := range workspaces {
		if w.Pth == workspace.Pth {
			updatedWorkspaces = append(updatedWorkspaces, workspace)
		} else {
			updatedWorkspaces = append(updatedWorkspaces, w)
		}
	}
	return updatedWorkspaces
}

// CreateStandaloneProjectsAndWorkspaces ...
func CreateStandaloneProjectsAndWorkspaces(projectFiles, workspaceFiles []string) ([]xcodeproj.ProjectModel, []xcodeproj.WorkspaceModel, error) {
	workspaces := []xcodeproj.WorkspaceModel{}
	for _, workspaceFile := range workspaceFiles {
		workspace, err := xcodeproj.NewWorkspace(workspaceFile, projectFiles...)
		if err != nil {
			return []xcodeproj.ProjectModel{}, []xcodeproj.WorkspaceModel{}, err
		}
		workspaces = append(workspaces, workspace)
	}

	standaloneProjects := []xcodeproj.ProjectModel{}
	for _, projectFile := range projectFiles {
		workspaceContains := false
		for _, workspace := range workspaces {
			_, found := FindProjectInList(projectFile, workspace.Projects)
			if found {
				workspaceContains = true
				break
			}
		}

		if !workspaceContains {
			project, err := xcodeproj.NewProject(projectFile)
			if err != nil {
				return []xcodeproj.ProjectModel{}, []xcodeproj.WorkspaceModel{}, err
			}
			standaloneProjects = append(standaloneProjects, project)
		}
	}

	return standaloneProjects, workspaces, nil
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

	for _, projectType := range projectTypes {
		switch projectType {
		case XcodeProjectTypeIOS:
			filters = append(filters, AllowIphoneosSDKFilter)
		case XcodeProjectTypeMacOS:
			filters = append(filters, AllowMacosxSDKFilter)
		}
	}

	return utility.FilterPaths(fileList, filters...)
}

// FilterRelevantWorkspaceFiles ...
func FilterRelevantWorkspaceFiles(fileList []string, projectTypes ...XcodeProjectType) ([]string, error) {
	filters := []utility.FilterFunc{
		AllowXCWorkspaceExtFilter,
		AllowIsDirectoryFilter,
		containsContentsXcworkspacedata,
		ForbidEmbeddedWorkspaceRegexpFilter,
		ForbidGitDirComponentFilter,
		ForbidPodsDirComponentFilter,
		ForbidCarthageDirComponentFilter,
		ForbidFramworkComponentWithExtensionFilter,
		ForbidCordovaLibDirComponentFilter,
		ForbidNodeModulesComponentFilter,
	}

	for _, projectType := range projectTypes {
		switch projectType {
		case XcodeProjectTypeIOS:
			filters = append(filters, AllowIphoneosSDKFilter)
		case XcodeProjectTypeMacOS:
			filters = append(filters, AllowMacosxSDKFilter)
		}
	}

	return utility.FilterPaths(fileList, filters...)
}

// FilterRelevantPodfiles ...
func FilterRelevantPodfiles(fileList []string) ([]string, error) {
	return utility.FilterPaths(fileList,
		AllowPodfileBaseFilter,
		ForbidGitDirComponentFilter,
		ForbidPodsDirComponentFilter,
		ForbidCarthageDirComponentFilter,
		ForbidFramworkComponentWithExtensionFilter,
		ForbidCordovaLibDirComponentFilter,
		ForbidNodeModulesComponentFilter)
}

// FilterRelevantCartFile ...
func FilterRelevantCartFile(fileList []string) ([]string, error) {
	return utility.FilterPaths(fileList,
		AllowCartfileBaseFilter,
		ForbidGitDirComponentFilter,
		ForbidPodsDirComponentFilter,
		ForbidCarthageDirComponentFilter,
		ForbidFramworkComponentWithExtensionFilter,
		ForbidCordovaLibDirComponentFilter,
		ForbidNodeModulesComponentFilter)
}
