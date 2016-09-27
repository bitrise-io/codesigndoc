package cmd

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/xamarin"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/project"
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

	xamSln, err := solution.New(xamarinCmd.SolutionFilePath, true)
	if err != nil {
		return printXamarinScanFinishedWithError("Failed to analyze Xamarin solution: %s", err)
	}
	log.Debugf("xamSln: %#v", xamSln)
	// filter only the iOS "app"" projects
	xamarinProjectsToChooseFrom := []project.Model{}
	for _, aXamarinProject := range xamSln.ProjectMap {
		switch aXamarinProject.ProjectType {
		case constants.ProjectTypeIos, constants.ProjectTypeTVOs, constants.ProjectTypeMac:
			if aXamarinProject.OutputType == "exe" {
				// possible project
				xamarinProjectsToChooseFrom = append(xamarinProjectsToChooseFrom, aXamarinProject)
			}
		default:
			continue
		}
	}
	log.Debugf("len(xamarinProjectsToChooseFrom): %#v", len(xamarinProjectsToChooseFrom))
	log.Debugf("xamarinProjectsToChooseFrom: %#v", xamarinProjectsToChooseFrom)

	// Xamarin Project
	selectedXamarinProject := project.Model{}
	{
		if xamarinProjectName != "" {
			// project specified via flag/param
			for _, aProj := range xamarinProjectsToChooseFrom {
				if xamarinProjectName == aProj.Name {
					selectedXamarinProject = aProj
					break
				}
			}
			if selectedXamarinProject.Name == "" {
				return printXamarinScanFinishedWithError(
					"Invalid Project specified (%s), either not found in the provided Solution or it can't be used for iOS Archive.",
					xamarinProjectName)
			}
		} else {
			// no project specified
			if len(xamarinProjectsToChooseFrom) == 1 {
				selectedXamarinProject = xamarinProjectsToChooseFrom[0]
			} else {
				projectNames := []string{}
				for _, aProj := range xamarinProjectsToChooseFrom {
					projectNames = append(projectNames, aProj.Name)
				}
				fmt.Println()
				answerValue, err := goinp.SelectFromStrings(
					`Select the Project Name you use for "Archive for Publishing" (usually ends with ".iOS", e.g.: MyProject.iOS)?`,
					projectNames,
				)
				if err != nil {
					return printXamarinScanFinishedWithError("Failed to select Project: %s", err)
				}
				log.Debugf("selected project: %v", answerValue)
				for _, aProj := range xamarinProjectsToChooseFrom {
					if answerValue == aProj.Name {
						selectedXamarinProject = aProj
						break
					}
				}
			}
		}
	}
	xamarinCmd.ProjectName = selectedXamarinProject.Name
	log.Debugf("xamarinCmd.ProjectName: %s", xamarinCmd.ProjectName)

	// Xamarin Configuration Name
	xamarinCmd.ConfigurationName = xamarinConfigurationName
	if xamarinCmd.ConfigurationName == "" {
		if selectedXamarinProject.Name == xamarinCmd.ProjectName {
			// print only the Project's Configs
		} else {
			// print all configs from the Solution, as the project was manually specified
		}
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
