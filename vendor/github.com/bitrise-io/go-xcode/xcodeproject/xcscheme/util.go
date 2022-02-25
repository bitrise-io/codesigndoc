package xcscheme

import (
	"path/filepath"
)

// FindSchemesIn ...
func FindSchemesIn(root string) (schemes []Scheme, err error) {
	//
	// Add the shared schemes to the list
	sharedPths, err := pathsByPattern(root, "xcshareddata", "xcschemes", "*.xcscheme")
	if err != nil {
		return nil, err
	}

	//
	// Add the non-shared user schemes to the list
	userPths, err := pathsByPattern(root, "xcuserdata", "*.xcuserdatad", "xcschemes", "*.xcscheme")
	if err != nil {
		return nil, err
	}

	for _, schemePaths := range []struct {
		schemes  []string
		isShared bool
	}{
		{sharedPths, true},
		{userPths, false},
	} {
		for _, pth := range schemePaths.schemes {
			scheme, err := Open(pth)
			if err != nil {
				return nil, err
			}

			scheme.IsShared = schemePaths.isShared
			schemes = append(schemes, scheme)
		}
	}

	return
}

func pathsByPattern(paths ...string) ([]string, error) {
	pattern := filepath.Join(paths...)
	return filepath.Glob(pattern)
}
