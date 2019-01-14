package cmd

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/bitrise-tools/xcode-project/xcscheme"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/xcode"
	"github.com/bitrise-tools/go-xcode/utility"
	"github.com/bitrise-tools/xcode-project/xcodeproj"
	"github.com/bitrise-tools/xcode-project/xcworkspace"
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

		schemes, schemeNames, err := scanSchemes(projectPath)
		if err != nil {
			return fmt.Errorf("failed to scan schemes, error: %s", err)
		}

		log.Debugf("schemes: %v", schemes)

		if len(schemes) == 0 {
			return ArchiveError{toolXcode, "no schemes found"}
		} else if len(schemes) == 1 {
			schemeToUse = schemes[0].Name
		} else {
			fmt.Println()
			selectedScheme, err := goinp.SelectFromStringsWithDefault("Select the Scheme you usually use in Xcode", 1, schemeNames)
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

func filterScheme(schemes []xcscheme.Scheme) []xcscheme.Scheme {
	var schemesWithTest []xcscheme.Scheme
	for _, scheme := range schemes {
		for _, e := range scheme.BuildAction.BuildActionEntries {
			if e.BuildForTesting == "YES" {
				schemesWithTest = append(schemesWithTest, scheme)
				break
			}
		}

	}
	return schemesWithTest
}

func findBuiltProject(pth, schemeName, configurationName string) (xcodeproj.XcodeProj, string, error) {
	var scheme xcscheme.Scheme
	var schemeContainerDir string

	if xcodeproj.IsXcodeProj(pth) {
		project, err := xcodeproj.Open(pth)
		if err != nil {
			return xcodeproj.XcodeProj{}, "", err
		}

		var ok bool
		scheme, ok = project.Scheme(schemeName)
		if !ok {
			return xcodeproj.XcodeProj{}, "", fmt.Errorf("no scheme found with name: %s in project: %s", schemeName, pth)
		}
		schemeContainerDir = filepath.Dir(pth)
	} else if xcworkspace.IsWorkspace(pth) {
		workspace, err := xcworkspace.Open(pth)
		if err != nil {
			return xcodeproj.XcodeProj{}, "", err
		}

		var containerProject string
		scheme, containerProject, err = workspace.Scheme(schemeName)
		if err != nil {
			return xcodeproj.XcodeProj{}, "", fmt.Errorf("no scheme found with name: %s in workspace: %s", schemeName, pth)
		}
		schemeContainerDir = filepath.Dir(containerProject)
	} else {
		return xcodeproj.XcodeProj{}, "", fmt.Errorf("unknown project extension: %s", filepath.Ext(pth))
	}

	if configurationName == "" {
		configurationName = scheme.TestAction.BuildConfiguration
	}

	if configurationName == "" {
		return xcodeproj.XcodeProj{}, "", fmt.Errorf("no configuration provided nor default defined for the scheme's (%s) archive action", schemeName)
	}

	var testEntry xcscheme.BuildActionEntry
	log.Warnf("scheme.BuildAction.BuildActionEntries: %+v", scheme.BuildAction.BuildActionEntries)
	for _, entry := range scheme.BuildAction.BuildActionEntries {
		if entry.BuildForTesting != "YES" || !strings.HasSuffix(entry.BuildableReference.BuildableName, ".xctest") {
			continue
		}
		testEntry = entry
		break
	}

	if testEntry.BuildableReference.BlueprintIdentifier == "" {
		return xcodeproj.XcodeProj{}, "", fmt.Errorf("archivable entry not found")
	}

	projectPth, err := testEntry.BuildableReference.ReferencedContainerAbsPath(schemeContainerDir)
	if err != nil {
		return xcodeproj.XcodeProj{}, "", err
	}

	project, err := xcodeproj.Open(projectPth)
	if err != nil {
		return xcodeproj.XcodeProj{}, "", err
	}

	project.Proj.Target()

	log.Warnf("Project: %v", project)
	return project, scheme.Name, nil
}

func scanSchemes(projectPath string) ([]xcscheme.Scheme, []string, error) {
	var schemes []xcscheme.Scheme
	if xcworkspace.IsWorkspace(projectPath) {
		workspace, err := xcworkspace.Open(projectPath)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to open workspace (%s), error: %s", projectPath, err)
		}

		schemesByContainer, err := workspace.Schemes()
		if err != nil {
			return nil, nil, ArchiveError{toolXcode, "failed to scan Schemes: " + err.Error()}
		}

		// Remove Cocoapod schemes
		for container, containerSchemes := range schemesByContainer {
			log.Warnf(container)
			if strings.ToLower(path.Base(container)) != "pods.xcodeproj" {
				schemes = append(schemes, containerSchemes...)
			}
		}
	} else {
		proj, err := xcodeproj.Open(projectPath)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to open project (%s), error: %s", projectPath, err)
		}

		schemes, err = proj.Schemes()
		if err != nil {
			return nil, nil, err
		}
	}

	schemes = filterScheme(schemes)

	for _, scheme := range schemes {
		_, _, err := findBuiltProject(projectPath, scheme.Name, "")
		if err != nil {
			return nil, nil, err
		}
	}

	var schemeNames []string
	for _, scheme := range schemes {
		schemeNames = append(schemeNames, scheme.Name)
	}

	return schemes, schemeNames, nil
}
