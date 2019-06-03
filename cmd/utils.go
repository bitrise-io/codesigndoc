package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/bitrise-io/bitrise-init/scanners/ios"
	"github.com/bitrise-io/bitrise-init/scanners/xamarin"
	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
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

// findProject scans the directory for Xcode Project (.xcworkspace / .xcodeproject) file first
// If can't find any, ask the user to drag-and-drop the file
func findXcodeProject() (string, error) {
	var projpth string

	projPaths, err := scanForProjectFiles(iOSProjectType)
	if err != nil {
		log.Printf("Failed: %s", err)
		fmt.Println()

		log.Infof("Provide the project file manually")
		askText := `Please drag-and-drop your Xcode Project (` + colorstring.Green(".xcodeproj") + `) or Workspace (` + colorstring.Green(".xcworkspace") + `) file, 
the one you usually open in Xcode, then hit Enter.
(Note: if you have a Workspace file you should most likely use that)`
		projpth, err = goinp.AskForPath(askText)
		if err != nil {
			return "", fmt.Errorf("failed to read input: %s", err)
		}

		return projpth, nil
	}

	if len(projPaths) == 1 {
		log.Printf("Found one project file: %s.", path.Base(projPaths[0]))
		return projPaths[0], nil
	}

	log.Printf("Found multiple project file: %s.", path.Base(projpth))
	projpth, err = goinp.SelectFromStringsWithDefault("Select the project file you want to scan", 1, projPaths)
	if err != nil {
		return "", fmt.Errorf("failed to select project file: %s", err)
	}

	return projpth, nil
}

// findSolution scans the directory for Xamarin.Solution file first
// If can't find any, ask the user to drag-and-drop the file
func findXamarinSolution() (string, error) {
	var solutionPth string
	solPaths, err := scanForProjectFiles(xamarinProjectType)
	if err != nil {
		log.Printf("Failed: %s", err)
		fmt.Println()

		log.Infof("Provide the solution file manually")
		askText := `Please drag-and-drop your Xamarin Solution (` + colorstring.Green(".sln") + `) file,
and then hit Enter`
		solutionPth, err = goinp.AskForPath(askText)
		if err != nil {
			return "", fmt.Errorf("failed to read input: %s", err)
		}

		return solutionPth, nil
	}

	if len(solPaths) == 1 {
		log.Printf("Found one solution file: %s.", path.Base(solPaths[0]))
		return solPaths[0], nil
	}

	log.Printf("Found multiple solution file: %s.", path.Base(solutionPth))
	solutionPth, err = goinp.SelectFromStringsWithDefault("Select the solution file you want to scan", 1, solPaths)
	if err != nil {
		return "", fmt.Errorf("failed to select solution file: %s", err)
	}

	return solutionPth, nil
}
