package project

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xamarin/constants"
	"github.com/bitrise-io/go-xamarin/utility"
)

// ConfigurationPlatformModel ...
type ConfigurationPlatformModel struct {
	Configuration string
	Platform      string
	OutputDir     string

	MtouchArchs []string
	BuildIpa    bool

	SignAndroid bool
}

// Model ...
type Model struct {
	Pth  string
	Name string // Set by solution analyze or set its path's filename without extension

	// Solution Configuration|Platform - Project Configuration|Platform map
	// !!! only set by solution analyze
	ConfigMap map[string]string

	ID            string
	SDK           constants.SDK
	TestFramework constants.TestFramework
	OutputType    string
	AssemblyName  string

	ReferredProjectIDs []string

	ManifestPth        string
	AndroidApplication bool

	Configs map[string]ConfigurationPlatformModel // Project Configuration|Platform - ConfigurationPlatformModel map
}

// New ...
func New(pth string) (Model, error) {
	return analyzeProject(pth)
}

func debugLog(err error, pth string) {
	log.Debugf("%v for project at %s", err, pth)
}

func analyzeTargetDefinition(projectModel Model, pth string) (Model, error) {
	projectDir := filepath.Dir(pth)
	var err error

	parsedProject, err := ParseProject(pth)
	if err != nil {
		return Model{}, err
	}

	for _, importedProject := range GetImportedProjects(parsedProject) {
		if !strings.Contains(importedProject, "$(MSBuild") {
			targetDefinitionPth := filepath.Join(projectDir, importedProject)

			if exist, err := pathutil.IsPathExists(targetDefinitionPth); err != nil {
				return Model{}, err
			} else if exist {
				projectFromTargetDefinition, err := analyzeTargetDefinition(projectModel, targetDefinitionPth)
				if err != nil {
					return Model{}, err
				}

				// Set properties became from solution analyze
				projectFromTargetDefinition.Name = projectModel.Name
				projectFromTargetDefinition.Pth = projectModel.Pth
				projectFromTargetDefinition.ConfigMap = projectModel.ConfigMap
				// ---

				projectModel = projectFromTargetDefinition
			}
		}
	}

	projectModel.ID, err = GetProjectGUID(parsedProject)
	if err != nil {
		debugLog(err, pth)
	}

	projectModel.OutputType, err = GetOutputType(parsedProject)
	if err != nil {
		debugLog(err, pth)
	}

	projectModel.AssemblyName, err = GetAssemblyName(parsedProject)
	if err != nil {
		debugLog(err, pth)
	}

	projectModel.TestFramework, err = GetTestFramework(parsedProject)
	if err != nil {
		debugLog(err, pth)
	}

	projectModel.SDK, err = GetResolvedProjectTypeGUIDs(parsedProject)
	if err != nil {
		debugLog(err, pth)
	}

	if projectModel.SDK == constants.SDKAndroid {
		projectModel.ManifestPth, err = GetResolvedAndroidManifestPath(parsedProject, projectDir)
		if err != nil {
			debugLog(err, pth)
		}

		projectModel.AndroidApplication, err = GetIsAndroidApplication(parsedProject)
		if err != nil {
			debugLog(err, pth)
		}
	}

	projectModel.ReferredProjectIDs = GetReferencedProjectIds(parsedProject)

	configPlatforms, err := GetPropertyGroupsConfiguration(parsedProject, projectDir, projectModel.SDK)
	if err != nil {
		debugLog(err, pth)
	}

	for _, configPlatform := range configPlatforms {
		projectModel.Configs[utility.ToConfig(configPlatform.Configuration, configPlatform.Platform)] = configPlatform
	}

	return projectModel, nil
}

func analyzeProject(pth string) (Model, error) {
	absPth, err := pathutil.AbsPath(pth)
	if err != nil {
		return Model{}, fmt.Errorf("failed to expand path (%s), error: %s", pth, err)
	}

	fileName := filepath.Base(absPth)
	ext := filepath.Ext(absPth)
	fileName = strings.TrimSuffix(fileName, ext)

	project := Model{
		Pth:           absPth,
		Name:          fileName,
		ConfigMap:     map[string]string{},
		Configs:       map[string]ConfigurationPlatformModel{},
		SDK:           constants.SDKUnknown,
		TestFramework: constants.TestFrameworkUnknown,
	}
	return analyzeTargetDefinition(project, absPth)
}
