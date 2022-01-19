package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/codesigndoc/codesigndocuitests"
	codesigndocutility "github.com/bitrise-io/codesigndoc/utility"
	"github.com/bitrise-io/codesigndoc/xcodeuitest"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/stringutil"
	"github.com/bitrise-io/go-xcode/utility"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcworkspace"
	"github.com/bitrise-io/goinp/goinp"
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
	xcodeUITestsCmd.Flags().StringVar(&paramXcodeDestination, "xcodebuild-destination", "", "The -destination option takes as its argument a destination specifier describing the device (or devices) to use as a destination. If a value is specified for this flag it'll be passed to xcodebuild.")
}

func scanXcodeUITestsProject(cmd *cobra.Command, args []string) error {
	absExportOutputDirPath, err := absOutputDir()
	if err != nil {
		return err
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

	if paramXcodeDestination != "" {
		xcodeUITestsCmd.DESTINATION = paramXcodeDestination
	} else {
		var project xcodeproj.XcodeProj
		var scheme xcscheme.Scheme

		if xcodeproj.IsXcodeProj(xcodeUITestsCmd.ProjectFilePath) {
			proj, err := xcodeproj.Open(xcodeUITestsCmd.ProjectFilePath)
			if err != nil {
				return fmt.Errorf("Failed to open project (%s), error: %s", xcodeUITestsCmd.ProjectFilePath, err)
			}

			projectScheme, _, err := project.Scheme(xcodeUITestsCmd.Scheme)
			if err != nil {
				return fmt.Errorf("failed to find scheme (%s) in project (%s), error: %s", xcodeUITestsCmd.Scheme, project.Path, err)
			}

			project = proj
			scheme = *projectScheme
		} else {
			workspace, err := xcworkspace.Open(xcodeUITestsCmd.ProjectFilePath)
			if err != nil {
				return err
			}

			projects, err := workspace.ProjectFileLocations()
			if err != nil {
				return err
			}

			for _, projectLocation := range projects {
				if exist, err := pathutil.IsPathExists(projectLocation); err != nil {
					return fmt.Errorf("failed to check if project exist at: %s, error: %s", projectLocation, err)
				} else if !exist {
					// at this point we are interested the schemes visible for the workspace
					continue
				}

				possibleProject, _ := xcodeproj.Open(projectLocation)
				projectScheme, _, _ := possibleProject.Scheme(xcodeUITestsCmd.Scheme)

				if projectScheme != nil {
					project = possibleProject
					scheme = *projectScheme

					break
				}
			}
		}

		platform, err := codesigndocutility.BuildableTargetPlatform(&project, &scheme, "", codesigndocutility.XcodeBuild{})
		if err == nil {
			destination := "generic/platform=" + string(platform)

			xcodeUITestsCmd.DESTINATION = destination

			fmt.Print("Setting -destination flag to: ", destination)
		}
	}

	fmt.Println()
	fmt.Println()
	log.Printf("🔦  Running an Xcode build-for-testing, to get all the required code signing settings...")
	xcodebuildOutputFilePath := filepath.Join(absExportOutputDirPath, "xcodebuild-output.log")

	buildForTestingPath, xcodebuildOutput, err := xcodeUITestsCmd.RunBuildForTesting()
	if writeFiles == codesign.WriteFilesAlways || writeFiles == codesign.WriteFilesFallback && err != nil { // save the xcodebuild output into a debug log file
		if err := os.MkdirAll(absExportOutputDirPath, 0700); err != nil {
			return fmt.Errorf("failed to create output directory, error: %s", err)
		}

		log.Infof("💡  "+colorstring.Yellow("Saving xcodebuild output into file")+": %s", xcodebuildOutputFilePath)
		if err := fileutil.WriteStringToFile(xcodebuildOutputFilePath, xcodebuildOutput); err != nil {
			log.Errorf("Failed to save xcodebuild output into file (%s), error: %s", xcodebuildOutputFilePath, err)
		}
	}
	if err != nil {
		log.Warnf("Last lines of the build log:")
		fmt.Println(stringutil.LastNLines(xcodebuildOutput, 15))

		log.Infof(colorstring.Yellow("Please check the build log to see what caused the error."))
		fmt.Println()

		log.Errorf("Xcode Build For Testing failed.")
		log.Infof(colorstring.Yellow("Open the project: ")+"%s", xcodeUITestsCmd.ProjectFilePath)
		log.Infof(colorstring.Yellow("and make sure that you can run Build For Testing, with the scheme: ")+"%s", xcodeUITestsCmd.Scheme)
		fmt.Println()

		return BuildForTestingError{toolXcode, err.Error()}
	}

	// If certificatesOnly is set, CollectCodesignFiles returns an empty slice for profiles
	certificatesToExport, profilesToExport, err := codesigndocuitests.CollectCodesignFiles(buildForTestingPath, certificatesOnly)
	if err != nil {
		return err
	}

	certificates, profiles, err := codesign.ExportCodesigningFiles(certificatesToExport, profilesToExport, isAskForPassword)
	if err != nil {
		return err
	}

	exportResult, err := codesign.UploadAndWriteCodesignFiles(certificates,
		profiles,
		codesign.WriteFilesConfig{
			WriteFiles:       writeFiles,
			AbsOutputDirPath: absExportOutputDirPath,
		},
		codesign.UploadConfig{
			PersonalAccessToken: personalAccessToken,
			AppSlug:             appSlug,
		})
	if err != nil {
		return err
	}

	printFinished(exportResult, absExportOutputDirPath)
	return nil
}
