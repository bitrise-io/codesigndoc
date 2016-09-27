package common

import "github.com/bitrise-tools/codesigndoc/provprofile"

// CodeSigningIdentityInfo ...
type CodeSigningIdentityInfo struct {
	Title string
}

// CodeSigningSettings ...
type CodeSigningSettings struct {
	Identities   []CodeSigningIdentityInfo
	ProvProfiles []provprofile.ProvisioningProfileInfo
	TeamIDs      []string
	AppBundleIDs []string
}
