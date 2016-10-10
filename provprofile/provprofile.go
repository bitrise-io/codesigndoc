package provprofile

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	plist "github.com/DHowett/go-plist"
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/maputil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/ryanuber/go-glob"
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
	AppID  string `plist:"application-identifier"`
	TeamID string `plist:"com.apple.developer.team-identifier"`
}

// ProvisioningProfileModel ...
type ProvisioningProfileModel struct {
	Entitlements   EntitlementsModel `plist:"Entitlements"`
	UUID           string            `plist:"UUID"`
	TeamName       string            `plist:"TeamName"`
	Name           string            `plist:"Name"`
	AppIDName      string            `plist:"AppIDName"`
	ExpirationDate time.Time         `plist:"ExpirationDate"`
}

// ProvisioningProfileFileInfoModel ...
type ProvisioningProfileFileInfoModel struct {
	Path                    string
	ProvisioningProfileInfo ProvisioningProfileModel
}

// ProvisioningProfileFileInfoModels ...
type ProvisioningProfileFileInfoModels []ProvisioningProfileFileInfoModel

// CollectTeamIDs ...
func (ppFileInfos ProvisioningProfileFileInfoModels) CollectTeamIDs() []string {
	teamIDsMap := map[string]interface{}{}
	for _, aProvProfileFileInfo := range ppFileInfos {
		teamIDsMap[aProvProfileFileInfo.ProvisioningProfileInfo.Entitlements.TeamID] = 1
	}
	return maputil.KeysOfStringInterfaceMap(teamIDsMap)
}

// CreateProvisioningProfileModelFromFile ...
func CreateProvisioningProfileModelFromFile(filePth string) (ProvisioningProfileModel, error) {
	profileContent, err := cmdex.NewCommand("security", "cms", "-D", "-i", filePth).RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return ProvisioningProfileModel{},
			fmt.Errorf("Failed to retrieve information about Provisioning Profile, error: %s",
				err)
	}

	var provProfileData ProvisioningProfileModel
	if err := plist.NewDecoder(strings.NewReader(profileContent)).Decode(&provProfileData); err != nil {
		return provProfileData,
			fmt.Errorf("Failed to parse Provisioning Profile content, error: %s", err)
	}

	if provProfileData.UUID == "" {
		return provProfileData,
			fmt.Errorf("No UUID found in the Provisioning Profile (%#v)", provProfileData)
	}

	return provProfileData, nil
}

// FindProvProfilesByAppID ...
// `appID`` supports "glob", e.g.: *.bundle.id will match any Prov Profile with ".bundle.id"
//   app ID suffix
func FindProvProfilesByAppID(appID string) ([]ProvisioningProfileFileInfoModel, error) {
	absProvProfileDirPath, err := pathutil.AbsPath(provProfileSystemDirPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to get Absolute path of Provisioning Profiles dir: %s", err)
	}

	pths, err := filepath.Glob(absProvProfileDirPath + "/*.mobileprovision")
	if err != nil {
		return nil, fmt.Errorf("Failed to perform *.mobileprovision search, error: %s", err)
	}

	provProfilePathsToReturn := []ProvisioningProfileFileInfoModel{}
	for _, aPth := range pths {
		provProfileData, err := CreateProvisioningProfileModelFromFile(aPth)
		if err != nil {
			return nil, fmt.Errorf("Failed to read Provisioning Profile infos from file (path: %s), error: %s",
				aPth, err)
		}

		if glob.Glob(appID, provProfileData.Entitlements.AppID) {
			provProfilePathsToReturn = append(provProfilePathsToReturn, ProvisioningProfileFileInfoModel{
				Path: aPth,
				ProvisioningProfileInfo: provProfileData,
			})
		}
	}

	return provProfilePathsToReturn, nil
}

// FindProvProfileByUUID ...
func FindProvProfileByUUID(provProfileUUID string) (ProvisioningProfileFileInfoModel, error) {
	absProvProfileDirPath, err := pathutil.AbsPath(provProfileSystemDirPath)
	if err != nil {
		return ProvisioningProfileFileInfoModel{}, fmt.Errorf("Failed to get Absolute path of Provisioning Profiles dir: %s", err)
	}

	// iOS / .mobileprovision
	{
		mobileProvPth := filepath.Join(absProvProfileDirPath, provProfileUUID+".mobileprovision")
		exist, err := pathutil.IsPathExists(mobileProvPth)
		if !exist || err != nil {
			log.Debugf("No mobileprovision file found at: %s | err: %s", mobileProvPth, err)
		} else {
			provProfileData, err := CreateProvisioningProfileModelFromFile(mobileProvPth)
			if err != nil {
				return ProvisioningProfileFileInfoModel{},
					fmt.Errorf("Failed to read Provisioning Profile infos from file (path: %s), error: %s",
						mobileProvPth, err)
			}
			return ProvisioningProfileFileInfoModel{
				Path: mobileProvPth,
				ProvisioningProfileInfo: provProfileData,
			}, nil
		}
	}

	// Mac / .provisionprofile
	{
		macProvProfPth := filepath.Join(absProvProfileDirPath, provProfileUUID+".provisionprofile")
		exist, err := pathutil.IsPathExists(macProvProfPth)
		if !exist || err != nil {
			log.Debugf("No provisionprofile file found at: %s | err: %s", macProvProfPth, err)
			return ProvisioningProfileFileInfoModel{}, fmt.Errorf("Failed to find Provisioning Profile with UUID: %s", provProfileUUID)
		}

		provProfileData, err := CreateProvisioningProfileModelFromFile(macProvProfPth)
		if err != nil {
			return ProvisioningProfileFileInfoModel{},
				fmt.Errorf("Failed to read Provisioning Profile infos from file (path: %s), error: %s",
					macProvProfPth, err)
		}
		return ProvisioningProfileFileInfoModel{
			Path: macProvProfPth,
			ProvisioningProfileInfo: provProfileData,
		}, nil
	}
}
