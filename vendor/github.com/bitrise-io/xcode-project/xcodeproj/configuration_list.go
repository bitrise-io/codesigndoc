package xcodeproj

import "github.com/bitrise-io/xcode-project/serialized"

// ConfigurationList ...
type ConfigurationList struct {
	ID                       string
	DefaultConfigurationName string
	BuildConfigurations      []BuildConfiguration
}

func parseConfigurationList(id string, objects serialized.Object) (ConfigurationList, error) {
	raw, err := objects.Object(id)
	if err != nil {
		return ConfigurationList{}, err
	}

	rawBuildConfigurations, err := raw.StringSlice("buildConfigurations")
	if err != nil {
		return ConfigurationList{}, err
	}

	var buildConfigurations []BuildConfiguration
	for _, rawID := range rawBuildConfigurations {
		buildConfiguration, err := parseBuildConfiguration(rawID, objects)
		if err != nil {
			return ConfigurationList{}, err
		}

		buildConfigurations = append(buildConfigurations, buildConfiguration)
	}

	var defaultConfigurationName string
	if aDefaultConfigurationName, err := raw.String("defaultConfigurationName"); err == nil {
		defaultConfigurationName = aDefaultConfigurationName
	}

	return ConfigurationList{
		ID: id,
		DefaultConfigurationName: defaultConfigurationName,
		BuildConfigurations:      buildConfigurations,
	}, nil
}
