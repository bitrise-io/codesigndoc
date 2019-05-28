package codesign

import (
	"github.com/bitrise-io/go-xcode/profileutil"
)

// FilterLatestProfiles renmoves older versions of the same profile
func FilterLatestProfiles(profiles []profileutil.ProvisioningProfileInfoModel) []profileutil.ProvisioningProfileInfoModel {
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
