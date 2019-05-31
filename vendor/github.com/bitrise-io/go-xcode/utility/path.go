package utility

import (
	"io/ioutil"
	"path/filepath"
	"strings"
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
