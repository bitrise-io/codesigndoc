package cmd

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-init/scanners/ios"
	"github.com/bitrise-core/bitrise-init/scanners/xamarin"
	"github.com/bitrise-core/bitrise-init/utility"
)

// projectType enum.
// Could be iOSProjectType = 0
// Or xamarinProjectType = 1
type projectType int

const (
	iOSProjectType projectType = iota
	xamarinProjectType
)

// Scans the root dir for the provided project files
func scanForProjectFiles(projType projectType) ([]string, error) {
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
		if projType == iOSProjectType {
			paths, err = ios.FilterRelevantWorkspaceFiles(fileList)
			if err != nil {
				return nil, fmt.Errorf("failed to search for workspace files, error: %s", err)
			}

			if len(paths) == 0 {
				paths, err = ios.FilterRelevantProjectFiles(fileList)
				if err != nil {
					return nil, fmt.Errorf("failed to search for project files, error: %s", err)
				}
			}
		} else if projType == xamarinProjectType {
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
