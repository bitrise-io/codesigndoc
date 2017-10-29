package cmd

import "github.com/spf13/cobra"

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a project's code signing settings",
	Long: `Scan a project's code signing settings,
and optionally export the require code signing files.`,
}

var (
	isAllowExport    = false
	isAskForPassword = false
	certificatsOnly  = false
)

func init() {
	RootCmd.AddCommand(scanCmd)
	scanCmd.PersistentFlags().BoolVar(&isAllowExport, "allow-export", false, "Automatically allow export of discovered files")
	scanCmd.PersistentFlags().BoolVar(&isAskForPassword, "ask-pass", false, "Ask for .p12 password, instead of using an empty password")
	scanCmd.PersistentFlags().BoolVar(&certificatsOnly, "certs-only", false, "Collect Certificates (Identities) only")
}
