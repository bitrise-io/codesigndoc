package steps

import (
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pointers"
	stepmanModels "github.com/bitrise-io/stepman/models"
)

func stepIDComposite(ID, version string) string {
	if version != "" {
		return ID + "@" + version
	}
	return ID
}

func stepListItem(stepIDComposite, title, runIf string, inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	step := stepmanModels.StepModel{}
	if title != "" {
		step.Title = pointers.NewStringPtr(title)
	}
	if runIf != "" {
		step.RunIf = pointers.NewStringPtr(runIf)
	}
	if len(inputs) > 0 {
		step.Inputs = inputs
	}

	return bitriseModels.StepListItemModel{
		stepIDComposite: step,
	}
}

// DefaultPrepareStepList ...
func DefaultPrepareStepList(isIncludeCache bool) []bitriseModels.StepListItemModel {
	stepList := []bitriseModels.StepListItemModel{
		ActivateSSHKeyStepListItem(),
		GitCloneStepListItem(),
	}

	if isIncludeCache {
		stepList = append(stepList, CachePullStepListItem())
	}

	return append(stepList, ScriptSteplistItem(ScriptDefaultTitle))
}

// DefaultDeployStepList ...
func DefaultDeployStepList(isIncludeCache bool) []bitriseModels.StepListItemModel {
	stepList := []bitriseModels.StepListItemModel{
		DeployToBitriseIoStepListItem(),
	}

	if isIncludeCache {
		stepList = append(stepList, CachePushStepListItem())
	}

	return stepList
}

// ActivateSSHKeyStepListItem ...
func ActivateSSHKeyStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(ActivateSSHKeyID, ActivateSSHKeyVersion)
	runIf := `{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}`
	return stepListItem(stepIDComposite, "", runIf)
}

// AndroidLintStepListItem ...
func AndroidLintStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(AndroidLintID, AndroidLintVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// AndroidUnitTestStepListItem ...
func AndroidUnitTestStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(AndroidUnitTestID, AndroidUnitTestVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// AndroidBuildStepListItem ...
func AndroidBuildStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(AndroidBuildID, AndroidBuildVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// GitCloneStepListItem ...
func GitCloneStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(GitCloneID, GitCloneVersion)
	return stepListItem(stepIDComposite, "", "")
}

// CachePullStepListItem ...
func CachePullStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CachePullID, CachePullVersion)
	return stepListItem(stepIDComposite, "", "")
}

// CachePushStepListItem ...
func CachePushStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CachePushID, CachePushVersion)
	return stepListItem(stepIDComposite, "", "")
}

// CertificateAndProfileInstallerStepListItem ...
func CertificateAndProfileInstallerStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CertificateAndProfileInstallerID, CertificateAndProfileInstallerVersion)
	return stepListItem(stepIDComposite, "", "")
}

// ChangeAndroidVersionCodeAndVersionNameStepListItem ...
func ChangeAndroidVersionCodeAndVersionNameStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(ChangeAndroidVersionCodeAndVersionNameID, ChangeAndroidVersionCodeAndVersionNameVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// DeployToBitriseIoStepListItem ...
func DeployToBitriseIoStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(DeployToBitriseIoID, DeployToBitriseIoVersion)
	return stepListItem(stepIDComposite, "", "")
}

// ScriptSteplistItem ...
func ScriptSteplistItem(title string, inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(ScriptID, ScriptVersion)
	return stepListItem(stepIDComposite, title, "", inputs...)
}

// SignAPKStepListItem ...
func SignAPKStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(SignAPKID, SignAPKVersion)
	return stepListItem(stepIDComposite, "", `{{getenv "BITRISEIO_ANDROID_KEYSTORE_URL" | ne ""}}`)
}

// InstallMissingAndroidToolsStepListItem ....
func InstallMissingAndroidToolsStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(InstallMissingAndroidToolsID, InstallMissingAndroidToolsVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// FastlaneStepListItem ...
func FastlaneStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(FastlaneID, FastlaneVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// CocoapodsInstallStepListItem ...
func CocoapodsInstallStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CocoapodsInstallID, CocoapodsInstallVersion)
	return stepListItem(stepIDComposite, "", "")
}

// CarthageStepListItem ...
func CarthageStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CarthageID, CarthageVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// RecreateUserSchemesStepListItem ...
func RecreateUserSchemesStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(RecreateUserSchemesID, RecreateUserSchemesVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// XcodeArchiveStepListItem ...
func XcodeArchiveStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XcodeArchiveID, XcodeArchiveVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// XcodeTestStepListItem ...
func XcodeTestStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XcodeTestID, XcodeTestVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// XamarinUserManagementStepListItem ...
func XamarinUserManagementStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XamarinUserManagementID, XamarinUserManagementVersion)
	runIf := ".IsCI"
	return stepListItem(stepIDComposite, "", runIf, inputs...)
}

// NugetRestoreStepListItem ...
func NugetRestoreStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(NugetRestoreID, NugetRestoreVersion)
	return stepListItem(stepIDComposite, "", "")
}

// XamarinComponentsRestoreStepListItem ...
func XamarinComponentsRestoreStepListItem() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XamarinComponentsRestoreID, XamarinComponentsRestoreVersion)
	return stepListItem(stepIDComposite, "", "")
}

// XamarinArchiveStepListItem ...
func XamarinArchiveStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XamarinArchiveID, XamarinArchiveVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// XcodeArchiveMacStepListItem ...
func XcodeArchiveMacStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XcodeArchiveMacID, XcodeArchiveMacVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// XcodeTestMacStepListItem ...
func XcodeTestMacStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(XcodeTestMacID, XcodeTestMacVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// CordovaArchiveStepListItem ...
func CordovaArchiveStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CordovaArchiveID, CordovaArchiveVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// IonicArchiveStepListItem ...
func IonicArchiveStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(IonicArchiveID, IonicArchiveVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// GenerateCordovaBuildConfigStepListItem ...
func GenerateCordovaBuildConfigStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(GenerateCordovaBuildConfigID, GenerateCordovaBuildConfigVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// JasmineTestRunnerStepListItem ...
func JasmineTestRunnerStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(JasmineTestRunnerID, JasmineTestRunnerVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// KarmaJasmineTestRunnerStepListItem ...
func KarmaJasmineTestRunnerStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(KarmaJasmineTestRunnerID, KarmaJasmineTestRunnerVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// NpmStepListItem ...
func NpmStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(NpmID, NpmVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// ExpoDetachStepListItem ...
func ExpoDetachStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(ExpoDetachID, ExpoDetachVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// YarnStepListItem ...
func YarnStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(YarnID, YarnVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// FlutterInstallStepListItem ...
func FlutterInstallStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(FlutterInstallID, FlutterInstallVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// FlutterTestStepListItem ...
func FlutterTestStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(FlutterTestID, FlutterTestVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// FlutterAnalyzeStepListItem ...
func FlutterAnalyzeStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(FlutterAnalyzeID, FlutterAnalyzeVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}

// FlutterBuildStepListItem ...
func FlutterBuildStepListItem(inputs ...envmanModels.EnvironmentItemModel) bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(FlutterBuildID, FlutterBuildVersion)
	return stepListItem(stepIDComposite, "", "", inputs...)
}
