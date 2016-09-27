package cmd

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/xamarin"
	"github.com/bitrise-tools/go-xamarin/solution"
	"github.com/spf13/cobra"
)

// xamarinCmd represents the xamarin command
var xamarinCmd = &cobra.Command{
	Use:   "xamarin",
	Short: "Xamarin project scanner",
	Long:  `Scan a Xamarin project`,

	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          scanXamarinProject,
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

func scanXamarinProject(cmd *cobra.Command, args []string) error {
	absExportOutputDirPath, err := initExportOutputDir()
	if err != nil {
		return printXamarinScanFinishedWithError("Failed to prepare Export directory: %s", err)
	}
	log.Info("absExportOutputDirPath: ", absExportOutputDirPath)

	// --- Inputs ---

	xamarinCmd := xamarin.CommandModel{}

	// Xamarin Solution Path
	xamarinCmd.SolutionFilePath = xamarinSolutionFilePath
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

	xamSln, err := solution.New(xamarinCmd.SolutionFilePath, false)
	if err != nil {
		return printXamarinScanFinishedWithError("Failed to analyze Xamarin solution: %s", err)
	}
	log.Infof("xamSln: %#v", xamSln)

	// Xamarin Project Name
	xamarinCmd.ProjectName = xamarinProjectName
	if xamarinCmd.ProjectName == "" {
		fmt.Println()
		answerValue, err := goinp.AskForString(
			`What's the name of the Project you use for "Archive for Publishing" (e.g.: MyProject.iOS)?`,
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
			`What's the name of the Configuration you use for "Archive for Publishing"?
Specify it if it's not "Release|iPhone", or hit Enter if it is`,
			"Release|iPhone",
		)
		if err != nil {
			return printXamarinScanFinishedWithError("Failed to read input: %s", err)
		}
		xamarinCmd.ConfigurationName = answerValue
	}

	fmt.Println()
	fmt.Println()
	log.Println(`🔦  Running a Build, to get all the required code signing settings...`)
	codeSigningSettings, logOutput, err := xamarinCmd.ScanCodeSigningSettings()
	// save the xamarin output into a debug log file
	logOutputFilePath := filepath.Join(absExportOutputDirPath, "xamarin-build-output.log")
	{
		log.Infof("  💡  "+colorstring.Yellow("Saving xamarin output into file")+": %s", logOutputFilePath)
		if logWriteErr := fileutil.WriteStringToFile(logOutputFilePath, logOutput); logWriteErr != nil {
			log.Errorf("Failed to save xamarin build output into file (%s), error: %s", logOutputFilePath, logWriteErr)
		} else if err != nil {
			log.Infoln(colorstring.Yellow("Please check the logfile (" + logOutputFilePath + ") to see what caused the error"))
			log.Infoln(colorstring.Red(`and make sure that you can "Archive for Publishing" this project from Xamarin!`))
			fmt.Println()
			log.Infoln("Open the project: ", xamarinCmd.SolutionFilePath)
			log.Infoln(`And do "Archive for Publishing", after selecting the Configuration: `, xamarinCmd.ConfigurationName)
			fmt.Println()
		}
	}
	if err != nil {
		return printXamarinScanFinishedWithError("Failed to detect code signing settings: %s", err)
	}
	log.Debugf("codeSigningSettings: %#v", codeSigningSettings)

	return exportCodeSigningFiles(absExportOutputDirPath, codeSigningSettings)
}
