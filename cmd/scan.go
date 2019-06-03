package cmd

import (
	"fmt"

	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/command"
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
		case "always":
			{
				writeFiles = codesign.WriteFilesAlways
			}
		case "fallback":
			{
				writeFiles = codesign.WriteFilesFallback
			}
		case "disable":
			{
				writeFiles = codesign.WriteFilesDisabled
			}
		default:
			{
				return fmt.Errorf("invalid value for %s flag. Valid values: 'always', 'fallback', 'disable'", writeFilesFlag)
			}
		}
		appSlug := cmd.Flag(appSlugFlag).Value.String()
		authToken := cmd.Flag(authTokenFlag).Value.String()
		if appSlug != "" && authToken == "" ||
			appSlug == "" && authToken != "" {
			return fmt.Errorf("both or none flags %s and %s are required to be set", appSlugFlag, authTokenFlag)
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
	scanCmd.PersistentFlags().String(writeFilesFlag, "always", `Set wether to export build logs and codesigning files to the ./codesigndoc_exports directory. Defaults to "always". Valid values: "always", "fallback", "disable".
- always: Writes artifacts in every case.
- fallback: Does not write artifacts if the automatic upload option is chosen interactively or by providing the auth-token and app-slug flag. Writes build log only on failure.
- disabled: Do not write any files to the export directory.`)
	// Flags used to automatically upload artifacts.
	scanCmd.PersistentFlags().StringVar(&personalAccessToken, authTokenFlag, "", `Bitrise personal access token. By default codesigndoc will ask for it interactively.
Will upload codesigning files automatically if provided. Requires the app-slug paramater to be also set.`)
	scanCmd.PersistentFlags().StringVar(&appSlug, appSlugFlag, "", `Bitrise app slug. By default codesigndoc will ask for it interactively.
Will upload codesigning files automatically if provided. Requires the auth-token parameter to be also set.`)
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

func printFinished(exportResult codesign.ExportReport, absOutputDir string) {
	if exportResult.CodesignFilesWritten {
		fmt.Println()
		log.Successf("Exports finished you can find the exported files at: %s", absOutputDir)

		if err := command.RunCommand("open", absOutputDir); err != nil {
			log.Errorf("Failed to open the export directory in Finder: %s", absOutputDir)
		} else {
			fmt.Println("Opened the directory in Finder.")
		}
	}

	fmt.Println()
	log.Successf("That's all.")

	if !exportResult.ProvisioningProfilesUploaded && !exportResult.CertificatesUploaded {
		log.Warnf("You just have to upload the found certificates (.p12) and provisioning profiles (.mobileprovision) and you'll be good to go!")
		fmt.Println()
	} else if !exportResult.CertificatesUploaded {
		log.Warnf("You just have to upload the found certificates (.p12) and you'll be good to go!")
		fmt.Println()
	}
}
