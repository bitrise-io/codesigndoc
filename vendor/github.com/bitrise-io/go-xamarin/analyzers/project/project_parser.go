package project

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xamarin/constants"
	"github.com/bitrise-io/go-xamarin/utility"
)

// Project is the struct for the csproj file.
type Project struct {
	XMLName        xml.Name        `xml:"Project"`
	Text           string          `xml:",chardata"`
	DefaultTargets string          `xml:"DefaultTargets,attr"`
	ToolsVersion   string          `xml:"ToolsVersion,attr"`
	Xmlns          string          `xml:"xmlns,attr"`
	PropertyGroups []PropertyGroup `xml:"PropertyGroup"`
	ItemGroups     []ItemGroup     `xml:"ItemGroup"`
	Imports        []Import        `xml:"Import"`
}

// Import the import values from the csproj file.
type Import struct {
	Text      string `xml:",chardata"`
	Project   string `xml:"Project,attr"`
	Label     string `xml:"Label,attr"`
	Condition string `xml:"Condition,attr"`
}

// PropertyGroup the property group from the csproj file.
type PropertyGroup struct {
	XMLName       xml.Name `xml:"PropertyGroup"`
	Text          string   `xml:",chardata"`
	Condition     string   `xml:"Condition,attr"`
	Configuration []struct {
		Text      string `xml:",chardata"`
		Condition string `xml:"Condition,attr"`
	} `xml:"Configuration"`
	Platform []struct {
		Text      string `xml:",chardata"`
		Condition string `xml:"Condition,attr"`
	} `xml:"Platform"`
	ProjectGUID               []string `xml:"ProjectGuid"`
	ProjectTypeGuids          []string `xml:"ProjectTypeGuids"`
	OutputType                []string `xml:"OutputType"`
	RootNamespace             []string `xml:"RootNamespace"`
	AssemblyName              []string `xml:"AssemblyName"`
	TargetFrameworkVersion    []string `xml:"TargetFrameworkVersion"`
	AndroidApplication        []string `xml:"AndroidApplication"`
	AndroidManifest           []string `xml:"AndroidManifest"`
	AndroidResgenFile         []string `xml:"AndroidResgenFile"`
	AndroidResgenClass        []string `xml:"AndroidResgenClass"`
	MonoAndroidResourcePrefix []string `xml:"MonoAndroidResourcePrefix"`
	MonoAndroidAssetsPrefix   []string `xml:"MonoAndroidAssetsPrefix"`
	DebugSymbols              []string `xml:"DebugSymbols"`
	DebugType                 []string `xml:"DebugType"`
	Optimize                  []string `xml:"Optimize"`
	OutputPath                []string `xml:"OutputPath"`
	DefineConstants           []string `xml:"DefineConstants"`
	ErrorReport               []string `xml:"ErrorReport"`
	WarningLevel              []string `xml:"WarningLevel"`
	AndroidLinkMode           []string `xml:"AndroidLinkMode"`
	AndroidManagedSymbols     []string `xml:"AndroidManagedSymbols"`
	AndroidUseSharedRuntime   []string `xml:"AndroidUseSharedRuntime"`
	MandroidI18n              []string `xml:"MandroidI18n"`
	MtouchArch                []string `xml:"MtouchArch"`
	AndroidSupportedAbis      []string `xml:"AndroidSupportedAbis"`
	BuildIpa                  []string `xml:"BuildIpa"`
	AndroidKeyStore           []string `xml:"AndroidKeyStore"`
}

// ItemGroup the item group from the csproj file.
type ItemGroup struct {
	XMLName   xml.Name `xml:"ItemGroup"`
	Text      string   `xml:",chardata"`
	Reference []struct {
		Text    string `xml:",chardata"`
		Include string `xml:"Include,attr"`
	} `xml:"Reference"`
	ProjectReferences []ProjReference `xml:"ProjectReference"`
	Compile           []struct {
		Text    string `xml:",chardata"`
		Include string `xml:"Include,attr"`
	} `xml:"Compile"`
	None []struct {
		Text    string `xml:",chardata"`
		Include string `xml:"Include,attr"`
	} `xml:"None"`
	AndroidResource []struct {
		Text    string `xml:",chardata"`
		Include string `xml:"Include,attr"`
	} `xml:"AndroidResource"`
}

// ProjReference the project reference from the csproj file.
type ProjReference struct {
	XMLName                 xml.Name `xml:"ProjectReference"`
	Text                    string   `xml:",chardata"`
	Include                 string   `xml:"Include,attr"`
	Project                 string   `xml:"Project"`
	Name                    string   `xml:"Name"`
	ReferenceOutputAssembly string   `xml:"ReferenceOutputAssembly"`
	Private                 string   `xml:"Private"`
}

const getterErrorMsg = "could not find %s"

// ParseProjectContent parses the given string content to Project struct.
func ParseProjectContent(content string) (Project, error) {
	var project Project
	if err := xml.Unmarshal([]byte(content), &project); err != nil {
		return Project{}, fmt.Errorf("failed to unmarshall conent. Error: %v", err)
	}
	return project, nil
}

// ParseProject parses the given project on path.
func ParseProject(path string) (Project, error) {
	projectDefinitionFileContent, err := fileutil.ReadStringFromFile(path)
	if err != nil {
		return Project{}, fmt.Errorf("failed to parse project at (%s), error: %s", path, err)
	}
	return ParseProjectContent(projectDefinitionFileContent)
}

// GetProjectGUID gets the guid from the given project.
func GetProjectGUID(project Project) (string, error) {
	for _, propertyGroup := range project.PropertyGroups {
		length := len(propertyGroup.ProjectGUID)
		if length > 0 {
			guid := propertyGroup.ProjectGUID[length-1]
			return strings.ToUpper(trimIDFixes(guid)), nil
		}
	}
	return "", fmt.Errorf(getterErrorMsg, "project guid")
}

// GetOutputType gets the output type from the given project.
func GetOutputType(project Project) (string, error) {
	for _, propertyGroup := range project.PropertyGroups {
		length := len(propertyGroup.OutputType)
		if length > 0 {
			return strings.ToLower(propertyGroup.OutputType[length-1]), nil
		}
	}
	return "", fmt.Errorf(getterErrorMsg, "project guid")
}

// GetAssemblyName gets the assembly name from the given project.
func GetAssemblyName(project Project) (string, error) {
	for _, propertyGroup := range project.PropertyGroups {
		length := len(propertyGroup.AssemblyName)
		if length > 0 {
			return propertyGroup.AssemblyName[length-1], nil
		}
	}
	return "", fmt.Errorf(getterErrorMsg, "assembly name")
}

// GetAndroidManifestPath gets the path for the Android manifest from the given project.
func GetAndroidManifestPath(project Project) (string, error) {
	for _, propertyGroup := range project.PropertyGroups {
		length := len(propertyGroup.AndroidManifest)
		if length > 0 {
			return propertyGroup.AndroidManifest[length-1], nil
		}
	}
	return "", fmt.Errorf(getterErrorMsg, "Android manifest path")
}

// GetResolvedAndroidManifestPath gets the resolved path for the Android manifest from the given project.
func GetResolvedAndroidManifestPath(project Project, projectDir string) (string, error) {
	relativePth, err := GetAndroidManifestPath(project)
	if err != nil {
		return "", err
	}
	relativePth = utility.FixWindowsPath(relativePth)
	return filepath.Join(projectDir, relativePth), nil
}

// GetIsAndroidApplication gets the bool value if the project is an Android project.
func GetIsAndroidApplication(project Project) (bool, error) {
	for _, propertyGroup := range project.PropertyGroups {
		length := len(propertyGroup.AndroidApplication)
		if length > 0 {
			return boolParse(propertyGroup.AndroidApplication[length-1]), nil
		}
	}
	return false, fmt.Errorf(getterErrorMsg, "Android application")
}

// GetConfiguration gets the configuration from the given property group.
func GetConfiguration(propertyGroup PropertyGroup) (string, error) {
	for _, configuration := range propertyGroup.Configuration {
		if configuration.Text != "" {
			return configuration.Text, nil
		}
	}
	return "", fmt.Errorf(getterErrorMsg, "configuration")
}

// GetPropertyGroupCondition gets the condition of the given property group.
func GetPropertyGroupCondition(propertyGroup PropertyGroup) (string, error) {
	if propertyGroup.Condition == "" {
		return "", fmt.Errorf(getterErrorMsg, "condition")
	}
	return propertyGroup.Condition, nil
}

// GetResolvedConfiguration gets the resolved configuration from the given property group.
func GetResolvedConfiguration(propertyGroup PropertyGroup) (string, error) {
	conditionText, err := GetPropertyGroupCondition(propertyGroup)
	if err != nil {
		return "", err
	}
	if matches := regexp.MustCompile(`'\$\(Configuration\)\|\$\(Platform\)'\s*==\s*'(?P<config>.*)\|(?P<platform>.*)'`).FindStringSubmatch(conditionText); len(matches) == 3 {
		return matches[1], nil
	}

	if matches := regexp.MustCompile(`'\$\(Configuration\)'\s*==\s*'(?P<config>.*)'`).FindStringSubmatch(conditionText); len(matches) == 2 {
		return matches[1], nil
	}

	return conditionText, nil
}

// GetResolvedPlatform gets the resolved platform from the given property group.
func GetResolvedPlatform(propertyGroup PropertyGroup) (string, error) {
	conditionText, err := GetPropertyGroupCondition(propertyGroup)
	if err != nil {
		return "", err
	}
	if matches := regexp.MustCompile(`'\$\(Configuration\)\|\$\(Platform\)'\s*==\s*'(?P<config>.*)\|(?P<platform>.*)'`).FindStringSubmatch(conditionText); len(matches) == 3 {
		return matches[2], nil
	}

	if matches := regexp.MustCompile(`'\$\(Platform\)'\s*==\s*'(?P<platform>.*)'`).FindStringSubmatch(conditionText); len(matches) == 2 {
		return matches[1], nil
	}

	return conditionText, nil
}

// GetPlatform gets the platform from the given property group.
func GetPlatform(propertyGroup PropertyGroup) (string, error) {
	for _, platform := range propertyGroup.Platform {
		if platform.Text != "" {
			return platform.Text, nil
		}
	}
	return "", fmt.Errorf(getterErrorMsg, "platform")
}

// GetOutputPath gets the output dir from the given property group.
func GetOutputPath(propertyGroup PropertyGroup) (string, error) {
	length := len(propertyGroup.OutputPath)
	if length > 0 {
		return propertyGroup.OutputPath[length-1], nil
	}
	return "", fmt.Errorf(getterErrorMsg, "output path")
}

// GetOutputDir gets the output dir from the property group.
func GetOutputDir(propertyGroup PropertyGroup, projectDir, configuration, platform string) (string, error) {
	relativePth, err := GetOutputPath(propertyGroup)
	if err != nil {
		return "", err
	}
	relativePth = utility.FixWindowsPath(relativePth)
	strings.Replace(relativePth, "$(Configuration)", configuration, -1)
	strings.Replace(relativePth, "$(Platform)", platform, -1)
	return filepath.Join(projectDir, relativePth), nil
}

// GetMtouchArch gets the MtouchArch from the given property group.
func GetMtouchArch(propertyGroup PropertyGroup) (string, error) {
	length := len(propertyGroup.MtouchArch)
	if length > 0 {
		return propertyGroup.MtouchArch[length-1], nil
	}
	return "", fmt.Errorf(getterErrorMsg, "mtouchArch")
}

// GetResolvedMtouchArch gets the resolved MtouchArch from the given property group.
func GetResolvedMtouchArch(propertyGroup PropertyGroup) ([]string, error) {
	mTouchArchs, err := GetMtouchArch(propertyGroup)
	if err != nil {
		return []string{}, err
	}
	return utility.SplitAndStripList(mTouchArchs, ","), nil
}

// GetBuildIpa gets the build IPA boolean from the given property group.
func GetBuildIpa(propertyGroup PropertyGroup) (bool, error) {
	length := len(propertyGroup.BuildIpa)
	if length > 0 {
		return boolParse(propertyGroup.BuildIpa[length-1]), nil
	}
	return false, fmt.Errorf(getterErrorMsg, "build IPA")
}

// GetAndroidKeyStore gets the Android keystore boolean from the given property group.
func GetAndroidKeyStore(propertyGroup PropertyGroup) (bool, error) {
	length := len(propertyGroup.AndroidKeyStore)
	if length > 0 {
		return boolParse(propertyGroup.AndroidKeyStore[length-1]), nil
	}
	return false, fmt.Errorf(getterErrorMsg, "Android keystore")
}

// GetProjectTypeGUIDs gets the project type GUIDs from the given project.
func GetProjectTypeGUIDs(project Project) (string, error) {
	for _, propertyGroup := range project.PropertyGroups {
		length := len(propertyGroup.ProjectTypeGuids)
		if length > 0 {
			return propertyGroup.ProjectTypeGuids[length-1], nil
		}
	}
	return "", fmt.Errorf(getterErrorMsg, "project GUIDs")
}

// GetResolvedProjectTypeGUIDs gets the project type GUIDs from the given project.
func GetResolvedProjectTypeGUIDs(project Project) (constants.SDK, error) {
	sdk := constants.SDKUnknown
	guidsLine, err := GetProjectTypeGUIDs(project)
	if err != nil {
		return sdk, err
	}
	projectTypeList := strings.Split(guidsLine, ";")
	for _, guid := range projectTypeList {
		guid = trimIDFixes(guid)

		sdk, err = constants.ParseProjectTypeGUID(guid)
		if err == nil {
			break
		}
	}
	return sdk, nil
}

// GetItemGroupIncludes gets the includes from the given item group.
func GetItemGroupIncludes(itemGroup ItemGroup) []string {
	var includes []string
	for _, reference := range itemGroup.Reference {
		if reference.Include != "" {
			includes = append(includes, reference.Include)
		}
	}
	return includes
}

// GetTestFramework gets the test framework for the given project.
func GetTestFramework(project Project) (constants.TestFramework, error) {
	var testFramework constants.TestFramework
	for _, itemGroup := range project.ItemGroups {
		includes := GetItemGroupIncludes(itemGroup)
		for _, include := range includes {
			switch include {
			case "Xamarin.UITest":
				return constants.TestFrameworkXamarinUITest, nil
			case "MonoTouch.NUnitLite":
				return constants.TestFrameworkNunitLiteTest, nil
			case "nunit", "nunit.framework":
				testFramework = constants.TestFrameworkNunitTest
			}
		}
	}
	if testFramework == "" {
		return constants.TestFrameworkUnknown, fmt.Errorf(getterErrorMsg, "testframework")
	}
	return testFramework, nil
}

// GetProjectReferences get the project references from the given project.
func GetProjectReferences(project Project) []ProjReference {
	var projectReferences []ProjReference
	for _, itemGroup := range project.ItemGroups {
		projectReferences = append(projectReferences, itemGroup.ProjectReferences...)
	}
	return projectReferences
}

// GetReferencedProjectIds gets the referenced project IDs.
func GetReferencedProjectIds(project Project) []string {
	projectReferences := GetProjectReferences(project)

	var projectIds []string
	for _, projectReference := range projectReferences {
		id := strings.ToUpper(projectReference.Project)
		id = trimIDFixes(id)
		projectIds = append(projectIds, id)
	}
	return projectIds
}

// GetImportedProjects gets the imported projects from a given project.
func GetImportedProjects(project Project) []string {
	var importedProjects []string
	for _, importItem := range project.Imports {
		importedProjects = append(importedProjects, utility.FixWindowsPath(importItem.Project))
	}
	return importedProjects
}

// GetPropertyGroupsConfiguration gets the configuration for each property group
func GetPropertyGroupsConfiguration(project Project, projectDir string, sdk constants.SDK) ([]ConfigurationPlatformModel, error) {
	var configModels []ConfigurationPlatformModel
	for _, propertyGroup := range project.PropertyGroups {
		var configModel ConfigurationPlatformModel
		var err error

		configModel.Configuration, err = GetResolvedConfiguration(propertyGroup)
		if err != nil {
			debugParseLog(err)
		}

		configModel.Platform, err = GetResolvedPlatform(propertyGroup)
		if err != nil {
			debugParseLog(err)
		}

		configModel.OutputDir, err = GetOutputDir(propertyGroup, projectDir, configModel.Configuration, configModel.Platform)
		if err != nil {
			debugParseLog(err)
		}
		if sdk == constants.SDKIOS || sdk == constants.SDKMacOS || sdk == constants.SDKTvOS {
			configModel.MtouchArchs, err = GetResolvedMtouchArch(propertyGroup)
			if err != nil {
				debugParseLog(err)
			}

			configModel.BuildIpa, err = GetBuildIpa(propertyGroup)
			if err != nil {
				debugParseLog(err)
			}
		}

		if sdk == constants.SDKAndroid {
			configModel.SignAndroid, err = GetAndroidKeyStore(propertyGroup)
			if err != nil {
				debugParseLog(err)
			}
		}

		configModels = append(configModels, configModel)
	}
	return configModels, nil
}

func boolParse(value string) bool {
	return strings.EqualFold(value, "true")
}

func trimIDFixes(id string) string {
	id = strings.TrimPrefix(id, "{")
	id = strings.TrimSuffix(id, "}")
	return id
}

func debugParseLog(err error) {
	log.Debugf("%v", err)
}
