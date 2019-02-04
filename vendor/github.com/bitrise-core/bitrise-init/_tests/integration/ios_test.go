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

func TestIOS(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__ios__")
	require.NoError(t, err)

	t.Log("ios-no-shared-schemes")
	{
		sampleAppDir := filepath.Join(tmpDir, "ios-no-shared-scheme")
		sampleAppURL := "https://github.com/bitrise-samples/ios-no-shared-schemes.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(iosNoSharedSchemesResultYML), strings.TrimSpace(result))
	}

	t.Log("ios-cocoapods-at-root")
	{
		sampleAppDir := filepath.Join(tmpDir, "ios-cocoapods-at-root")
		sampleAppURL := "https://github.com/bitrise-samples/ios-cocoapods-at-root.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(iosCocoapodsAtRootResultYML), strings.TrimSpace(result))
	}

	t.Log("sample-apps-ios-watchkit")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-ios-watchkit")
		sampleAppURL := "https://github.com/bitrise-io/sample-apps-ios-watchkit.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsIosWatchkitResultYML), strings.TrimSpace(result))
	}

	t.Log("sample-apps-carthage")
	{
		//
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-carthage")
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-carthage.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsCarthageResultYML), strings.TrimSpace(result))
	}
}

var iosNoSharedSchemesVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.RecreateUserSchemesVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.RecreateUserSchemesVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
}

var iosNoSharedSchemesResultYML = fmt.Sprintf(`options:
  ios:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      BitriseXcode7Sample.xcodeproj:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          BitriseXcode7Sample:
            title: ipa export method
            env_key: BITRISE_EXPORT_METHOD
            value_map:
              ad-hoc:
                config: ios-test-missing-shared-schemes-config
              app-store:
                config: ios-test-missing-shared-schemes-config
              development:
                config: ios-test-missing-shared-schemes-config
              enterprise:
                config: ios-test-missing-shared-schemes-config
configs:
  ios:
    ios-test-missing-shared-schemes-config: |
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
          - recreate-user-schemes@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
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
          - recreate-user-schemes@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  ios:
  - |-
    No shared schemes found for project: BitriseXcode7Sample.xcodeproj.
    Automatically generated schemes may differ from the ones in your project.
    Make sure to <a href="http://devcenter.bitrise.io/ios/frequent-ios-issues/#xcode-scheme-not-found">share your schemes</a> for the expected behaviour.
`, iosNoSharedSchemesVersions...)

var iosCocoapodsAtRootVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
}

var iosCocoapodsAtRootResultYML = fmt.Sprintf(`options:
  ios:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      iOSMinimalCocoaPodsSample.xcworkspace:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          iOSMinimalCocoaPodsSample:
            title: ipa export method
            env_key: BITRISE_EXPORT_METHOD
            value_map:
              ad-hoc:
                config: ios-pod-test-config
              app-store:
                config: ios-pod-test-config
              development:
                config: ios-pod-test-config
              enterprise:
                config: ios-pod-test-config
configs:
  ios:
    ios-pod-test-config: |
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
          - cocoapods-install@%s: {}
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
          - cocoapods-install@%s: {}
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  ios: []
`, iosCocoapodsAtRootVersions...)

var sampleAppsIosWatchkitVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

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

var sampleAppsIosWatchkitResultYML = fmt.Sprintf(`options:
  ios:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      watch-test.xcodeproj:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          Complication - watch-test WatchKit App:
            title: ipa export method
            env_key: BITRISE_EXPORT_METHOD
            value_map:
              ad-hoc:
                config: ios-config
              app-store:
                config: ios-config
              development:
                config: ios-config
              enterprise:
                config: ios-config
          Glance - watch-test WatchKit App:
            title: ipa export method
            env_key: BITRISE_EXPORT_METHOD
            value_map:
              ad-hoc:
                config: ios-config
              app-store:
                config: ios-config
              development:
                config: ios-config
              enterprise:
                config: ios-config
          Notification - watch-test WatchKit App:
            title: ipa export method
            env_key: BITRISE_EXPORT_METHOD
            value_map:
              ad-hoc:
                config: ios-config
              app-store:
                config: ios-config
              development:
                config: ios-config
              enterprise:
                config: ios-config
          watch-test:
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
          watch-test WatchKit App:
            title: ipa export method
            env_key: BITRISE_EXPORT_METHOD
            value_map:
              ad-hoc:
                config: ios-config
              app-store:
                config: ios-config
              development:
                config: ios-config
              enterprise:
                config: ios-config
configs:
  ios:
    ios-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
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
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
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
  ios: []
`, sampleAppsIosWatchkitVersions...)

var sampleAppsCarthageVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CarthageVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CarthageVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
}

var sampleAppsCarthageResultYML = fmt.Sprintf(`options:
  ios:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      sample-apps-carthage.xcodeproj:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          sample-apps-carthage:
            title: ipa export method
            env_key: BITRISE_EXPORT_METHOD
            value_map:
              ad-hoc:
                config: ios-carthage-test-config
              app-store:
                config: ios-carthage-test-config
              development:
                config: ios-carthage-test-config
              enterprise:
                config: ios-carthage-test-config
configs:
  ios:
    ios-carthage-test-config: |
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
          - carthage@%s:
              inputs:
              - carthage_command: bootstrap
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
          - carthage@%s:
              inputs:
              - carthage_command: bootstrap
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  ios: []
`, sampleAppsCarthageVersions...)
