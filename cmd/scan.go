package cmd

import (
	"fmt"

	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/spf13/cobra"
)

const (
	appSlugFlag    = "app-slug"
	authTokenFlag  = "auth-token"
	writeFilesFlag = "write-files"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a project's code signing settings",
	Long: `Scan a project's code signing settings,
and export the require code signing files.`,
	TraverseChildren: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		switch cmd.Flag(writeFilesFlag).Value.String() {
		case "":
			{
				writeFiles = codesign.WriteFilesAlways
			}
		case string(codesign.WriteFilesAlways):
			{
				writeFiles = codesign.WriteFilesAlways
			}
		case string(codesign.WriteFilesFallback):
			{
				writeFiles = codesign.WriteFilesFallback
			}
		case string(codesign.WriteFilesDisabled):
			{
				writeFiles = codesign.WriteFilesDisabled
			}
		default:
			{
				return fmt.Errorf("invalid value for write-files paramter. Valid values: 'always','fallback','disabled'")
			}
		}
		log.Printf("File output level: %s", writeFiles)

		appSlug := cmd.Flag(appSlugFlag).Value.String()
		authToken := cmd.Flag(authTokenFlag).Value.String()
		if appSlug != "" && authToken == "" ||
			appSlug == "" && authToken != "" {
			return fmt.Errorf("both %s and %s are required to be set for automatic upload", appSlugFlag, authTokenFlag)
		}
		return nil
	},
}

var (
	isAskForPassword bool
	certificatesOnly bool
	writeFiles       codesign.WriteFilesLevel

	personalAccessToken string
	appSlug             string
)

func init() {
	RootCmd.AddCommand(scanCmd)
	scanCmd.PersistentFlags().BoolVar(&isAskForPassword, "ask-pass", false, "Ask for .p12 password, instead of using an empty password")
	scanCmd.PersistentFlags().BoolVar(&certificatesOnly, "certs-only", false, "Collect Certificates (Identities) only")
	var writeFilesRaw string
	scanCmd.PersistentFlags().StringVar(&writeFilesRaw, writeFilesFlag, "", "Set wether to export artifacts to a local directory.")
	// Flags used to automatically upload artifacts
	scanCmd.PersistentFlags().StringVar(&personalAccessToken, authTokenFlag, "", "Personal access token. Requires the app-slug paramater to be also set. Will upload codesigning files automatically if provided.")
	scanCmd.PersistentFlags().StringVar(&appSlug, appSlugFlag, "", "App Slug. Requires the auth-token parameter to be also set. Will upload codesigning files automatically if provided.")
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
If you can, and you get a valid .ipa/.app file if you export from ` + string(e.tool) + `,
` + colorstring.Red("please create an issue") + ` on GitHub at: https://github.com/bitrise-io/codesigndoc/issues
with as many details & logs as you can share!
------------------------------

` + colorstring.Redf("Error: %s", e.msg)
}

// BuildForTestingError ...
type BuildForTestingError struct {
	tool Tool
	msg  string
}

// Error ...
func (e BuildForTestingError) Error() string {
	return colorstring.Redf("Error: %s", e.msg) + `

------------------------------` + `
First of all, check the selected scheme in ` + string(e.tool) + `:
- Make sure, you have enabled at least one UITest target for test run in the selected scheme's build option.
- Make sure that the UITest target is added (and enabled) in the selected scheme's test option.

After this ` + colorstring.Red("please make sure that you can run build-for-testing for your app from "+e.tool+".") + `
codesigndoc only works if you can run build-for-testing for your app from ` + string(e.tool) + `.
For this run a ` + colorstring.Red("clean") + ` in your ` + string(e.tool) + `, after that, run a ` + colorstring.Red("build-for-testing") + ` for your app in ` + string(e.tool) + `.
If you can, and you get a valid *-Runner.app file, ` + colorstring.Red("please create an issue") + ` on GitHub at: https://github.com/bitrise-io/codesigndoc/issues
with as many details & logs as you can share!
------------------------------
`
}

func printFinished(provProfilesUploaded bool, certsUploaded bool) {
	fmt.Println()
	log.Successf("That's all.")

	if !provProfilesUploaded && !certsUploaded {
		log.Warnf("You just have to upload the found certificates (.p12) and provisioning profiles (.mobileprovision) and you'll be good to go!")
		fmt.Println()
	} else if !certsUploaded {
		log.Warnf("You just have to upload the found certificates (.p12) and you'll be good to go!")
		fmt.Println()
	}
}
