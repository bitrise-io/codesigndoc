package xcode

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/progress"
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
func (xccmd CommandModel) GenerateArchive() (string, string, error) {
	xcoutput := ""
	var err error

	tmpDir, err := pathutil.NormalizedOSTempDirPath("__codesigndoc__")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp dir for archives, error: %s", err)
	}
	tmpArchivePath := filepath.Join(tmpDir, xccmd.Scheme+".xcarchive")

	progress.SimpleProgress(".", 1*time.Second, func() {
		xcoutput, err = xccmd.RunXcodebuildCommand("clean", "archive", "-archivePath", tmpArchivePath)
	})
	fmt.Println()

	if err != nil {
		return "", xcoutput, err
	}
	return tmpArchivePath, xcoutput, nil
}

func (xccmd CommandModel) xcodeProjectOrWorkspaceParam() (string, error) {
	if strings.HasSuffix(xccmd.ProjectFilePath, "xcworkspace") {
		return "-workspace", nil
	} else if strings.HasSuffix(xccmd.ProjectFilePath, "xcodeproj") {
		return "-project", nil
	}
	return "", fmt.Errorf("invalid project/workspace file, the extension should be either .xcworkspace or .xcodeproj ; (file path: %s)", xccmd.ProjectFilePath)
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

	if xccmd.SDK != "" {
		baseArgs = append(baseArgs, "-sdk", xccmd.SDK)
	}

	if xccmd.CodeSignIdentity != "" {
		baseArgs = append(baseArgs, `CODE_SIGN_IDENTITY=`+xccmd.CodeSignIdentity)
	}
	return append(baseArgs, xcodebuildActionArgs...), nil
}

// RunXcodebuildCommand ...
func (xccmd CommandModel) RunXcodebuildCommand(xcodebuildActionArgs ...string) (string, error) {
	xcodeCmdParamsToRun, err := xccmd.transformToXcodebuildParams(xcodebuildActionArgs...)
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

// ScanSchemes ...
func (xccmd CommandModel) ScanSchemes() ([]string, error) {
	xcoutput, err := xccmd.RunXcodebuildCommand("-list")
	if err != nil {
		return []string{}, fmt.Errorf("error: %s | xcodebuild output: %s", err, xcoutput)
	}

	parsedSchemes := parseSchemesFromXcodeOutput(xcoutput)
	return parsedSchemes, nil
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
		if regexp.MustCompile(`^[[:space:]]*Schemes:$`).MatchString(line) {
			isSchemeDelimiterFound = true
		}
	}
	return foundSchemes
}
