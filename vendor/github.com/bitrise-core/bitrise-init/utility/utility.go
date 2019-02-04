package utility

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

// PackagesModel ...
type PackagesModel struct {
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func parsePackagesJSONContent(content string) (PackagesModel, error) {
	var packages PackagesModel
	if err := json.Unmarshal([]byte(content), &packages); err != nil {
		return PackagesModel{}, err
	}
	return packages, nil
}

// ParsePackagesJSON ...
func ParsePackagesJSON(packagesJSONPth string) (PackagesModel, error) {
	content, err := fileutil.ReadStringFromFile(packagesJSONPth)
	if err != nil {
		return PackagesModel{}, err
	}
	return parsePackagesJSONContent(content)
}

// RelPath ...
func RelPath(basePth, pth string) (string, error) {
	absBasePth, err := pathutil.AbsPath(basePth)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(absBasePth, "/private/var") {
		absBasePth = strings.TrimPrefix(absBasePth, "/private")
	}

	absPth, err := pathutil.AbsPath(pth)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(absPth, "/private/var") {
		absPth = strings.TrimPrefix(absPth, "/private")
	}

	return filepath.Rel(absBasePth, absPth)
}

// ListPathInDirSortedByComponents ...
func ListPathInDirSortedByComponents(searchDir string, relPath bool) ([]string, error) {
	searchDir, err := filepath.Abs(searchDir)
	if err != nil {
		return []string{}, err
	}

	fileList := []string{}

	if err := filepath.Walk(searchDir, func(path string, _ os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if relPath {
			rel, err := filepath.Rel(searchDir, path)
			if err != nil {
				return err
			}
			path = rel
		}

		fileList = append(fileList, path)

		return nil
	}); err != nil {
		return []string{}, err
	}
	return SortPathsByComponents(fileList)
}

// FilterPaths ...
func FilterPaths(fileList []string, filters ...FilterFunc) ([]string, error) {
	filtered := []string{}

	for _, pth := range fileList {
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

// FilterFunc ...
type FilterFunc func(string) (bool, error)

// BaseFilter ...
func BaseFilter(base string, allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		b := filepath.Base(pth)
		return (allowed == strings.EqualFold(base, b)), nil
	}
}

// ExtensionFilter ...
func ExtensionFilter(ext string, allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		e := filepath.Ext(pth)
		return (allowed == strings.EqualFold(ext, e)), nil
	}
}

// RegexpFilter ...
func RegexpFilter(pattern string, allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		re := regexp.MustCompile(pattern)
		found := re.FindString(pth) != ""
		return (allowed == found), nil
	}
}

// ComponentFilter ...
func ComponentFilter(component string, allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		found := false
		pathComponents := strings.Split(pth, string(filepath.Separator))
		for _, c := range pathComponents {
			if c == component {
				found = true
			}
		}
		return (allowed == found), nil
	}
}

// ComponentWithExtensionFilter ...
func ComponentWithExtensionFilter(ext string, allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		found := false
		pathComponents := strings.Split(pth, string(filepath.Separator))
		for _, c := range pathComponents {
			e := filepath.Ext(c)
			if e == ext {
				found = true
			}
		}
		return (allowed == found), nil
	}
}

// IsDirectoryFilter ...
func IsDirectoryFilter(allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		fileInf, err := os.Lstat(pth)
		if err != nil {
			return false, err
		}
		if fileInf == nil {
			return false, errors.New("no file info available")
		}
		return (allowed == fileInf.IsDir()), nil
	}
}

// InDirectoryFilter ...
func InDirectoryFilter(dir string, allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		in := (filepath.Dir(pth) == dir)
		return (allowed == in), nil
	}
}

// FileContains ...
func FileContains(pth, str string) (bool, error) {
	content, err := fileutil.ReadStringFromFile(pth)
	if err != nil {
		return false, err
	}

	return strings.Contains(content, str), nil
}
