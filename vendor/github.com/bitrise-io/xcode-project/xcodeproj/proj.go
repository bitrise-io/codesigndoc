package xcodeproj

import (
	"github.com/bitrise-io/xcode-project/serialized"
)

// Proj ...
type Proj struct {
	ID                     string
	BuildConfigurationList ConfigurationList
	Targets                []Target
	Attributes             ProjectAtributes
}

func parseProj(id string, objects serialized.Object) (Proj, error) {
	rawPBXProj, err := objects.Object(id)
	if err != nil {
		return Proj{}, err
	}

	projectAttributes, err := parseProjectAttributes(rawPBXProj)
	if err != nil {
		return Proj{}, err
	}

	buildConfigurationListID, err := rawPBXProj.String("buildConfigurationList")
	if err != nil {
		return Proj{}, err
	}

	buildConfigurationList, err := parseConfigurationList(buildConfigurationListID, objects)
	if err != nil {
		return Proj{}, err
	}

	rawTargets, err := rawPBXProj.StringSlice("targets")
	if err != nil {
		return Proj{}, err
	}

	var targets []Target
	for i := range rawTargets {
		target, err := parseTarget(rawTargets[i], objects)
		if err != nil {
			return Proj{}, err
		}
		targets = append(targets, target)
	}

	return Proj{
		ID:                     id,
		BuildConfigurationList: buildConfigurationList,
		Targets:                targets,
		Attributes:             projectAttributes,
	}, nil
}

// Target ...
func (p Proj) Target(id string) (Target, bool) {
	for _, target := range p.Targets {
		if target.ID == id {
			return target, true
		}
	}
	return Target{}, false
}

// TargetByName ...
func (p Proj) TargetByName(name string) (Target, bool) {
	for _, target := range p.Targets {
		if target.Name == name {
			return target, true
		}
	}
	return Target{}, false
}
