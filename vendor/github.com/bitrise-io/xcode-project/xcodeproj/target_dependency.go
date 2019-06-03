package xcodeproj

import "github.com/bitrise-io/xcode-project/serialized"

// TargetDependency ...
type TargetDependency struct {
	ID     string
	Target Target
}

func parseTargetDependency(id string, objects serialized.Object) (TargetDependency, error) {
	rawTargetDependency, err := objects.Object(id)
	if err != nil {
		return TargetDependency{}, err
	}

	targetID, err := rawTargetDependency.String("target")
	if err != nil {
		return TargetDependency{}, err
	}

	target, err := parseTarget(targetID, objects)
	if err != nil {
		return TargetDependency{}, err
	}

	return TargetDependency{
		ID:     id,
		Target: target,
	}, nil
}
