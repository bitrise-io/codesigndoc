package flutter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/ios"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/xcode-project/xcworkspace"
	yaml "gopkg.in/yaml.v2"
)

const (
	scannerName                = "flutter"
	configName                 = "flutter-config"
	projectLocationInputKey    = "project_location"
	platformInputKey           = "platform"
	defaultIOSConfiguration    = "Release"
	projectLocationInputEnvKey = "BITRISE_FLUTTER_PROJECT_LOCATION"
	projectLocationInputTitle  = "Project Location"
	projectTypeInputEnvKey     = "BITRISE_FLUTTER_PROJECT_TYPE"
	projectTypeInputTitle      = "Project Type"
	testsInputTitle            = "Run tests found in the project"
	platformInputTitle         = "Platform"
)

var (
	projectTypes = []string{
		"app",
		"plugin",
		"package",
	}
	platforms = []string{
		"none",
		"android",
		"ios",
		"both",
	}
)

//------------------
// ScannerInterface
//------------------

// Scanner ...
type Scanner struct {
	projects []project
}

type project struct {
	path              string
	xcodeProjectPaths map[string][]string
	hasTest           bool
	hasIosProject     bool
	hasAndroidProject bool
}

type pubspec struct {
	Name string `yaml:"name"`
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return scannerName
}

func findProjectLocations(searchDir string) ([]string, error) {
	fileList, err := utility.ListPathInDirSortedByComponents(searchDir, true)
	if err != nil {
		return nil, err
	}

	filters := []utility.FilterFunc{
		utility.BaseFilter("pubspec.yaml", true),
	}

	paths, err := utility.FilterPaths(fileList, filters...)
	if err != nil {
		return nil, err
	}

	for i, path := range paths {
		paths[i] = filepath.Dir(path)
	}

	return paths, nil
}

func findWorkspaceLocations(projectLocation string) ([]string, error) {
	fileList, err := utility.ListPathInDirSortedByComponents(projectLocation, true)
	if err != nil {
		return nil, err
	}

	for i, file := range fileList {
		fileList[i] = filepath.Join(projectLocation, file)
	}

	filters := []utility.FilterFunc{
		ios.AllowXCWorkspaceExtFilter,
		ios.AllowIsDirectoryFilter,
		ios.ForbidEmbeddedWorkspaceRegexpFilter,
		ios.ForbidGitDirComponentFilter,
		ios.ForbidPodsDirComponentFilter,
		ios.ForbidCarthageDirComponentFilter,
		ios.ForbidFramworkComponentWithExtensionFilter,
		ios.ForbidCordovaLibDirComponentFilter,
		ios.ForbidNodeModulesComponentFilter,
	}

	return utility.FilterPaths(fileList, filters...)
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	log.TInfof("Search for project(s)")
	projectLocations, err := findProjectLocations(searchDir)
	if err != nil {
		return false, err
	}

	log.TPrintf("Paths containing pubspec.yaml(%d):", len(projectLocations))
	for _, p := range projectLocations {
		log.TPrintf("- %s", p)
	}
	log.TPrintf("")

	log.TInfof("Fetching pubspec.yaml files")
projects:
	for _, projectLocation := range projectLocations {
		var proj project

		pubspecPath := filepath.Join(projectLocation, "pubspec.yaml")
		pubspecFile, err := os.Open(pubspecPath)
		if err != nil {
			log.TErrorf("Failed to open pubspec.yaml file at: %s, error: %s", pubspecPath, err)
			return false, err
		}

		var ps pubspec
		if err := yaml.NewDecoder(pubspecFile).Decode(&ps); err != nil {
			log.TErrorf("Failed to decode yaml pubspec.yaml file at: %s, error: %s", pubspecPath, err)
			return false, err
		}

		testsDirPath := filepath.Join(projectLocation, "test")
		if exists, err := pathutil.IsDirExists(testsDirPath); err == nil && exists {
			if files, err := ioutil.ReadDir(testsDirPath); err == nil && len(files) > 0 {
				for _, file := range files {
					if strings.HasSuffix(file.Name(), "_test.dart") {
						proj.hasTest = true
						break
					}
				}
			}
		}

		iosProjPath := filepath.Join(projectLocation, "ios", "Runner.xcworkspace")
		if exists, err := pathutil.IsPathExists(iosProjPath); err == nil && exists {
			proj.hasIosProject = true
		}

		androidProjPath := filepath.Join(projectLocation, "android", "build.gradle")
		if exists, err := pathutil.IsPathExists(androidProjPath); err == nil && exists {
			proj.hasAndroidProject = true
		}

		log.TPrintf("- Project name: %s", ps.Name)
		log.TPrintf("  Path: %s", projectLocation)
		log.TPrintf("  HasTest: %t", proj.hasTest)
		log.TPrintf("  HasAndroidProject: %t", proj.hasAndroidProject)
		log.TPrintf("  HasIosProject: %t", proj.hasIosProject)

		proj.path = projectLocation

		if proj.hasIosProject {
			if workspaceLocations, err := findWorkspaceLocations(filepath.Join(projectLocation, "ios")); err != nil {
				log.TWarnf("Failed to check path at: %s, error: %s", filepath.Join(projectLocation, "ios"), err)
			} else {
				log.TPrintf("  XCWorkspaces(%d):", len(workspaceLocations))

				for _, workspaceLocation := range workspaceLocations {
					log.TPrintf("    Path: %s", workspaceLocation)
					ws, err := xcworkspace.Open(workspaceLocation)
					if err != nil {
						continue projects
					}
					schemeMap, err := ws.Schemes()
					if err != nil {
						continue projects
					}

					proj.xcodeProjectPaths = map[string][]string{}

					for _, schemes := range schemeMap {
						if len(schemes) > 0 {
							log.TPrintf("    Schemes(%d):", len(schemes))
						}
						for _, scheme := range schemes {
							log.TPrintf("    - %s", scheme.Name)
							proj.xcodeProjectPaths[workspaceLocation] = append(proj.xcodeProjectPaths[workspaceLocation], scheme.Name)
						}
					}
				}
			}
		}

		scanner.projects = append(scanner.projects, proj)
	}

	if len(scanner.projects) == 0 {
		return false, nil
	}

	return true, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return []string{
		string(ios.XcodeProjectTypeIOS),
		android.ScannerName,
	}
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, error) {
	flutterProjectLocationOption := models.NewOption(projectLocationInputTitle, projectLocationInputEnvKey)

	for _, project := range scanner.projects {
		if project.hasTest {
			flutterProjectHasTestOption := models.NewOption(testsInputTitle, "")
			flutterProjectLocationOption.AddOption(project.path, flutterProjectHasTestOption)

			for _, v := range []string{"yes", "no"} {
				cfg := configName
				if v == "yes" {
					cfg += "-test"
				}

				if project.hasIosProject || project.hasAndroidProject {
					if project.hasIosProject {
						projectPathOption := models.NewOption(ios.ProjectPathInputTitle, ios.ProjectPathInputEnvKey)
						flutterProjectHasTestOption.AddOption(v, projectPathOption)

						for xcodeWorkspacePath, schemes := range project.xcodeProjectPaths {
							schemeOption := models.NewOption(ios.SchemeInputTitle, ios.SchemeInputEnvKey)
							projectPathOption.AddOption(xcodeWorkspacePath, schemeOption)

							for _, scheme := range schemes {
								exportMethodOption := models.NewOption(ios.IosExportMethodInputTitle, ios.ExportMethodInputEnvKey)
								schemeOption.AddOption(scheme, exportMethodOption)

								for _, exportMethod := range ios.IosExportMethods {
									configOption := models.NewConfigOption(cfg + "-app-" + getBuildablePlatform(project.hasAndroidProject, project.hasIosProject))
									exportMethodOption.AddConfig(exportMethod, configOption)
								}
							}
						}
					} else {
						configOption := models.NewConfigOption(cfg + "-app-" + getBuildablePlatform(project.hasAndroidProject, project.hasIosProject))
						flutterProjectHasTestOption.AddOption(v, configOption)
					}
				} else {
					configOption := models.NewConfigOption(cfg)
					flutterProjectHasTestOption.AddOption(v, configOption)
				}
			}
		} else {
			cfg := configName

			if project.hasIosProject || project.hasAndroidProject {
				if project.hasIosProject {
					projectPathOption := models.NewOption(ios.ProjectPathInputTitle, ios.ProjectPathInputEnvKey)
					flutterProjectLocationOption.AddOption(project.path, projectPathOption)

					for xcodeWorkspacePath, schemes := range project.xcodeProjectPaths {
						schemeOption := models.NewOption(ios.SchemeInputTitle, ios.SchemeInputEnvKey)
						projectPathOption.AddOption(xcodeWorkspacePath, schemeOption)

						for _, scheme := range schemes {
							exportMethodOption := models.NewOption(ios.IosExportMethodInputTitle, ios.ExportMethodInputEnvKey)
							schemeOption.AddOption(scheme, exportMethodOption)

							for _, exportMethod := range ios.IosExportMethods {
								configOption := models.NewConfigOption(cfg + "-app-" + getBuildablePlatform(project.hasAndroidProject, project.hasIosProject))
								exportMethodOption.AddConfig(exportMethod, configOption)
							}
						}
					}
				} else {
					configOption := models.NewConfigOption(cfg + "-app-" + getBuildablePlatform(project.hasAndroidProject, project.hasIosProject))
					flutterProjectLocationOption.AddOption(project.path, configOption)
				}
			} else {
				configOption := models.NewConfigOption(cfg)
				flutterProjectLocationOption.AddOption(project.path, configOption)
			}
		}
	}

	return *flutterProjectLocationOption, nil, nil
}

func getBuildablePlatform(hasAndroidProject, hasIosProject bool) string {
	switch {
	case hasAndroidProject && !hasIosProject:
		return "android"
	case !hasAndroidProject && hasIosProject:
		return "ios"
	default:
		return "both"
	}
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionNode {
	flutterProjectLocationOption := models.NewOption(projectLocationInputTitle, projectLocationInputEnvKey)

	flutterProjectHasTestOption := models.NewOption(testsInputTitle, "")
	flutterProjectLocationOption.AddOption("_", flutterProjectHasTestOption)

	for _, v := range []string{"yes", "no"} {
		cfg := configName
		if v == "yes" {
			cfg += "-test"
		}
		flutterPlatformOption := models.NewOption(platformInputTitle, "")
		flutterProjectHasTestOption.AddOption(v, flutterPlatformOption)

		for _, platform := range platforms {
			if platform != "none" {
				if platform != "android" {
					projectPathOption := models.NewOption(ios.ProjectPathInputTitle, ios.ProjectPathInputEnvKey)
					flutterPlatformOption.AddOption(platform, projectPathOption)

					schemeOption := models.NewOption(ios.SchemeInputTitle, ios.SchemeInputEnvKey)
					projectPathOption.AddOption("_", schemeOption)

					exportMethodOption := models.NewOption(ios.IosExportMethodInputTitle, ios.ExportMethodInputEnvKey)
					schemeOption.AddOption("_", exportMethodOption)

					for _, exportMethod := range ios.IosExportMethods {
						configOption := models.NewConfigOption(cfg + "-app-" + platform)
						exportMethodOption.AddConfig(exportMethod, configOption)
					}
				} else {
					configOption := models.NewConfigOption(cfg + "-app-" + platform)
					flutterPlatformOption.AddConfig(platform, configOption)
				}
			} else {
				configOption := models.NewConfigOption(cfg)
				flutterPlatformOption.AddConfig(platform, configOption)
			}
		}
	}

	return *flutterProjectLocationOption
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	return scanner.DefaultConfigs()
}

// DefaultConfigs ...
func (scanner Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configs := models.BitriseConfigMap{}

	for _, variant := range []struct {
		configID string
		test     bool
		deploy   bool
		platform string
	}{
		{test: false, deploy: false, configID: configName},
		{test: true, deploy: false, configID: configName + "-test"},
		{test: false, deploy: true, platform: "both", configID: configName + "-app-both"},
		{test: true, deploy: true, platform: "both", configID: configName + "-test-app-both"},
		{test: false, deploy: true, platform: "android", configID: configName + "-app-android"},
		{test: true, deploy: true, platform: "android", configID: configName + "-test-app-android"},
		{test: false, deploy: true, platform: "ios", configID: configName + "-app-ios"},
		{test: true, deploy: true, platform: "ios", configID: configName + "-test-app-ios"},
	} {
		configBuilder := models.NewDefaultConfigBuilder()

		// primary

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.FlutterInstallStepListItem())

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.FlutterAnalyzeStepListItem(
			envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
		))

		if variant.test {
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.FlutterTestStepListItem(
				envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
			))
		}

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

		// deploy

		if variant.deploy {
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(false)...)

			if variant.platform != "android" {
				configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())
			}

			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterInstallStepListItem())

			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterAnalyzeStepListItem(
				envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
			))

			if variant.test {
				configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterTestStepListItem(
					envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
				))
			}

			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterBuildStepListItem(
				envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
				envmanModels.EnvironmentItemModel{platformInputKey: variant.platform},
			))

			if variant.platform != "android" {
				configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(
					envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
					envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
					envmanModels.EnvironmentItemModel{ios.ExportMethodInputKey: "$" + ios.ExportMethodInputEnvKey},
					envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: defaultIOSConfiguration},
				))
			}

			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)
		}

		config, err := configBuilder.Generate(scannerName)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		configs[variant.configID] = string(data)
	}

	return configs, nil
}
