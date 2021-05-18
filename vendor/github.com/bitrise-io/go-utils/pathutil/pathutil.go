package pathutil

import (
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// RevokableChangeDir ...
func RevokableChangeDir(dir string) (func() error, error) {
	origDir, err := CurrentWorkingDirectoryAbsolutePath()
	if err != nil {
		return nil, err
	}

	revokeFn := func() error {
		return os.Chdir(origDir)
	}

	return revokeFn, os.Chdir(dir)
}

// ChangeDirForFunction ...
func ChangeDirForFunction(dir string, fn func()) error {
	revokeFn, err := RevokableChangeDir(dir)
	if err != nil {
		return err
	}

	fn()

	return revokeFn()
}

// IsRelativePath ...
func IsRelativePath(pth string) bool {
	if strings.HasPrefix(pth, "./") {
		return true
	}

	if strings.HasPrefix(pth, "/") {
		return false
	}

	if strings.HasPrefix(pth, "$") {
		return false
	}

	return true
}

// EnsureDirExist ...
func EnsureDirExist(dir string) error {
	exist, err := IsDirExists(dir)
	if !exist || err != nil {
		return os.MkdirAll(dir, 0777)
	}
	return nil
}

func genericIsPathExists(pth string) (os.FileInfo, bool, error) {
	if pth == "" {
		return nil, false, errors.New("No path provided")
	}
	fileInf, err := os.Lstat(pth)
	if err == nil {
		return fileInf, true, nil
	}
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	return fileInf, false, err
}

// IsPathExists ...
func IsPathExists(pth string) (bool, error) {
	_, isExists, err := genericIsPathExists(pth)
	return isExists, err
}

// PathCheckAndInfos ...
// Returns:
// 1. file info or nil
// 2. bool, indicating whether the path exists
// 3. error, if any error happens during the check
func PathCheckAndInfos(pth string) (os.FileInfo, bool, error) {
	return genericIsPathExists(pth)
}

// IsDirExists ...
func IsDirExists(pth string) (bool, error) {
	fileInf, isExists, err := genericIsPathExists(pth)
	if err != nil {
		return false, err
	}
	if !isExists {
		return false, nil
	}
	if fileInf == nil {
		return false, errors.New("No file info available")
	}
	return fileInf.IsDir(), nil
}

// AbsPath expands ENV vars and the ~ character
//	then call Go's Abs
func AbsPath(pth string) (string, error) {
	if pth == "" {
		return "", errors.New("No Path provided")
	}

	pth, err := ExpandTilde(pth)
	if err != nil {
		return "", err
	}

	return filepath.Abs(os.ExpandEnv(pth))
}

// ExpandTilde ...
func ExpandTilde(pth string) (string, error) {
	if pth == "" {
		return "", errors.New("No Path provided")
	}

	if strings.HasPrefix(pth, "~") {
		pth = strings.TrimPrefix(pth, "~")

		if len(pth) == 0 || strings.HasPrefix(pth, "/") {
			return os.ExpandEnv("$HOME" + pth), nil
		}

		splitPth := strings.Split(pth, "/")
		username := splitPth[0]

		usr, err := user.Lookup(username)
		if err != nil {
			return "", err
		}

		pathInUsrHome := strings.Join(splitPth[1:], "/")

		return filepath.Join(usr.HomeDir, pathInUsrHome), nil
	}

	return pth, nil
}

// CurrentWorkingDirectoryAbsolutePath ...
func CurrentWorkingDirectoryAbsolutePath() (string, error) {
	return filepath.Abs("./")
}

// UserHomeDir ...
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// NormalizedOSTempDirPath ...
// Creates a temp dir, and returns its path.
// If tmpDirNamePrefix is provided it'll be used
//  as the tmp dir's name prefix.
// Normalized: it's guaranteed that the path won't end with '/'.
func NormalizedOSTempDirPath(tmpDirNamePrefix string) (retPth string, err error) {
	retPth, err = ioutil.TempDir("", tmpDirNamePrefix)
	if strings.HasSuffix(retPth, "/") {
		retPth = retPth[:len(retPth)-1]
	}
	return
}

// GetFileName returns the name of the file from a given path or the name of the directory if it is a directory
func GetFileName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// ListPathInDirSortedByComponents ...
func ListPathInDirSortedByComponents(searchDir string, relPath bool) ([]string, error) {
	searchDir, err := filepath.Abs(searchDir)
	if err != nil {
		return []string{}, err
	}

	var fileList []string

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
