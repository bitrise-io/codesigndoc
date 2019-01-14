package cmd

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/xcode"
	"github.com/bitrise-tools/go-xcode/utility"
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
	// absExportOutputDirPath, err := initExportOutputDir()
	// if err != nil {
	// 	return fmt.Errorf("failed to prepare Export directory: %s", err)
	// }

	// Output tools versions
	xcodebuildVersion, err := utility.GetXcodeVersion()
	if err != nil {
		return fmt.Errorf("failed to get Xcode (xcodebuild) version, error: %s", err)
	}
	fmt.Println()
	log.Infof("%s: %s (%s)", colorstring.Green("Xcode (xcodebuild) version"), xcodebuildVersion.Version, xcodebuildVersion.BuildVersion)
	fmt.Println()

	// xcodebuildOutput := ""
	xcodeUITestsCmd := xcode.CommandModel{}

	projectPath := paramXcodeProjectFilePath
	if projectPath == "" {
		askText := `Please drag-and-drop your Xcode Project (` + colorstring.Green(".xcodeproj") + `) or Workspace (` + colorstring.Green(".xcworkspace") + `) file, 
the one you usually open in Xcode, then hit Enter.
(Note: if you have a Workspace file you should most likely use that)`
		projpth, err := goinp.AskForPath(askText)
		if err != nil {
			return fmt.Errorf("failed to read input: %s", err)
		}

		projectPath = strings.Trim(strings.TrimSpace(projpth), "'\"")
	}
	log.Debugf("projectPath: %s", projectPath)
	xcodeUITestsCmd.ProjectFilePath = projectPath

	schemeToUse := paramXcodeScheme
	if schemeToUse == "" {
		fmt.Println()
		log.Printf("🔦  Scanning Schemes ...")
		schemes, err := xcodeUITestsCmd.ScanSchemes()
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
	xcodeUITestsCmd.Scheme = schemeToUse

	if paramXcodebuildSDK != "" {
		xcodeUITestsCmd.SDK = paramXcodebuildSDK
	}

	return nil
}
