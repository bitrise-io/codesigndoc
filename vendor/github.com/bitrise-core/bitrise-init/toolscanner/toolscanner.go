package toolscanner

import (
	"github.com/bitrise-core/bitrise-init/models"
	bitriseModels "github.com/bitrise-io/bitrise/models"
)

// ProjectTypeEnvKey is the name of the enviroment variable used to substitute the project type for
// automation tool scanner's config
const (
	ProjectTypeUserTitle = "Project type"
	// The key is used in the options decision tree model.
	// If empty, it will not be inserted into the bitrise.yml
	ProjectTypeEnvKey = ""
)

// AddProjectTypeToConfig returns the config filled in with every detected project type, that could be selected
func AddProjectTypeToConfig(configName string, config bitriseModels.BitriseDataModel, detectedProjectTypes []string) map[string]bitriseModels.BitriseDataModel {
	configMapWithProjecTypes := map[string]bitriseModels.BitriseDataModel{}
	for _, projectType := range detectedProjectTypes {
		configWithProjectType := config
		configWithProjectType.ProjectType = projectType
		configMapWithProjecTypes[appendProjectTypeToConfigName(configName, projectType)] = configWithProjectType
	}
	return configMapWithProjecTypes
}

// AddProjectTypeToOptions adds a project type question to automation tool scanners's option tree
func AddProjectTypeToOptions(scannerOptionTree models.OptionNode, detectedProjectTypes []string) models.OptionNode {
	optionsTreeWithProjectTypeRoot := models.NewOption(ProjectTypeUserTitle, ProjectTypeEnvKey)
	for _, projectType := range detectedProjectTypes {
		optionsTreeWithProjectTypeRoot.AddOption(projectType,
			appendProjectTypeToConfig(scannerOptionTree, projectType))
	}
	return *optionsTreeWithProjectTypeRoot
}

func appendProjectTypeToConfigName(configName string, projectType string) string {
	return configName + "_" + projectType
}

func appendProjectTypeToConfig(options models.OptionNode, projectType string) *models.OptionNode {
	var appendToConfigNames func(*models.OptionNode)
	appendToConfigNames = func(node *models.OptionNode) {
		if (*node).IsConfigOption() || (*node).ChildOptionMap == nil {
			(*node).Config = appendProjectTypeToConfigName((*node).Config, projectType)
			return
		}
		for _, child := range (*node).ChildOptionMap {
			appendToConfigNames(child)
		}
	}
	optionsWithProjectType := options.Copy()
	appendToConfigNames(optionsWithProjectType)
	return optionsWithProjectType
}
