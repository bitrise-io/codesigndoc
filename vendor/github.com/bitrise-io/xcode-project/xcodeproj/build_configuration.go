package xcodeproj

import "github.com/bitrise-io/xcode-project/serialized"

// BuildConfiguration ..
type BuildConfiguration struct {
	ID            string
	Name          string
	BuildSettings serialized.Object
}

func parseBuildConfiguration(id string, objects serialized.Object) (BuildConfiguration, error) {
	raw, err := objects.Object(id)
	if err != nil {
		return BuildConfiguration{}, err
	}

	name, err := raw.String("name")
	if err != nil {
		return BuildConfiguration{}, err
	}

	buildSettings, err := raw.Object("buildSettings")
	if err != nil {
		return BuildConfiguration{}, err
	}

	return BuildConfiguration{
		ID:            id,
		Name:          name,
		BuildSettings: buildSettings,
	}, nil
}
