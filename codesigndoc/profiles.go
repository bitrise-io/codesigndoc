package codesigndoc

import (
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

// profileExportFileName creates a file name for the given profile with pattern: uuid.escaped_profile_name.[mobileprovision|provisionprofile]
func profileExportFileName(info profileutil.ProvisioningProfileInfoModel, path string) string {
	replaceRexp, err := regexp.Compile("[^A-Za-z0-9_.-]")
	if err != nil {
		log.Warnf("Invalid regex, error: %s", err)
		return ""
	}
	safeTitle := replaceRexp.ReplaceAllString(info.Name, "")
	extension := ".mobileprovision"
	if strings.HasSuffix(path, ".provisionprofile") {
		extension = ".provisionprofile"
	}

	return info.UUID + "." + safeTitle + extension
}

// filterLatestProfiles renmoves older versions of the same profile
func filterLatestProfiles(profiles []profileutil.ProvisioningProfileInfoModel) []profileutil.ProvisioningProfileInfoModel {
	profilesByBundleIDAndName := map[string][]profileutil.ProvisioningProfileInfoModel{}
	for _, profile := range profiles {
		bundleID := profile.BundleID
		name := profile.Name
		bundleIDAndName := bundleID + name
		profs, ok := profilesByBundleIDAndName[bundleIDAndName]
		if !ok {
			profs = []profileutil.ProvisioningProfileInfoModel{}
		}
		profs = append(profs, profile)
		profilesByBundleIDAndName[bundleIDAndName] = profs
	}

	filteredProfiles := []profileutil.ProvisioningProfileInfoModel{}
	for _, profiles := range profilesByBundleIDAndName {
		var latestProfile *profileutil.ProvisioningProfileInfoModel
		for _, profile := range profiles {
			if latestProfile == nil || profile.ExpirationDate.After(latestProfile.ExpirationDate) {
				latestProfile = &profile
			}
		}
		filteredProfiles = append(filteredProfiles, *latestProfile)
	}
	return filteredProfiles
}
