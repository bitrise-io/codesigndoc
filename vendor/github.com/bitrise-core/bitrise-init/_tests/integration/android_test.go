package integration

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/command/git"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func gitClone(t *testing.T, dir, uri string) {
	fmt.Printf("cloning into: %s\n", dir)
	g, err := git.New(dir)
	require.NoError(t, err)
	require.NoError(t, g.Clone(uri).Run())
}

func TestAndroid(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__android__")
	require.NoError(t, err)

	t.Log("sample-apps-android-sdk22")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-android-sdk22")
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-android-sdk22.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "sample-apps-android-sdk22", strings.TrimSpace(sampleAppsAndroid22ResultYML), strings.TrimSpace(result), sampleAppsAndroid22Versions...)
	}

	t.Log("android-non-executable-gradlew")
	{
		sampleAppDir := filepath.Join(tmpDir, "android-non-executable-gradlew")
		sampleAppURL := "https://github.com/bitrise-samples/android-non-executable-gradlew.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "android-non-executable-gradlew", strings.TrimSpace(androidNonExecutableGradlewResultYML), strings.TrimSpace(result), androidNonExecutableGradlewVersions...)
	}

	t.Log("android-sdk22-no-gradlew")
	{
		sampleAppDir := filepath.Join(tmpDir, "android-sdk22-no-gradlew")
		sampleAppURL := "https://github.com/bitrise-samples/android-sdk22-no-gradlew.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		_, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.EqualError(t, err, "exit status 1")

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "android-sdk22-no-gradlew", strings.TrimSpace(sampleAppsSDK22NoGradlewResultYML), strings.TrimSpace(result))
	}

	t.Log("android-sdk22-subdir")
	{
		sampleAppDir := filepath.Join(tmpDir, "android-sdk22-subdir")
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-android-sdk22-subdir"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		_, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "android-sdk22-subdir", strings.TrimSpace(sampleAppsAndroidSDK22SubdirResultYML), strings.TrimSpace(result), sampleAppsAndroidSDK22SubdirVersions...)
	}
}

var sampleAppsAndroidSDK22SubdirVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.ChangeAndroidVersionCodeAndVersionNameVersion,
	steps.AndroidLintVersion,
	steps.AndroidUnitTestVersion,
	steps.AndroidBuildVersion,
	steps.SignAPKVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidLintVersion,
	steps.AndroidUnitTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
}

var sampleAppsAndroidSDK22SubdirResultYML = fmt.Sprintf(`options:
  android:
    title: The root directory of an Android project
    env_key: PROJECT_LOCATION
    value_map:
      src:
        title: Module
        env_key: MODULE
        value_map:
          app:
            title: Variant
            env_key: VARIANT
            value_map:
              "":
                config: android-config
configs:
  android:
    android-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: android
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        deploy:
          description: |
            ## How to get a signed APK

            This workflow contains the **Sign APK** step. To sign your APK all you have to do is to:

            1. Click on **Code Signing** tab
            1. Find the **ANDROID KEYSTORE FILE** section
            1. Click or drop your file on the upload file field
            1. Fill the displayed 3 input fields:
             1. **Keystore password**
             1. **Keystore alias**
             1. **Private key password**
            1. Click on **[Save metadata]** button

            That's it! From now on, **Sign APK** step will receive your uploaded files.

            ## To run this workflow

            If you want to run this workflow manually:

            1. Open the app's build list page
            2. Click on **[Start/Schedule a Build]** button
            3. Select **deploy** in **Workflow** dropdown input
            4. Click **[Start Build]** button

            Or if you need this workflow to be started by a GIT event:

            1. Click on **Triggers** tab
            2. Setup your desired event (push/tag/pull) and select **deploy** workflow
            3. Click on **[Done]** and then **[Save]** buttons

            The next change in your repository that matches any of your trigger map event will start **deploy** workflow.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - change-android-versioncode-and-versionname@%s:
              inputs:
              - build_gradle_path: $PROJECT_LOCATION/$MODULE/build.gradle
          - android-lint@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-unit-test@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - sign-apk@%s:
              run_if: '{{getenv "BITRISEIO_ANDROID_KEYSTORE_URL" | ne ""}}'
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
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-lint@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-unit-test@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  android: []
`, sampleAppsAndroidSDK22SubdirVersions...)

var sampleAppsSDK22NoGradlewResultYML = `warnings:
  android:
  - "<b>No Gradle Wrapper (gradlew) found.</b> \nUsing a Gradle Wrapper (gradlew)
    is required, as the wrapper is what makes sure\nthat the right Gradle version
    is installed and used for the build. More info/guide: <a>https://docs.gradle.org/current/userguide/gradle_wrapper.html</a>"
errors:
  general:
  - No known platform detected
`

var sampleAppsAndroid22Versions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.ChangeAndroidVersionCodeAndVersionNameVersion,
	steps.AndroidLintVersion,
	steps.AndroidUnitTestVersion,
	steps.AndroidBuildVersion,
	steps.SignAPKVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidLintVersion,
	steps.AndroidUnitTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
}

var sampleAppsAndroid22ResultYML = fmt.Sprintf(`options:
  android:
    title: The root directory of an Android project
    env_key: PROJECT_LOCATION
    value_map:
      .:
        title: Module
        env_key: MODULE
        value_map:
          app:
            title: Variant
            env_key: VARIANT
            value_map:
              "":
                config: android-config
configs:
  android:
    android-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: android
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        deploy:
          description: |
            ## How to get a signed APK

            This workflow contains the **Sign APK** step. To sign your APK all you have to do is to:

            1. Click on **Code Signing** tab
            1. Find the **ANDROID KEYSTORE FILE** section
            1. Click or drop your file on the upload file field
            1. Fill the displayed 3 input fields:
             1. **Keystore password**
             1. **Keystore alias**
             1. **Private key password**
            1. Click on **[Save metadata]** button

            That's it! From now on, **Sign APK** step will receive your uploaded files.

            ## To run this workflow

            If you want to run this workflow manually:

            1. Open the app's build list page
            2. Click on **[Start/Schedule a Build]** button
            3. Select **deploy** in **Workflow** dropdown input
            4. Click **[Start Build]** button

            Or if you need this workflow to be started by a GIT event:

            1. Click on **Triggers** tab
            2. Setup your desired event (push/tag/pull) and select **deploy** workflow
            3. Click on **[Done]** and then **[Save]** buttons

            The next change in your repository that matches any of your trigger map event will start **deploy** workflow.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - change-android-versioncode-and-versionname@%s:
              inputs:
              - build_gradle_path: $PROJECT_LOCATION/$MODULE/build.gradle
          - android-lint@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-unit-test@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - sign-apk@%s:
              run_if: '{{getenv "BITRISEIO_ANDROID_KEYSTORE_URL" | ne ""}}'
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
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-lint@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-unit-test@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  android: []
`, sampleAppsAndroid22Versions...)

var androidNonExecutableGradlewVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.ChangeAndroidVersionCodeAndVersionNameVersion,
	steps.AndroidLintVersion,
	steps.AndroidUnitTestVersion,
	steps.AndroidBuildVersion,
	steps.SignAPKVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidLintVersion,
	steps.AndroidUnitTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
}

var androidNonExecutableGradlewResultYML = fmt.Sprintf(`options:
  android:
    title: The root directory of an Android project
    env_key: PROJECT_LOCATION
    value_map:
      .:
        title: Module
        env_key: MODULE
        value_map:
          app:
            title: Variant
            env_key: VARIANT
            value_map:
              "":
                config: android-config
configs:
  android:
    android-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: android
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        deploy:
          description: |
            ## How to get a signed APK

            This workflow contains the **Sign APK** step. To sign your APK all you have to do is to:

            1. Click on **Code Signing** tab
            1. Find the **ANDROID KEYSTORE FILE** section
            1. Click or drop your file on the upload file field
            1. Fill the displayed 3 input fields:
             1. **Keystore password**
             1. **Keystore alias**
             1. **Private key password**
            1. Click on **[Save metadata]** button

            That's it! From now on, **Sign APK** step will receive your uploaded files.

            ## To run this workflow

            If you want to run this workflow manually:

            1. Open the app's build list page
            2. Click on **[Start/Schedule a Build]** button
            3. Select **deploy** in **Workflow** dropdown input
            4. Click **[Start Build]** button

            Or if you need this workflow to be started by a GIT event:

            1. Click on **Triggers** tab
            2. Setup your desired event (push/tag/pull) and select **deploy** workflow
            3. Click on **[Done]** and then **[Save]** buttons

            The next change in your repository that matches any of your trigger map event will start **deploy** workflow.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - change-android-versioncode-and-versionname@%s:
              inputs:
              - build_gradle_path: $PROJECT_LOCATION/$MODULE/build.gradle
          - android-lint@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-unit-test@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - sign-apk@%s:
              run_if: '{{getenv "BITRISEIO_ANDROID_KEYSTORE_URL" | ne ""}}'
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
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-lint@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-unit-test@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  android: []
`, androidNonExecutableGradlewVersions...)
