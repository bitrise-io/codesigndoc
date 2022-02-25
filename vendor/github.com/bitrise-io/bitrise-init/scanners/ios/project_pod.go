package ios

import (
	"fmt"

	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
)

type podWorkspace struct {
	workspacePath     string
	workspaceProjects []container
}

func newPodWorkspace(path string, projects []container) podWorkspace {
	return podWorkspace{
		workspacePath:     path,
		workspaceProjects: projects,
	}
}

func (w podWorkspace) isWorkspace() bool {
	return true
}

func (w podWorkspace) path() string {
	return w.workspacePath
}

func (w podWorkspace) schemes() (map[string][]xcscheme.Scheme, error) {
	projectToSchemes := make(map[string][]xcscheme.Scheme)

	for _, p := range w.workspaceProjects {
		innerSchemes, err := p.schemes()
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}

		for path, schemes := range innerSchemes {
			projectToSchemes[path] = schemes
		}
	}

	return projectToSchemes, nil
}

func (w podWorkspace) projects() ([]xcodeproj.XcodeProj, []string, error) {
	var innerProjects []xcodeproj.XcodeProj
	for _, p := range w.workspaceProjects {
		projects, _, err := p.projects()
		if err != nil {
			return nil, nil, err
		}

		innerProjects = append(innerProjects, projects...)
	}

	return innerProjects, nil, nil
}

func (w podWorkspace) projectPaths() ([]string, error) {
	var paths []string

	for _, p := range w.workspaceProjects {
		paths = append(paths, p.path())
	}

	return paths, nil
}
