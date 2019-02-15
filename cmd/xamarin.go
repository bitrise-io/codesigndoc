package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/codesigndoc"
	"github.com/bitrise-tools/codesigndoc/xamarin"
	"github.com/bitrise-tools/go-xamarin/analyzers/project"
	"github.com/bitrise-tools/go-xamarin/analyzers/solution"
	"github.com/bitrise-tools/go-xamarin/builder"
	"github.com/bitrise-tools/go-xamarin/constants"
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
	paramXamarinSolutionFilePath  = ""
	paramXamarinConfigurationName = ""
)

func init() {
	scanCmd.AddCommand(xamarinCmd)

	xamarinCmd.Flags().StringVar(&paramXamarinSolutionFilePath, "file", "", `Xamarin Solution file path`)
	xamarinCmd.Flags().StringVar(&paramXamarinConfigurationName, "config", "", `Xamarin Configuration Name (e.g.: "Release|iPhone")`)
}

func archivableSolutionConfigNames(projectsByID map[string]project.Model) []string {
	archivableSolutionConfigNameSet := map[string]bool{}
	for _, project := range projectsByID {
		if project.SDK != constants.SDKIOS || project.OutputType != "exe" {
			continue
		}

		var archivableProjectConfigNames []string
		for name, config := range project.Configs {
			if builder.IsDeviceArch(config.MtouchArchs...) {
				archivableProjectConfigNames = append(archivableProjectConfigNames, name)
			}

		}

		for solutionConfigName, projectConfigName := range project.ConfigMap {
			for _, archivableProjectConfigName := range archivableProjectConfigNames {
				if archivableProjectConfigName == projectConfigName {
					archivableSolutionConfigNameSet[solutionConfigName] = true
				}
			}
		}
	}

	archivableSolutionConfigNames := []string{}
	for configName := range archivableSolutionConfigNameSet {
		archivableSolutionConfigNames = append(archivableSolutionConfigNames, configName)
	}
	sort.Strings(archivableSolutionConfigNames)

	return archivableSolutionConfigNames
}

func scanXamarinProject(cmd *cobra.Command, args []string) error {
	absExportOutputDirPath, err := initExportOutputDir()
	if err != nil {
		return fmt.Errorf("failed to prepare Export directory: %s", err)
	}

	xamarinCmd := xamarin.CommandModel{}
	// --- Inputs ---

	// Xamarin Solution Path
	xamarinCmd.SolutionFilePath = paramXamarinSolutionFilePath
	if xamarinCmd.SolutionFilePath == "" {
		fmt.Println()
		log.Infof("Scan the directory for solution files")
		log.Warnf("You can specify the Xamarin Solution file to scan with the --file flag.")

		//
		// Scan the directory for Xamarin.Solution file first
		// If can't find any, ask the user to drag-and-drop the file
		xamarinCmd.SolutionFilePath, err = findXamarinSolution()
		if err != nil {
			return err
		}
	}
	log.Debugf("xamSolutionPth: %s", xamarinCmd.SolutionFilePath)

	xamSln, err := solution.New(xamarinCmd.SolutionFilePath, true)
	if err != nil {
		return fmt.Errorf("failed to analyze Xamarin solution: %s", err)
	}

	if enableVerboseLog {
		b, err := json.MarshalIndent(xamSln, "", "\t")
		if err == nil {
			log.Debugf("xamarin solution:\n%s", b)
		}
	}

	archivableSolutionConfigNames := archivableSolutionConfigNames(xamSln.ProjectMap)

	if len(archivableSolutionConfigNames) < 1 {
		return ArchiveError{toolXamarin, `no acceptable Configuration found in the provided Solution and Project, or none can be used for iOS "Archive for Publishing".`}
	}

	// Xamarin Configuration Name
	selectedXamarinConfigurationName := ""
	{
		if paramXamarinConfigurationName != "" {
			// configuration specified via flag/param
			for _, configName := range archivableSolutionConfigNames {
				if paramXamarinConfigurationName == configName {
					selectedXamarinConfigurationName = configName
					break
				}
			}
			if selectedXamarinConfigurationName == "" {
				return ArchiveError{toolXamarin, fmt.Sprintf("invalid Configuration specified (%s), either not found in the provided Solution and Project or it can't be used for iOS Archive.", paramXamarinConfigurationName)}
			}
		} else {
			// no configuration CLI param specified
			if len(archivableSolutionConfigNames) == 1 {
				selectedXamarinConfigurationName = archivableSolutionConfigNames[0]
			} else {
				fmt.Println()
				answerValue, err := goinp.SelectFromStringsWithDefault(`Select the Configuration Name you use for "Archive for Publishing" (usually Release|iPhone)`, 1, archivableSolutionConfigNames)
				if err != nil {
					return fmt.Errorf("failed to select Configuration: %s", err)
				}
				log.Debugf("selected configuration: %v", answerValue)
				selectedXamarinConfigurationName = answerValue
			}
		}
	}
	if selectedXamarinConfigurationName == "" {
		return ArchiveError{toolXamarin, `no acceptable Configuration found (it was empty) in the provided Solution and Project, or none can be used for iOS "Archive for Publishing".`}
	}
	if err := xamarinCmd.SetConfigurationPlatformCombination(selectedXamarinConfigurationName); err != nil {
		return fmt.Errorf("failed to set Configuration Platform combination for the command, error: %s", err)
	}

	fmt.Println()
	fmt.Println()
	log.Printf(`🔦  Running a Build, to get all the required code signing settings...`)
	archivePath, xamLogOut, err := xamarinCmd.GenerateArchive()
	logOutput := xamLogOut
	// save the xamarin output into a debug log file
	logOutputFilePath := filepath.Join(absExportOutputDirPath, "xamarin-build-output.log")
	{
		log.Infof("💡  "+colorstring.Yellow("Saving xamarin output into file")+": %s", logOutputFilePath)
		if logWriteErr := fileutil.WriteStringToFile(logOutputFilePath, logOutput); logWriteErr != nil {
			log.Errorf("Failed to save xamarin build output into file (%s), error: %s", logOutputFilePath, logWriteErr)
		} else if err != nil {
			log.Warnf("Please check the logfile (%s) to see what caused the error", logOutputFilePath)
			log.Warnf(`and make sure that you can "Archive for Publishing" this project from Xamarin!`)
			fmt.Println()
			log.Infof("Open the project: %s", xamarinCmd.SolutionFilePath)
			log.Infof(`And do "Archive for Publishing", after selecting the Configuration+Platform: %s|%s`, xamarinCmd.Configuration, xamarinCmd.Platform)
			fmt.Println()
		}
	}
	if err != nil {
		return ArchiveError{toolXamarin, "failed to run xamarin build command: " + err.Error()}
	}

	certsUploaded, provProfilesUploaded, err := codesigndoc.ExportCodesignFiles(archivePath, absExportOutputDirPath, certificatesOnly, isAskForPassword)
	if err != nil {
		return err
	}

	printFinished(provProfilesUploaded, certsUploaded)
	return nil
}
