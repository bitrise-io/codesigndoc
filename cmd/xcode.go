package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/codesigndoc/codesigndoc"
	"github.com/bitrise-io/codesigndoc/utility"
	"github.com/bitrise-io/codesigndoc/xcode"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcworkspace"
	"github.com/bitrise-io/goinp/goinp"
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
	paramXcodeProjectFilePath string
	paramXcodeScheme          string
	paramXcodebuildSDK        string
	paramXcodeDestination     string
)

func init() {
	scanCmd.AddCommand(xcodeCmd)

	xcodeCmd.Flags().StringVar(&paramXcodeProjectFilePath, "file", "", "Xcode Project/Workspace file path")
	xcodeCmd.Flags().StringVar(&paramXcodeScheme, "scheme", "", "Xcode Scheme")
	xcodeCmd.Flags().StringVar(&paramXcodebuildSDK, "xcodebuild-sdk", "", "xcodebuild -sdk param. If a value is specified for this flag it'll be passed to xcodebuild as the value of the -sdk flag. For more info about the values please see xcodebuild's -sdk flag docs. Example value: iphoneos")
	xcodeCmd.Flags().StringVar(&paramXcodeDestination, "xcodebuild-destination", "", "The -destination option takes as its argument a destination specifier describing the device (or devices) to use as a destination. If a value is specified for this flag it'll be passed to xcodebuild.")
}

func absOutputDir() (string, error) {
	confExportOutputDirPath := "./codesigndoc_exports"
	absExportOutputDirPath, err := pathutil.AbsPath(confExportOutputDirPath)
	log.Debugf("absExportOutputDirPath: %s", absExportOutputDirPath)
	if err != nil {
		return absExportOutputDirPath, fmt.Errorf("Failed to determine absolute path of export dir: %s", confExportOutputDirPath)
	}
	return absExportOutputDirPath, nil
}

func scanXcodeProject(_ *cobra.Command, _ []string) error {
	absExportOutputDirPath, err := absOutputDir()
	if err != nil {
		return err
	}

	xcodeCmd := xcode.CommandModel{}

	projectPath := paramXcodeProjectFilePath
	if projectPath == "" {
		log.Infof("Scan the directory for project files")
		log.Warnf("You can specify the Xcode project/workscape file to scan with the --file flag.")

		// Scan the directory for Xcode Project (.xcworkspace / .xcodeproject) file first
		// If can't find any, ask the user to drag-and-drop the file
		projpth, err := findXcodeProject()
		if err != nil {
			return err
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

	if paramXcodeDestination != "" {
		xcodeCmd.DESTINATION = paramXcodeDestination
	} else {
		var project xcodeproj.XcodeProj
		var scheme xcscheme.Scheme

		if xcodeproj.IsXcodeProj(xcodeCmd.ProjectFilePath) {
			proj, err := xcodeproj.Open(xcodeCmd.ProjectFilePath)
			if err != nil {
				return fmt.Errorf("Failed to open project (%s), error: %s", xcodeCmd.ProjectFilePath, err)
			}

			projectScheme, _, err := proj.Scheme(xcodeCmd.Scheme)

			if err != nil {
				return fmt.Errorf("failed to find scheme (%s) in project (%s), error: %s", xcodeCmd.Scheme, proj.Path, err)
			}

			project = proj
			scheme = *projectScheme
		} else {
			workspace, err := xcworkspace.Open(xcodeCmd.ProjectFilePath)
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
				projectScheme, _, _ := possibleProject.Scheme(xcodeCmd.Scheme)

				if projectScheme != nil {
					project = possibleProject
					scheme = *projectScheme

					break
				}
			}
		}

		platform, err := utility.BuildableTargetPlatform(&project, &scheme, "", utility.XcodeBuild{})
		if err == nil {
			destination := "generic/platform=" + string(platform)

			xcodeCmd.DESTINATION = destination

			fmt.Print("Setting -destination flag to: ", destination)
		}
	}

	writeBuildLogs := func(xcodebuildOutput string) error {
		if writeFiles == codesign.WriteFilesAlways || writeFiles == codesign.WriteFilesFallback && err != nil { // save the xcodebuild output into a debug log file
			xcodebuildOutputFilePath := filepath.Join(absExportOutputDirPath, "xcodebuild-output.log")
			if err := os.MkdirAll(absExportOutputDirPath, 0700); err != nil {
				return fmt.Errorf("failed to create output directory, error: %s", err)
			}

			log.Infof("💡  "+colorstring.Yellow("Saving xcodebuild output into file")+": %s", xcodebuildOutputFilePath)
			if err := fileutil.WriteStringToFile(xcodebuildOutputFilePath, xcodebuildOutput); err != nil {
				return fmt.Errorf("Failed to save xcodebuild output into file (%s), error: %s", xcodebuildOutputFilePath, err)
			}
		}
		return nil
	}

	archivePath, err := codesigndoc.BuildXcodeArchive(xcodeCmd, writeBuildLogs)
	if err != nil {
		return ArchiveError{toolXcode, err.Error()}
	}

	certificates, profiles, err := codesigndoc.CodesigningFilesForXCodeProject(archivePath, certificatesOnly, isAskForPassword)
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
