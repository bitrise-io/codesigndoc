package ios

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"path/filepath"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xcode/xcodeproj"
)

const (
	defaultConfigNameFormat = "default-%s-config"
	configNameFormat        = "%s%s-config"
)

const (
	// ProjectPathInputKey ...
	ProjectPathInputKey = "project_path"
	// ProjectPathInputEnvKey ...
	ProjectPathInputEnvKey = "BITRISE_PROJECT_PATH"
	// ProjectPathInputTitle ...
	ProjectPathInputTitle = "Project (or Workspace) path"
)

const (
	// SchemeInputKey ...
	SchemeInputKey = "scheme"
	// SchemeInputEnvKey ...
	SchemeInputEnvKey = "BITRISE_SCHEME"
	// SchemeInputTitle ...
	SchemeInputTitle = "Scheme name"
)

const (
	// ExportMethodInputKey ...
	ExportMethodInputKey = "export_method"
	// ExportMethodInputEnvKey ...
	ExportMethodInputEnvKey = "BITRISE_EXPORT_METHOD"
	// IosExportMethodInputTitle ...
	IosExportMethodInputTitle = "ipa export method"
	// MacExportMethodInputTitle ...
	MacExportMethodInputTitle = "Application export method\nNOTE: `none` means: Export a copy of the application without re-signing."
)

// IosExportMethods ...
var IosExportMethods = []string{"app-store", "ad-hoc", "enterprise", "development"}

// MacExportMethods ...
var MacExportMethods = []string{"app-store", "developer-id", "development", "none"}

const (
	// ConfigurationInputKey ...
	ConfigurationInputKey = "configuration"
)

const (
	// CarthageCommandInputKey ...
	CarthageCommandInputKey = "carthage_command"
)

const cartfileBase = "Cartfile"
const cartfileResolvedBase = "Cartfile.resolved"

// AllowCartfileBaseFilter ...
var AllowCartfileBaseFilter = utility.BaseFilter(cartfileBase, true)

// ConfigDescriptor ...
type ConfigDescriptor struct {
	HasPodfile           bool
	CarthageCommand      string
	HasTest              bool
	MissingSharedSchemes bool
}

// NewConfigDescriptor ...
func NewConfigDescriptor(hasPodfile bool, carthageCommand string, hasXCTest bool, missingSharedSchemes bool) ConfigDescriptor {
	return ConfigDescriptor{
		HasPodfile:           hasPodfile,
		CarthageCommand:      carthageCommand,
		HasTest:              hasXCTest,
		MissingSharedSchemes: missingSharedSchemes,
	}
}

// ConfigName ...
func (descriptor ConfigDescriptor) ConfigName(projectType XcodeProjectType) string {
	qualifiers := ""
	if descriptor.HasPodfile {
		qualifiers += "-pod"
	}
	if descriptor.CarthageCommand != "" {
		qualifiers += "-carthage"
	}
	if descriptor.HasTest {
		qualifiers += "-test"
	}
	if descriptor.MissingSharedSchemes {
		qualifiers += "-missing-shared-schemes"
	}
	return fmt.Sprintf(configNameFormat, string(projectType), qualifiers)
}

// HasCartfileInDirectoryOf ...
func HasCartfileInDirectoryOf(pth string) bool {
	dir := filepath.Dir(pth)
	cartfilePth := filepath.Join(dir, cartfileBase)
	exist, err := pathutil.IsPathExists(cartfilePth)
	if err != nil {
		return false
	}
	return exist
}

// HasCartfileResolvedInDirectoryOf ...
func HasCartfileResolvedInDirectoryOf(pth string) bool {
	dir := filepath.Dir(pth)
	cartfileResolvedPth := filepath.Join(dir, cartfileResolvedBase)
	exist, err := pathutil.IsPathExists(cartfileResolvedPth)
	if err != nil {
		return false
	}
	return exist
}

// Detect ...
func Detect(projectType XcodeProjectType, searchDir string) (bool, error) {
	fileList, err := utility.ListPathInDirSortedByComponents(searchDir, true)
	if err != nil {
		return false, err
	}

	log.TInfof("Filter relevant Xcode project files")

	relevantXcodeprojectFiles, err := FilterRelevantProjectFiles(fileList, projectType)
	if err != nil {
		return false, err
	}

	log.TPrintf("%d Xcode %s project files found", len(relevantXcodeprojectFiles), string(projectType))
	for _, xcodeprojectFile := range relevantXcodeprojectFiles {
		log.TPrintf("- %s", xcodeprojectFile)
	}

	if len(relevantXcodeprojectFiles) == 0 {
		log.TPrintf("platform not detected")
		return false, nil
	}

	log.TSuccessf("Platform detected")

	return true, nil
}

func printMissingSharedSchemesAndGenerateWarning(projectPth, defaultGitignorePth string, targets []xcodeproj.TargetModel) string {
	isXcshareddataGitignored := false
	if exist, err := pathutil.IsPathExists(defaultGitignorePth); err != nil {
		log.TWarnf("Failed to check if .gitignore file exists at: %s, error: %s", defaultGitignorePth, err)
	} else if exist {
		isGitignored, err := utility.FileContains(defaultGitignorePth, "xcshareddata")
		if err != nil {
			log.TWarnf("Failed to check if xcshareddata gitignored, error: %s", err)
		} else {
			isXcshareddataGitignored = isGitignored
		}
	}

	log.TPrintf("")
	log.TErrorf("No shared schemes found, adding recreate-user-schemes step...")
	log.TErrorf("The newly generated schemes may differ from the ones in your project.")

	message := `No shared schemes found for project: ` + projectPth + `.` + "\n"

	if isXcshareddataGitignored {
		log.TErrorf("Your gitignore file (%s) contains 'xcshareddata', maybe shared schemes are gitignored?", defaultGitignorePth)
		log.TErrorf("If not, make sure to share your schemes, to have the expected behaviour.")

		message += `Your gitignore file (` + defaultGitignorePth + `) contains 'xcshareddata', maybe shared schemes are gitignored?` + "\n"
	} else {
		log.TErrorf("Make sure to share your schemes, to have the expected behaviour.")
	}

	message += `Automatically generated schemes may differ from the ones in your project.
Make sure to <a href="http://devcenter.bitrise.io/ios/frequent-ios-issues/#xcode-scheme-not-found">share your schemes</a> for the expected behaviour.`

	log.TPrintf("")

	log.TWarnf("%d user schemes will be generated", len(targets))
	for _, target := range targets {
		log.TWarnf("- %s", target.Name)
	}

	log.TPrintf("")

	return message
}

func detectCarthageCommand(projectPth string) (string, string) {
	carthageCommand := ""
	warning := ""

	if HasCartfileInDirectoryOf(projectPth) {
		if HasCartfileResolvedInDirectoryOf(projectPth) {
			carthageCommand = "bootstrap"
		} else {
			dir := filepath.Dir(projectPth)
			cartfilePth := filepath.Join(dir, "Cartfile")

			warning = fmt.Sprintf(`Cartfile found at (%s), but no Cartfile.resolved exists in the same directory.
It is <a href="https://github.com/Carthage/Carthage/blob/master/Documentation/Artifacts.md#cartfileresolved">strongly recommended to commit this file to your repository</a>`, cartfilePth)

			carthageCommand = "update"
		}
	}

	return carthageCommand, warning
}

// GenerateOptions ...
func GenerateOptions(projectType XcodeProjectType, searchDir string) (models.OptionNode, []ConfigDescriptor, models.Warnings, error) {
	warnings := models.Warnings{}

	fileList, err := utility.ListPathInDirSortedByComponents(searchDir, true)
	if err != nil {
		return models.OptionNode{}, []ConfigDescriptor{}, models.Warnings{}, err
	}

	// Separate workspaces and standalon projects
	projectFiles, err := FilterRelevantProjectFiles(fileList, projectType)
	if err != nil {
		return models.OptionNode{}, []ConfigDescriptor{}, models.Warnings{}, err
	}

	workspaceFiles, err := FilterRelevantWorkspaceFiles(fileList, projectType)
	if err != nil {
		return models.OptionNode{}, []ConfigDescriptor{}, models.Warnings{}, err
	}

	standaloneProjects, workspaces, err := CreateStandaloneProjectsAndWorkspaces(projectFiles, workspaceFiles)
	if err != nil {
		return models.OptionNode{}, []ConfigDescriptor{}, models.Warnings{}, err
	}

	exportMethodInputTitle := ""
	exportMethods := []string{}
	if projectType == XcodeProjectTypeIOS {
		exportMethodInputTitle = IosExportMethodInputTitle
		exportMethods = IosExportMethods
	} else {
		exportMethodInputTitle = MacExportMethodInputTitle
		exportMethods = MacExportMethods
	}

	// Create cocoapods workspace-project mapping
	log.TInfof("Searching for Podfile")

	podfiles, err := FilterRelevantPodfiles(fileList)
	if err != nil {
		return models.OptionNode{}, []ConfigDescriptor{}, models.Warnings{}, err
	}

	log.TPrintf("%d Podfiles detected", len(podfiles))

	for _, podfile := range podfiles {
		log.TPrintf("- %s", podfile)

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfile, projectFiles)
		if err != nil {
			warning := fmt.Sprintf("Failed to determine cocoapods project-workspace mapping, error: %s", err)
			warnings = append(warnings, warning)
			log.Warnf(warning)
			continue
		}

		aStandaloneProjects, aWorkspaces, err := MergePodWorkspaceProjectMap(workspaceProjectMap, standaloneProjects, workspaces)
		if err != nil {
			warning := fmt.Sprintf("Failed to create cocoapods project-workspace mapping, error: %s", err)
			warnings = append(warnings, warning)
			log.Warnf(warning)
			continue
		}

		standaloneProjects = aStandaloneProjects
		workspaces = aWorkspaces
	}

	// Carthage
	log.TInfof("Searching for Cartfile")

	cartfiles, err := FilterRelevantCartFile(fileList)
	if err != nil {
		return models.OptionNode{}, []ConfigDescriptor{}, models.Warnings{}, err
	}

	log.TPrintf("%d Cartfiles detected", len(cartfiles))
	for _, file := range cartfiles {
		log.TPrintf("- %s", file)
	}

	// Create config descriptors & options
	configDescriptors := []ConfigDescriptor{}

	defaultGitignorePth := filepath.Join(searchDir, ".gitignore")

	projectPathOption := models.NewOption(ProjectPathInputTitle, ProjectPathInputEnvKey)

	// Standalon Projects
	for _, project := range standaloneProjects {
		log.TInfof("Inspecting standalone project file: %s", project.Pth)

		schemeOption := models.NewOption(SchemeInputTitle, SchemeInputEnvKey)
		projectPathOption.AddOption(project.Pth, schemeOption)

		carthageCommand, warning := detectCarthageCommand(project.Pth)
		if warning != "" {
			warnings = append(warnings, warning)
		}

		log.TPrintf("%d shared schemes detected", len(project.SharedSchemes))

		if len(project.SharedSchemes) == 0 {
			message := printMissingSharedSchemesAndGenerateWarning(project.Pth, defaultGitignorePth, project.Targets)
			if message != "" {
				warnings = append(warnings, message)
			}

			for _, target := range project.Targets {

				exportMethodOption := models.NewOption(exportMethodInputTitle, ExportMethodInputEnvKey)
				schemeOption.AddOption(target.Name, exportMethodOption)

				for _, exportMethod := range exportMethods {
					configDescriptor := NewConfigDescriptor(false, carthageCommand, target.HasXCTest, true)
					configDescriptors = append(configDescriptors, configDescriptor)

					configOption := models.NewConfigOption(configDescriptor.ConfigName(projectType))
					exportMethodOption.AddConfig(exportMethod, configOption)
				}
			}
		} else {
			for _, scheme := range project.SharedSchemes {
				log.TPrintf("- %s", scheme.Name)

				exportMethodOption := models.NewOption(exportMethodInputTitle, ExportMethodInputEnvKey)
				schemeOption.AddOption(scheme.Name, exportMethodOption)

				for _, exportMethod := range exportMethods {
					configDescriptor := NewConfigDescriptor(false, carthageCommand, scheme.HasXCTest, false)
					configDescriptors = append(configDescriptors, configDescriptor)

					configOption := models.NewConfigOption(configDescriptor.ConfigName(projectType))
					exportMethodOption.AddConfig(exportMethod, configOption)

				}
			}
		}
	}

	// Workspaces
	for _, workspace := range workspaces {
		log.TInfof("Inspecting workspace file: %s", workspace.Pth)

		schemeOption := models.NewOption(SchemeInputTitle, SchemeInputEnvKey)
		projectPathOption.AddOption(workspace.Pth, schemeOption)

		carthageCommand, warning := detectCarthageCommand(workspace.Pth)
		if warning != "" {
			warnings = append(warnings, warning)
		}

		sharedSchemes := workspace.GetSharedSchemes()
		log.TPrintf("%d shared schemes detected", len(sharedSchemes))

		if len(sharedSchemes) == 0 {
			targets := workspace.GetTargets()

			message := printMissingSharedSchemesAndGenerateWarning(workspace.Pth, defaultGitignorePth, targets)
			if message != "" {
				warnings = append(warnings, message)
			}

			for _, target := range targets {
				exportMethodOption := models.NewOption(exportMethodInputTitle, ExportMethodInputEnvKey)
				schemeOption.AddOption(target.Name, exportMethodOption)

				for _, exportMethod := range exportMethods {
					configDescriptor := NewConfigDescriptor(workspace.IsPodWorkspace, carthageCommand, target.HasXCTest, true)
					configDescriptors = append(configDescriptors, configDescriptor)

					configOption := models.NewConfigOption(configDescriptor.ConfigName(projectType))
					exportMethodOption.AddConfig(exportMethod, configOption)
				}
			}
		} else {
			for _, scheme := range sharedSchemes {
				log.TPrintf("- %s", scheme.Name)

				exportMethodOption := models.NewOption(exportMethodInputTitle, ExportMethodInputEnvKey)
				schemeOption.AddOption(scheme.Name, exportMethodOption)

				for _, exportMethod := range exportMethods {
					configDescriptor := NewConfigDescriptor(workspace.IsPodWorkspace, carthageCommand, scheme.HasXCTest, false)
					configDescriptors = append(configDescriptors, configDescriptor)

					configOption := models.NewConfigOption(configDescriptor.ConfigName(projectType))
					exportMethodOption.AddConfig(exportMethod, configOption)
				}
			}
		}
	}

	configDescriptors = RemoveDuplicatedConfigDescriptors(configDescriptors, projectType)

	if len(configDescriptors) == 0 {
		log.TErrorf("No valid %s config found", string(projectType))
		return models.OptionNode{}, []ConfigDescriptor{}, warnings, fmt.Errorf("No valid %s config found", string(projectType))
	}

	return *projectPathOption, configDescriptors, warnings, nil
}

// GenerateDefaultOptions ...
func GenerateDefaultOptions(projectType XcodeProjectType) models.OptionNode {
	projectPathOption := models.NewOption(ProjectPathInputTitle, ProjectPathInputEnvKey)

	schemeOption := models.NewOption(SchemeInputTitle, SchemeInputEnvKey)
	projectPathOption.AddOption("_", schemeOption)

	exportMethodInputTitle := ""
	exportMethods := []string{}
	if projectType == XcodeProjectTypeIOS {
		exportMethodInputTitle = IosExportMethodInputTitle
		exportMethods = IosExportMethods
	} else {
		exportMethodInputTitle = MacExportMethodInputTitle
		exportMethods = MacExportMethods
	}

	exportMethodOption := models.NewOption(exportMethodInputTitle, ExportMethodInputEnvKey)
	schemeOption.AddOption("_", exportMethodOption)

	for _, exportMethod := range exportMethods {
		configOption := models.NewConfigOption(fmt.Sprintf(defaultConfigNameFormat, string(projectType)))
		exportMethodOption.AddConfig(exportMethod, configOption)
	}

	return *projectPathOption
}

// GenerateConfigBuilder ...
func GenerateConfigBuilder(projectType XcodeProjectType, hasPodfile, hasTest, missingSharedSchemes bool, carthageCommand string, isIncludeCache bool) models.ConfigBuilderModel {
	configBuilder := models.NewDefaultConfigBuilder()

	// CI
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

	if missingSharedSchemes {
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.RecreateUserSchemesStepListItem(
			envmanModels.EnvironmentItemModel{ProjectPathInputKey: "$" + ProjectPathInputEnvKey},
		))
	}

	if hasPodfile {
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CocoapodsInstallStepListItem())
	}

	if carthageCommand != "" {
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CarthageStepListItem(
			envmanModels.EnvironmentItemModel{CarthageCommandInputKey: carthageCommand},
		))
	}

	xcodeStepInputModels := []envmanModels.EnvironmentItemModel{
		envmanModels.EnvironmentItemModel{ProjectPathInputKey: "$" + ProjectPathInputEnvKey},
		envmanModels.EnvironmentItemModel{SchemeInputKey: "$" + SchemeInputEnvKey},
	}
	xcodeArchiveStepInputModels := append(xcodeStepInputModels, envmanModels.EnvironmentItemModel{ExportMethodInputKey: "$" + ExportMethodInputEnvKey})

	if hasTest {
		switch projectType {
		case XcodeProjectTypeIOS:
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XcodeTestStepListItem(xcodeStepInputModels...))
		case XcodeProjectTypeMacOS:
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XcodeTestMacStepListItem(xcodeStepInputModels...))
		}
	} else {
		switch projectType {
		case XcodeProjectTypeIOS:
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XcodeArchiveStepListItem(xcodeArchiveStepInputModels...))
		case XcodeProjectTypeMacOS:
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XcodeArchiveMacStepListItem(xcodeArchiveStepInputModels...))
		}
	}

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(isIncludeCache)...)

	if hasTest {
		// CD
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

		if missingSharedSchemes {
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.RecreateUserSchemesStepListItem(
				envmanModels.EnvironmentItemModel{ProjectPathInputKey: "$" + ProjectPathInputEnvKey},
			))
		}

		if hasPodfile {
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CocoapodsInstallStepListItem())
		}

		if carthageCommand != "" {
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CarthageStepListItem(
				envmanModels.EnvironmentItemModel{CarthageCommandInputKey: carthageCommand},
			))
		}

		switch projectType {
		case XcodeProjectTypeIOS:
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeTestStepListItem(xcodeStepInputModels...))
		case XcodeProjectTypeMacOS:
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeTestMacStepListItem(xcodeStepInputModels...))
		}

		switch projectType {
		case XcodeProjectTypeIOS:
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(xcodeArchiveStepInputModels...))
		case XcodeProjectTypeMacOS:
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveMacStepListItem(xcodeArchiveStepInputModels...))
		}

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(isIncludeCache)...)
	}

	return *configBuilder
}

// RemoveDuplicatedConfigDescriptors ...
func RemoveDuplicatedConfigDescriptors(configDescriptors []ConfigDescriptor, projectType XcodeProjectType) []ConfigDescriptor {
	descritorNameMap := map[string]ConfigDescriptor{}
	for _, descriptor := range configDescriptors {
		name := descriptor.ConfigName(projectType)
		descritorNameMap[name] = descriptor
	}

	descriptors := []ConfigDescriptor{}
	for _, descriptor := range descritorNameMap {
		descriptors = append(descriptors, descriptor)
	}

	return descriptors
}

// GenerateConfig ...
func GenerateConfig(projectType XcodeProjectType, configDescriptors []ConfigDescriptor, isIncludeCache bool) (models.BitriseConfigMap, error) {
	bitriseDataMap := models.BitriseConfigMap{}
	for _, descriptor := range configDescriptors {
		configBuilder := GenerateConfigBuilder(projectType, descriptor.HasPodfile, descriptor.HasTest, descriptor.MissingSharedSchemes, descriptor.CarthageCommand, isIncludeCache)

		config, err := configBuilder.Generate(string(projectType))
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		bitriseDataMap[descriptor.ConfigName(projectType)] = string(data)
	}

	return bitriseDataMap, nil
}

// GenerateDefaultConfig ...
func GenerateDefaultConfig(projectType XcodeProjectType, isIncludeCache bool) (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)

	// CI
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CertificateAndProfileInstallerStepListItem())
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.RecreateUserSchemesStepListItem(
		envmanModels.EnvironmentItemModel{ProjectPathInputKey: "$" + ProjectPathInputEnvKey},
	))

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CocoapodsInstallStepListItem())

	xcodeTestStepInputModels := []envmanModels.EnvironmentItemModel{
		envmanModels.EnvironmentItemModel{ProjectPathInputKey: "$" + ProjectPathInputEnvKey},
		envmanModels.EnvironmentItemModel{SchemeInputKey: "$" + SchemeInputEnvKey},
	}
	xcodeArchiveStepInputModels := append(xcodeTestStepInputModels, envmanModels.EnvironmentItemModel{ExportMethodInputKey: "$" + ExportMethodInputEnvKey})

	switch projectType {
	case XcodeProjectTypeIOS:
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XcodeTestStepListItem(xcodeTestStepInputModels...))
	case XcodeProjectTypeMacOS:
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XcodeTestMacStepListItem(xcodeTestStepInputModels...))
	}
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(true)...)

	// CD
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.RecreateUserSchemesStepListItem(
		envmanModels.EnvironmentItemModel{ProjectPathInputKey: "$" + ProjectPathInputEnvKey},
	))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CocoapodsInstallStepListItem())

	switch projectType {
	case XcodeProjectTypeIOS:
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeTestStepListItem(xcodeTestStepInputModels...))

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(xcodeArchiveStepInputModels...))
	case XcodeProjectTypeMacOS:
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeTestMacStepListItem(xcodeTestStepInputModels...))
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveMacStepListItem(xcodeArchiveStepInputModels...))
	}

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(true)...)

	config, err := configBuilder.Generate(string(projectType))
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	return models.BitriseConfigMap{
		fmt.Sprintf(defaultConfigNameFormat, string(projectType)): string(data),
	}, nil
}
