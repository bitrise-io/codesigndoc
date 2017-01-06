package xcode

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/maputil"
	"github.com/bitrise-io/go-utils/progress"
	"github.com/bitrise-io/go-utils/readerutil"
	"github.com/bitrise-io/go-utils/regexputil"
	"github.com/bitrise-tools/codesigndoc/common"
	"github.com/bitrise-tools/codesigndoc/provprofile"
)

// CommandModel ...
type CommandModel struct {
	// ProjectFilePath - might be a `xcodeproj` or `xcworkspace`
	ProjectFilePath  string
	Scheme           string
	CodeSignIdentity string
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

func parseCodeSigningSettingsFromXcodeOutput(xcodeOutput string) (common.CodeSigningSettings, error) {
	logReader := bufio.NewReader(strings.NewReader(xcodeOutput))

	identitiesMap := map[string]common.CodeSigningIdentityInfo{}
	provProfilesMap := map[string]provprofile.ProvisioningProfileInfo{}
	teamIDsMap := map[string]interface{}{}
	appIDsMap := map[string]interface{}{}

	// scan log line by line
	{
		line, readErr := readerutil.ReadLongLine(logReader)
		for ; readErr == nil; line, readErr = readerutil.ReadLongLine(logReader) {
			// Team ID
			if rexp := regexp.MustCompile(`^[[:space:]]*"com.apple.developer.team-identifier" = (?P<teamid>[a-zA-Z0-9]+);$`); rexp.MatchString(line) {
				results, isFound := regexputil.NamedFindStringSubmatch(rexp, line)
				if !isFound {
					log.Error("Failed to scan TeamID: not found in the logs")
					continue
				}
				teamIDsMap[results["teamid"]] = 1
			}

			// App Bundle ID
			if rexp := regexp.MustCompile(`^[[:space:]]*"application-identifier" = "(?P<appbundleid>.+)";$`); rexp.MatchString(line) {
				results, isFound := regexputil.NamedFindStringSubmatch(rexp, line)
				if !isFound {
					log.Error("Failed to scan App Bundle ID: not found in the logs")
					continue
				}
				appIDsMap[results["appbundleid"]] = 1
			}

			// Signing Identity
			if rexp := regexp.MustCompile(`^[[:space:]]*Signing Identity:[[:space:]]*"(?P<title>.+)"$`); rexp.MatchString(line) {
				results, isFound := regexputil.NamedFindStringSubmatch(rexp, line)
				if !isFound {
					log.Error("Failed to scan Signing Identity title: not found in the logs")
					continue
				}
				codeSigningID := common.CodeSigningIdentityInfo{Title: results["title"]}
				identitiesMap[codeSigningID.Title] = codeSigningID
			}
			// Prov. Profile - title line
			if rexp := regexp.MustCompile(`^[[:space:]]*Provisioning Profile:[[:space:]]*"(?P<title>.+)"$`); rexp.MatchString(line) {
				results, isFound := regexputil.NamedFindStringSubmatch(rexp, line)
				if !isFound {
					log.Error("Failed to scan Provisioning Profile title: not found in the logs")
					continue
				}
				tmpProvProfile := provprofile.ProvisioningProfileInfo{Title: results["title"]}

				// read next line
				line, readErr = readerutil.ReadLongLine(logReader)
				if readErr != nil {
					continue
				}
				if line == "" {
					log.Error("Failed to scan Provisioning Profile UUID: no more lines to scan")
					continue
				}
				provProfileUUIDLine := line

				rexp = regexp.MustCompile(`^[[:space:]]*\((?P<uuid>[a-zA-Z0-9-]{36})\)`)
				results, isFound = regexputil.NamedFindStringSubmatch(rexp, provProfileUUIDLine)
				if !isFound {
					log.Errorf("Failed to scan Provisioning Profile UUID: pattern not found | line was: %s", provProfileUUIDLine)
					continue
				}
				tmpProvProfile.UUID = results["uuid"]
				provProfilesMap[tmpProvProfile.Title] = tmpProvProfile
			}
		}
		if readErr != nil && readErr != io.EOF {
			return common.CodeSigningSettings{}, fmt.Errorf("Failed to scan log output, error: %s", readErr)
		}
	}

	identities := []common.CodeSigningIdentityInfo{}
	for _, v := range identitiesMap {
		identities = append(identities, v)
	}
	provProfiles := []provprofile.ProvisioningProfileInfo{}
	for _, v := range provProfilesMap {
		provProfiles = append(provProfiles, v)
	}
	teamIDs := maputil.KeysOfStringInterfaceMap(teamIDsMap)
	appIDs := maputil.KeysOfStringInterfaceMap(appIDsMap)

	return common.CodeSigningSettings{
		Identities:   identities,
		ProvProfiles: provProfiles,
		TeamIDs:      teamIDs,
		AppIDs:       appIDs,
	}, nil
}

// GenerateLog : generates the log for subsequent "Scan" call
func (xccmd CommandModel) GenerateLog() (string, error) {
	xcoutput := ""
	var err error

	progress.SimpleProgress(".", 1*time.Second, func() {
		xcoutput, err = xccmd.RunXcodebuildCommand("clean", "archive")
	})
	fmt.Println()

	if err != nil {
		return xcoutput, fmt.Errorf("Failed to Archive, error: %s", err)
	}
	return xcoutput, nil
}

// ScanCodeSigningSettings ...
func (xccmd CommandModel) ScanCodeSigningSettings(logToScan string) (common.CodeSigningSettings, error) {
	return parseCodeSigningSettingsFromXcodeOutput(logToScan)
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

	log.Infof("$ xcodebuild %s", cmdex.PrintableCommandArgs(true, xcodeCmdParamsToRun))
	fmt.Print("Running and analyzing log ...")
	xcoutput, err := cmdex.RunCommandAndReturnCombinedStdoutAndStderr("xcodebuild", xcodeCmdParamsToRun...)
	if err != nil {
		return xcoutput, fmt.Errorf("Failed to run xcodebuild command, error: %s", err)
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
