package ios

import (
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/pathfilters"
	"github.com/bitrise-io/go-xcode/xcodeproj"
)

// XcodeProjectType ...
type XcodeProjectType string

const (
	// XcodeProjectTypeIOS ...
	XcodeProjectTypeIOS XcodeProjectType = "ios"
	// XcodeProjectTypeMacOS ...
	XcodeProjectTypeMacOS XcodeProjectType = "macos"
)

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
	filters := []pathutil.FilterFunc{
		pathfilters.AllowXcodeProjExtFilter,
		pathfilters.AllowIsDirectoryFilter,
		pathfilters.ForbidEmbeddedWorkspaceRegexpFilter,
		pathfilters.ForbidGitDirComponentFilter,
		pathfilters.ForbidPodsDirComponentFilter,
		pathfilters.ForbidCarthageDirComponentFilter,
		pathfilters.ForbidFramworkComponentWithExtensionFilter,
		pathfilters.ForbidCordovaLibDirComponentFilter,
		pathfilters.ForbidNodeModulesComponentFilter,
	}

	for _, projectType := range projectTypes {
		switch projectType {
		case XcodeProjectTypeIOS:
			filters = append(filters, pathfilters.AllowIphoneosSDKFilter)
		case XcodeProjectTypeMacOS:
			filters = append(filters, pathfilters.AllowMacosxSDKFilter)
		}
	}

	return pathutil.FilterPaths(fileList, filters...)
}

// FilterRelevantWorkspaceFiles ...
func FilterRelevantWorkspaceFiles(fileList []string, projectTypes ...XcodeProjectType) ([]string, error) {
	filters := []pathutil.FilterFunc{
		pathfilters.AllowXCWorkspaceExtFilter,
		pathfilters.AllowIsDirectoryFilter,
		pathfilters.AllowWorkspaceWithContentsFile,
		pathfilters.ForbidEmbeddedWorkspaceRegexpFilter,
		pathfilters.ForbidGitDirComponentFilter,
		pathfilters.ForbidPodsDirComponentFilter,
		pathfilters.ForbidCarthageDirComponentFilter,
		pathfilters.ForbidFramworkComponentWithExtensionFilter,
		pathfilters.ForbidCordovaLibDirComponentFilter,
		pathfilters.ForbidNodeModulesComponentFilter,
	}

	for _, projectType := range projectTypes {
		switch projectType {
		case XcodeProjectTypeIOS:
			filters = append(filters, pathfilters.AllowIphoneosSDKFilter)
		case XcodeProjectTypeMacOS:
			filters = append(filters, pathfilters.AllowMacosxSDKFilter)
		}
	}

	return pathutil.FilterPaths(fileList, filters...)
}

// FilterRelevantPodfiles ...
func FilterRelevantPodfiles(fileList []string) ([]string, error) {
	return pathutil.FilterPaths(fileList,
		AllowPodfileBaseFilter,
		pathfilters.ForbidGitDirComponentFilter,
		pathfilters.ForbidPodsDirComponentFilter,
		pathfilters.ForbidCarthageDirComponentFilter,
		pathfilters.ForbidFramworkComponentWithExtensionFilter,
		pathfilters.ForbidCordovaLibDirComponentFilter,
		pathfilters.ForbidNodeModulesComponentFilter)
}

// FilterRelevantCartFile ...
func FilterRelevantCartFile(fileList []string) ([]string, error) {
	return pathutil.FilterPaths(fileList,
		AllowCartfileBaseFilter,
		pathfilters.ForbidGitDirComponentFilter,
		pathfilters.ForbidPodsDirComponentFilter,
		pathfilters.ForbidCarthageDirComponentFilter,
		pathfilters.ForbidFramworkComponentWithExtensionFilter,
		pathfilters.ForbidCordovaLibDirComponentFilter,
		pathfilters.ForbidNodeModulesComponentFilter)
}
