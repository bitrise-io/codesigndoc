package xamarin

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/bitrise-io/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/log"
)

const scannerName = "xamarin"

const (
	defaultConfigName = "default-xamarin-config"
)

const (
	xamarinSolutionInputKey     = "xamarin_solution"
	xamarinSolutionInputEnvKey  = "BITRISE_PROJECT_PATH"
	xamarinSolutionInputTitle   = "Path to the Xamarin Solution file"
	xamarinSolutionInputSummary = "Your solution file has to contain all the solution configurations you wish to use on Bitrise. A solution configuration specifies how projects in the solution are to be built and deployed."
)

const (
	xamarinConfigurationInputKey     = "xamarin_configuration"
	xamarinConfigurationInputEnvKey  = "BITRISE_XAMARIN_CONFIGURATION"
	xamarinConfigurationInputTitle   = "Xamarin solution configuration"
	xamarinConfigurationInputSummary = "The Xamarin solution configuration that you wish to run in your first build. You can change this at any time in your Workflows."
)

const (
	xamarinPlatformInputKey     = "xamarin_platform"
	xamarinPlatformInputEnvKey  = "BITRISE_XAMARIN_PLATFORM"
	xamarinPlatformInputTitle   = "Xamarin solution platform"
	xamarinPlatformInputSummary = ""
)

const (
	xamarinIosLicenceInputKey     = "xamarin_ios_license"
	xamarinAndroidLicenceInputKey = "xamarin_android_license"
	xamarinMacLicenseInputKey     = "xamarin_mac_license"
)

func configName(hasNugetPackages, hasXamarinComponents bool) string {
	name := "xamarin-"
	if hasNugetPackages {
		name = name + "nuget-"
	}
	if hasXamarinComponents {
		name = name + "components-"
	}
	return name + "config"
}

//--------------------------------------------------
// Scanner
//--------------------------------------------------

// Scanner ...
type Scanner struct {
	FileList      []string
	SolutionFiles []string

	HasNugetPackages     bool
	HasXamarinComponents bool

	HasIosProject     bool
	HasAndroidProject bool
	HasMacProject     bool
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return scannerName
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	fileList, err := utility.ListPathInDirSortedByComponents(searchDir, true)
	if err != nil {
		return false, fmt.Errorf("failed to search for files in (%s), error: %s", searchDir, err)
	}
	scanner.FileList = fileList

	// Search for solution file
	log.TInfof("Searching for solution files")

	solutionFiles, err := FilterSolutionFiles(fileList)
	if err != nil {
		return false, fmt.Errorf("failed to search for solution files, error: %s", err)
	}

	scanner.SolutionFiles = solutionFiles

	log.TPrintf("%d solution files detected", len(solutionFiles))
	for _, file := range solutionFiles {
		log.TPrintf("- %s", file)
	}

	if len(solutionFiles) == 0 {
		log.TPrintf("platform not detected")
		return false, nil
	}

	log.TSuccessf("Platform detected")

	return true, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return []string{}
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	log.TInfof("Searching for NuGet packages & Xamarin Components")

	warnings := models.Warnings{}

	for _, file := range scanner.FileList {
		// Search for nuget packages
		if !scanner.HasNugetPackages {
			baseName := filepath.Base(file)
			if baseName == "packages.config" {
				scanner.HasNugetPackages = true
			}
		}

		// If adding a component:
		// /Components/[COMPONENT_NAME]/ dir added
		// ItemGroup/XamarinComponentReference added to the project
		// packages.config added to the project's folder
		if !scanner.HasXamarinComponents {
			componentsExp := regexp.MustCompile(".*Components/.+")
			if result := componentsExp.FindString(file); result != "" {
				scanner.HasXamarinComponents = true
			}
		}

		if scanner.HasNugetPackages && scanner.HasXamarinComponents {
			break
		}
	}

	if scanner.HasNugetPackages {
		log.TPrintf("Nuget packages found")
	} else {
		log.TPrintf("NO Nuget packages found")
	}

	if scanner.HasXamarinComponents {
		log.TPrintf("Xamarin Components found")
	} else {
		log.TPrintf("NO Xamarin Components found")
	}

	// Check for solution configs
	validSolutionMap := map[string]map[string][]string{}
	for _, solutionFile := range scanner.SolutionFiles {
		log.TInfof("Inspecting solution file: %s", solutionFile)

		configs, err := GetSolutionConfigs(solutionFile)
		if err != nil {
			log.TWarnf("Failed to get solution configs, error: %s", err)
			warnings = append(warnings, fmt.Sprintf("Failed to get solution (%s) configs, error: %s", solutionFile, err))
			continue
		}

		if len(configs) > 0 {
			log.TPrintf("%d configurations found", len(configs))
			for config, platforms := range configs {
				log.TPrintf("- %s with platforms: %v", config, platforms)
			}

			validSolutionMap[solutionFile] = configs
		} else {
			log.TWarnf("No config found for %s", solutionFile)
			warnings = append(warnings, fmt.Sprintf("No configs found for solution: %s", solutionFile))
		}
	}

	if len(validSolutionMap) == 0 {
		log.TErrorf("No valid solution file found")
		return models.OptionNode{}, warnings, nil, errors.New("No valid solution file found")
	}

	// Check for solution projects
	xamarinSolutionOption := models.NewOption(xamarinSolutionInputTitle, xamarinSolutionInputSummary, xamarinSolutionInputEnvKey, models.TypeSelector)

	for solutionFile, configMap := range validSolutionMap {
		xamarinConfigurationOption := models.NewOption(xamarinConfigurationInputTitle, xamarinConfigurationInputSummary, xamarinConfigurationInputEnvKey, models.TypeSelector)
		xamarinSolutionOption.AddOption(solutionFile, xamarinConfigurationOption)

		for config, platforms := range configMap {
			xamarinPlatformOption := models.NewOption(xamarinPlatformInputTitle, xamarinPlatformInputSummary, xamarinPlatformInputEnvKey, models.TypeSelector)
			xamarinConfigurationOption.AddOption(config, xamarinPlatformOption)

			for _, platform := range platforms {
				configOption := models.NewConfigOption(configName(scanner.HasNugetPackages, scanner.HasXamarinComponents), nil)
				xamarinPlatformOption.AddConfig(platform, configOption)
			}
		}
	}

	return *xamarinSolutionOption, warnings, nil, nil
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionNode {
	xamarinSolutionOption := models.NewOption(xamarinSolutionInputTitle, xamarinSolutionInputSummary, xamarinSolutionInputEnvKey, models.TypeUserInput)

	xamarinConfigurationOption := models.NewOption(xamarinConfigurationInputTitle, xamarinConfigurationInputSummary, xamarinConfigurationInputEnvKey, models.TypeUserInput)
	xamarinSolutionOption.AddOption("", xamarinConfigurationOption)

	xamarinPlatformOption := models.NewOption(xamarinPlatformInputTitle, xamarinPlatformInputSummary, xamarinPlatformInputEnvKey, models.TypeUserInput)
	xamarinConfigurationOption.AddOption("", xamarinPlatformOption)

	configOption := models.NewConfigOption(defaultConfigName, nil)
	xamarinPlatformOption.AddConfig("", configOption)

	return *xamarinSolutionOption
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

	// XamarinUserManagement
	if scanner.HasXamarinComponents {
		inputs := []envmanModels.EnvironmentItemModel{}
		if scanner.HasIosProject {
			inputs = append(inputs, envmanModels.EnvironmentItemModel{xamarinIosLicenceInputKey: "yes"})
		}
		if scanner.HasAndroidProject {
			inputs = append(inputs, envmanModels.EnvironmentItemModel{xamarinAndroidLicenceInputKey: "yes"})
		}
		if scanner.HasMacProject {
			inputs = append(inputs, envmanModels.EnvironmentItemModel{xamarinMacLicenseInputKey: "yes"})
		}

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XamarinUserManagementStepListItem(inputs...))
	}

	// NugetRestore
	if scanner.HasNugetPackages {
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NugetRestoreStepListItem())
	}

	// XamarinComponentsRestore
	if scanner.HasXamarinComponents {
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XamarinComponentsRestoreStepListItem())
	}

	// XamarinArchive
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XamarinArchiveStepListItem(
		envmanModels.EnvironmentItemModel{xamarinSolutionInputKey: "$" + xamarinSolutionInputEnvKey},
		envmanModels.EnvironmentItemModel{xamarinConfigurationInputKey: "$" + xamarinConfigurationInputEnvKey},
		envmanModels.EnvironmentItemModel{xamarinPlatformInputKey: "$" + xamarinPlatformInputEnvKey},
	))

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

	config, err := configBuilder.Generate(scannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	return models.BitriseConfigMap{
		configName(scanner.HasNugetPackages, scanner.HasXamarinComponents): string(data),
	}, nil
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CertificateAndProfileInstallerStepListItem())
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XamarinUserManagementStepListItem())

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NugetRestoreStepListItem())
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XamarinComponentsRestoreStepListItem())

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XamarinArchiveStepListItem(
		envmanModels.EnvironmentItemModel{xamarinSolutionInputKey: "$" + xamarinSolutionInputEnvKey},
		envmanModels.EnvironmentItemModel{xamarinConfigurationInputKey: "$" + xamarinConfigurationInputEnvKey},
		envmanModels.EnvironmentItemModel{xamarinPlatformInputKey: "$" + xamarinPlatformInputEnvKey},
	))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

	config, err := configBuilder.Generate(scannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	return models.BitriseConfigMap{
		defaultConfigName: string(data),
	}, nil
}
