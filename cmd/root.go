package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	paramLogLevel = "info"
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
		// Log level
		logLevel, err := log.ParseLevel(paramLogLevel)
		if err != nil {
			return fmt.Errorf("Failed to parse log level: %s", err)
		}
		log.SetLevel(logLevel)
		log.Debugf("Loglevel: %s", logLevel)

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
	RootCmd.PersistentFlags().StringVarP(&paramLogLevel,
		"loglevel", "l",
		"info",
		"Log level (options: debug, info, warn, error, fatal, panic).")
}
