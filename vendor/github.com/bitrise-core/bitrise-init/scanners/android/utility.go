package android

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pathutil"
)

// Constants ...
const (
	ScannerName       = "android"
	ConfigName        = "android-config"
	DefaultConfigName = "default-android-config"

	ProjectLocationInputKey    = "project_location"
	ProjectLocationInputEnvKey = "PROJECT_LOCATION"
	ProjectLocationInputTitle  = "The root directory of an Android project"

	ModuleBuildGradlePathInputKey = "build_gradle_path"

	VariantInputKey    = "variant"
	VariantInputEnvKey = "VARIANT"
	VariantInputTitle  = "Variant"

	ModuleInputKey    = "module"
	ModuleInputEnvKey = "MODULE"
	ModuleInputTitle  = "Module"

	GradlewPathInputKey    = "gradlew_path"
	GradlewPathInputEnvKey = "GRADLEW_PATH"
	GradlewPathInputTitle  = "Gradlew file path"
)

func walk(src string, fn func(path string, info os.FileInfo) error) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == src {
			return nil
		}
		return fn(path, info)
	})
}

func checkFiles(path string, files ...string) (bool, error) {
	for _, file := range files {
		exists, err := pathutil.IsPathExists(filepath.Join(path, file))
		if err != nil {
			return false, err
		}
		if !exists {
			return false, nil
		}
	}
	return true, nil
}

func walkMultipleFiles(searchDir string, files ...string) (matches []string, err error) {
	match, err := checkFiles(searchDir, files...)
	if err != nil {
		return nil, err
	}
	if match {
		matches = append(matches, searchDir)
	}
	return matches, walk(searchDir, func(path string, info os.FileInfo) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			match, err := checkFiles(path, files...)
			if err != nil {
				return err
			}
			if match {
				matches = append(matches, path)
			}
		}
		return nil
	})
}

func checkGradlew(projectDir string) error {
	gradlewPth := filepath.Join(projectDir, "gradlew")
	exist, err := pathutil.IsPathExists(gradlewPth)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(`<b>No Gradle Wrapper (gradlew) found.</b> 
Using a Gradle Wrapper (gradlew) is required, as the wrapper is what makes sure
that the right Gradle version is installed and used for the build. More info/guide: <a>https://docs.gradle.org/current/userguide/gradle_wrapper.html</a>`)
	}
	return nil
}

func (scanner *Scanner) generateConfigBuilder() models.ConfigBuilderModel {
	configBuilder := models.NewDefaultConfigBuilder()

	projectLocationEnv, gradlewPath, moduleEnv, variantEnv := "$"+ProjectLocationInputEnvKey, "$"+ProjectLocationInputEnvKey+"/gradlew", "$"+ModuleInputEnvKey, "$"+VariantInputEnvKey

	//-- primary
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(true)...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
		envmanModels.EnvironmentItemModel{GradlewPathInputKey: gradlewPath},
	))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.AndroidLintStepListItem(
		envmanModels.EnvironmentItemModel{
			ProjectLocationInputKey: projectLocationEnv,
		},
		envmanModels.EnvironmentItemModel{
			ModuleInputKey: moduleEnv,
		},
		envmanModels.EnvironmentItemModel{
			VariantInputKey: variantEnv,
		},
	))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.AndroidUnitTestStepListItem(
		envmanModels.EnvironmentItemModel{
			ProjectLocationInputKey: projectLocationEnv,
		},
		envmanModels.EnvironmentItemModel{
			ModuleInputKey: moduleEnv,
		},
		envmanModels.EnvironmentItemModel{
			VariantInputKey: variantEnv,
		},
	))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(true)...)

	//-- deploy
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(true)...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
		envmanModels.EnvironmentItemModel{GradlewPathInputKey: gradlewPath},
	))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.ChangeAndroidVersionCodeAndVersionNameStepListItem(
		envmanModels.EnvironmentItemModel{ModuleBuildGradlePathInputKey: filepath.Join(projectLocationEnv, moduleEnv, "build.gradle")},
	))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidLintStepListItem(
		envmanModels.EnvironmentItemModel{
			ProjectLocationInputKey: projectLocationEnv,
		},
		envmanModels.EnvironmentItemModel{
			ModuleInputKey: moduleEnv,
		},
		envmanModels.EnvironmentItemModel{
			VariantInputKey: variantEnv,
		},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidUnitTestStepListItem(
		envmanModels.EnvironmentItemModel{
			ProjectLocationInputKey: projectLocationEnv,
		},
		envmanModels.EnvironmentItemModel{
			ModuleInputKey: moduleEnv,
		},
		envmanModels.EnvironmentItemModel{
			VariantInputKey: variantEnv,
		},
	))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
		envmanModels.EnvironmentItemModel{
			ProjectLocationInputKey: projectLocationEnv,
		},
		envmanModels.EnvironmentItemModel{
			ModuleInputKey: moduleEnv,
		},
		envmanModels.EnvironmentItemModel{
			VariantInputKey: variantEnv,
		},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.SignAPKStepListItem())
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(true)...)

	configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployWorkflowDescription)

	return *configBuilder
}
