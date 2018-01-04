package builder

import (
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/analyzers/solution"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/utility"
	"github.com/stretchr/testify/require"
)

func TestValidateSolutionPth(t *testing.T) {
	t.Log("it validates solution path")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		solutionPth := filepath.Join(tmpDir, "solution.sln")
		require.NoError(t, fileutil.WriteStringToFile(solutionPth, "solution"))
		require.NoError(t, validateSolutionPth(solutionPth))
	}

	t.Log("it fails if file not exist")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		solutionPth := filepath.Join(tmpDir, "solution.sln")
		require.Error(t, validateSolutionPth(solutionPth))
	}

	t.Log("it fails if path is not solution path")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		projectPth := filepath.Join(tmpDir, "project.csproj")
		require.Error(t, validateSolutionPth(projectPth))
	}
}

func TestValidateSolutionConfig(t *testing.T) {
	t.Log("it validates if solution config exist")
	{
		configuration := "Release"
		platform := "iPhone"
		config := utility.ToConfig(configuration, platform)

		solution := solution.Model{
			ConfigMap: map[string]string{
				config: config,
			},
		}

		require.NoError(t, validateSolutionConfig(solution, configuration, platform))
	}

	t.Log("it fails if solution config not exist")
	{
		configuration := "Release"
		platform := "iPhone"
		config := utility.ToConfig(configuration, platform)

		solution := solution.Model{
			ConfigMap: map[string]string{
				config: config,
			},
		}

		require.Error(t, validateSolutionConfig(solution, configuration, "Any CPU"))
	}
}

func TestWhitelistAllows(t *testing.T) {
	t.Log("empty whitelist means allow any project type")
	{
		whitelist := []constants.SDK{}
		require.Equal(t, true, whitelistAllows(constants.SDKIOS, whitelist...))
	}

	t.Log("it allows project type that exists in whitelist")
	{
		whitelist := []constants.SDK{constants.SDKIOS}
		require.Equal(t, true, whitelistAllows(constants.SDKIOS, whitelist...))
	}

	t.Log("it allows project type that exists in whitelist")
	{
		whitelist := []constants.SDK{constants.SDKAndroid, constants.SDKIOS}
		require.Equal(t, true, whitelistAllows(constants.SDKIOS, whitelist...))
	}

	t.Log("it allows project type that exists in whitelist")
	{
		whitelist := []constants.SDK{constants.SDKAndroid, constants.SDKIOS}
		require.Equal(t, true, whitelistAllows(constants.SDKAndroid, whitelist...))
	}

	t.Log("it does not allows project type that does not exists in whitelist")
	{
		whitelist := []constants.SDK{constants.SDKIOS}
		require.Equal(t, false, whitelistAllows(constants.SDKAndroid, whitelist...))
	}
}

func TestIsArchitectureArchiveablet(t *testing.T) {
	t.Log("default architectures is armv7")
	{
		require.Equal(t, true, IsDeviceArch())
	}

	t.Log("arm architectures are archivables")
	{
		require.Equal(t, true, IsDeviceArch("armv7"))
	}

	t.Log("it is case insensitive")
	{
		require.Equal(t, true, IsDeviceArch("ARM7"))
	}

	t.Log("x86 architectures are not archivables")
	{
		require.Equal(t, false, IsDeviceArch("x86"))
	}
}

func TestIsPlatformAnyCPU(t *testing.T) {
	t.Log("true for Any CPU")
	{
		require.Equal(t, true, isPlatformAnyCPU("Any CPU"))
	}

	t.Log("true for AnyCPU")
	{
		require.Equal(t, true, isPlatformAnyCPU("AnyCPU"))
	}

	t.Log("false for other platforms")
	{
		require.Equal(t, false, isPlatformAnyCPU("iPhone"))
	}
}

func TestAndroidPackageNameFromManifestContent(t *testing.T) {
	t.Log("it finds package name in manifest")
	{
		packageName, err := androidPackageNameFromManifestContent(manifestFileContent)
		require.NoError(t, err)
		require.Equal(t, "hu.bitrise.test", packageName)
	}
}
