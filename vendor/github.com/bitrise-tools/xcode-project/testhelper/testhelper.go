package testhelper

import (
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

var clonedRespos = map[string]string{}

// GitCloneIntoTmpDir ...
func GitCloneIntoTmpDir(t *testing.T, repo string) string {
	if tmpDir, ok := clonedRespos[repo]; ok {
		return tmpDir
	}

	tmpDir, err := pathutil.NormalizedOSTempDirPath("__xcode-proj__")
	require.NoError(t, err)

	cmd := command.New("git", "clone", repo, tmpDir)
	require.NoError(t, cmd.Run())

	clonedRespos[repo] = tmpDir

	return tmpDir
}

// CreateTmpFile ...
func CreateTmpFile(t *testing.T, name, content string) string {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__xcode-proj__")
	require.NoError(t, err)

	pth := filepath.Join(tmpDir, name)
	require.NoError(t, fileutil.WriteStringToFile(pth, content))
	return pth
}
