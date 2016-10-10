package cmd

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/xcode"
	"github.com/spf13/cobra"
)

// xcodeCmd represents the xcode command
var xcodeCmd = &cobra.Command{
	Use:   "xcode",
	Short: "Xcode project scanner",
	Long:  `Scan an Xcode project`,

	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          scanXcodeProject,
}

var (
	paramXcodeProjectFilePath        = ""
	paramXcodeScheme                 = ""
	paramXcodebuildOutputLogFilePath = ""
)

func init() {
	scanCmd.AddCommand(xcodeCmd)

	xcodeCmd.Flags().StringVar(&paramXcodeProjectFilePath,
		"file", "",
		"Xcode Project/Workspace file path")
	xcodeCmd.Flags().StringVar(&paramXcodeScheme,
		"scheme", "",
		"Xcode Scheme")
	xcodeCmd.Flags().StringVar(&paramXcodebuildOutputLogFilePath,
		"xcodebuild-log", "",
		"xcodebuild output log (file path). If specified it will be used instead of running xcodebuild")
}

func printXcodeScanFinishedWithError(format string, args ...interface{}) error {
	return printFinishedWithError("Xcode", format, args...)
}

func scanXcodeProject(cmd *cobra.Command, args []string) error {
	absExportOutputDirPath, err := initExportOutputDir()
	if err != nil {
		return printXcodeScanFinishedWithError("Failed to prepare Export directory: %s", err)
	}

	xcodebuildOutput := ""
	xcodeCmd := xcode.CommandModel{}
	if paramXcodebuildOutputLogFilePath != "" {
		xcLog, err := fileutil.ReadStringFromFile(paramXcodebuildOutputLogFilePath)
		if err != nil {
			return printXcodeScanFinishedWithError("Failed to read log from the specified log file, error: %s", err)
		}
		xcodebuildOutput = xcLog
	} else {
		projectPath := paramXcodeProjectFilePath
		if projectPath == "" {
			askText := `Please drag-and-drop your Xcode Project (` + colorstring.Green(".xcodeproj") + `)
   or Workspace (` + colorstring.Green(".xcworkspace") + `) file, the one you usually open in Xcode,
   then hit Enter.

  (Note: if you have a Workspace file you should most likely use that)`
			fmt.Println()
			projpth, err := goinp.AskForPath(askText)
			if err != nil {
				return printXcodeScanFinishedWithError("Failed to read input: %s", err)
			}
			projectPath = projpth
		}
		log.Debugf("projectPath: %s", projectPath)
		xcodeCmd.ProjectFilePath = projectPath

		schemeToUse := paramXcodeScheme
		if schemeToUse == "" {
			fmt.Println()
			fmt.Println()
			log.Println("🔦  Scanning Schemes ...")
			schemes, err := xcodeCmd.ScanSchemes()
			if err != nil {
				return printXcodeScanFinishedWithError("Failed to scan Schemes: %s", err)
			}
			log.Debugf("schemes: %v", schemes)

			fmt.Println()
			selectedScheme, err := goinp.SelectFromStrings("Select the Scheme you usually use in Xcode", schemes)
			if err != nil {
				return printXcodeScanFinishedWithError("Failed to select Scheme: %s", err)
			}
			log.Debugf("selected scheme: %v", selectedScheme)
			schemeToUse = selectedScheme
		}
		xcodeCmd.Scheme = schemeToUse

		fmt.Println()
		fmt.Println()
		log.Println("🔦  Running an Xcode Archive, to get all the required code signing settings...")
		xcLog, err := xcodeCmd.GenerateLog()
		xcodebuildOutput = xcLog
		// save the xcodebuild output into a debug log file
		xcodebuildOutputFilePath := filepath.Join(absExportOutputDirPath, "xcodebuild-output.log")
		{
			log.Infof("  💡  "+colorstring.Yellow("Saving xcodebuild output into file")+": %s", xcodebuildOutputFilePath)
			if logWriteErr := fileutil.WriteStringToFile(xcodebuildOutputFilePath, xcodebuildOutput); logWriteErr != nil {
				log.Errorf("Failed to save xcodebuild output into file (%s), error: %s", xcodebuildOutputFilePath, logWriteErr)
			} else if err != nil {
				log.Infoln(colorstring.Yellow("Please check the logfile (" + xcodebuildOutputFilePath + ") to see what caused the error"))
				log.Infoln(colorstring.Red("and make sure that you can Archive this project from Xcode!"))
				fmt.Println()
				log.Infoln("Open the project:", xcodeCmd.ProjectFilePath)
				log.Infoln("and Archive, using the Scheme:", xcodeCmd.Scheme)
				fmt.Println()
			}
		}
		if err != nil {
			return printXcodeScanFinishedWithError("Failed to run Xcode Archive: %s", err)
		}
	}

	codeSigningSettings, err := xcodeCmd.ScanCodeSigningSettings(xcodebuildOutput)
	if err != nil {
		return printXcodeScanFinishedWithError("Failed to detect code signing settings: %s", err)
	}
	log.Debugf("codeSigningSettings: %#v", codeSigningSettings)

	return exportCodeSigningFiles("Xcode", absExportOutputDirPath, codeSigningSettings)
}
