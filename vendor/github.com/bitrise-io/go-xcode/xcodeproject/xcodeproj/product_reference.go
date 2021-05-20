package xcodeproj

import "github.com/bitrise-io/go-xcode/xcodeproject/serialized"

// ProductReference ...
type ProductReference struct {
	Path string
}

func parseProductReference(id string, objects serialized.Object) (ProductReference, error) {
	raw, err := objects.Object(id)
	if err != nil {
		return ProductReference{}, err
	}

	pth, err := raw.String("path")
	if err != nil {
		return ProductReference{}, err
	}

	return ProductReference{
		Path: pth,
	}, nil
}
