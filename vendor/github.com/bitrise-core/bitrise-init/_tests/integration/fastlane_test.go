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

func TestFastlane(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("fastlane")
	require.NoError(t, err)

	t.Log("fastlane")
	{
		sampleAppDir := filepath.Join(tmpDir, "__fastlane__")
		sampleAppURL := "https://github.com/bitrise-samples/fastlane.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(fastlaneResultYML), strings.TrimSpace(result))
	}
}

var fastlaneVersions = []interface{}{
	// fastlane
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FastlaneVersion,
	steps.DeployToBitriseIoVersion,

	// ios
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
}

var fastlaneResultYML = fmt.Sprintf(`options:
  fastlane:
    title: Project type
    value_map:
      ios:
        title: Working directory
        env_key: FASTLANE_WORK_DIR
        value_map:
          BitriseFastlaneSample:
            title: Fastlane lane
            env_key: FASTLANE_LANE
            value_map:
              ios test:
                config: fastlane-config_ios
  ios:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      BitriseFastlaneSample/BitriseFastlaneSample.xcodeproj:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          BitriseFastlaneSample:
            title: ipa export method
            env_key: BITRISE_EXPORT_METHOD
            value_map:
              ad-hoc:
                config: ios-test-config
              app-store:
                config: ios-test-config
              development:
                config: ios-test-config
              enterprise:
                config: ios-test-config
configs:
  fastlane:
    fastlane-config_ios: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      app:
        envs:
        - FASTLANE_XCODE_LIST_TIMEOUT: "120"
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
          - fastlane@%s:
              inputs:
              - lane: $FASTLANE_LANE
              - work_dir: $FASTLANE_WORK_DIR
          - deploy-to-bitrise-io@%s: {}
  ios:
    ios-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
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
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - xcode-archive@%s:
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
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  fastlane: []
  ios: []
`, fastlaneVersions...)
