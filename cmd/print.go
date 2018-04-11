package cmd

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xcode/export"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

func printFinished(provProfilesUploaded bool, certsUploaded bool) {
	fmt.Println()
	log.Successf("That's all.")

	if !provProfilesUploaded && !certsUploaded {
		log.Warnf("You just have to upload the found certificates (.p12) and provisioning profiles (.mobileprovision) and you'll be good to go!")
		fmt.Println()
	}
}

// PrintIOSCodesignGroup ...
func printIOSCodesignGroup(group export.IosCodeSignGroup) {
	printCodesignGroup(group.BundleIDProfileMap, group.Certificate.TeamName, group.Certificate.TeamID, group.Certificate.CommonName, group.Certificate.Serial)
}

func printMacOsCodesignGroup(group export.MacCodeSignGroup) {
	printCodesignGroup(group.BundleIDProfileMap, group.Certificate.TeamName, group.Certificate.TeamID, group.Certificate.CommonName, group.Certificate.Serial)
}

func printCodesignGroup(bundleIDProfileMap map[string]profileutil.ProvisioningProfileInfoModel, teamName string, teamID string, commonName string, serial string) {
	fmt.Printf("%s %s (%s)\n", colorstring.Green("development team:"), teamName, teamID)
	fmt.Printf("%s %s [%s]\n", colorstring.Green("codesign identity:"), commonName, serial)
	idx := -1
	for bundleID, profile := range bundleIDProfileMap {
		idx++
		if idx == 0 {
			fmt.Printf("%s %s -> %s\n", colorstring.Greenf("provisioning profiles:"), profile.Name, bundleID)
		} else {
			fmt.Printf("%s%s -> %s\n", strings.Repeat(" ", len("provisioning profiles: ")), profile.Name, bundleID)
		}
	}
}
