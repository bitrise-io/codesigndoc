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

func TestMacOS(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__macos__")
	require.NoError(t, err)

	t.Log("sample-apps-osx-10-11")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-osx-10-11")
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-osx-10-11.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsOSX1011ResultYML), strings.TrimSpace(result))
	}
}

var sampleAppsOSX1011Versions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeTestMacVersion,
	steps.XcodeArchiveMacVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeTestMacVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
}

var sampleAppsOSX1011ResultYML = fmt.Sprintf(`options:
  macos:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      sample-apps-osx-10-11.xcodeproj:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          sample-apps-osx-10-11:
            title: |-
              Application export method
              NOTE: `+"`none`"+` means: Export a copy of the application without re-signing.
            env_key: BITRISE_EXPORT_METHOD
            value_map:
              app-store:
                config: macos-test-config
              developer-id:
                config: macos-test-config
              development:
                config: macos-test-config
              none:
                config: macos-test-config
configs:
  macos:
    macos-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: macos
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        deploy:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - xcode-archive-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  macos: []
`, sampleAppsOSX1011Versions...)
