package ionic

import (
	"fmt"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/cordova"
	"github.com/bitrise-core/bitrise-init/scanners/ios"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

const scannerName = "ionic"

const (
	configName        = "ionic-config"
	defaultConfigName = "default-ionic-config"
)

// Step Inputs
const (
	workDirInputKey    = "workdir"
	workDirInputTitle  = "Directory of Ionic Config.xml"
	workDirInputEnvKey = "IONIC_WORK_DIR"
)

const (
	platformInputKey    = "platform"
	platformInputTitle  = "Platform to use in ionic-cli commands"
	platformInputEnvKey = "IONIC_PLATFORM"
)

const (
	targetInputKey = "target"
	targetEmulator = "emulator"
)

//------------------
// ScannerInterface
//------------------

// Scanner ...
type Scanner struct {
	cordovaConfigPth    string
	relCordovaConfigDir string
	searchDir           string
	hasKarmaJasmineTest bool
	hasJasmineTest      bool
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

	// Search for config.xml file
	log.TInfof("Searching for config.xml file")

	configXMLPth, err := cordova.FilterRootConfigXMLFile(fileList)
	if err != nil {
		return false, fmt.Errorf("failed to search for config.xml file, error: %s", err)
	}

	log.TPrintf("config.xml: %s", configXMLPth)

	if configXMLPth == "" {
		log.TPrintf("platform not detected")
		return false, nil
	}

	widget, err := cordova.ParseConfigXML(configXMLPth)
	if err != nil {
		log.TPrintf("can not parse config.xml as a Cordova widget, error: %s", err)
		log.TPrintf("platform not detected")
		return false, nil
	}

	// ensure it is a cordova widget
	if !strings.Contains(widget.XMLNSCDV, "cordova.apache.org") {
		log.TPrintf("config.xml propert: xmlns:cdv does not contain cordova.apache.org")
		log.TPrintf("platform not detected")
		return false, nil
	}

	// ensure it is an ionic project
	projectBaseDir := filepath.Dir(configXMLPth)

	ionicProjectExist, err := pathutil.IsPathExists(filepath.Join(projectBaseDir, "ionic.project"))
	if err != nil {
		return false, fmt.Errorf("failed to check if project is an ionic project, error: %s", err)
	}

	ionicConfigExist, err := pathutil.IsPathExists(filepath.Join(projectBaseDir, "ionic.config.json"))
	if err != nil {
		return false, fmt.Errorf("failed to check if project is an ionic project, error: %s", err)
	}

	if !ionicProjectExist && !ionicConfigExist {
		log.Printf("no ionic.project file nor ionic.config.json found, seems to be a cordova project")
		return false, nil
	}

	log.TSuccessf("Platform detected")

	scanner.cordovaConfigPth = configXMLPth
	scanner.searchDir = searchDir

	return true, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return []string{
		string(ios.XcodeProjectTypeIOS),
		string(ios.XcodeProjectTypeMacOS),
		cordova.ScannerName,
		android.ScannerName,
	}
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, error) {
	warnings := models.Warnings{}
	projectRootDir := filepath.Dir(scanner.cordovaConfigPth)

	packagesJSONPth := filepath.Join(projectRootDir, "package.json")
	packages, err := utility.ParsePackagesJSON(packagesJSONPth)
	if err != nil {
		return models.OptionNode{}, warnings, err
	}

	// Search for karma/jasmine tests
	log.TPrintf("Searching for karma/jasmine test")

	karmaTestDetected := false

	karmaJasmineDependencyFound := false
	for dependency := range packages.Dependencies {
		if strings.Contains(dependency, "karma-jasmine") {
			karmaJasmineDependencyFound = true
		}
	}
	if !karmaJasmineDependencyFound {
		for dependency := range packages.DevDependencies {
			if strings.Contains(dependency, "karma-jasmine") {
				karmaJasmineDependencyFound = true
			}
		}
	}
	log.TPrintf("karma-jasmine dependency found: %v", karmaJasmineDependencyFound)

	if karmaJasmineDependencyFound {
		karmaConfigJSONPth := filepath.Join(projectRootDir, "karma.conf.js")
		if exist, err := pathutil.IsPathExists(karmaConfigJSONPth); err != nil {
			return models.OptionNode{}, warnings, err
		} else if exist {
			karmaTestDetected = true
		}
	}
	log.TPrintf("karma.conf.js found: %v", karmaTestDetected)

	scanner.hasKarmaJasmineTest = karmaTestDetected
	// ---

	// Search for jasmine tests
	jasminTestDetected := false

	if !karmaTestDetected {
		log.TPrintf("Searching for jasmine test")

		jasmineDependencyFound := false
		for dependency := range packages.Dependencies {
			if strings.Contains(dependency, "jasmine") {
				jasmineDependencyFound = true
				break
			}
		}
		if !jasmineDependencyFound {
			for dependency := range packages.DevDependencies {
				if strings.Contains(dependency, "jasmine") {
					jasmineDependencyFound = true
					break
				}
			}
		}
		log.TPrintf("jasmine dependency found: %v", jasmineDependencyFound)

		if jasmineDependencyFound {
			jasmineConfigJSONPth := filepath.Join(projectRootDir, "spec", "support", "jasmine.json")
			if exist, err := pathutil.IsPathExists(jasmineConfigJSONPth); err != nil {
				return models.OptionNode{}, warnings, err
			} else if exist {
				jasminTestDetected = true
			}
		}

		log.TPrintf("jasmine.json found: %v", jasminTestDetected)

		scanner.hasJasmineTest = jasminTestDetected
	}
	// ---

	// Get relative config.xml dir
	cordovaConfigDir := filepath.Dir(scanner.cordovaConfigPth)
	relCordovaConfigDir, err := utility.RelPath(scanner.searchDir, cordovaConfigDir)
	if err != nil {
		return models.OptionNode{}, warnings, fmt.Errorf("Failed to get relative config.xml dir path, error: %s", err)
	}
	if relCordovaConfigDir == "." {
		// config.xml placed in the search dir, no need to change-dir in the workflows
		relCordovaConfigDir = ""
	}
	scanner.relCordovaConfigDir = relCordovaConfigDir
	// ---

	// Options
	var rootOption *models.OptionNode

	platforms := []string{"ios", "android", "ios,android"}

	if relCordovaConfigDir != "" {
		rootOption = models.NewOption(workDirInputTitle, workDirInputEnvKey)

		projectTypeOption := models.NewOption(platformInputTitle, platformInputEnvKey)
		rootOption.AddOption(relCordovaConfigDir, projectTypeOption)

		for _, platform := range platforms {
			configOption := models.NewConfigOption(configName)
			projectTypeOption.AddConfig(platform, configOption)
		}
	} else {
		rootOption = models.NewOption(platformInputTitle, platformInputEnvKey)

		for _, platform := range platforms {
			configOption := models.NewConfigOption(configName)
			rootOption.AddConfig(platform, configOption)
		}
	}
	// ---

	return *rootOption, warnings, nil
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionNode {
	workDirOption := models.NewOption(workDirInputTitle, workDirInputEnvKey)

	projectTypeOption := models.NewOption(platformInputTitle, platformInputEnvKey)
	workDirOption.AddOption("_", projectTypeOption)

	platforms := []string{
		"ios",
		"android",
		"ios,android",
	}
	for _, platform := range platforms {
		configOption := models.NewConfigOption(defaultConfigName)
		projectTypeOption.AddConfig(platform, configOption)
	}

	return *workDirOption
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)

	workdirEnvList := []envmanModels.EnvironmentItemModel{}
	if scanner.relCordovaConfigDir != "" {
		workdirEnvList = append(workdirEnvList, envmanModels.EnvironmentItemModel{workDirInputKey: "$" + workDirInputEnvKey})
	}

	if scanner.hasJasmineTest || scanner.hasKarmaJasmineTest {
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))

		// CI
		if scanner.hasKarmaJasmineTest {
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.KarmaJasmineTestRunnerStepListItem(workdirEnvList...))
		} else if scanner.hasJasmineTest {
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.JasmineTestRunnerStepListItem(workdirEnvList...))
		}
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

		// CD
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(false)...)
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))

		if scanner.hasKarmaJasmineTest {
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.KarmaJasmineTestRunnerStepListItem(workdirEnvList...))
		} else if scanner.hasJasmineTest {
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.JasmineTestRunnerStepListItem(workdirEnvList...))
		}

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.GenerateCordovaBuildConfigStepListItem())

		ionicArchiveEnvs := []envmanModels.EnvironmentItemModel{
			envmanModels.EnvironmentItemModel{platformInputKey: "$" + platformInputEnvKey},
			envmanModels.EnvironmentItemModel{targetInputKey: targetEmulator},
		}
		if scanner.relCordovaConfigDir != "" {
			ionicArchiveEnvs = append(ionicArchiveEnvs, envmanModels.EnvironmentItemModel{workDirInputKey: "$" + workDirInputEnvKey})
		}
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.IonicArchiveStepListItem(ionicArchiveEnvs...))
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)

		config, err := configBuilder.Generate(scannerName)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		return models.BitriseConfigMap{
			configName: string(data),
		}, nil
	}

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CertificateAndProfileInstallerStepListItem())
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.GenerateCordovaBuildConfigStepListItem())

	ionicArchiveEnvs := []envmanModels.EnvironmentItemModel{
		envmanModels.EnvironmentItemModel{platformInputKey: "$" + platformInputEnvKey},
		envmanModels.EnvironmentItemModel{targetInputKey: targetEmulator},
	}
	if scanner.relCordovaConfigDir != "" {
		ionicArchiveEnvs = append(ionicArchiveEnvs, envmanModels.EnvironmentItemModel{workDirInputKey: "$" + workDirInputEnvKey})
	}
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.IonicArchiveStepListItem(ionicArchiveEnvs...))
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
		configName: string(data),
	}, nil
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(
		envmanModels.EnvironmentItemModel{"command": "install"},
		envmanModels.EnvironmentItemModel{workDirInputKey: "$" + workDirInputEnvKey}))

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.GenerateCordovaBuildConfigStepListItem())
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.IonicArchiveStepListItem(
		envmanModels.EnvironmentItemModel{workDirInputKey: "$" + workDirInputEnvKey},
		envmanModels.EnvironmentItemModel{platformInputKey: "$" + platformInputEnvKey},
		envmanModels.EnvironmentItemModel{targetInputKey: targetEmulator}))

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
