package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/xcode"
	"github.com/codegangsta/cli"
)

func scan(c *cli.Context) {
	projectPath := c.String(FileParamKey)
	if projectPath == "" {
		askText := `Please drang-and-drop your Xcode Project (` + colorstring.Green(".xcodeproj") + `)
   or Workspace (` + colorstring.Green(".xcworkspace") + `) file, the one you usually open in Xcode,
   then hit Enter.

  (Note: if you have a Workspace file you should use that)`
		fmt.Println()
		projpth, err := goinp.AskForString(askText)
		if err != nil {
			log.Fatalf("Failed to read input: %s", err)
		}
		projectPath = projpth
	}
	log.Debugf("projectPath: %s", projectPath)
	xcodeCmd := xcode.CommandModel{
		ProjectFilePath: projectPath,
	}

	schemeToUse := c.String(SchemeParamKey)
	if schemeToUse == "" {
		log.Println("Scanning Schemes ...")
		schemes, err := xcodeCmd.ScanSchemes()
		if err != nil {
			log.Fatalf("Failed to scan Schemes: %s", err)
		}
		log.Debugf("schemes: %v", schemes)

		fmt.Println()
		selectedScheme, err := goinp.SelectFromStrings("Select the Scheme you usually use in Xcode", schemes)
		if err != nil {
			log.Fatalf("Failed to select Scheme: %s", err)
		}
		log.Debugf("selected scheme: %v", selectedScheme)
		schemeToUse = selectedScheme
	}
	xcodeCmd.Scheme = schemeToUse

	fmt.Println()
	fmt.Println("Running an Xcode Archive, to get all the required code signing settings...")
	if err := xcodeCmd.ScanCodeSigningSettings(); err != nil {
		log.Fatalf("Failed to detect code signing settings: %s", err)
	}
}
