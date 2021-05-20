package utility

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/bitrise-io/go-utils/pathutil"
)

// ListEntries ...
func ListEntries(dir string, filters ...pathutil.FilterFunc) ([]string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return []string{}, err
	}

	entries, err := ioutil.ReadDir(absDir)
	if err != nil {
		return []string{}, err
	}

	var paths []string
	for _, entry := range entries {
		pth := filepath.Join(absDir, entry.Name())
		paths = append(paths, pth)
	}

	return pathutil.FilterPaths(paths, filters...)
}

// FindFileInAppDir ...
func FindFileInAppDir(appDir, fileName string) (string, error) {
	filePth := filepath.Join(appDir, fileName)
	if exist, err := pathutil.IsPathExists(filePth); err != nil {
		return "", err
	} else if exist {
		return filePth, nil
	}
	// ---

	// It's somewhere else - let's find it!
	apps, err := ListEntries(appDir, pathutil.ExtensionFilter(".app", true))
	if err != nil {
		return "", err
	}

	for _, app := range apps {
		pths, err := ListEntries(app, pathutil.BaseFilter(fileName, true))
		if err != nil {
			return "", err
		}

		if len(pths) > 0 {
			return pths[0], nil
		}
	}
	// ---

	return "", fmt.Errorf("failed to find %s", fileName)
}
