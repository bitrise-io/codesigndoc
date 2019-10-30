package utility

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
)

// FilterFunc ...
type FilterFunc func(pth string) (bool, error)

// FilterPaths ...
func FilterPaths(paths []string, filters ...FilterFunc) ([]string, error) {
	filtered := []string{}

	for _, pth := range paths {
		allowed := true
		for _, filter := range filters {
			if allows, err := filter(pth); err != nil {
				return []string{}, err
			} else if !allows {
				allowed = false
				break
			}
		}
		if allowed {
			filtered = append(filtered, pth)
		}
	}

	return filtered, nil
}

// ListEntries ...
func ListEntries(dir string, filters ...FilterFunc) ([]string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return []string{}, err
	}

	entries, err := ioutil.ReadDir(absDir)
	if err != nil {
		return []string{}, err
	}

	paths := []string{}
	for _, entry := range entries {
		pth := filepath.Join(absDir, entry.Name())
		paths = append(paths, pth)
	}

	return FilterPaths(paths, filters...)
}

// ExtensionFilter ...
func ExtensionFilter(ext string, allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		e := filepath.Ext(pth)
		return (allowed == strings.EqualFold(ext, e)), nil
	}
}

// BaseFilter ...
func BaseFilter(base string, allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		b := filepath.Base(pth)
		return (allowed == strings.EqualFold(base, b)), nil
	}
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
	apps, err := ListEntries(appDir, ExtensionFilter(".app", true))
	if err != nil {
		return "", err
	}

	for _, app := range apps {
		pths, err := ListEntries(app, BaseFilter(fileName, true))
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
