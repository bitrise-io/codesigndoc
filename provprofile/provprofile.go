package provprofile

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	plist "github.com/DHowett/go-plist"
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/maputil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/pkg/errors"
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
	AppID string `plist:"application-identifier"`
}

// ProvisioningProfileModel ...
type ProvisioningProfileModel struct {
	Entitlements    EntitlementsModel `plist:"Entitlements"`
	UUID            string            `plist:"UUID"`
	TeamName        string            `plist:"TeamName"`
	Name            string            `plist:"Name"`
	AppIDName       string            `plist:"AppIDName"`
	TeamIdentifiers []string          `plist:"TeamIdentifier"`
	ExpirationDate  time.Time         `plist:"ExpirationDate"`
}

// TeamID ...
func (provProfile ProvisioningProfileModel) TeamID() (string, error) {
	if len(provProfile.TeamIdentifiers) == 0 {
		return "", errors.New("No TeamIdentifier specified")
	}
	if len(provProfile.TeamIdentifiers) != 1 {
		return "", errors.New("More than one TeamIdentifier specified")
	}

	teamID := provProfile.TeamIdentifiers[0]
	if len(teamID) == 0 {
		return "", errors.New("An empty item specified for TeamIdentifier")
	}

	return teamID, nil
}

// ProvisioningProfileFileInfoModel ...
type ProvisioningProfileFileInfoModel struct {
	Path                    string
	ProvisioningProfileInfo ProvisioningProfileModel
}

// ProvisioningProfileFileInfoModels ...
type ProvisioningProfileFileInfoModels []ProvisioningProfileFileInfoModel

// CollectTeamIDs ...
func (ppFileInfos ProvisioningProfileFileInfoModels) CollectTeamIDs() ([]string, error) {
	teamIDsMap := map[string]interface{}{}
	for _, aProvProfileFileInfo := range ppFileInfos {
		teamID, err := aProvProfileFileInfo.ProvisioningProfileInfo.TeamID()
		if err != nil {
			return []string{}, fmt.Errorf("Team ID error for profile (uuid: %s), error: %s", aProvProfileFileInfo.ProvisioningProfileInfo.UUID, err)
		}
		teamIDsMap[teamID] = 1
	}
	return maputil.KeysOfStringInterfaceMap(teamIDsMap), nil
}

// CreateProvisioningProfileModelFromFile ...
func CreateProvisioningProfileModelFromFile(filePth string) (ProvisioningProfileModel, error) {
	profileContent, err := command.New("security", "cms", "-D", "-i", filePth).RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return ProvisioningProfileModel{},
			fmt.Errorf("Failed to retrieve information about Provisioning Profile, error: %s",
				err)
	}

	log.Debugln("profileContent: ", profileContent)

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

func walkProvProfiles(ppWalkFn func(provProfile ProvisioningProfileFileInfoModel) error) error {
	absProvProfileDirPath, err := pathutil.AbsPath(provProfileSystemDirPath)
	if err != nil {
		return errors.Wrap(err, "Failed to get Absolute path of Provisioning Profiles dir")
	}

	pths, err := filepath.Glob(absProvProfileDirPath + "/*.mobileprovision")
	if err != nil {
		return errors.Wrap(err, "Failed to perform *.mobileprovision search")
	}

	for _, aPth := range pths {
		provProfileData, err := CreateProvisioningProfileModelFromFile(aPth)
		if err != nil {
			return errors.Wrapf(err, "Failed to read Provisioning Profile infos from file (path: %s)",
				aPth)
		}

		if time.Now().After(provProfileData.ExpirationDate) {
			log.Warnf(colorstring.Yellow(" (!) ")+"Provisioning Profile %s "+colorstring.Yellow("expired")+" at %s. Skipping.",
				colorstring.Blue(provProfileData.UUID), colorstring.Blue(provProfileData.ExpirationDate))
			log.Warnf("     If you want to delete this Provisioning Profile you can find it at: %s",
				colorstring.Yellow(aPth))
			continue
		}

		if err := ppWalkFn(ProvisioningProfileFileInfoModel{Path: aPth, ProvisioningProfileInfo: provProfileData}); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
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

		if time.Now().After(provProfileData.ExpirationDate) {
			log.Warnf(colorstring.Yellow(" (!) ")+"Provisioning Profile %s "+colorstring.Yellow("expired")+" at %s. Skipping.",
				colorstring.Blue(provProfileData.UUID), colorstring.Blue(provProfileData.ExpirationDate))
			log.Warnf("     If you want to delete this Provisioning Profile you can find it at: %s",
				colorstring.Yellow(aPth))
			continue
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
		isFound := false
		mobileProvPth := filepath.Join(absProvProfileDirPath, provProfileUUID+".mobileprovision")
		exist, err := pathutil.IsPathExists(mobileProvPth)
		if err != nil {
			log.Debugf("No mobileprovision file found at: %s | err: %s", mobileProvPth, err)
		} else if !exist {
			log.Debugf("Not found at path (%s), doing a full search by content ...", mobileProvPth)
			// try by content
			err := walkProvProfiles(func(ppf ProvisioningProfileFileInfoModel) error {
				if ppf.ProvisioningProfileInfo.UUID == provProfileUUID {
					isFound = true
					mobileProvPth = ppf.Path
				}
				return nil
			})
			if err != nil {
				log.Debugf("Error during Prov Profile walk: %+v", errors.WithStack(err))
			}
			if !isFound {
				log.Debugf("Prov Profile not found by UUID (walk)")
			}
		} else {
			isFound = true
		}

		if isFound {
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
