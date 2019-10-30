package xcodeproj

import (
	"io/ioutil"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/xcode-project/serialized"
	"howett.net/plist"
)

// ReadPlistFile returns a parsed object representing a plist file residing at path
// and a format identifier specifying the plist file format.
// Error is returned if:
// - file at path cannot be read
// - file is not a valid plist file
// Format IDs:
// - XMLFormat      = 1
// - BinaryFormat   = 2
// - OpenStepFormat = 3
// - GNUStepFormat  = 4
func ReadPlistFile(path string) (serialized.Object, int, error) {
	codeSignEntitlementsContent, err := fileutil.ReadBytesFromFile(path)
	if err != nil {
		return nil, 0, err
	}

	var codeSignEntitlements serialized.Object
	format, err := plist.Unmarshal(codeSignEntitlementsContent, &codeSignEntitlements)
	if err != nil {
		return nil, 0, err
	}

	return codeSignEntitlements, format, nil
}

// WritePlistFile writes a parsed object representing a plist file according the
// format identifier specified.
// Error is returned if:
// - file at path cannot be written
// - format id is incorrect
// Valid format IDs:
// - XMLFormat      = 1
// - BinaryFormat   = 2
// - OpenStepFormat = 3
// - GNUStepFormat  = 4
func WritePlistFile(path string, entitlements serialized.Object, format int) error {
	marshalled, err := plist.Marshal(entitlements, format)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, marshalled, 0644)
}
