package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/spf13/cobra"
)

// xamarinCmd represents the xamarin command
var xamarinCmd = &cobra.Command{
	Use:   "xamarin",
	Short: "Xamarin project scanner",
	Long:  `Scan a Xamarin project`,
	RunE:  scanXamarinProject,
}

var (
	xamarinSolutionFilePath  = ""
	xamarinProjectName       = ""
	xamarinConfigurationName = ""
)

func init() {
	scanCmd.AddCommand(xamarinCmd)

	xamarinCmd.Flags().StringVar(&xamarinSolutionFilePath,
		"file", "",
		`Xamarin Solution file path`)
	xamarinCmd.Flags().StringVar(&xamarinProjectName,
		"project", "",
		`Xamarin iOS Project Name (e.g.: "MyProject.iOS")`)
	xamarinCmd.Flags().StringVar(&xamarinConfigurationName,
		"config", "",
		`Xamarin Configuration Name (e.g.: "Release|iPhone")`)
}

func printXamarinScanFinishedWithError(format string, args ...interface{}) error {
	return printFinishedWithError("Xamarin", format, args...)
}

// XamarinCommandModel ...
type XamarinCommandModel struct {
	SolutionFilePath  string
	ProjectName       string
	ConfigurationName string
}

func scanXamarinProject(cmd *cobra.Command, args []string) error {
	absExportOutputDirPath, err := initExportOutputDir()
	if err != nil {
		return printXamarinScanFinishedWithError("Failed to prepare Export directory: %s", err)
	}
	log.Info("absExportOutputDirPath: ", absExportOutputDirPath)

	// --- Inputs ---

	xamarinCmd := XamarinCommandModel{}

	// Xamarin Solution Path
	xamarinCmd.SolutionFilePath = xcodeProjectFilePath
	if xamarinCmd.SolutionFilePath == "" {
		askText := `Please drag-and-drop your Xamarin Solution (` + colorstring.Green(".sln") + `)
   file here, and then hit Enter`
		fmt.Println()
		projpth, err := goinp.AskForPath(askText)
		if err != nil {
			return printXamarinScanFinishedWithError("Failed to read input: %s", err)
		}
		xamarinCmd.SolutionFilePath = projpth
	}
	log.Debugf("xamSolutionPth: %s", xamarinCmd.SolutionFilePath)

	// Xamarin Project Name
	xamarinCmd.ProjectName = xamarinProjectName
	if xamarinCmd.ProjectName == "" {
		fmt.Println()
		answerValue, err := goinp.AskForString(
			`What's the name of the Project you use for creating your iOS apps (e.g.: MyProject.iOS)?`,
		)
		if err != nil {
			return printXamarinScanFinishedWithError("Failed to read input: %s", err)
		}
		xamarinCmd.ProjectName = answerValue
	}

	// Xamarin Configuration Name
	xamarinCmd.ConfigurationName = xamarinConfigurationName
	if xamarinCmd.ConfigurationName == "" {
		fmt.Println()
		answerValue, err := goinp.AskForStringWithDefault(
			`What's the name of the Configuration you use for creating your iOS apps?
Specify it if it's not "Release|iPhone", or hit Enter if it is`,
			"Release|iPhone",
		)
		if err != nil {
			return printXamarinScanFinishedWithError("Failed to read input: %s", err)
		}
		xamarinCmd.ConfigurationName = answerValue
	}

	log.Infof("xamarinCmd: %#v", xamarinCmd)

	return nil
}
