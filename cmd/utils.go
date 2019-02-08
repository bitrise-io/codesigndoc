package cmd

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-init/scanners/ios"
	"github.com/bitrise-core/bitrise-init/scanners/xamarin"
	"github.com/bitrise-core/bitrise-init/utility"
)

// ProjectType ...
type ProjectType int

const (
	iOSProjectType ProjectType = iota
	xamarinProjectType
)

// Scans the root dir for the provided project files
func scanForProjectFiles(projectType ProjectType) ([]string, error) {
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
		if projectType == iOSProjectType {
			paths, err = ios.FilterRelevantWorkspaceFiles(fileList)
			if err != nil {
				return nil, fmt.Errorf("failed to search for solution files, error: %s", err)
			}

			if len(paths) == 0 {
				paths, err = ios.FilterRelevantProjectFiles(fileList)
				if err != nil {
					return nil, fmt.Errorf("failed to search for solution files, error: %s", err)
				}
			}
		} else if projectType == xamarinProjectType {
			paths, err = xamarin.FilterSolutionFiles(fileList)
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
