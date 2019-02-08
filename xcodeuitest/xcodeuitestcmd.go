package xcodeuitest

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/progress"
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

	// SDK: if defined it'll be passed as the -sdk flag to xcodebuild.
	// For more info about the possible values please see xcodebuild's docs about the -sdk flag.
	// Only passed to xcodebuild if not empty!
	SDK string
}

// RunBuildForTesting runs the build-for-tesing xcode command
func (xcuitestcmd CommandModel) RunBuildForTesting() (string, string, error) {
	xcoutput := ""
	var err error

	tmpDir, err := pathutil.NormalizedOSTempDirPath("__codesigndoc__")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp dir for archives, error: %s", err)
	}
	tmpBuildPath := filepath.Join(tmpDir, xcuitestcmd.Scheme)

	progress.SimpleProgress(".", 1*time.Second, func() {
		xcoutput, err = xcuitestcmd.RunXcodebuildCommand("clean", "build-for-testing", "CONFIGURATION_BUILD_DIR="+tmpBuildPath)
	})
	fmt.Println()

	if err != nil {
		return "", xcoutput, err
	}
	return tmpBuildPath, xcoutput, nil
}

func (xcuitestcmd CommandModel) xcodeProjectOrWorkspaceParam() (string, error) {
	if strings.HasSuffix(xcuitestcmd.ProjectFilePath, "xcworkspace") {
		return "-workspace", nil
	} else if strings.HasSuffix(xcuitestcmd.ProjectFilePath, "xcodeproj") {
		return "-project", nil
	}
	return "", fmt.Errorf("invalid project/workspace file, the extension should be either .xcworkspace or .xcodeproj ; (file path: %s)", xcuitestcmd.ProjectFilePath)
}

func (xcuitestcmd CommandModel) transformToXcodebuildParams(xcodebuildActionArgs ...string) ([]string, error) {
	projParam, err := xcuitestcmd.xcodeProjectOrWorkspaceParam()
	if err != nil {
		return []string{}, err
	}

	baseArgs := []string{projParam, xcuitestcmd.ProjectFilePath}
	if xcuitestcmd.Scheme != "" {
		baseArgs = append(baseArgs, "-scheme", xcuitestcmd.Scheme)
	}

	if xcuitestcmd.SDK != "" {
		baseArgs = append(baseArgs, "-sdk", xcuitestcmd.SDK)
	}

	return append(baseArgs, xcodebuildActionArgs...), nil
}

// RunXcodebuildCommand TODO comment
func (xcuitestcmd CommandModel) RunXcodebuildCommand(xcodebuildActionArgs ...string) (string, error) {
	xcodeCmdParamsToRun, err := xcuitestcmd.transformToXcodebuildParams(xcodebuildActionArgs...)
	if err != nil {
		return "", err
	}

	log.Infof("$ xcodebuild %s", command.PrintableCommandArgs(true, xcodeCmdParamsToRun))
	xcoutput, err := command.RunCommandAndReturnCombinedStdoutAndStderr("xcodebuild", xcodeCmdParamsToRun...)
	if err != nil {
		return xcoutput, fmt.Errorf("failed to run xcodebuild command, error: %s", err)
	}

	log.Debugf("xcoutput: %s", xcoutput)
	return xcoutput, nil
}

// ScanSchemes scans the provided project or workspace for schemes.
// Returns the schemes of the provided project's or the provided workspace's + project's schemes, the names of the schemes,
// the schemes with UITest target and the schemes' which has UITest.
func (xcuitestcmd CommandModel) ScanSchemes() (schemes []xcscheme.Scheme, schemesWitUITests []xcscheme.Scheme, err error) {
	if xcworkspace.IsWorkspace(xcuitestcmd.ProjectFilePath) {
		workspace, err := xcworkspace.Open(xcuitestcmd.ProjectFilePath)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to open workspace (%s), error: %s", xcuitestcmd.ProjectFilePath, err)
		}

		schemesByContainer, err := workspace.Schemes()
		if err != nil {
			return nil, nil, err
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
			return nil, nil, fmt.Errorf("Failed to open project (%s), error: %s", xcuitestcmd.ProjectFilePath, err)
		}

		schemes, err = proj.Schemes()
		if err != nil {
			return nil, nil, err
		}
	}

	// Check every scheme if it has UITest target in the .pbxproj file or not.
	{
		var proj xcodeproj.Proj
		for _, scheme := range schemes {
			xcproj, _, err := findBuiltProject(xcuitestcmd.ProjectFilePath, scheme.Name, "")
			if err != nil {
				log.Debugf("failed to find build project, error: %s", err)
				continue
			}

			log.Printf("Check if the (%s) scheme is valid for UI testing", scheme.Name)
			proj = xcproj.Proj
			if schemesHasUITest(scheme, proj) {
				schemesWitUITests = append(schemesWitUITests, scheme)
			}

		}
	}

	return
}

// Return true if the provided scheme has BuildActionEntry with BuildForTesting = YES and one of this entries has target with *ui-testing product type.
func schemesHasUITest(scheme xcscheme.Scheme, proj xcodeproj.Proj) bool {
	var testables []xcscheme.TestableReference

	for _, testable := range scheme.TestAction.Testables {
		if testable.Skipped == "NO" {
			testables = append(testables, testable)
		}
	}

	if len(testables) == 0 {
		log.Printf("No enabled test targets found in (%s) => Skipping...\n", scheme.Name)
		return false
	}

	log.Printf("Found enabled test target in (%s)", scheme.Name)

	for _, entry := range testables {
		for _, target := range proj.Targets {
			if target.ID == entry.BuildableReference.BlueprintIdentifier {
				if strings.HasSuffix(target.ProductType, "ui-testing") {
					log.Printf("Found enabled test target (%s) => Add it to the list...\n", scheme.Name)
					return true
				}
			}
		}
	}

	log.Printf("No UITest target found in (%s) => Skipping...\n", scheme.Name)
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
			return xcodeproj.XcodeProj{}, "", err
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

	var testEntry xcscheme.TestableReference
	for _, testable := range scheme.TestAction.Testables {
		if testable.Skipped == "NO" {
			testEntry = testable
		}
	}

	if testEntry.BuildableReference.BlueprintIdentifier == "" {
		return xcodeproj.XcodeProj{}, "", fmt.Errorf("testable entry not found")
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
