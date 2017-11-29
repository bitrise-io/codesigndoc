package export

import (
	"sort"

	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
	glob "github.com/ryanuber/go-glob"
)

// IosCodeSignGroup ...
type IosCodeSignGroup struct {
	Certificate        certificateutil.CertificateInfoModel
	BundleIDProfileMap map[string]profileutil.ProvisioningProfileInfoModel
}

func createSingleWildcardGroups(group SelectableCodeSignGroup, alreadyUsedProfileUUIDMap map[string]bool) []IosCodeSignGroup {
	groups := []IosCodeSignGroup{}

	certificate := group.Certificate
	bundleIDProfilesMap := group.BundleIDProfilesMap

	bundleIDs := []string{}
	profiles := []profileutil.ProvisioningProfileInfoModel{}
	for bundleID, matchingProfiles := range bundleIDProfilesMap {
		bundleIDs = append(bundleIDs, bundleID)
		profiles = append(profiles, matchingProfiles...)
	}

	for _, profile := range profiles {
		if alreadyUsedProfileUUIDMap[profile.UUID] {
			continue
		}

		matchesForAllBundleID := true
		for _, bundleID := range bundleIDs {
			if !glob.Glob(profile.BundleID, bundleID) {
				matchesForAllBundleID = false
				break
			}
		}
		if matchesForAllBundleID {
			bundleIDProfileMap := map[string]profileutil.ProvisioningProfileInfoModel{}
			for _, bundleID := range bundleIDs {
				bundleIDProfileMap[bundleID] = profile
			}

			group := IosCodeSignGroup{
				Certificate:        certificate,
				BundleIDProfileMap: bundleIDProfileMap,
			}
			groups = append(groups, group)

			alreadyUsedProfileUUIDMap[profile.UUID] = true
		}
	}
	return groups
}

func createXcodeManagedGroups(group SelectableCodeSignGroup, alreadyUsedProfileUUIDMap map[string]bool) []IosCodeSignGroup {
	groups := []IosCodeSignGroup{}

	certificate := group.Certificate
	bundleIDProfilesMap := group.BundleIDProfilesMap

	bundleIDs := []string{}
	profiles := []profileutil.ProvisioningProfileInfoModel{}
	for bundleID, matchingProfiles := range bundleIDProfilesMap {
		bundleIDs = append(bundleIDs, bundleID)
		profiles = append(profiles, matchingProfiles...)
	}

	// collect xcode managed profiles
	xcodeManagedProfiles := []profileutil.ProvisioningProfileInfoModel{}
	for _, profile := range profiles {
		if !alreadyUsedProfileUUIDMap[profile.UUID] && profile.IsXcodeManaged() {
			xcodeManagedProfiles = append(xcodeManagedProfiles, profile)
		}
	}
	sort.Sort(ByBundleIDLength(xcodeManagedProfiles))

	// map profiles to bundle ids + remove the already used profiles
	bundleIDMannagedProfilesMap := map[string][]profileutil.ProvisioningProfileInfoModel{}
	for _, bundleID := range bundleIDs {
		for _, profile := range xcodeManagedProfiles {
			if !glob.Glob(profile.BundleID, bundleID) {
				continue
			}

			matchingProfiles := bundleIDMannagedProfilesMap[bundleID]
			if matchingProfiles == nil {
				matchingProfiles = []profileutil.ProvisioningProfileInfoModel{}
			}
			matchingProfiles = append(matchingProfiles, profile)
			bundleIDMannagedProfilesMap[bundleID] = matchingProfiles
		}
	}

	if len(bundleIDMannagedProfilesMap) == len(bundleIDs) {
		// if only one profile can sign a bundle id, remove it from bundleIDMannagedProfilesMap
		alreadyUsedManagedProfileMap := map[string]bool{}
		for _, profiles := range bundleIDMannagedProfilesMap {
			if len(profiles) == 1 {
				profile := profiles[0]
				alreadyUsedManagedProfileMap[profile.UUID] = true
			}
		}

		bundleIDMannagedProfileMap := map[string]profileutil.ProvisioningProfileInfoModel{}
		for bundleID, profiles := range bundleIDMannagedProfilesMap {
			if len(profiles) == 1 {
				bundleIDMannagedProfileMap[bundleID] = profiles[0]
			} else {
				remainingProfiles := []profileutil.ProvisioningProfileInfoModel{}
				for _, profile := range profiles {
					if !alreadyUsedManagedProfileMap[profile.UUID] {
						remainingProfiles = append(remainingProfiles, profile)
					}
				}
				if len(remainingProfiles) == 1 {
					bundleIDMannagedProfileMap[bundleID] = remainingProfiles[0]
				}
			}
		}

		// create code sign group
		if len(bundleIDMannagedProfileMap) == len(bundleIDs) {
			for _, profile := range bundleIDMannagedProfileMap {
				alreadyUsedProfileUUIDMap[profile.UUID] = true
			}

			group := IosCodeSignGroup{
				Certificate:        certificate,
				BundleIDProfileMap: bundleIDMannagedProfileMap,
			}
			groups = append(groups, group)
		}
	}

	return groups
}

func createNotXcodeManagedGroups(group SelectableCodeSignGroup, alreadyUsedProfileUUIDMap map[string]bool) []IosCodeSignGroup {
	groups := []IosCodeSignGroup{}

	certificate := group.Certificate
	bundleIDProfilesMap := group.BundleIDProfilesMap

	bundleIDs := []string{}
	profiles := []profileutil.ProvisioningProfileInfoModel{}
	for bundleID, matchingProfiles := range bundleIDProfilesMap {
		bundleIDs = append(bundleIDs, bundleID)
		profiles = append(profiles, matchingProfiles...)
	}

	// collect xcode managed profiles
	notXcodeManagedProfiles := []profileutil.ProvisioningProfileInfoModel{}
	for _, profile := range profiles {
		if !alreadyUsedProfileUUIDMap[profile.UUID] && !profile.IsXcodeManaged() {
			notXcodeManagedProfiles = append(notXcodeManagedProfiles, profile)
		}
	}
	sort.Sort(ByBundleIDLength(notXcodeManagedProfiles))

	// map profiles to bundle ids + remove the already used profiles
	bundleIDNotMannagedProfilesMap := map[string][]profileutil.ProvisioningProfileInfoModel{}
	for _, bundleID := range bundleIDs {
		for _, profile := range notXcodeManagedProfiles {
			if !glob.Glob(profile.BundleID, bundleID) {
				continue
			}

			matchingProfiles := bundleIDNotMannagedProfilesMap[bundleID]
			if matchingProfiles == nil {
				matchingProfiles = []profileutil.ProvisioningProfileInfoModel{}
			}
			matchingProfiles = append(matchingProfiles, profile)
			bundleIDNotMannagedProfilesMap[bundleID] = matchingProfiles
		}
	}

	if len(bundleIDNotMannagedProfilesMap) == len(bundleIDs) {
		// if only one profile can sign a bundle id, remove it from bundleIDNotMannagedProfilesMap
		alreadyUsedNotManagedProfileMap := map[string]bool{}
		for _, profiles := range bundleIDNotMannagedProfilesMap {
			if len(profiles) == 1 {
				profile := profiles[0]
				alreadyUsedNotManagedProfileMap[profile.UUID] = true
			}
		}

		bundleIDNotMannagedProfileMap := map[string]profileutil.ProvisioningProfileInfoModel{}
		for bundleID, profiles := range bundleIDNotMannagedProfilesMap {
			if len(profiles) == 1 {
				bundleIDNotMannagedProfileMap[bundleID] = profiles[0]
			} else {
				remainingProfiles := []profileutil.ProvisioningProfileInfoModel{}
				for _, profile := range profiles {
					if !alreadyUsedNotManagedProfileMap[profile.UUID] {
						remainingProfiles = append(remainingProfiles, profile)
					}
				}
				if len(remainingProfiles) == 1 {
					bundleIDNotMannagedProfileMap[bundleID] = remainingProfiles[0]
				}
			}
		}

		// create code sign group
		if len(bundleIDNotMannagedProfileMap) == len(bundleIDs) {
			for _, profile := range bundleIDNotMannagedProfileMap {
				alreadyUsedProfileUUIDMap[profile.UUID] = true
			}

			codeSignGroup := IosCodeSignGroup{
				Certificate:        certificate,
				BundleIDProfileMap: bundleIDNotMannagedProfileMap,
			}
			groups = append(groups, codeSignGroup)
		}
	}

	return groups
}

func createRemainingGroups(group SelectableCodeSignGroup, alreadyUsedProfileUUIDMap map[string]bool) []IosCodeSignGroup {
	groups := []IosCodeSignGroup{}

	certificate := group.Certificate
	bundleIDProfilesMap := group.BundleIDProfilesMap

	bundleIDs := []string{}
	profiles := []profileutil.ProvisioningProfileInfoModel{}
	for bundleID, matchingProfiles := range bundleIDProfilesMap {
		bundleIDs = append(bundleIDs, bundleID)
		profiles = append(profiles, matchingProfiles...)
	}

	if len(alreadyUsedProfileUUIDMap) != len(profiles) {
		bundleIDProfileMap := map[string]profileutil.ProvisioningProfileInfoModel{}
		for _, bundleID := range bundleIDs {
			for _, profile := range profiles {
				if alreadyUsedProfileUUIDMap[profile.UUID] {
					continue
				}

				if !glob.Glob(profile.BundleID, bundleID) {
					continue
				}

				bundleIDProfileMap[bundleID] = profile
				break
			}
		}

		if len(bundleIDProfileMap) == len(bundleIDs) {
			group := IosCodeSignGroup{
				Certificate:        certificate,
				BundleIDProfileMap: bundleIDProfileMap,
			}
			groups = append(groups, group)
		}
	}

	return groups
}

// CreateIosCodeSignGroups ...
func CreateIosCodeSignGroups(selectableGroups []SelectableCodeSignGroup) []IosCodeSignGroup {
	alreadyUsedProfileUUIDMap := map[string]bool{}

	singleWildcardGroups := []IosCodeSignGroup{}
	xcodeManagedGroups := []IosCodeSignGroup{}
	notXcodeManagedGroups := []IosCodeSignGroup{}
	remainingGroups := []IosCodeSignGroup{}

	for _, selectableGroup := range selectableGroups {
		// create groups with single wildcard profiles
		singleWildcardGroups = append(singleWildcardGroups, createSingleWildcardGroups(selectableGroup, alreadyUsedProfileUUIDMap)...)

		// create groups with xcode managed profiles
		xcodeManagedGroups = append(xcodeManagedGroups, createXcodeManagedGroups(selectableGroup, alreadyUsedProfileUUIDMap)...)

		// create groups with NOT xcode managed profiles
		notXcodeManagedGroups = append(notXcodeManagedGroups, createNotXcodeManagedGroups(selectableGroup, alreadyUsedProfileUUIDMap)...)

		// if there are remaining profiles we create a not exact group by using the first matching profile for every bundle id
		remainingGroups = append(remainingGroups, createRemainingGroups(selectableGroup, alreadyUsedProfileUUIDMap)...)
	}

	codeSignGroups := []IosCodeSignGroup{}
	codeSignGroups = append(codeSignGroups, notXcodeManagedGroups...)
	codeSignGroups = append(codeSignGroups, xcodeManagedGroups...)
	codeSignGroups = append(codeSignGroups, singleWildcardGroups...)
	codeSignGroups = append(codeSignGroups, remainingGroups...)

	return codeSignGroups
}
