package xcodeproj

import (
	"fmt"

	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
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
		return Proj{}, fmt.Errorf("failed to access object with id %s: %s", id, err)
	}

	projectAttributes, err := parseProjectAttributes(rawPBXProj)
	if err != nil {
		return Proj{}, fmt.Errorf("failed to parse project attributes: %s", err)
	}

	buildConfigurationListID, err := rawPBXProj.String("buildConfigurationList")
	if err != nil {
		return Proj{}, fmt.Errorf("failed to access build configuration list: %s", err)
	}

	buildConfigurationList, err := parseConfigurationList(buildConfigurationListID, objects)
	if err != nil {
		return Proj{}, fmt.Errorf("failed to parse build configuration list: %s", err)
	}

	rawTargets, err := rawPBXProj.StringSlice("targets")
	if err != nil {
		return Proj{}, fmt.Errorf("failed to access targets: %s", err)
	}

	var targets []Target
	for _, targetID := range rawTargets {
		// rawTargets can contain more target IDs than the project configuration has
		hasTargetNode, err := hasTargetNode(targetID, objects)
		if err != nil {
			return Proj{}, fmt.Errorf("failed to access target object with id %s: %s", targetID, err)
		}

		if !hasTargetNode {
			continue
		}

		target, err := parseTarget(targetID, objects)
		if err != nil {
			return Proj{}, fmt.Errorf("failed to parse target with id: %s: %s", targetID, err)
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

func hasTargetNode(id string, objects serialized.Object) (bool, error) {
	if _, err := objects.Object(id); err != nil {
		if serialized.IsKeyNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
