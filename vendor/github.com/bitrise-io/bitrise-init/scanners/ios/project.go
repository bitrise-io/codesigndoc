package ios

import (
	"fmt"

	"github.com/bitrise-io/go-utils/pathutil"
	xcodeproject "github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcworkspace"
)

type container interface {
	isWorkspace() bool
	// path returns container absolute path
	path() string
	// schemes returns schemes mapped to the project or workspace path
	schemes() (map[string][]xcscheme.Scheme, error)
	projects() ([]xcodeproject.XcodeProj, []string, error)
	projectPaths() ([]string, error)
}

// projectContainer ...
type projectContainer struct {
	project xcodeproject.XcodeProj
}

func newProject(path string) (projectContainer, error) {
	if !xcodeproject.IsXcodeProj(path) {
		return projectContainer{}, fmt.Errorf("path (%s) is not a Project", path)
	}

	project, err := xcodeproject.Open(path)
	if err != nil {
		return projectContainer{}, fmt.Errorf("failed to open Project (%s): %v", path, err)
	}

	return projectContainer{
		project: project,
	}, nil
}

func (p projectContainer) isWorkspace() bool {
	return false
}

func (p projectContainer) path() string {
	return p.project.Path
}

func (p projectContainer) schemes() (map[string][]xcscheme.Scheme, error) {
	projectSchemes, err := p.project.Schemes()
	if err != nil {
		return nil, fmt.Errorf("failed to list Schemes in Project (%s): %v", p.project.Path, err)
	}

	containerToSchemes := make(map[string][]xcscheme.Scheme)
	containerToSchemes[p.project.Path] = projectSchemes

	return containerToSchemes, nil
}

func (p projectContainer) projects() ([]xcodeproject.XcodeProj, []string, error) {
	return []xcodeproject.XcodeProj{p.project}, []string{}, nil
}

func (p projectContainer) projectPaths() ([]string, error) {
	return []string{p.project.Path}, nil
}

// workspaceContainer ...
type workspaceContainer struct {
	workspace xcworkspace.Workspace
}

func newWorkspace(path string) (workspaceContainer, error) {
	if !xcworkspace.IsWorkspace(path) {
		return workspaceContainer{}, fmt.Errorf("path (%s) is not a Workspace", path)
	}

	workspace, err := xcworkspace.Open(path)
	if err != nil {
		return workspaceContainer{}, fmt.Errorf("failed to open Workspace (%s): %v", path, err)
	}

	return workspaceContainer{
		workspace: workspace,
	}, nil
}

func (w workspaceContainer) isWorkspace() bool {
	return true
}

func (w workspaceContainer) path() string {
	return w.workspace.Path
}

func (w workspaceContainer) schemes() (map[string][]xcscheme.Scheme, error) {
	containerToSchemes, err := w.workspace.Schemes()
	if err != nil {
		return nil, fmt.Errorf("failed to list Schemes in Workspace (%s): %v", w.workspace.Path, err)
	}

	return containerToSchemes, nil
}

func (w workspaceContainer) projects() ([]xcodeproject.XcodeProj, []string, error) {
	projPaths, err := w.workspace.ProjectFileLocations()
	if err != nil {
		return nil, nil, err
	}

	var projects []xcodeproject.XcodeProj
	var missingProjects []string
	for _, projPath := range projPaths {
		if exist, err := pathutil.IsPathExists(projPath); err != nil {
			return nil, nil, fmt.Errorf("failed to list Projects in the Workspace (%s), can not check if path (%s) exists: %v", w.workspace.Path, projPath, err)
		} else if !exist {
			missingProjects = append(missingProjects, projPath)
			continue
		}

		project, err := xcodeproject.Open(projPath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open Project (%s) in the Workspace (%s): %v", projPath, w.workspace.Path, err)
		}

		projects = append(projects, project)
	}

	return projects, missingProjects, nil
}

func (w workspaceContainer) projectPaths() ([]string, error) {
	return w.workspace.ProjectFileLocations()
}

func numberOfSharedSchemes(containerToSchemes map[string][]xcscheme.Scheme) int {
	var count int
	for _, schemes := range containerToSchemes {
		for _, scheme := range schemes {
			if scheme.IsShared {
				count++
			}
		}
	}

	return count
}
