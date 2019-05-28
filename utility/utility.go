package utility

import (
	"regexp"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/profileutil"
)

// ProfileExportFileNameNoPath creates a file name for the given profile with pattern: uuid.escaped_profile_name.[mobileprovision|provisionprofile]
func ProfileExportFileNameNoPath(info profileutil.ProvisioningProfileInfoModel) string {
	replaceRexp, err := regexp.Compile("[^A-Za-z0-9_.-]")
	if err != nil {
		log.Warnf("Invalid regex, error: %s", err)
		return ""
	}
	safeTitle := replaceRexp.ReplaceAllString(info.Name, "")
	extension := ".mobileprovision"
	if info.Type == profileutil.ProfileTypeMacOs {
		extension = ".provisionprofile"
	}

	return info.UUID + "." + safeTitle + extension
}
