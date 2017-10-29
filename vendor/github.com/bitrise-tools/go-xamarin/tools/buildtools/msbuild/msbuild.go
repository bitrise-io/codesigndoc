package msbuild

import (
	"fmt"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/tools/buildtools/xbuild"
)

// New ...
func New(solutionPth, projectPth string) (*xbuild.Model, error) {
	absSolutionPth, err := pathutil.AbsPath(solutionPth)
	if err != nil {
		return nil, fmt.Errorf("Failed to expand path (%s), error: %s", solutionPth, err)
	}

	absProjectPth := ""
	if projectPth != "" {
		absPth, err := pathutil.AbsPath(projectPth)
		if err != nil {
			return nil, fmt.Errorf("Failed to expand path (%s), error: %s", projectPth, err)
		}
		absProjectPth = absPth
	}

	return &xbuild.Model{SolutionPth: absSolutionPth, ProjectPth: absProjectPth, BuildTool: constants.MsbuildPath}, nil
}
