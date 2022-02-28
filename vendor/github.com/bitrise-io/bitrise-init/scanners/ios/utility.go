package ios

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/analytics"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
)

const (
	defaultConfigNameFormat = "default-%s-config"
	configNameFormat        = "%s%s-config"
	iconFailureTag          = "icon_lookup"
)

const (
	// ProjectPathInputKey ...
	ProjectPathInputKey = "project_path"
	// ProjectPathInputEnvKey ...
	ProjectPathInputEnvKey = "BITRISE_PROJECT_PATH"
	// ProjectPathInputTitle ...
	ProjectPathInputTitle = "Project or Workspace path"
	// ProjectPathInputSummary ...
	ProjectPathInputSummary = "The location of your Xcode project or Xcode workspace files, stored as an Environment Variable. In your Workflows, you can specify paths relative to this path."
)

const (
	// SchemeInputKey ...
	SchemeInputKey = "scheme"
	// SchemeInputEnvKey ...
	SchemeInputEnvKey = "BITRISE_SCHEME"
	// SchemeInputTitle ...
	SchemeInputTitle = "Scheme name"
	// SchemeInputSummary ...
	SchemeInputSummary = "An Xcode scheme defines a collection of targets to build, a configuration to use when building, and a collection of tests to execute. Only shared schemes are detected automatically but you can use any scheme as a target on Bitrise. You can change the scheme at any time in your Env Vars."
)

const (
	// DistributionMethodInputKey ...
	DistributionMethodInputKey = "distribution_method"
	// DistributionMethodEnvKey ...
	DistributionMethodEnvKey = "BITRISE_DISTRIBUTION_METHOD"
	// DistributionMethodInputTitle ...
	DistributionMethodInputTitle = "Distribution method"
	// DistributionMethodInputSummary ...
	DistributionMethodInputSummary = "The export method used to create an .ipa file in your builds, stored as an Environment Variable. You can change this at any time, or even create several .ipa files with different export methods in the same build."
)

const (
	// ExportMethodInputKey ...
	ExportMethodInputKey = "export_method"
	// ExportMethodEnvKey ...
	ExportMethodEnvKey = "BITRISE_EXPORT_METHOD"
	// ExportMethodInputTitle ...
	ExportMethodInputTitle = "Application export method\nNOTE: `none` means: Export a copy of the application without re-signing."
	// ExportMethodInputSummary ...
	ExportMethodInputSummary = "The export method used to create an .app file in your builds, stored as an Environment Variable. You can change this at any time, or even create several .app files with different export methods in the same build."
)

// XCConfigContentInputKey ...
const XCConfigContentInputKey = "xcconfig_content"

// IosExportMethods ...
var IosExportMethods = []string{"app-store", "ad-hoc", "enterprise", "development"}

const (
	// ExportXCArchiveProductInputKey ...
	ExportXCArchiveProductInputKey = "product"

	// ExportXCArchiveProductInputAppClipValue ...
	ExportXCArchiveProductInputAppClipValue = "app-clip"
)

// MacExportMethods ...
var MacExportMethods = []string{"app-store", "developer-id", "development", "none"}

const (
	// ConfigurationInputKey ...
	ConfigurationInputKey = "configuration"
)

const (
	// AutomaticCodeSigningInputKey ...
	AutomaticCodeSigningInputKey = "automatic_code_signing"
	// AutomaticCodeSigningInputAPIKeyValue ...
	AutomaticCodeSigningInputAPIKeyValue = "api-key"
)

const (
	// CarthageCommandInputKey ...
	CarthageCommandInputKey = "carthage_command"
)

const cartfileBase = "Cartfile"
const cartfileResolvedBase = "Cartfile.resolved"

// AllowCartfileBaseFilter ...
var AllowCartfileBaseFilter = pathutil.BaseFilter(cartfileBase, true)

// Scheme is an Xcode project scheme or target
type Scheme struct {
	Name       string
	Missing    bool
	HasXCTests bool
	HasAppClip bool

	Icons models.Icons
}

// Project is an Xcode project on the filesystem
type Project struct {
	RelPath string
	// Is it a standalone project or a workspace?
	IsWorkspace    bool
	IsPodWorkspace bool

	// Carthage command to run: bootstrap/update
	CarthageCommand string
	Warnings        models.Warnings

	Schemes []Scheme
}

// DetectResult ...
type DetectResult struct {
	Projects []Project
	Warnings models.Warnings
}

type containers struct {
	standaloneProjects []container
	workspaces         []container
	podWorkspacePaths  []string
}

// ConfigDescriptor ...
type ConfigDescriptor struct {
	HasPodfile           bool
	CarthageCommand      string
	HasTest              bool
	HasAppClip           bool
	ExportMethod         string
	MissingSharedSchemes bool
}

// NewConfigDescriptor ...
func NewConfigDescriptor(hasPodfile bool, carthageCommand string, hasXCTest, hasAppClip bool, exportMethod string, missingSharedSchemes bool) ConfigDescriptor {
	return ConfigDescriptor{
		HasPodfile:           hasPodfile,
		CarthageCommand:      carthageCommand,
		HasTest:              hasXCTest,
		HasAppClip:           hasAppClip,
		ExportMethod:         exportMethod,
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
	if descriptor.HasAppClip {
		qualifiers += fmt.Sprintf("-app-clip-%s", descriptor.ExportMethod)
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

func fileContains(pth, str string) (bool, error) {
	content, err := fileutil.ReadStringFromFile(pth)
	if err != nil {
		return false, err
	}

	return strings.Contains(content, str), nil
}

func printMissingSharedSchemesAndGenerateWarning(projectRelPth, defaultGitignorePth string) string {
	isXcshareddataGitignored := false
	if exist, err := pathutil.IsPathExists(defaultGitignorePth); err != nil {
		log.TWarnf("Failed to check if .gitignore file exists at: %s, error: %s", defaultGitignorePth, err)
	} else if exist {
		isGitignored, err := fileContains(defaultGitignorePth, "xcshareddata")
		if err != nil {
			log.TWarnf("Failed to check if xcshareddata gitignored, error: %s", err)
		} else {
			isXcshareddataGitignored = isGitignored
		}
	}

	log.TPrintf("")
	log.TErrorf("No shared schemes found, adding recreate-user-schemes Step...")
	log.TErrorf("The newly generated schemes may differ from the ones in your project.")

	message := `No shared schemes found for project: ` + projectRelPth + `.` + "\n"

	if isXcshareddataGitignored {
		log.TErrorf("Your gitignore file (%s) contains 'xcshareddata', maybe shared schemes are gitignored?", defaultGitignorePth)
		log.TErrorf("If not, make sure to share your schemes, to have the expected behaviour.")

		message += `Your gitignore file (` + defaultGitignorePth + `) contains 'xcshareddata', maybe shared schemes are gitignored?` + "\n"
	} else {
		log.TErrorf("Make sure to share your schemes, to have the expected behaviour.")
	}

	message += `Automatically generated schemes may differ from the ones in your project.
Make sure to <a href="https://support.bitrise.io/hc/en-us/articles/4405779956625">share your schemes</a> for the expected behaviour.`

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

func relPathForLog(searchDir string, path string) string {
	relPath, err := filepath.Rel(searchDir, path)
	if err != nil {
		log.TWarnf("failed to get relative path: %s", err)
		return ""
	}

	return relPath
}

// ParseProjects collects available iOS/macOS projects
func ParseProjects(projectType XcodeProjectType, searchDir string, excludeAppIcon, suppressPodFileParseError bool) (DetectResult, error) {
	var (
		projects []Project
		warnings models.Warnings
	)

	fileList, err := pathutil.ListPathInDirSortedByComponents(searchDir, false)
	if err != nil {
		return DetectResult{}, err
	}

	// Separate workspaces and standalone projects
	log.TInfof("Filtering relevant Xcode project files")
	projectFiles, err := FilterRelevantProjectFiles(fileList, projectType)
	if err != nil {
		return DetectResult{}, err
	}

	log.TPrintf("%d Xcode %s project files found", len(projectFiles), string(projectType))
	for _, xcodeprojectFile := range projectFiles {
		log.TPrintf("- %s", relPathForLog(searchDir, xcodeprojectFile))
	}

	if len(projectFiles) == 0 {
		log.TPrintf("Platform not detected")
		return DetectResult{}, nil
	}

	log.TSuccessf("Platform detected")

	workspaceFiles, err := FilterRelevantWorkspaceFiles(fileList, projectType)
	if err != nil {
		return DetectResult{}, err
	}

	detectedContainers, err := createStandaloneProjectsAndWorkspaces(projectFiles, workspaceFiles)
	if err != nil {
		return DetectResult{}, err
	}

	// Create cocoapods workspace-project mapping
	log.TInfof("Searching for Podfile")

	podfiles, err := FilterRelevantPodfiles(fileList)
	if err != nil {
		return DetectResult{}, err
	}

	log.TPrintf("%d Podfiles detected", len(podfiles))

	for _, podfile := range podfiles {
		log.TPrintf("- %s", relPathForLog(searchDir, podfile))

		podfileParser := podfileParser{
			podfilePth:                podfile,
			suppressPodFileParseError: suppressPodFileParseError,
		}

		workspaceProjectMap, err := podfileParser.GetWorkspaceProjectMap(projectFiles)
		if err != nil {
			warning := fmt.Sprintf("Failed to determine cocoapods project-workspace mapping, error: %s", err)
			warnings = append(warnings, warning)
			log.Warnf(warning)

			continue
		}

		detectedContainers, err = mergePodWorkspaceProjectMap(workspaceProjectMap, detectedContainers)
		if err != nil {
			warning := fmt.Sprintf("Failed to create cocoapods project-workspace mapping, error: %s", err)
			warnings = append(warnings, warning)
			log.Warnf(warning)

			continue
		}
	}

	// Carthage
	log.TInfof("Searching for Cartfile")

	cartfiles, err := FilterRelevantCartFile(fileList)
	if err != nil {
		return DetectResult{
			Warnings: warnings,
		}, err
	}

	log.TPrintf("%d Cartfiles detected", len(cartfiles))
	for _, file := range cartfiles {
		log.TPrintf("- %s", relPathForLog(searchDir, file))
	}

	defaultGitignorePth := filepath.Join(searchDir, ".gitignore")

	for _, container := range append(detectedContainers.standaloneProjects, detectedContainers.workspaces...) {
		var (
			projectWarnings []string
			detectedSchemes []Scheme
		)

		containerPath := container.path()
		containerRelPath, err := filepath.Rel(searchDir, containerPath)
		if err != nil {
			return DetectResult{Warnings: warnings}, fmt.Errorf("failed to get relative path: %s", err)
		}

		log.TInfof("Inspecting file: %s", containerRelPath)
		carthageCommand, warning := detectCarthageCommand(containerPath)
		if warning != "" {
			projectWarnings = append(projectWarnings, warning)
		}

		projectToSchemes, err := container.schemes()
		if err != nil {
			return DetectResult{}, fmt.Errorf("failed to read Schemes: %s", err)
		}

		numSharedShemes := numberOfSharedSchemes(projectToSchemes)
		log.TPrintf("%d shared schemes detected", numSharedShemes)
		shouldRecreateSchemes := numSharedShemes == 0

		containerProjects, missingProjects, err := container.projects()
		if err != nil {
			return DetectResult{}, fmt.Errorf("%s", err)
		}

		for _, missingProject := range missingProjects {
			log.Warnf("Skipping Project (%s), as it is not present", relPathForLog(searchDir, missingProject))
		}

		if shouldRecreateSchemes {
			message := printMissingSharedSchemesAndGenerateWarning(containerRelPath, defaultGitignorePth)
			if message != "" {
				projectWarnings = append(projectWarnings, message)
			}
		}

		for _, project := range containerProjects {
			var sharedSchemes []xcscheme.Scheme
			for _, s := range projectToSchemes[project.Path] {
				if s.IsShared {
					sharedSchemes = append(sharedSchemes, s)
				}
			}

			if shouldRecreateSchemes {
				sharedSchemes = project.ReCreateSchemes()
			}

			for _, scheme := range sharedSchemes {
				log.TPrintf("- %s", scheme.Name)

				var icons models.Icons
				if !excludeAppIcon {
					if icons, err = lookupIconByScheme(project, scheme, searchDir); err != nil {
						log.Warnf("could not get icons for app: %s, error: %s", containerRelPath, err)
						analytics.LogInfo(iconFailureTag, analytics.DetectorErrorData(string(XcodeProjectTypeIOS), err), "Failed to lookup ios icons")
					}
				}

				detectedSchemes = append(detectedSchemes, Scheme{
					Name:       scheme.Name,
					Missing:    shouldRecreateSchemes,
					HasXCTests: scheme.IsTestable(),
					HasAppClip: schemeHasAppClipTarget(project, scheme),
					Icons:      icons,
				})
			}
		}

		projects = append(projects, Project{
			RelPath:         containerRelPath,
			IsWorkspace:     container.isWorkspace(),
			IsPodWorkspace:  sliceutil.IsStringInSlice(containerPath, detectedContainers.podWorkspacePaths),
			Schemes:         detectedSchemes,
			CarthageCommand: carthageCommand,
			Warnings:        projectWarnings,
		})
	}

	return DetectResult{
		Projects: projects,
		Warnings: warnings,
	}, nil
}

// GenerateOptions ...
func GenerateOptions(projectType XcodeProjectType, result DetectResult) (models.OptionNode, []ConfigDescriptor, models.Icons, models.Warnings, error) {
	var (
		exportMethodInputTitle   string
		exportMethodInputSummary string
		exportMethodEnvKey       string
		exportMethods            []string
	)

	if projectType == XcodeProjectTypeIOS {
		exportMethodInputTitle = DistributionMethodInputTitle
		exportMethodInputSummary = DistributionMethodInputSummary
		exportMethodEnvKey = DistributionMethodEnvKey
		exportMethods = IosExportMethods
	} else {
		exportMethodInputTitle = ExportMethodInputTitle
		exportMethodInputSummary = ExportMethodInputSummary
		exportMethodEnvKey = ExportMethodEnvKey
		exportMethods = MacExportMethods
	}

	var (
		allWarnings         = result.Warnings
		iconsForAllProjects models.Icons
		configDescriptors   []ConfigDescriptor
	)

	projectPathOption := models.NewOption(ProjectPathInputTitle, ProjectPathInputSummary, ProjectPathInputEnvKey, models.TypeSelector)
	for _, project := range result.Projects {
		allWarnings = append(allWarnings, project.Warnings...)

		schemeOption := models.NewOption(SchemeInputTitle, SchemeInputSummary, SchemeInputEnvKey, models.TypeSelector)
		projectPathOption.AddOption(project.RelPath, schemeOption)

		for _, scheme := range project.Schemes {
			exportMethodOption := models.NewOption(exportMethodInputTitle, exportMethodInputSummary, exportMethodEnvKey, models.TypeSelector)
			schemeOption.AddOption(scheme.Name, exportMethodOption)

			iconsForAllProjects = append(iconsForAllProjects, scheme.Icons...)

			iconIDs := []string{}
			for _, icon := range scheme.Icons {
				iconIDs = append(iconIDs, icon.Filename)
			}

			for _, exportMethod := range exportMethods {
				// Whether app-clip export Step is added later depends on the used export method
				configDescriptor := NewConfigDescriptor(project.IsPodWorkspace, project.CarthageCommand, scheme.HasXCTests, scheme.HasAppClip, exportMethod, scheme.Missing)
				configDescriptors = append(configDescriptors, configDescriptor)
				configOption := models.NewConfigOption(configDescriptor.ConfigName(projectType), iconIDs)

				exportMethodOption.AddConfig(exportMethod, configOption)
			}
		}
	}

	configDescriptors = RemoveDuplicatedConfigDescriptors(configDescriptors, projectType)
	if len(configDescriptors) == 0 {
		log.TErrorf("No valid %s config found", string(projectType))
		return models.OptionNode{}, []ConfigDescriptor{}, nil, allWarnings, fmt.Errorf("No valid %s config found", string(projectType))
	}

	return *projectPathOption, configDescriptors, iconsForAllProjects, allWarnings, nil
}

// GenerateDefaultOptions ...
func GenerateDefaultOptions(projectType XcodeProjectType) models.OptionNode {
	projectPathOption := models.NewOption(ProjectPathInputTitle, ProjectPathInputSummary, ProjectPathInputEnvKey, models.TypeUserInput)

	schemeOption := models.NewOption(SchemeInputTitle, SchemeInputSummary, SchemeInputEnvKey, models.TypeUserInput)
	projectPathOption.AddOption("", schemeOption)

	var exportMethodInputTitle string
	var exportMethodInputSummary string
	var exportMethodEnvKey string
	var exportMethods []string

	if projectType == XcodeProjectTypeIOS {
		exportMethodInputTitle = DistributionMethodInputTitle
		exportMethodInputSummary = DistributionMethodInputSummary
		exportMethodEnvKey = DistributionMethodEnvKey
		exportMethods = IosExportMethods
	} else {
		exportMethodInputTitle = ExportMethodInputTitle
		exportMethodInputSummary = ExportMethodInputSummary
		exportMethodEnvKey = ExportMethodEnvKey
		exportMethods = MacExportMethods
	}

	exportMethodOption := models.NewOption(exportMethodInputTitle, exportMethodInputSummary, exportMethodEnvKey, models.TypeSelector)
	schemeOption.AddOption("", exportMethodOption)

	for _, exportMethod := range exportMethods {
		configOption := models.NewConfigOption(fmt.Sprintf(defaultConfigNameFormat, string(projectType)), nil)
		exportMethodOption.AddConfig(exportMethod, configOption)
	}

	return *projectPathOption
}

// GenerateConfigBuilder ...
func GenerateConfigBuilder(
	projectType XcodeProjectType,
	isPrivateRepository,
	hasPodfile,
	hasTest,
	hasAppClip,
	missingSharedSchemes,
	includeCache bool,
	carthageCommand,
	exportMethod string,
) models.ConfigBuilderModel {
	configBuilder := models.NewDefaultConfigBuilder()

	params := workflowSetupParams{
		projectType:          projectType,
		configBuilder:        configBuilder,
		isPrivateRepository:  isPrivateRepository,
		includeCache:         includeCache,
		missingSharedSchemes: missingSharedSchemes,
		hasTests:             hasTest,
		hasAppClip:           hasAppClip,
		hasPodfile:           hasPodfile,
		carthageCommand:      carthageCommand,
		exportMethod:         exportMethod,
	}

	createPrimaryWorkflow(params)
	createDeployWorkflow(params)

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
func GenerateConfig(projectType XcodeProjectType, configDescriptors []ConfigDescriptor, isPrivateRepository bool) (models.BitriseConfigMap, error) {
	bitriseDataMap := models.BitriseConfigMap{}
	for _, descriptor := range configDescriptors {
		configBuilder := GenerateConfigBuilder(
			projectType,
			isPrivateRepository,
			descriptor.HasPodfile,
			descriptor.HasTest,
			descriptor.HasAppClip,
			descriptor.MissingSharedSchemes,
			true,
			descriptor.CarthageCommand,
			descriptor.ExportMethod)

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
func GenerateDefaultConfig(projectType XcodeProjectType) (models.BitriseConfigMap, error) {
	configBuilder := GenerateConfigBuilder(
		projectType,
		true,
		true,
		true,
		false,
		true,
		true,
		"",
		"")

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
