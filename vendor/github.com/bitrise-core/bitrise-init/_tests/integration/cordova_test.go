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

func TestCordova(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__cordova__")
	require.NoError(t, err)

	t.Log("sample-apps-cordova-with-jasmine")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-cordova-with-jasmine")
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-cordova-with-jasmine.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsCordovaWithJasmineResultYML), strings.TrimSpace(result))
	}

	t.Log("sample-apps-cordova-with-karma-jasmine")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-cordova-with-karma-jasmine")
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-cordova-with-karma-jasmine.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsCordovaWithKarmaJasmineResultYML), strings.TrimSpace(result))
	}
}

var sampleAppsCordovaWithJasmineVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.NpmVersion,
	steps.JasmineTestRunnerVersion,
	steps.GenerateCordovaBuildConfigVersion,
	steps.CordovaArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.NpmVersion,
	steps.JasmineTestRunnerVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsCordovaWithJasmineResultYML = fmt.Sprintf(`options:
  cordova:
    title: Platform to use in cordova-cli commands
    env_key: CORDOVA_PLATFORM
    value_map:
      android:
        config: cordova-config
      ios:
        config: cordova-config
      ios,android:
        config: cordova-config
configs:
  cordova:
    cordova-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: cordova
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
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - npm@%s:
              inputs:
              - command: install
          - jasmine-runner@%s: {}
          - generate-cordova-build-configuration@%s: {}
          - cordova-archive@%s:
              inputs:
              - platform: $CORDOVA_PLATFORM
              - target: emulator
          - deploy-to-bitrise-io@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - npm@%s:
              inputs:
              - command: install
          - jasmine-runner@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  cordova: []
`, sampleAppsCordovaWithJasmineVersions...)

var sampleAppsCordovaWithKarmaJasmineVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.NpmVersion,
	steps.KarmaJasmineTestRunnerVersion,
	steps.GenerateCordovaBuildConfigVersion,
	steps.CordovaArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.NpmVersion,
	steps.KarmaJasmineTestRunnerVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsCordovaWithKarmaJasmineResultYML = fmt.Sprintf(`options:
  cordova:
    title: Platform to use in cordova-cli commands
    env_key: CORDOVA_PLATFORM
    value_map:
      android:
        config: cordova-config
      ios:
        config: cordova-config
      ios,android:
        config: cordova-config
configs:
  cordova:
    cordova-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: cordova
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
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - npm@%s:
              inputs:
              - command: install
          - karma-jasmine-runner@%s: {}
          - generate-cordova-build-configuration@%s: {}
          - cordova-archive@%s:
              inputs:
              - platform: $CORDOVA_PLATFORM
              - target: emulator
          - deploy-to-bitrise-io@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - npm@%s:
              inputs:
              - command: install
          - karma-jasmine-runner@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  cordova: []
  `, sampleAppsCordovaWithKarmaJasmineVersions...)
