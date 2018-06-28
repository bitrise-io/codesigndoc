package cmd

import (
	"fmt"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a project's code signing settings",
	Long: `Scan a project's code signing settings,
and export the require code signing files.`,
	TraverseChildren: true,
}

var (
	isAllowExport    = false
	isAskForPassword = false
	certificatesOnly = false
)

func init() {
	RootCmd.AddCommand(scanCmd)
	scanCmd.PersistentFlags().BoolVar(&isAllowExport, "allow-export", false, "Automatically allow export of discovered files")
	scanCmd.PersistentFlags().BoolVar(&isAskForPassword, "ask-pass", false, "Ask for .p12 password, instead of using an empty password")
	scanCmd.PersistentFlags().BoolVar(&certificatesOnly, "certs-only", false, "Collect Certificates (Identities) only")
}

// Tool ...
type Tool string

const (
	toolXcode   Tool = "Xcode"
	toolXamarin Tool = "Visual Studio"
)

// ArchiveError ...
type ArchiveError struct {
	tool Tool
	msg  string
}

// Error ...
func (e ArchiveError) Error() string {
	return `
------------------------------` + `
First of all ` + colorstring.Red("please make sure that you can Archive your app from "+e.tool+".") + `
codesigndoc only works if you can archive your app from ` + string(e.tool) + `.
If you can, and you get a valid IPA/.app file if you export from ` + string(e.tool) + `,
` + colorstring.Red("please create an issue") + ` on GitHub at: https://github.com/bitrise-tools/codesigndoc/issues
with as many details & logs as you can share!
------------------------------

` + colorstring.Redf("Error: %s", e.msg)
}

func printFinished(provProfilesUploaded bool, certsUploaded bool) {
	fmt.Println()
	log.Successf("That's all.")

	if !provProfilesUploaded && !certsUploaded {
		log.Warnf("You just have to upload the found certificates (.p12) and provisioning profiles (.mobileprovision) and you'll be good to go!")
		fmt.Println()
	}
}
