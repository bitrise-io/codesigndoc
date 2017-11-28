package cmd

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xcode/export"
)

func printFinishedWithError(toolName, format string, args ...interface{}) error {
	fmt.Println()
	fmt.Println("------------------------------")
	fmt.Println("First of all " + colorstring.Red("please make sure that you can Archive your app from "+toolName+"."))
	fmt.Println("codesigndoc only works if you can archive your app from " + toolName + ".")
	fmt.Println("If you can, and you get a valid IPA file if you export from " + toolName + ",")
	fmt.Println(colorstring.Red("please create an issue") + " on GitHub at: https://github.com/bitrise-tools/codesigndoc/issues")
	fmt.Println("with as many details & logs as you can share!")
	fmt.Println("------------------------------")
	fmt.Println()

	return fmt.Errorf(colorstring.Red("Error: ")+format, args...)
}

func printFinished() {
	fmt.Println()
	log.Successf("That's all.")
	log.Warnf("You just have to upload the found code signing files (.p12 and .mobileprovision) and you'll be good to go!")
	fmt.Println()
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
