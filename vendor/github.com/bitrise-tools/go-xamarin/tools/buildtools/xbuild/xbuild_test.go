package xbuild

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/testutil"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Log("it create new xbuild model")
	{
		currentDir, err := pathutil.CurrentWorkingDirectoryAbsolutePath()
		require.NoError(t, err)

		xbuild, err := New("solution.sln", "")
		require.NoError(t, err)
		require.NotNil(t, xbuild)

		require.Equal(t, constants.XbuildPath, xbuild.BuildTool)
		require.Equal(t, filepath.Join(currentDir, "solution.sln"), xbuild.SolutionPth)
		require.Equal(t, "", xbuild.configuration)
		require.Equal(t, "", xbuild.platform)
		require.Equal(t, "", xbuild.target)

		require.Equal(t, false, xbuild.buildIpa)
		require.Equal(t, false, xbuild.archiveOnBuild)

		require.Equal(t, 0, len(xbuild.customOptions))
	}

	t.Log("it create new xbuild model")
	{
		currentDir, err := pathutil.CurrentWorkingDirectoryAbsolutePath()
		require.NoError(t, err)

		xbuild, err := New("solution.sln", "project.csproj")
		require.NoError(t, err)
		require.NotNil(t, xbuild)

		require.Equal(t, constants.XbuildPath, xbuild.BuildTool)
		require.Equal(t, filepath.Join(currentDir, "solution.sln"), xbuild.SolutionPth)
		require.Equal(t, filepath.Join(currentDir, "project.csproj"), xbuild.ProjectPth)
		require.Equal(t, "", xbuild.configuration)
		require.Equal(t, "", xbuild.platform)
		require.Equal(t, "", xbuild.target)

		require.Equal(t, false, xbuild.buildIpa)
		require.Equal(t, false, xbuild.archiveOnBuild)

		require.Equal(t, 0, len(xbuild.customOptions))
	}
}

func TestSetProperties(t *testing.T) {
	t.Log("it sets target")
	{
		xbuild, err := New("/solution.sln", "")
		require.NoError(t, err)
		require.NotNil(t, xbuild)
		require.Equal(t, "", xbuild.target)

		xbuild.SetTarget("Build")
		require.Equal(t, "Build", xbuild.target)
	}

	t.Log("it sets configuration")
	{
		xbuild, err := New("/solution.sln", "")
		require.NoError(t, err)
		require.NotNil(t, xbuild)
		require.Equal(t, "", xbuild.configuration)

		xbuild.SetConfiguration("Release")
		require.Equal(t, "Release", xbuild.configuration)
	}

	t.Log("it sets platform")
	{
		xbuild, err := New("/solution.sln", "")
		require.NoError(t, err)
		require.NotNil(t, xbuild)
		require.Equal(t, "", xbuild.platform)

		xbuild.SetPlatform("iPhone")
		require.Equal(t, "iPhone", xbuild.platform)
	}

	t.Log("it sets build ipa")
	{
		xbuild, err := New("/solution.sln", "")
		require.NoError(t, err)
		require.NotNil(t, xbuild)
		require.Equal(t, false, xbuild.buildIpa)

		xbuild.SetBuildIpa(true)
		require.Equal(t, true, xbuild.buildIpa)
	}

	t.Log("it sets archive on build")
	{
		xbuild, err := New("/solution.sln", "")
		require.NoError(t, err)
		require.NotNil(t, xbuild)
		require.Equal(t, false, xbuild.archiveOnBuild)

		xbuild.SetArchiveOnBuild(true)
		require.Equal(t, true, xbuild.archiveOnBuild)
	}

	t.Log("it appends custom options")
	{
		xbuild, err := New("/solution.sln", "")
		require.NoError(t, err)
		require.NotNil(t, xbuild)
		require.Equal(t, 0, len(xbuild.customOptions))

		customOptions := []string{"/verbosity:minimal", "/nologo"}
		xbuild.SetCustomOptions(customOptions...)
		testutil.EqualSlicesWithoutOrder(t, customOptions, xbuild.customOptions)
	}
}

func TestBuildCommandSlice(t *testing.T) {
	t.Log("solution-dir test")
	{
		currentDir, err := pathutil.CurrentWorkingDirectoryAbsolutePath()
		require.NoError(t, err)

		xbuild, err := New("./test/solution.sln", "./test/ios/project.csproj")
		require.NoError(t, err)
		desired := []string{constants.XbuildPath, filepath.Join(currentDir, "test/ios/project.csproj"), fmt.Sprintf("/p:SolutionDir=%s", filepath.Join(currentDir, "test"))}
		require.Equal(t, desired, xbuild.buildCommandSlice())
	}

	t.Log("solution-dir test")
	{
		xbuild, err := New("/Users/Develop/test/solution.sln", "/Users/Develop/test/test/ios/project.csproj")
		require.NoError(t, err)
		desired := []string{constants.XbuildPath, "/Users/Develop/test/test/ios/project.csproj", "/p:SolutionDir=/Users/Develop/test"}
		require.Equal(t, desired, xbuild.buildCommandSlice())
	}

	t.Log("it build command slice from model")
	{
		xbuild, err := New("/solution.sln", "")
		require.NoError(t, err)
		desired := []string{constants.XbuildPath, "/solution.sln", "/p:SolutionDir=/"}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetTarget("Build")
		desired = []string{constants.XbuildPath, "/solution.sln", "/target:Build", "/p:SolutionDir=/"}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetConfiguration("Release")
		desired = []string{constants.XbuildPath, "/solution.sln", "/target:Build", "/p:SolutionDir=/", "/p:Configuration=Release"}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetPlatform("iPhone")
		desired = []string{constants.XbuildPath, "/solution.sln", "/target:Build", "/p:SolutionDir=/", "/p:Configuration=Release", "/p:Platform=iPhone"}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetArchiveOnBuild(true)
		desired = []string{constants.XbuildPath, "/solution.sln", "/target:Build", "/p:SolutionDir=/", "/p:Configuration=Release", "/p:Platform=iPhone", "/p:ArchiveOnBuild=true"}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetBuildIpa(true)
		desired = []string{constants.XbuildPath, "/solution.sln", "/target:Build", "/p:SolutionDir=/", "/p:Configuration=Release", "/p:Platform=iPhone", "/p:ArchiveOnBuild=true", "/p:BuildIpa=true"}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetCustomOptions("/nologo")
		desired = []string{constants.XbuildPath, "/solution.sln", "/target:Build", "/p:SolutionDir=/", "/p:Configuration=Release", "/p:Platform=iPhone", "/p:ArchiveOnBuild=true", "/p:BuildIpa=true", "/nologo"}
		require.Equal(t, desired, xbuild.buildCommandSlice())
	}
}

func TestPrintableCommand(t *testing.T) {
	t.Log("solution-dir test")
	{
		currentDir, err := pathutil.CurrentWorkingDirectoryAbsolutePath()
		require.NoError(t, err)

		xbuild, err := New("./test/solution.sln", "./test/ios/project.csproj")
		require.NoError(t, err)
		desired := fmt.Sprintf(`"%s" "%s" "%s"`, constants.XbuildPath, filepath.Join(currentDir, "test/ios/project.csproj"), fmt.Sprintf("/p:SolutionDir=%s", filepath.Join(currentDir, "test")))
		require.Equal(t, desired, xbuild.PrintableCommand())
	}

	t.Log("solution-dir test")
	{
		xbuild, err := New("/Users/Develop/test/solution.sln", "/Users/Develop/test/test/ios/project.csproj")
		require.NoError(t, err)
		desired := fmt.Sprintf(`"%s" "/Users/Develop/test/test/ios/project.csproj" "/p:SolutionDir=/Users/Develop/test"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())
	}

	t.Log("it creates printable command")
	{
		xbuild, err := New("/solution.sln", "")
		require.NoError(t, err)
		desired := fmt.Sprintf(`"%s" "/solution.sln" "/p:SolutionDir=/"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetTarget("Build")
		desired = fmt.Sprintf(`"%s" "/solution.sln" "/target:Build" "/p:SolutionDir=/"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetConfiguration("Release")
		desired = fmt.Sprintf(`"%s" "/solution.sln" "/target:Build" "/p:SolutionDir=/" "/p:Configuration=Release"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetPlatform("iPhone")
		desired = fmt.Sprintf(`"%s" "/solution.sln" "/target:Build" "/p:SolutionDir=/" "/p:Configuration=Release" "/p:Platform=iPhone"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetArchiveOnBuild(true)
		desired = fmt.Sprintf(`"%s" "/solution.sln" "/target:Build" "/p:SolutionDir=/" "/p:Configuration=Release" "/p:Platform=iPhone" "/p:ArchiveOnBuild=true"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetBuildIpa(true)
		desired = fmt.Sprintf(`"%s" "/solution.sln" "/target:Build" "/p:SolutionDir=/" "/p:Configuration=Release" "/p:Platform=iPhone" "/p:ArchiveOnBuild=true" "/p:BuildIpa=true"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetCustomOptions("/nologo")
		desired = fmt.Sprintf(`"%s" "/solution.sln" "/target:Build" "/p:SolutionDir=/" "/p:Configuration=Release" "/p:Platform=iPhone" "/p:ArchiveOnBuild=true" "/p:BuildIpa=true" "/nologo"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())
	}
}
