package xcarchive

import (
	"path/filepath"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/plistutil"
)

// IsMacOS try to find the Contents dir under the .app/.
// If its finds it the archive is MacOs. If it does not the archive is iOS.
func IsMacOS(archPath string) (bool, error) {
	log.Debugf("Checking archive is MacOS or iOS")
	infoPlistPath := filepath.Join(archPath, "Info.plist")

	plist, err := plistutil.NewPlistDataFromFile(infoPlistPath)
	if err != nil {
		return false, err
	}

	appProperties, found := plist.GetMapStringInterface("ApplicationProperties")
	if !found {
		return false, err
	}

	applicationPath, found := appProperties.GetString("ApplicationPath")
	if !found {
		return false, err
	}

	applicationPath = filepath.Join(archPath, "Products", applicationPath)
	contentsPath := filepath.Join(applicationPath, "Contents")

	exist, err := pathutil.IsDirExists(contentsPath)
	if err != nil {
		return false, err
	}

	return exist, nil
}
