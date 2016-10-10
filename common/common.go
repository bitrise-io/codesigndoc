package common

import (
	"regexp"

	"github.com/bitrise-io/go-utils/regexputil"
	"github.com/bitrise-tools/codesigndoc/provprofile"
)

// CodeSigningIdentityInfo ...
type CodeSigningIdentityInfo struct {
	Title string
}

// CodeSigningSettings ...
type CodeSigningSettings struct {
	Identities   []CodeSigningIdentityInfo
	ProvProfiles []provprofile.ProvisioningProfileInfo
	TeamIDs      []string
	// Full AppIDs, in the form: TEAMID.BUNDLEID
	AppIDs []string
}

// BundleIDFromAppID ...
// if the App ID is not in the form "TEAMID.BUNDLE.ID" then an empty string will be returned
func BundleIDFromAppID(appID string) string {
	rexp := regexp.MustCompile(`^(?P<teamID>[a-zA-Z0-9]+)\.(?P<bundleID>.+)$`)
	results, isFound := regexputil.NamedFindStringSubmatch(rexp, appID)
	if isFound && results["bundleID"] != "" && results["bundleID"] != "*" {
		return results["bundleID"]
	}
	return ""
}
