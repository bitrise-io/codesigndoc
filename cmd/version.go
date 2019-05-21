package cmd

import (
	"fmt"
	"runtime"

	"github.com/bitrise-io/codesigndoc/version"
	"github.com/spf13/cobra"
)

var (
	isFullVersionPrint = false
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version number",
	Long:  `Prints version number`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.VERSION)

		if isFullVersionPrint {
			fmt.Println()
			fmt.Println("go: " + runtime.Version())
			fmt.Println("arch: " + runtime.GOARCH)
			fmt.Println("os: " + runtime.GOOS)
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&isFullVersionPrint, "full", false, "Full version")
}
