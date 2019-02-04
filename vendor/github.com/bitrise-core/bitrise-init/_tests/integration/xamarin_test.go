package integration

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestXamarin(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__xamarin__")
	require.NoError(t, err)

	t.Log("xamarin-sample-app")
	{
		sampleAppDir := filepath.Join(tmpDir, "xamarin-sample-app")
		sampleAppURL := "https://github.com/bitrise-samples/xamarin-sample-app.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(xamarinSampleAppResultYML), strings.TrimSpace(result))
	}

	t.Log("sample-apps-xamarin-ios")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-xamarin-ios")
		sampleAppURL := "https://github.com/bitrise-io/sample-apps-xamarin-ios.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsXamarinIosResultYML), strings.TrimSpace(result))
	}

	t.Log("sample-apps-xamarin-android")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-xamarin-android")
		sampleAppURL := "https://github.com/bitrise-io/sample-apps-xamarin-android.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsXamarinAndroidResultYML), strings.TrimSpace(result))
	}
}

var xamarinSampleAppVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XamarinUserManagementVersion,
	steps.NugetRestoreVersion,
	steps.XamarinComponentsRestoreVersion,
	steps.XamarinArchiveVersion,
	steps.DeployToBitriseIoVersion,
}

var xamarinSampleAppResultYML = fmt.Sprintf(`options:
  xamarin:
    title: Path to the Xamarin Solution file
    env_key: BITRISE_PROJECT_PATH
    value_map:
      XamarinSampleApp.sln:
        title: Xamarin solution configuration
        env_key: BITRISE_XAMARIN_CONFIGURATION
        value_map:
          Debug:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-components-config
              iPhone:
                config: xamarin-nuget-components-config
              iPhoneSimulator:
                config: xamarin-nuget-components-config
          Release:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-components-config
              iPhone:
                config: xamarin-nuget-components-config
              iPhoneSimulator:
                config: xamarin-nuget-components-config
configs:
  xamarin:
    xamarin-nuget-components-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: xamarin
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - xamarin-user-management@%s:
              run_if: .IsCI
          - nuget-restore@%s: {}
          - xamarin-components-restore@%s: {}
          - xamarin-archive@%s:
              inputs:
              - xamarin_solution: $BITRISE_PROJECT_PATH
              - xamarin_configuration: $BITRISE_XAMARIN_CONFIGURATION
              - xamarin_platform: $BITRISE_XAMARIN_PLATFORM
          - deploy-to-bitrise-io@%s: {}
warnings:
  xamarin: []
`, xamarinSampleAppVersions...)

var sampleAppsXamarinIosVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.NugetRestoreVersion,
	steps.XamarinArchiveVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsXamarinIosResultYML = fmt.Sprintf(`options:
  xamarin:
    title: Path to the Xamarin Solution file
    env_key: BITRISE_PROJECT_PATH
    value_map:
      CreditCardValidator.iOS.sln:
        title: Xamarin solution configuration
        env_key: BITRISE_XAMARIN_CONFIGURATION
        value_map:
          Debug:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-config
              iPhone:
                config: xamarin-nuget-config
              iPhoneSimulator:
                config: xamarin-nuget-config
          Release:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-config
              iPhone:
                config: xamarin-nuget-config
              iPhoneSimulator:
                config: xamarin-nuget-config
configs:
  xamarin:
    xamarin-nuget-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: xamarin
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - nuget-restore@%s: {}
          - xamarin-archive@%s:
              inputs:
              - xamarin_solution: $BITRISE_PROJECT_PATH
              - xamarin_configuration: $BITRISE_XAMARIN_CONFIGURATION
              - xamarin_platform: $BITRISE_XAMARIN_PLATFORM
          - deploy-to-bitrise-io@%s: {}
warnings:
  xamarin: []
`, sampleAppsXamarinIosVersions...)

var sampleAppsXamarinAndroidVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.NugetRestoreVersion,
	steps.XamarinArchiveVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsXamarinAndroidResultYML = fmt.Sprintf(`options:
  xamarin:
    title: Path to the Xamarin Solution file
    env_key: BITRISE_PROJECT_PATH
    value_map:
      CreditCardValidator.Droid.sln:
        title: Xamarin solution configuration
        env_key: BITRISE_XAMARIN_CONFIGURATION
        value_map:
          Debug:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-config
          Release:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-config
configs:
  xamarin:
    xamarin-nuget-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: xamarin
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - nuget-restore@%s: {}
          - xamarin-archive@%s:
              inputs:
              - xamarin_solution: $BITRISE_PROJECT_PATH
              - xamarin_configuration: $BITRISE_XAMARIN_CONFIGURATION
              - xamarin_platform: $BITRISE_XAMARIN_PLATFORM
          - deploy-to-bitrise-io@%s: {}
warnings:
  xamarin: []
`, sampleAppsXamarinAndroidVersions...)
