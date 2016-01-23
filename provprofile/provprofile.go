package provprofile

import (
	"fmt"
	"path"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/going/toolkit/log"
)

const (
	provProfileSystemDirPath = "~/Library/MobileDevice/Provisioning Profiles"
)

// ProvisioningProfileInfo ...
type ProvisioningProfileInfo struct {
	Title string
	UUID  string
}

// FindProvProfileFile ...
func FindProvProfileFile(provProfileInfo ProvisioningProfileInfo) (string, error) {
	absProvProfileDirPath, err := pathutil.AbsPath(provProfileSystemDirPath)
	if err != nil {
		return "", fmt.Errorf("Failed to get Absolute path of Provisioning Profiles dir: %s", err)
	}

	mobileProvPth := path.Join(absProvProfileDirPath, provProfileInfo.UUID+".mobileprovision")
	exist, err := pathutil.IsPathExists(mobileProvPth)
	if !exist || err != nil {
		log.Debugf("No mobileprovision file found at: %s | err: %s", mobileProvPth, err)
	} else {
		return mobileProvPth, nil
	}

	macProvProfPth := path.Join(absProvProfileDirPath, provProfileInfo.UUID+".provisionprofile")
	pathutil.IsPathExists(macProvProfPth)
	if !exist || err != nil {
		log.Debugf("No provisionprofile file found at: %s | err: %s", macProvProfPth, err)
		return "", fmt.Errorf("Failed to find Provisioning Profile with UUID: %s", provProfileInfo.UUID)
	}

	return macProvProfPth, nil
}

// // PrintFileInfo ...
// func PrintFileInfo(pth string) error {
// 	decoder := plist.NewDecoder(strings.NewReader(pth))
// 	dataMap := map[string]string{}
// 	if err := decoder.Decode(&dataMap); err != nil {
// 		return fmt.Errorf("Failed to decode Plist data: %s", err)
// 	}
// 	fmt.Printf("dataMap: %#v\n", dataMap)
// 	return nil
// }
