package utility

import (
	"os"
	"strings"
	"testing"

	"path/filepath"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestListPathInDirSortedByComponents(t *testing.T) {
	t.Log()
	{
		files, err := ListPathInDirSortedByComponents("./", true)
		require.NoError(t, err)
		require.NotEqual(t, 0, len(files))
	}

	t.Log()
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("__lis_path_test__")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.RemoveAll(tmpDir))
		}()

		pths := []string{
			filepath.Join(tmpDir, "testdir/testfile"),
			filepath.Join(tmpDir, "testdir/testdir/testfile"),
		}

		for _, pth := range pths {
			dir := filepath.Dir(pth)
			require.NoError(t, os.MkdirAll(dir, 0700))

			require.NoError(t, fileutil.WriteStringToFile(pth, "test"))
		}

		expected := []string{
			".",
			"testdir",
			"testdir/testdir",
			"testdir/testfile",
			"testdir/testdir/testfile",
		}

		files, err := ListPathInDirSortedByComponents(tmpDir, true)
		require.NoError(t, err)
		require.Equal(t, expected, files)
	}

	t.Log()
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("__lis_path_test1__")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.RemoveAll(tmpDir))
		}()

		pths := []string{
			filepath.Join(tmpDir, "testdir/testfile"),
			filepath.Join(tmpDir, "testdir/testdir/testfile"),
		}

		for _, pth := range pths {
			dir := filepath.Dir(pth)
			require.NoError(t, os.MkdirAll(dir, 0700))

			require.NoError(t, fileutil.WriteStringToFile(pth, "test"))
		}

		expected := []string{
			tmpDir,
			filepath.Join(tmpDir, "testdir"),
			filepath.Join(tmpDir, "testdir/testdir"),
			filepath.Join(tmpDir, "testdir/testfile"),
			filepath.Join(tmpDir, "testdir/testdir/testfile"),
		}

		files, err := ListPathInDirSortedByComponents(tmpDir, false)
		require.NoError(t, err)
		require.Equal(t, expected, files)
	}
}

func TestFilterPaths(t *testing.T) {
	t.Log("without any filter")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filtered, err := FilterPaths(paths)
		require.NoError(t, err)
		require.Equal(t, paths, filtered)
	}

	t.Log("with filter")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := func(pth string) (bool, error) {
			return strings.Contains(pth, "vagrant"), nil
		}
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/vagrant/test"}, filtered)
	}
}

func TestBaseFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"path/to/my/gradlew",
			"path/to/my/gradlew/file",
		}
		filter := BaseFilter("gradlew", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/gradlew"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"path/to/my/gradlew",
			"path/to/my/gradlew/file",
		}
		filter := BaseFilter("gradlew", false)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/gradlew/file"}, filtered)
	}
}

func TestExtensionFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"path/to/my/project.xcodeproj",
			"path/to/my/project.xcworkspace",
		}
		filter := ExtensionFilter(".xcodeproj", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/project.xcodeproj"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"path/to/my/project.xcodeproj",
			"path/to/my/project.xcworkspace",
		}
		filter := ExtensionFilter(".xcodeproj", false)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/project.xcworkspace"}, filtered)
	}
}

func TestRegexpFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"path/to/my/project.xcodeproj",
			"path/to/my/project.xcworkspace",
		}
		filter := RegexpFilter(".*.xcodeproj", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/project.xcodeproj"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"path/to/my/project.xcodeproj",
			"path/to/my/project.xcworkspace",
		}
		filter := RegexpFilter(".*.xcodeproj", false)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/project.xcworkspace"}, filtered)
	}
}

func TestComponentFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := ComponentFilter("bitrise", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/bitrise/test"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := ComponentFilter("bitrise", false)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/vagrant/test"}, filtered)
	}
}

func TestComponentWithExtensionFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"/Users/bitrise.framework/test",
			"/Users/vagrant/test",
		}
		filter := ComponentWithExtensionFilter(".framework", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/bitrise.framework/test"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"/Users/bitrise.framework/test",
			"/Users/vagrant/test",
		}
		filter := ComponentWithExtensionFilter(".framework", false)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/vagrant/test"}, filtered)
	}
}

func TestIsDirectoryFilter(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__bitrise-init__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	tmpFile := filepath.Join(tmpDir, "file.txt")
	require.NoError(t, fileutil.WriteStringToFile(tmpFile, ""))

	t.Log("allow")
	{
		paths := []string{
			tmpDir,
			tmpFile,
		}
		filter := IsDirectoryFilter(true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{tmpDir}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			tmpDir,
			tmpFile,
		}
		filter := IsDirectoryFilter(false)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{tmpFile}, filtered)
	}
}

func TestInDirectoryFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := InDirectoryFilter("/Users/bitrise", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/bitrise/test"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := InDirectoryFilter("/Users/bitrise", false)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/vagrant/test"}, filtered)
	}
}
