package cmd

import (
	"fmt"

	"github.com/bitrise-tools/codesigndoc/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version number",
	Long:  `Prints version number`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.VERSION)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
