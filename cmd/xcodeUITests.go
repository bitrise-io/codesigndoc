package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-io/codesigndoc/codesigndocuitests"
	"github.com/bitrise-io/codesigndoc/xcodeuitest"
	"github.com/bitrise-io/go-xcode/utility"
	"github.com/spf13/cobra"
)

var xcodeUITestsCmd = &cobra.Command{
	Use:   "xcodeuitests",
	Short: "Xcode project scanner for UI tests",
	Long:  `Scan an Xcode project for UI test targets`,

	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          scanXcodeUITestsProject,
}

func init() {
	scanCmd.AddCommand(xcodeUITestsCmd)

	xcodeUITestsCmd.Flags().StringVar(&paramXcodeProjectFilePath, "file", "", "Xcode Project/Workspace file path")
	xcodeUITestsCmd.Flags().StringVar(&paramXcodeScheme, "scheme", "", "Xcode Scheme")
	xcodeUITestsCmd.Flags().StringVar(&paramXcodebuildSDK, "xcodebuild-sdk", "", "xcodebuild -sdk param. If a value is specified for this flag it'll be passed to xcodebuild as the value of the -sdk flag. For more info about the values please see xcodebuild's -sdk flag docs. Example value: iphoneos")
}

func scanXcodeUITestsProject(cmd *cobra.Command, args []string) error {
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

	projectPath := paramXcodeProjectFilePath
	if projectPath == "" {
		log.Infof("Scan the directory for project files")
		log.Warnf("You can specify the Xcode project/workscape file to scan with the --file flag.")

		//
		// Scan the directory for Xcode Project (.xcworkspace / .xcodeproject) file first
		// If can't find any, ask the user to drag-and-drop the file
		projpth, err := findXcodeProject()
		if err != nil {
			return err
		}

		projectPath = strings.Trim(strings.TrimSpace(projpth), "'\"")
	}
	log.Debugf("projectPath: %s", projectPath)
	xcodeUITestsCmd := xcodeuitest.CommandModel{ProjectFilePath: projectPath}

	schemeToUse := paramXcodeScheme
	if schemeToUse == "" {
		fmt.Println()
		log.Printf("🔦  Scanning Schemes ...")

		schemes, schemesWitUITests, err := xcodeUITestsCmd.ScanSchemes()
		if err != nil {
			return fmt.Errorf("failed to scan schemes, error: %s", err)
		}

		log.Debugf("schemes: %v", schemes)

		if len(schemesWitUITests) == 0 {
			return BuildForTestingError{toolXcode, "no schemes found with UITest target enabled:"}
		} else if len(schemesWitUITests) == 1 {
			log.Infof("Only one scheme found with UITest target enabled:")
			log.Printf(schemesWitUITests[0].Name)
			schemeToUse = schemesWitUITests[0].Name
		} else {
			fmt.Println()
			log.Infof("Schemes with UITest target enabled:")

			// Iterate trough the scheme arrays and get the scheme names
			var schemesWitUITestNames []string
			{
				for _, schemeWithUITest := range schemesWitUITests {
					schemesWitUITestNames = append(schemesWitUITestNames, schemeWithUITest.Name)
				}
			}

			selectedScheme, err := goinp.SelectFromStringsWithDefault("Select the Scheme you usually use in Xcode", 1, schemesWitUITestNames)
			if err != nil {
				return fmt.Errorf("failed to select Scheme: %s", err)
			}
			schemeToUse = selectedScheme
		}

		log.Debugf("selected scheme: %v", schemeToUse)
	}
	xcodeUITestsCmd.Scheme = schemeToUse

	if paramXcodebuildSDK != "" {
		xcodeUITestsCmd.SDK = paramXcodebuildSDK
	}

	fmt.Println()
	log.Printf("🔦  Running an Xcode build-for-testing, to get all the required code signing settings...")
	buildForTestingPath, buildLog, err := xcodeUITestsCmd.RunBuildForTesting()
	xcodebuildOutput := buildLog
	// save the xcodebuild output into a debug log file
	xcodebuildOutputFilePath := filepath.Join(absExportOutputDirPath, "xcodebuild-output.log")
	{
		log.Infof("💡  "+colorstring.Yellow("Saving xcodebuild output into file")+": %s", xcodebuildOutputFilePath)
		if logWriteErr := fileutil.WriteStringToFile(xcodebuildOutputFilePath, xcodebuildOutput); logWriteErr != nil {
			log.Errorf("Failed to save xcodebuild output into file (%s), error: %s", xcodebuildOutputFilePath, logWriteErr)
		} else if err != nil {
			log.Warnf("Please check the logfile (%s) to see what caused the error", xcodebuildOutputFilePath)
			log.Warnf("and make sure that you can run Build for testing against the project from Xcode!")
			fmt.Println()
			log.Printf("Open the project: %s", xcodeUITestsCmd.ProjectFilePath)
			fmt.Println()
		}
	}
	if err != nil {
		return BuildForTestingError{toolXcode, err.Error()}
	}

	certsUploaded, provProfilesUploaded, err := codesigndocuitests.ExportCodesignFiles(buildForTestingPath, absExportOutputDirPath, certificatesOnly, isAskForPassword)
	if err != nil {
		return err
	}

	printFinished(provProfilesUploaded, certsUploaded)
	return nil
}
