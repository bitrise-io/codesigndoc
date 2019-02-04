package expo

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/ios"
	"github.com/bitrise-core/bitrise-init/scanners/reactnative"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/xcode-project/serialized"
	yaml "gopkg.in/yaml.v2"
)

const (
	configName        = "react-native-expo-config"
	defaultConfigName = "default-" + configName
)

const deployWorkflowDescription = `## Configure Android part of the deploy workflow

To generate a signed APK:

1. Open the **Workflow** tab of your project on Bitrise.io
1. Add **Sign APK step right after Android Build step**
1. Click on **Code Signing** tab
1. Find the **ANDROID KEYSTORE FILE** section
1. Click or drop your file on the upload file field
1. Fill the displayed 3 input fields:
1. **Keystore password**
1. **Keystore alias**
1. **Private key password**
1. Click on **[Save metadata]** button

That's it! From now on, **Sign APK** step will receive your uploaded files.

## Configure iOS part of the deploy workflow

To generate IPA:

1. Open the **Workflow** tab of your project on Bitrise.io
1. Click on **Code Signing** tab
1. Find the **PROVISIONING PROFILE** section
1. Click or drop your file on the upload file field
1. Find the **CODE SIGNING IDENTITY** section
1. Click or drop your file on the upload file field
1. Click on **Workflows** tab
1. Select deploy workflow
1. Select **Xcode Archive & Export for iOS** step
1. Open **Force Build Settings** input group
1. Specify codesign settings
Set **Force code signing with Development Team**, **Force code signing with Code Signing Identity**  
and **Force code signing with Provisioning Profile** inputs regarding to the uploaded codesigning files
1. Specify manual codesign style
If the codesigning files, are generated manually on the Apple Developer Portal,  
you need to explicitly specify to use manual coedsign settings  
(as ejected rn projects have xcode managed codesigning turned on).  
To do so, add 'CODE_SIGN_STYLE="Manual"' to 'Additional options for xcodebuild call' input

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
`

// Name ...
const Name = "react-native-expo"

// Scanner ...
type Scanner struct {
	searchDir      string
	packageJSONPth string
	usesExpoKit    bool
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return Name
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	scanner.searchDir = searchDir

	log.TInfof("Collect package.json files")

	packageJSONPths, err := reactnative.CollectPackageJSONFiles(searchDir)
	if err != nil {
		return false, err
	}

	if len(packageJSONPths) == 0 {
		return false, nil
	}

	log.TPrintf("%d package.json file detected", len(packageJSONPths))
	for _, pth := range packageJSONPths {
		log.TPrintf("- %s", pth)
	}
	log.TPrintf("")

	log.TInfof("Filter package.json files with expo dependency")

	relevantPackageJSONPths := []string{}
	for _, packageJSONPth := range packageJSONPths {
		packages, err := utility.ParsePackagesJSON(packageJSONPth)
		if err != nil {
			log.Warnf("Failed to parse package json file: %s, skipping...", packageJSONPth)
			continue
		}

		_, found := packages.Dependencies["expo"]
		if !found {
			continue
		}

		// app.json file is a required part of react native projects and it exists next to the root package.json file
		appJSONPth := filepath.Join(filepath.Dir(packageJSONPth), "app.json")
		if exist, err := pathutil.IsPathExists(appJSONPth); err != nil {
			log.Warnf("Failed to check if app.json file exist at: %s, skipping package json file: %s, error: %s", appJSONPth, packageJSONPth, err)
			continue
		} else if !exist {
			log.Warnf("No app.json file exist at: %s, skipping package json file: %s", appJSONPth, packageJSONPth)
			continue
		}

		relevantPackageJSONPths = append(relevantPackageJSONPths, packageJSONPth)
	}

	log.TPrintf("%d package.json file detected with expo dependency", len(relevantPackageJSONPths))
	for _, pth := range relevantPackageJSONPths {
		log.TPrintf("- %s", pth)
	}
	log.TPrintf("")

	if len(relevantPackageJSONPths) == 0 {
		return false, nil
	} else if len(relevantPackageJSONPths) > 1 {
		log.TWarnf("Multiple package.json file found, using: %s\n", relevantPackageJSONPths[0])
	}

	scanner.packageJSONPth = relevantPackageJSONPths[0]
	return true, nil
}

func appJSONIssue(appJSONPth, reason, explanation string) string {
	return fmt.Sprintf("app.json file (%s) %s\n%s", appJSONPth, reason, explanation)
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, error) {
	warnings := models.Warnings{}

	// we need to know if the project uses the Expo Kit,
	// since its usage differentiates the eject process and the config options
	usesExpoKit := false

	fileList, err := utility.ListPathInDirSortedByComponents(scanner.searchDir, true)
	if err != nil {
		return models.OptionNode{}, warnings, err
	}

	filters := []utility.FilterFunc{
		utility.ExtensionFilter(".js", true),
		utility.ComponentFilter("node_modules", false),
	}
	sourceFiles, err := utility.FilterPaths(fileList, filters...)
	if err != nil {
		return models.OptionNode{}, warnings, err
	}

	re := regexp.MustCompile(`import .* from 'expo'`)

SourceFileLoop:
	for _, sourceFile := range sourceFiles {
		f, err := os.Open(sourceFile)
		if err != nil {
			return models.OptionNode{}, warnings, err
		}
		defer func() {
			if cerr := f.Close(); cerr != nil {
				log.Warnf("Failed to close: %s, error: %s", f.Name(), err)
			}
		}()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if match := re.FindString(scanner.Text()); match != "" {
				usesExpoKit = true
				break SourceFileLoop
			}
		}
		if err := scanner.Err(); err != nil {
			return models.OptionNode{}, warnings, err
		}
	}

	scanner.usesExpoKit = usesExpoKit
	log.TPrintf("Uses ExpoKit: %v", usesExpoKit)

	// ensure app.json contains the required information (for non interactive eject)
	// and predict the ejected project name
	var projectName string

	rootDir := filepath.Dir(scanner.packageJSONPth)
	appJSONPth := filepath.Join(rootDir, "app.json")
	appJSON, err := fileutil.ReadStringFromFile(appJSONPth)
	if err != nil {
		return models.OptionNode{}, warnings, err
	}
	var app serialized.Object
	if err := json.Unmarshal([]byte(appJSON), &app); err != nil {
		return models.OptionNode{}, warnings, err
	}

	if usesExpoKit {
		// if the project uses Expo Kit app.json needs to contain expo/ios/bundleIdentifier and expo/android/package entries
		// to be able to eject in non interactive mode
		errorMessage := `If the project uses Expo Kit the app.json file needs to contain:
- expo/name
- expo/ios/bundleIdentifier
- expo/android/package
entries.`

		expoObj, err := app.Object("expo")
		if err != nil {
			return models.OptionNode{}, warnings, errors.New(appJSONIssue(appJSONPth, "missing expo entry", errorMessage))
		}
		projectName, err = expoObj.String("name")
		if err != nil || projectName == "" {
			return models.OptionNode{}, warnings, errors.New(appJSONIssue(appJSONPth, "missing or empty expo/name entry", errorMessage))
		}

		iosObj, err := expoObj.Object("ios")
		if err != nil {
			return models.OptionNode{}, warnings, errors.New(appJSONIssue(appJSONPth, "missing expo/ios entry", errorMessage))
		}
		bundleID, err := iosObj.String("bundleIdentifier")
		if err != nil || bundleID == "" {
			return models.OptionNode{}, warnings, errors.New(appJSONIssue(appJSONPth, "missing or empty expo/ios/bundleIdentifier entry", errorMessage))
		}

		androidObj, err := expoObj.Object("android")
		if err != nil {
			return models.OptionNode{}, warnings, errors.New(appJSONIssue(appJSONPth, "missing expo/android entry", errorMessage))
		}
		packageName, err := androidObj.String("package")
		if err != nil || packageName == "" {
			return models.OptionNode{}, warnings, errors.New(appJSONIssue(appJSONPth, "missing or empty expo/android/package entry", errorMessage))
		}
	} else {
		// if the project does not use Expo Kit app.json needs to contain name and displayName entries
		// to be able to eject in non interactive mode
		errorMessage := `The app.json file needs to contain:
- name
- displayName
entries.`

		projectName, err = app.String("name")
		if err != nil || projectName == "" {
			return models.OptionNode{}, warnings, errors.New(appJSONIssue(appJSONPth, "missing or empty name entry", errorMessage))
		}
		displayName, err := app.String("displayName")
		if err != nil || displayName == "" {
			return models.OptionNode{}, warnings, errors.New(appJSONIssue(appJSONPth, "missing or empty displayName entry", errorMessage))
		}
	}

	log.TPrintf("Project name: %v", projectName)

	// ios options
	projectPathOption := models.NewOption(ios.ProjectPathInputTitle, ios.ProjectPathInputEnvKey)
	schemeOption := models.NewOption(ios.SchemeInputTitle, ios.SchemeInputEnvKey)

	if usesExpoKit {
		projectName = strings.ToLower(regexp.MustCompile(`(?i:[^a-z0-9_\-])`).ReplaceAllString(projectName, "-"))
		projectPathOption.AddOption(filepath.Join("./", "ios", projectName+".xcworkspace"), schemeOption)
	} else {
		projectPathOption.AddOption(filepath.Join("./", "ios", projectName+".xcodeproj"), schemeOption)
	}

	developmentTeamOption := models.NewOption("iOS Development team", "BITRISE_IOS_DEVELOPMENT_TEAM")
	schemeOption.AddOption(projectName, developmentTeamOption)

	exportMethodOption := models.NewOption(ios.IosExportMethodInputTitle, ios.ExportMethodInputEnvKey)
	developmentTeamOption.AddOption("_", exportMethodOption)

	// android options
	packageJSONDir := filepath.Dir(scanner.packageJSONPth)
	relPackageJSONDir, err := utility.RelPath(scanner.searchDir, packageJSONDir)
	if err != nil {
		return models.OptionNode{}, warnings, fmt.Errorf("Failed to get relative package.json dir path, error: %s", err)
	}
	if relPackageJSONDir == "." {
		// package.json placed in the search dir, no need to change-dir in the workflows
		relPackageJSONDir = ""
	}

	var moduleOption *models.OptionNode
	if relPackageJSONDir == "" {
		projectLocationOption := models.NewOption(android.ProjectLocationInputTitle, android.ProjectLocationInputEnvKey)
		for _, exportMethod := range ios.IosExportMethods {
			exportMethodOption.AddOption(exportMethod, projectLocationOption)
		}

		moduleOption = models.NewOption(android.ModuleInputTitle, android.ModuleInputEnvKey)
		projectLocationOption.AddOption("./android", moduleOption)
	} else {
		workDirOption := models.NewOption("Project root directory (the directory of the project app.json/package.json file)", "WORKDIR")
		for _, exportMethod := range ios.IosExportMethods {
			exportMethodOption.AddOption(exportMethod, workDirOption)
		}

		projectLocationOption := models.NewOption(android.ProjectLocationInputTitle, android.ProjectLocationInputEnvKey)
		workDirOption.AddOption(relPackageJSONDir, projectLocationOption)

		moduleOption = models.NewOption(android.ModuleInputTitle, android.ModuleInputEnvKey)
		projectLocationOption.AddOption(filepath.Join(relPackageJSONDir, "android"), moduleOption)
	}

	buildVariantOption := models.NewOption(android.VariantInputTitle, android.VariantInputEnvKey)
	moduleOption.AddOption("app", buildVariantOption)

	// expo options
	if scanner.usesExpoKit {
		userNameOption := models.NewOption("Expo username", "EXPO_USERNAME")
		buildVariantOption.AddOption("Release", userNameOption)

		passwordOption := models.NewOption("Expo password", "EXPO_PASSWORD")
		userNameOption.AddOption("_", passwordOption)

		configOption := models.NewConfigOption(configName)
		passwordOption.AddConfig("_", configOption)
	} else {
		configOption := models.NewConfigOption(configName)
		buildVariantOption.AddConfig("Release", configOption)
	}

	return *projectPathOption, warnings, nil
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	configMap := models.BitriseConfigMap{}

	// determine workdir
	packageJSONDir := filepath.Dir(scanner.packageJSONPth)
	relPackageJSONDir, err := utility.RelPath(scanner.searchDir, packageJSONDir)
	if err != nil {
		return models.BitriseConfigMap{}, fmt.Errorf("Failed to get relative package.json dir path, error: %s", err)
	}
	if relPackageJSONDir == "." {
		// package.json placed in the search dir, no need to change-dir in the workflows
		relPackageJSONDir = ""
	}
	log.TPrintf("Working directory: %v", relPackageJSONDir)

	workdirEnvList := []envmanModels.EnvironmentItemModel{}
	if relPackageJSONDir != "" {
		workdirEnvList = append(workdirEnvList, envmanModels.EnvironmentItemModel{reactnative.WorkDirInputKey: relPackageJSONDir})
	}

	// determine dependency manager step
	hasYarnLockFile := false
	if exist, err := pathutil.IsPathExists(filepath.Join(relPackageJSONDir, "yarn.lock")); err != nil {
		log.Warnf("Failed to check if yarn.lock file exists in the workdir: %s", err)
		log.TPrintf("Dependency manager: npm")
	} else if exist {
		log.TPrintf("Dependency manager: yarn")
		hasYarnLockFile = true
	} else {
		log.TPrintf("Dependency manager: npm")
	}

	// find test script in package.json file
	b, err := fileutil.ReadBytesFromFile(scanner.packageJSONPth)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}
	var packageJSON serialized.Object
	if err := json.Unmarshal([]byte(b), &packageJSON); err != nil {
		return models.BitriseConfigMap{}, err
	}

	hasTest := false
	if scripts, err := packageJSON.Object("scripts"); err == nil {
		if _, err := scripts.String("test"); err == nil {
			hasTest = true
		}
	}
	log.TPrintf("test script found in package.json: %v", hasTest)

	if !hasTest {
		// if the project has no test script defined,
		// we can only provide deploy like workflow,
		// so that is going to be the primary workflow

		configBuilder := models.NewDefaultConfigBuilder()
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)

		if hasYarnLockFile {
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.YarnStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))
		} else {
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))
		}

		projectDir := relPackageJSONDir
		if relPackageJSONDir == "" {
			projectDir = "./"
		}
		if scanner.usesExpoKit {
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.ExpoDetachStepListItem(
				envmanModels.EnvironmentItemModel{"project_path": projectDir},
				envmanModels.EnvironmentItemModel{"user_name": "$EXPO_USERNAME"},
				envmanModels.EnvironmentItemModel{"password": "$EXPO_PASSWORD"},
				envmanModels.EnvironmentItemModel{"run_publish": "yes"},
			))
		} else {
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.ExpoDetachStepListItem(
				envmanModels.EnvironmentItemModel{"project_path": projectDir},
			))
		}

		// android build
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
			envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.ProjectLocationInputEnvKey + "/gradlew"},
		))
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.AndroidBuildStepListItem(
			envmanModels.EnvironmentItemModel{android.ProjectLocationInputKey: "$" + android.ProjectLocationInputEnvKey},
			envmanModels.EnvironmentItemModel{android.ModuleInputKey: "$" + android.ModuleInputEnvKey},
			envmanModels.EnvironmentItemModel{android.VariantInputKey: "$" + android.VariantInputEnvKey},
		))

		// ios build
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

		if scanner.usesExpoKit {
			// in case of expo kit rn project expo eject generates an ios project with Podfile
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CocoapodsInstallStepListItem())
		}

		xcodeArchiveInputs := []envmanModels.EnvironmentItemModel{
			envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
			envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
			envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: "Release"},
			envmanModels.EnvironmentItemModel{ios.ExportMethodInputKey: "$" + ios.ExportMethodInputEnvKey},
			envmanModels.EnvironmentItemModel{"force_team_id": "$BITRISE_IOS_DEVELOPMENT_TEAM"},
		}
		if !scanner.usesExpoKit {
			// in case of plain rn project new xcode build system needs to be turned off
			xcodeArchiveInputs = append(xcodeArchiveInputs, envmanModels.EnvironmentItemModel{"xcodebuild_options": "-UseModernBuildSystem=NO"})
		}
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XcodeArchiveStepListItem(xcodeArchiveInputs...))

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)
		configBuilder.SetWorkflowDescriptionTo(models.PrimaryWorkflowID, deployWorkflowDescription)

		bitriseDataModel, err := configBuilder.Generate(Name)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(bitriseDataModel)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		configMap[configName] = string(data)

		return configMap, nil
	}

	// primary workflow
	configBuilder := models.NewDefaultConfigBuilder()
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)
	if hasYarnLockFile {
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.YarnStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.YarnStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "test"})...))
	} else {
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "test"})...))
	}
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

	// deploy workflow
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(false)...)
	if hasYarnLockFile {
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.YarnStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))
	} else {
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))
	}

	projectDir := relPackageJSONDir
	if relPackageJSONDir == "" {
		projectDir = "./"
	}
	if scanner.usesExpoKit {
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.ExpoDetachStepListItem(
			envmanModels.EnvironmentItemModel{"project_path": projectDir},
			envmanModels.EnvironmentItemModel{"user_name": "$EXPO_USERNAME"},
			envmanModels.EnvironmentItemModel{"password": "$EXPO_PASSWORD"},
			envmanModels.EnvironmentItemModel{"run_publish": "yes"},
		))
	} else {
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.ExpoDetachStepListItem(
			envmanModels.EnvironmentItemModel{"project_path": projectDir},
		))
	}

	// android build
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
		envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.ProjectLocationInputEnvKey + "/gradlew"},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
		envmanModels.EnvironmentItemModel{android.ProjectLocationInputKey: "$" + android.ProjectLocationInputEnvKey},
		envmanModels.EnvironmentItemModel{android.ModuleInputKey: "$" + android.ModuleInputEnvKey},
		envmanModels.EnvironmentItemModel{android.VariantInputKey: "$" + android.VariantInputEnvKey},
	))

	// ios build
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

	if scanner.usesExpoKit {
		// in case of expo kit rn project expo eject generates an ios project with Podfile
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CocoapodsInstallStepListItem())
	}

	xcodeArchiveInputs := []envmanModels.EnvironmentItemModel{
		envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: "Release"},
		envmanModels.EnvironmentItemModel{ios.ExportMethodInputKey: "$" + ios.ExportMethodInputEnvKey},
		envmanModels.EnvironmentItemModel{"force_team_id": "$BITRISE_IOS_DEVELOPMENT_TEAM"},
	}
	if !scanner.usesExpoKit {
		xcodeArchiveInputs = append(xcodeArchiveInputs, envmanModels.EnvironmentItemModel{"xcodebuild_options": "-UseModernBuildSystem=NO"})
	}
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(xcodeArchiveInputs...))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)
	configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployWorkflowDescription)

	bitriseDataModel, err := configBuilder.Generate(Name)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(bitriseDataModel)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	configMap[configName] = string(data)

	return configMap, nil
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionNode {
	expoKitOption := models.NewOption("Project uses Expo Kit (any js file imports expo dependency)?", "USES_EXPO_KIT")

	// with Expo Kit
	{
		// ios options
		workspacePathOption := models.NewOption("The iOS workspace path generated ny the 'expo eject' process", ios.ProjectPathInputEnvKey)
		expoKitOption.AddOption("yes", workspacePathOption)

		schemeOption := models.NewOption("The iOS scheme name generated by the 'expo eject' process", ios.SchemeInputEnvKey)
		workspacePathOption.AddOption("_", schemeOption)

		exportMethodOption := models.NewOption(ios.IosExportMethodInputTitle, ios.ExportMethodInputEnvKey)
		schemeOption.AddOption("_", exportMethodOption)

		// android options
		workDirOption := models.NewOption("Project root directory (the directory of the project app.json/package.json file)", "WORKDIR")
		for _, exportMethod := range ios.IosExportMethods {
			exportMethodOption.AddOption(exportMethod, workDirOption)
		}

		projectLocationOption := models.NewOption(android.ProjectLocationInputTitle, android.ProjectLocationInputEnvKey)
		workDirOption.AddOption("_", projectLocationOption)

		moduleOption := models.NewOption(android.ModuleInputTitle, android.ModuleInputEnvKey)
		projectLocationOption.AddOption("./android", moduleOption)

		buildVariantOption := models.NewOption(android.VariantInputTitle, android.VariantInputEnvKey)
		moduleOption.AddOption("app", buildVariantOption)

		// Expo CLI options
		userNameOption := models.NewOption("Expo username", "EXPO_USERNAME")
		buildVariantOption.AddOption("Release", userNameOption)

		passwordOption := models.NewOption("Expo password", "EXPO_PASSWORD")
		userNameOption.AddOption("_", passwordOption)

		configOption := models.NewConfigOption("react-native-expo-expo-kit-default-config")
		passwordOption.AddConfig("_", configOption)
	}

	// without Expo Kit
	{
		// ios options
		projectPathOption := models.NewOption("The iOS project path generated ny the 'expo eject' process", ios.ProjectPathInputEnvKey)
		expoKitOption.AddOption("no", projectPathOption)

		schemeOption := models.NewOption("The iOS scheme name generated by the 'expo eject' process", ios.SchemeInputEnvKey)
		projectPathOption.AddOption("_", schemeOption)

		exportMethodOption := models.NewOption(ios.IosExportMethodInputTitle, ios.ExportMethodInputEnvKey)
		schemeOption.AddOption("_", exportMethodOption)

		// android options
		workDirOption := models.NewOption("Project root directory (the directory of the project app.json/package.json file)", "WORKDIR")
		for _, exportMethod := range ios.IosExportMethods {
			exportMethodOption.AddOption(exportMethod, workDirOption)
		}

		projectLocationOption := models.NewOption(android.ProjectLocationInputTitle, android.ProjectLocationInputEnvKey)
		workDirOption.AddOption("_", projectLocationOption)

		moduleOption := models.NewOption(android.ModuleInputTitle, android.ModuleInputEnvKey)
		projectLocationOption.AddOption("./android", moduleOption)

		buildVariantOption := models.NewOption(android.VariantInputTitle, android.VariantInputEnvKey)
		moduleOption.AddOption("app", buildVariantOption)

		configOption := models.NewConfigOption("react-native-expo-plain-default-config")
		buildVariantOption.AddConfig("Release", configOption)
	}

	return *expoKitOption
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configMap := models.BitriseConfigMap{}

	// with Expo Kit
	{
		// primary workflow
		configBuilder := models.NewDefaultConfigBuilder()
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{reactnative.WorkDirInputKey: "$WORKDIR"}, envmanModels.EnvironmentItemModel{"command": "install"}))
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{reactnative.WorkDirInputKey: "$WORKDIR"}, envmanModels.EnvironmentItemModel{"command": "test"}))
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

		// deploy workflow
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(false)...)
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{reactnative.WorkDirInputKey: "$WORKDIR"}, envmanModels.EnvironmentItemModel{"command": "install"}))

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.ExpoDetachStepListItem(
			envmanModels.EnvironmentItemModel{"project_path": "$WORKDIR"},
			envmanModels.EnvironmentItemModel{"user_name": "$EXPO_USERNAME"},
			envmanModels.EnvironmentItemModel{"password": "$EXPO_PASSWORD"},
			envmanModels.EnvironmentItemModel{"run_publish": "yes"},
		))

		// android build
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
			envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.ProjectLocationInputEnvKey + "/gradlew"},
		))
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
			envmanModels.EnvironmentItemModel{android.ProjectLocationInputKey: "$" + android.ProjectLocationInputEnvKey},
			envmanModels.EnvironmentItemModel{android.ModuleInputKey: "$" + android.ModuleInputEnvKey},
			envmanModels.EnvironmentItemModel{android.VariantInputKey: "$" + android.VariantInputEnvKey},
		))

		// ios build
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CocoapodsInstallStepListItem())

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(
			envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
			envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
			envmanModels.EnvironmentItemModel{ios.ExportMethodInputKey: "$" + ios.ExportMethodInputEnvKey},
			envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: "Release"},
		))

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)
		configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployWorkflowDescription)

		bitriseDataModel, err := configBuilder.Generate(Name)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(bitriseDataModel)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		configMap["default-react-native-expo-expo-kit-config"] = string(data)
	}

	{
		// primary workflow
		configBuilder := models.NewDefaultConfigBuilder()
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{reactnative.WorkDirInputKey: "$WORKDIR"}, envmanModels.EnvironmentItemModel{"command": "install"}))
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{reactnative.WorkDirInputKey: "$WORKDIR"}, envmanModels.EnvironmentItemModel{"command": "test"}))
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

		// deploy workflow
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(false)...)
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{reactnative.WorkDirInputKey: "$WORKDIR"}, envmanModels.EnvironmentItemModel{"command": "install"}))

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.ExpoDetachStepListItem(envmanModels.EnvironmentItemModel{"project_path": "$WORKDIR"}))

		// android build
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
			envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.ProjectLocationInputEnvKey + "/gradlew"},
		))
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
			envmanModels.EnvironmentItemModel{android.ProjectLocationInputKey: "$" + android.ProjectLocationInputEnvKey},
			envmanModels.EnvironmentItemModel{android.ModuleInputKey: "$" + android.ModuleInputEnvKey},
			envmanModels.EnvironmentItemModel{android.VariantInputKey: "$" + android.VariantInputEnvKey},
		))

		// ios build
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(
			envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
			envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
			envmanModels.EnvironmentItemModel{ios.ExportMethodInputKey: "$" + ios.ExportMethodInputEnvKey},
			envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: "Release"},
		))

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)
		configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployWorkflowDescription)

		bitriseDataModel, err := configBuilder.Generate(Name)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(bitriseDataModel)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		configMap["default-react-native-expo-plain-config"] = string(data)
	}

	return configMap, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return []string{
		reactnative.Name,
		string(ios.XcodeProjectTypeIOS),
		string(ios.XcodeProjectTypeMacOS),
		android.ScannerName,
	}
}
