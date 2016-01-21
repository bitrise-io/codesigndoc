package xcode

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
)

// CommandModel ...
type CommandModel struct {
	// ProjectFilePath - might be a `xcodeproj` or `xcworkspace`
	ProjectFilePath string
	Scheme          string
}

func parseSchemesFromXcodeOutput(xcodeOutput string) []string {
	scanner := bufio.NewScanner(strings.NewReader(xcodeOutput))

	foundSchemes := []string{}
	isSchemeDelimiterFound := false
	for scanner.Scan() {
		line := scanner.Text()
		if isSchemeDelimiterFound {
			foundSchemes = append(foundSchemes, strings.TrimSpace(line))
		}
		if regexp.MustCompile(`^[ ]*Schemes:$`).MatchString(line) {
			isSchemeDelimiterFound = true
		}
	}
	return foundSchemes
}

func findCodeSigningSettingsFromXcodeOutput(xcodeOutput string) {
	scanner := bufio.NewScanner(strings.NewReader(xcodeOutput))

	for scanner.Scan() {
		line := scanner.Text()
		if regexp.MustCompile(`^[ ]*Signing Identity: .*`).MatchString(line) {
			fmt.Println("-> line: ", line)
		}
		if regexp.MustCompile(`^[ ]*Provisioning Profile: .*`).MatchString(line) {
			fmt.Println("-> line: ", line)
		}
	}
}

// ScanCodeSigningSettings ...
func (xccmd CommandModel) ScanCodeSigningSettings() error {
	xcoutput, err := xccmd.RunXcodebuildCommand("clean", "archive")
	if err != nil {
		return fmt.Errorf("Failed to Archive: %s | full output: %s", err, xcoutput)
	}
	findCodeSigningSettingsFromXcodeOutput(xcoutput)
	return nil
}

func (xccmd CommandModel) xcodeProjectOrWorkspaceParam() (string, error) {
	if strings.HasSuffix(xccmd.ProjectFilePath, "xcworkspace") {
		return "-workspace", nil
	} else if strings.HasSuffix(xccmd.ProjectFilePath, "xcodeproj") {
		return "-project", nil
	}
	return "", fmt.Errorf("Invalid project/workspace file, the extension should be either .xcworkspace or .xcodeproj ; (file path: %s)", xccmd.ProjectFilePath)
}

func (xccmd CommandModel) transformToXcodebuildParams(xcodebuildActionArgs ...string) ([]string, error) {
	projParam, err := xccmd.xcodeProjectOrWorkspaceParam()
	if err != nil {
		return []string{}, err
	}

	baseArgs := []string{projParam, xccmd.ProjectFilePath}
	if xccmd.Scheme != "" {
		baseArgs = append(baseArgs, "-scheme", xccmd.Scheme)
	}
	return append(baseArgs, xcodebuildActionArgs...), nil
}

// RunXcodebuildCommand ...
func (xccmd CommandModel) RunXcodebuildCommand(xcodebuildActionArgs ...string) (string, error) {
	xcodeCmdParamsToRun, err := xccmd.transformToXcodebuildParams(xcodebuildActionArgs...)
	if err != nil {
		return "", err
	}

	log.Debugf("$ xcodebuild %s", cmdex.PrintableCommandArgs(true, xcodeCmdParamsToRun))
	xcoutput, err := cmdex.RunCommandInDirAndReturnCombinedStdoutAndStderr("", "xcodebuild", xcodeCmdParamsToRun...)
	if err != nil {
		return "", fmt.Errorf("Failed to run 'xcodebuild -list': %s | error: %s", xcoutput, err)
	}

	log.Debugf("xcoutput: %s", xcoutput)
	return xcoutput, nil
}

// ScanSchemes ...
func (xccmd CommandModel) ScanSchemes() ([]string, error) {
	xcoutput, err := xccmd.RunXcodebuildCommand("-list")
	if err != nil {
		return []string{}, err
	}

	parsedSchemes := parseSchemesFromXcodeOutput(xcoutput)
	return parsedSchemes, nil
}
