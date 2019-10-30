package xcodeproj

import (
	"fmt"

	"github.com/bitrise-io/xcode-project/serialized"
)

// ProjectAtributes ...
//
// **Deprecated**: use the func (p XcodeProj) Attributes() (serialized.Object, error) method instead
type ProjectAtributes struct {
	TargetAttributes serialized.Object
}

//
// **Deprecated**: use the func (p XcodeProj) Attributes() (serialized.Object, error) method instead
func parseProjectAttributes(rawPBXProj serialized.Object) (ProjectAtributes, error) {
	var attributes ProjectAtributes
	attributesObject, err := rawPBXProj.Object("attributes")
	if err != nil {
		return ProjectAtributes{}, err
	}

	attributes.TargetAttributes, err = parseTargetAttributes(attributesObject)
	if err != nil && !serialized.IsKeyNotFoundError(err) {
		return ProjectAtributes{}, err
	}

	return attributes, nil
}

//
// **Deprecated**: use the func (p XcodeProj) TargetAttributes() (serialized.Object, error) method instead
func parseTargetAttributes(attributesObject serialized.Object) (serialized.Object, error) {
	return attributesObject.Object("TargetAttributes")
}

// Attributes ...
func (p XcodeProj) Attributes() (serialized.Object, error) {
	objects, err := p.RawProj.Object("objects")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project attributes, the objects of the project are not found, error: %s", err)
	}

	object, err := objects.Object(p.Proj.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project attributes, the project objects wit ID (%s) is not found, error: %s", p.Proj.ID, err)
	}

	return object.Object("attributes")
}

// TargetAttributes ...
func (p XcodeProj) TargetAttributes() (serialized.Object, error) {
	attributes, err := p.Attributes()
	if err != nil {
		return nil, err
	}
	return attributes.Object("TargetAttributes")
}
