package cmd

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/spf13/cobra"
)

var (
	enableVerboseLog = false
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "codesigndoc",
	Short: "Your friendly iOS Code Signing Doctor",
	Long: `Your friendly iOS Code Signing Doctor

Using this tool is as easy as running "codesigndoc scan xcode/xamarin" and following the guide it prints.

At the end of the process you'll have all the code signing files
(.p12 Identity file including the Certificate and Private Key,
and the required Provisioning Profiles) required to do a successful Archive of your iOS project.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.SetEnableDebugLog(enableVerboseLog)
		log.Debugf("EnableDebugLog: %v", enableVerboseLog)

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&enableVerboseLog, "verbose", "v", false, "Enable verbose logging")
}
