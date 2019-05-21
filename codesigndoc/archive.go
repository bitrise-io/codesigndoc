package codesigndoc

import (
	"github.com/bitrise-io/go-xcode/plistutil"
	"github.com/bitrise-io/go-xcode/profileutil"
)

// Archive ...
type Archive interface {
	BundleIDEntitlementsMap() map[string]plistutil.PlistData
	IsXcodeManaged() bool
	SigningIdentity() string
	BundleIDProfileInfoMap() map[string]profileutil.ProvisioningProfileInfoModel
}
