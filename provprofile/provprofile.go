package provprofile

import (
	"fmt"
	"path/filepath"
	"strings"

	plist "github.com/DHowett/go-plist"
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/pathutil"
)

const (
	provProfileSystemDirPath = "~/Library/MobileDevice/Provisioning Profiles"
)

// ProvisioningProfileInfo ...
type ProvisioningProfileInfo struct {
	Title string
	UUID  string
}

// EntitlementsModel ...
type EntitlementsModel struct {
	AppID string `plist:"application-identifier"`
}

// ProvisioningProfileModel ...
type ProvisioningProfileModel struct {
	Entitlements EntitlementsModel `plist:"Entitlements"`
}

// FindProvProfilesFileByAppID ...
func FindProvProfilesFileByAppID(appID string) ([]string, error) {
	absProvProfileDirPath, err := pathutil.AbsPath(provProfileSystemDirPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to get Absolute path of Provisioning Profiles dir: %s", err)
	}

	pths, err := filepath.Glob(absProvProfileDirPath + "/*.mobileprovision")
	if err != nil {
		return nil, fmt.Errorf("Failed to perform *.mobileprovision search, error: %s", err)
	}

	provProfilePathsToReturn := []string{}
	for _, aPth := range pths {
		profileContent, err := cmdex.NewCommand("security", "cms", "-D", "-i", aPth).RunAndReturnTrimmedCombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve information about Provisioning Profile (path: %s), error: %s",
				aPth, err)
		}

		var provProfileData ProvisioningProfileModel
		if err := plist.NewDecoder(strings.NewReader(profileContent)).Decode(&provProfileData); err != nil {
			return nil, fmt.Errorf("Failed to parse Provisioning Profile (path: %s), error: %s", aPth, err)
		}
		if provProfileData.Entitlements.AppID == appID {
			provProfilePathsToReturn = append(provProfilePathsToReturn, aPth)
		}
	}

	return provProfilePathsToReturn, nil
}

// FindProvProfileFileByUUID ...
func FindProvProfileFileByUUID(provProfileUUID string) (string, error) {
	absProvProfileDirPath, err := pathutil.AbsPath(provProfileSystemDirPath)
	if err != nil {
		return "", fmt.Errorf("Failed to get Absolute path of Provisioning Profiles dir: %s", err)
	}

	mobileProvPth := filepath.Join(absProvProfileDirPath, provProfileUUID+".mobileprovision")
	exist, err := pathutil.IsPathExists(mobileProvPth)
	if !exist || err != nil {
		log.Debugf("No mobileprovision file found at: %s | err: %s", mobileProvPth, err)
	} else {
		return mobileProvPth, nil
	}

	macProvProfPth := filepath.Join(absProvProfileDirPath, provProfileUUID+".provisionprofile")
	exist, err = pathutil.IsPathExists(macProvProfPth)
	if !exist || err != nil {
		log.Debugf("No provisionprofile file found at: %s | err: %s", macProvProfPth, err)
		return "", fmt.Errorf("Failed to find Provisioning Profile with UUID: %s", provProfileUUID)
	}

	return macProvProfPth, nil
}
