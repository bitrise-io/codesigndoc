package cmd

import "github.com/bitrise-tools/go-xcode/plistutil"

// Archive ...
type Archive interface {
	BundleIDEntitlementsMap() map[string]plistutil.PlistData
	IsXcodeManaged() bool
}
