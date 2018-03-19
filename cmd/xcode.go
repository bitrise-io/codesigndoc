package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/xcode"
	"github.com/bitrise-tools/go-xcode/utility"
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
	paramXcodeProjectFilePath = ""
	paramXcodeScheme          = ""
	paramXcodebuildSDK        = ""
)

func init() {
	scanCmd.AddCommand(xcodeCmd)

	xcodeCmd.Flags().StringVar(&paramXcodeProjectFilePath, "file", "", "Xcode Project/Workspace file path")
	xcodeCmd.Flags().StringVar(&paramXcodeScheme, "scheme", "", "Xcode Scheme")
	xcodeCmd.Flags().StringVar(&paramXcodebuildSDK, "xcodebuild-sdk", "", "xcodebuild -sdk param. If a value is specified for this flag it'll be passed to xcodebuild as the value of the -sdk flag. For more info about the values please see xcodebuild's -sdk flag docs. Example value: iphoneos")
}

func scanXcodeProject(cmd *cobra.Command, args []string) error {
	absExportOutputDirPath, err := initExportOutputDir()
	if err != nil {
		return fmt.Errorf("failed to prepare Export directory: %s", err)
	}

	// Output tools versions
	xcodebuildVersion, err := utility.GetXcodeVersion()
	if err != nil {
		return fmt.Errorf("failed to get Xcode (xcodebuild) version, error: %s", err)
	}
	fmt.Println()
	log.Infof("%s: %s (%s)", colorstring.Green("Xcode (xcodebuild) version"), xcodebuildVersion.Version, xcodebuildVersion.BuildVersion)
	fmt.Println()

	xcodebuildOutput := ""
	xcodeCmd := xcode.CommandModel{}

	projectPath := paramXcodeProjectFilePath
	if projectPath == "" {
		askText := `Please drag-and-drop your Xcode Project (` + colorstring.Green(".xcodeproj") + `) or Workspace (` + colorstring.Green(".xcworkspace") + `) file, 
the one you usually open in Xcode, then hit Enter.

(Note: if you have a Workspace file you should most likely use that)`
		fmt.Println()
		projpth, err := goinp.AskForPath(askText)
		if err != nil {
			return fmt.Errorf("failed to read input: %s", err)
		}

		projectPath = strings.Trim(strings.TrimSpace(projpth), "'\"")
	}
	log.Debugf("projectPath: %s", projectPath)
	xcodeCmd.ProjectFilePath = projectPath

	schemeToUse := paramXcodeScheme
	if schemeToUse == "" {
		fmt.Println()
		log.Printf("🔦  Scanning Schemes ...")
		schemes, err := xcodeCmd.ScanSchemes()
		if err != nil {
			return ArchiveError{toolXcode, "failed to scan Schemes: " + err.Error()}
		}
		log.Debugf("schemes: %v", schemes)

		if len(schemes) == 0 {
			return ArchiveError{toolXcode, "no schemes found"}
		} else if len(schemes) == 1 {
			schemeToUse = schemes[0]
		} else {
			fmt.Println()
			selectedScheme, err := goinp.SelectFromStringsWithDefault("Select the Scheme you usually use in Xcode", 1, schemes)
			if err != nil {
				return fmt.Errorf("failed to select Scheme: %s", err)
			}
			schemeToUse = selectedScheme
		}

		log.Debugf("selected scheme: %v", schemeToUse)
	}
	xcodeCmd.Scheme = schemeToUse

	if paramXcodebuildSDK != "" {
		xcodeCmd.SDK = paramXcodebuildSDK
	}

	fmt.Println()
	fmt.Println()
	log.Printf("🔦  Running an Xcode Archive, to get all the required code signing settings...")
	archivePath, buildLog, err := xcodeCmd.GenerateArchive()
	xcodebuildOutput = buildLog
	// save the xcodebuild output into a debug log file
	xcodebuildOutputFilePath := filepath.Join(absExportOutputDirPath, "xcodebuild-output.log")
	{
		log.Infof("💡  "+colorstring.Yellow("Saving xcodebuild output into file")+": %s", xcodebuildOutputFilePath)
		if logWriteErr := fileutil.WriteStringToFile(xcodebuildOutputFilePath, xcodebuildOutput); logWriteErr != nil {
			log.Errorf("Failed to save xcodebuild output into file (%s), error: %s", xcodebuildOutputFilePath, logWriteErr)
		} else if err != nil {
			log.Warnf("Please check the logfile (%s) to see what caused the error", xcodebuildOutputFilePath)
			log.Warnf("and make sure that you can Archive this project from Xcode!")
			fmt.Println()
			log.Printf("Open the project: %s", xcodeCmd.ProjectFilePath)
			log.Printf("and Archive, using the Scheme: %s", xcodeCmd.Scheme)
			fmt.Println()
		}
	}
	if err != nil {
		return ArchiveError{toolXcode, err.Error()}
	}

	return exportCodesignFiles("Xcode", archivePath, absExportOutputDirPath)
}
