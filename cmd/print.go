package cmd

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xcode/export"
)

func printFinished(provProfilesUploaded bool, certsUploaded bool) {
	fmt.Println()
	log.Successf("That's all.")

	if !provProfilesUploaded && !certsUploaded {
		log.Warnf("You just have to upload the found certificates (.p12) and provisioning profiles (.mobileprovision) and you'll be good to go!")
		fmt.Println()
	}
}

func printCodesignGroup(group export.IosCodeSignGroup) {
	fmt.Printf("%s %s (%s)\n", colorstring.Green("development team:"), group.Certificate.TeamName, group.Certificate.TeamID)
	fmt.Printf("%s %s [%s]\n", colorstring.Green("codesign identity:"), group.Certificate.CommonName, group.Certificate.Serial)
	idx := -1
	for bundleID, profile := range group.BundleIDProfileMap {
		idx++
		if idx == 0 {
			fmt.Printf("%s %s -> %s\n", colorstring.Greenf("provisioning profiles:"), profile.Name, bundleID)
		} else {
			fmt.Printf("%s%s -> %s\n", strings.Repeat(" ", len("provisioning profiles: ")), profile.Name, bundleID)
		}
	}
}
