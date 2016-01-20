package cli

import (
	"log"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/xcode"
	"github.com/codegangsta/cli"
)

func getAvailableSchemes(xcode.CommandModel) {

}

func scan(c *cli.Context) {
	askText := `Please drang-and-drop your Xcode Project (` + colorstring.Green(".xcodeproj") + `)
 or Workspace (` + colorstring.Green(".xcworkspace") + `) file, the one you usually open in Xcode,
 then hit Enter.

(Note: if you have a Workspace file you should use that)`
	projectPath, err := goinp.AskForString(askText)
	if err != nil {
		log.Fatalf("Failed to read input: %s", err)
	}
	xcodeCmd := xcode.CommandModel{ProjectFilePath: projectPath}

	getAvailableSchemes(xcodeCmd)
}
