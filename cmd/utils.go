package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectType ...
type ProjectType string

const (
	// iOS
	xcodeWorkspace ProjectType = "*.xcworkspace"
	xcodeProject   ProjectType = "*.xcodeproj"
)

// Scans the root dir for the provided project files
// If none of them in the root dir, then it will return an error
func scanForProjectFiles(projectFiles []ProjectType) (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for _, projectFile := range projectFiles {
		pathPattern := filepath.Join(root, string(projectFile))
		paths, err := filepath.Glob(pathPattern)
		if err != nil {
			return "", err
		}

		switch len(paths) {
		case 0:
			continue
		case 1:
			return paths[0], nil
		default:
			return "", fmt.Errorf("multiple project file (%s) found in the root (%s) directory: %s", projectFile, root, strings.Join(paths, "\n"))
		}
	}

	return "", fmt.Errorf("no project file found: %s", root)

}
