package xcodeuitest

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/bitrise-tools/xcode-project/xcodeproj"
	"github.com/bitrise-tools/xcode-project/xcscheme"
	"github.com/bitrise-tools/xcode-project/xcworkspace"
)

// CommandModel ...
type CommandModel struct {
	// --- Required ---

	// ProjectFilePath - might be a `xcodeproj` or `xcworkspace`
	ProjectFilePath string

	// --- Optional ---

	// Scheme will be passed to xcodebuild as the -scheme flag's value
	// Only passed to xcodebuild if not empty!
	Scheme string

	// CodeSignIdentity will be passed to xcodebuild as an CODE_SIGN_IDENTITY= argument.
	// Only passed to xcodebuild if not empty!
	CodeSignIdentity string

	// SDK: if defined it'll be passed as the -sdk flag to xcodebuild.
	// For more info about the possible values please see xcodebuild's docs about the -sdk flag.
	// Only passed to xcodebuild if not empty!
	SDK string
}

// GenerateArchive : generates the archive for subsequent "Scan"
// func (xcuicmd CommandModel) GenerateArchive() (string, string, error) {
// 	xcoutput := ""
// 	var err error

// 	tmpDir, err := pathutil.NormalizedOSTempDirPath("__codesigndoc__")
// 	if err != nil {
// 		return "", "", fmt.Errorf("failed to create temp dir for archives, error: %s", err)
// 	}
// 	tmpArchivePath := filepath.Join(tmpDir, xcuicmd.Scheme+".xcarchive")

// 	progress.SimpleProgress(".", 1*time.Second, func() {
// 		xcoutput, err = xcuicmd.RunXcodebuildCommand("clean", "archive", "-archivePath", tmpArchivePath)
// 	})
// 	fmt.Println()

// 	if err != nil {
// 		return "", xcoutput, err
// 	}
// 	return tmpArchivePath, xcoutput, nil
// }

// ScanSchemes TODO
func (xcuitestcmd CommandModel) ScanSchemes() (schemes []xcscheme.Scheme, schemesWitUITests []xcscheme.Scheme, schemeNames []string, schemesWitUITestNames []string, error error) {
	if xcworkspace.IsWorkspace(xcuitestcmd.ProjectFilePath) {
		workspace, err := xcworkspace.Open(xcuitestcmd.ProjectFilePath)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("Failed to open workspace (%s), error: %s", xcuitestcmd.ProjectFilePath, err)
		}

		schemesByContainer, err := workspace.Schemes()
		if err != nil {
			return nil, nil, nil, nil, err
		}

		// Remove Cocoapod schemes
		for container, containerSchemes := range schemesByContainer {
			if strings.ToLower(path.Base(container)) != "pods.xcodeproj" {
				schemes = append(schemes, containerSchemes...)
			}
		}
	} else {
		proj, err := xcodeproj.Open(xcuitestcmd.ProjectFilePath)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("Failed to open project (%s), error: %s", xcuitestcmd.ProjectFilePath, err)
		}

		schemes, err = proj.Schemes()
		if err != nil {
			return
		}
	}

	// Check every scheme if it has UITest target in the .pbxproj file or not.
	{
		var proj xcodeproj.Proj
		for _, scheme := range schemes {
			xcproj, _, err := findBuiltProject(xcuitestcmd.ProjectFilePath, scheme.Name, "")
			if err != nil {
				continue
			}

			proj = xcproj.Proj
			if schemesHasUITest(scheme, proj) {
				schemesWitUITests = append(schemesWitUITests, scheme)
			}

		}
	}

	// Iterate trough the scheme arrays and get the scheme names
	{
		for _, scheme := range schemes {
			schemeNames = append(schemeNames, scheme.Name)
		}

		for _, schemeWithUITest := range schemesWitUITests {
			schemesWitUITestNames = append(schemesWitUITestNames, schemeWithUITest.Name)
		}
	}

	return
}

func schemesHasUITest(scheme xcscheme.Scheme, proj xcodeproj.Proj) bool {
	var testEntries []xcscheme.BuildActionEntry

	for _, entry := range scheme.BuildAction.BuildActionEntries {
		if entry.BuildForTesting != "YES" || !strings.HasSuffix(entry.BuildableReference.BuildableName, ".xctest") {
			continue
		}
		testEntries = append(testEntries, entry)
	}

	for _, entry := range testEntries {
		for _, target := range proj.Targets {
			if target.ID == entry.BuildableReference.BlueprintIdentifier {
				if strings.HasSuffix(target.ProductType, "ui-testing") {
					return true
				}
			}
		}
	}
	return false

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

	return project, scheme.Name, nil
}
