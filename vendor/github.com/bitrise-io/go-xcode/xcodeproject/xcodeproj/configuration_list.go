package xcodeproj

import (
	"fmt"

	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
)

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
		ID:                       id,
		DefaultConfigurationName: defaultConfigurationName,
		BuildConfigurations:      buildConfigurations,
	}, nil
}

// BuildConfigurationList ...
func (p XcodeProj) BuildConfigurationList(targetID string) (serialized.Object, error) {
	objects, err := p.RawProj.Object("objects")
	if err != nil {
		return nil, fmt.Errorf("failed to read project: %s", err)
	}

	object, err := objects.Object(targetID)
	if err != nil {
		return nil, fmt.Errorf("failed to read target (%s) object: %s", targetID, err)
	}
	buildConfigurationListID, err := object.String("buildConfigurationList")
	if err != nil {
		return nil, fmt.Errorf("failed to read target (%s) build configuration list: %s", targetID, err)
	}

	return objects.Object(buildConfigurationListID)
}

// BuildConfigurations ...
func (p XcodeProj) BuildConfigurations(buildConfigurationList serialized.Object) ([]serialized.Object, error) {
	buildConfigurationIDList, err := buildConfigurationList.StringSlice("buildConfigurations")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the buildConfigurations attributes of the the buildConfigurationList (%v), error: %s", buildConfigurationList, err)
	}

	var buildConfigurations []serialized.Object
	for _, id := range buildConfigurationIDList {
		objects, err := p.RawProj.Object("objects")
		if err != nil {
			return nil, fmt.Errorf("failed to fetch target buildConfigurations, the objects of the project are not found, error: %s", err)
		}

		buildConfiguration, err := objects.Object(id)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch target buildConfiguration objects with ID (%s), error: %s", id, err)
		}
		buildConfigurations = append(buildConfigurations, buildConfiguration)
	}
	return buildConfigurations, nil
}
