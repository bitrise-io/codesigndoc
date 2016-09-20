package xcode

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/maputil"
	"github.com/bitrise-io/go-utils/progress"
	"github.com/bitrise-io/go-utils/regexputil"
	"github.com/bitrise-tools/codesigndoc/provprofile"
)

// CommandModel ...
type CommandModel struct {
	// ProjectFilePath - might be a `xcodeproj` or `xcworkspace`
	ProjectFilePath  string
	Scheme           string
	CodeSignIdentity string
}

// CodeSigningIdentityInfo ...
type CodeSigningIdentityInfo struct {
	Title string
}

// CodeSigningSettings ...
type CodeSigningSettings struct {
	Identities   []CodeSigningIdentityInfo
	ProvProfiles []provprofile.ProvisioningProfileInfo
	TeamIDs      []string
	AppBundleIDs []string
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

func parseCodeSigningSettingsFromXcodeOutput(xcodeOutput string) CodeSigningSettings {
	scanner := bufio.NewScanner(strings.NewReader(xcodeOutput))

	identitiesMap := map[string]CodeSigningIdentityInfo{}
	provProfilesMap := map[string]provprofile.ProvisioningProfileInfo{}
	teamIDsMap := map[string]interface{}{}
	appBundleIDsMap := map[string]interface{}{}
	for scanner.Scan() {
		line := scanner.Text()

		// Team ID
		if rexp := regexp.MustCompile(`^[ ]*"com.apple.developer.team-identifier" = (?P<teamid>[a-zA-Z0-9]+);$`); rexp.MatchString(line) {
			results, err := regexputil.NamedFindStringSubmatch(rexp, line)
			if err != nil {
				log.Errorf("Failed to scan TeamID: %s", err)
				continue
			}
			teamIDsMap[results["teamid"]] = 1
		}

		// App Bundle ID
		if rexp := regexp.MustCompile(`^[ ]*"application-identifier" = "(?P<appbundleid>.+)";$`); rexp.MatchString(line) {
			results, err := regexputil.NamedFindStringSubmatch(rexp, line)
			if err != nil {
				log.Errorf("Failed to scan App Bundle ID: %s", err)
				continue
			}
			appBundleIDsMap[results["appbundleid"]] = 1
		}

		// Signing Identity
		if rexp := regexp.MustCompile(`^[ ]*Signing Identity:[ ]*"(?P<title>.+)"`); rexp.MatchString(line) {
			results, err := regexputil.NamedFindStringSubmatch(rexp, line)
			if err != nil {
				log.Errorf("Failed to scan Signing Identity title: %s", err)
				continue
			}
			codeSigningID := CodeSigningIdentityInfo{Title: results["title"]}
			identitiesMap[codeSigningID.Title] = codeSigningID
		}
		// Prov. Profile - title line
		if rexp := regexp.MustCompile(`^[ ]*Provisioning Profile:[ ]*"(?P<title>.+)"`); rexp.MatchString(line) {
			results, err := regexputil.NamedFindStringSubmatch(rexp, line)
			if err != nil {
				log.Errorf("Failed to scan Provisioning Profile title: %s", err)
				continue
			}
			tmpProvProfile := provprofile.ProvisioningProfileInfo{Title: results["title"]}
			if !scanner.Scan() {
				log.Error("Failed to scan Provisioning Profile UUID: no more lines to scan")
				continue
			}
			provProfileUUIDLine := scanner.Text()

			rexp = regexp.MustCompile(`^[ ]*\((?P<uuid>[a-zA-Z0-9-]{36})\)`)
			results, err = regexputil.NamedFindStringSubmatch(rexp, provProfileUUIDLine)
			if err != nil {
				log.Errorf("Failed to scan Provisioning Profile UUID: %s | line was: %s", err, provProfileUUIDLine)
				continue
			}
			tmpProvProfile.UUID = results["uuid"]
			provProfilesMap[tmpProvProfile.Title] = tmpProvProfile
		}
	}

	identities := []CodeSigningIdentityInfo{}
	for _, v := range identitiesMap {
		identities = append(identities, v)
	}
	provProfiles := []provprofile.ProvisioningProfileInfo{}
	for _, v := range provProfilesMap {
		provProfiles = append(provProfiles, v)
	}
	teamIDs := maputil.KeysOfStringInterfaceMap(teamIDsMap)
	appBundleIDs := maputil.KeysOfStringInterfaceMap(appBundleIDsMap)

	return CodeSigningSettings{
		Identities:   identities,
		ProvProfiles: provProfiles,
		TeamIDs:      teamIDs,
		AppBundleIDs: appBundleIDs,
	}
}

// ScanCodeSigningSettings ...
func (xccmd CommandModel) ScanCodeSigningSettings() (CodeSigningSettings, string, error) {
	xcoutput := ""
	var err error

	progress.SimpleProgress(".", 1*time.Second, func() {
		xcoutput, err = xccmd.RunXcodebuildCommand("clean", "archive")
	})
	fmt.Println()

	if err != nil {
		return CodeSigningSettings{}, xcoutput, fmt.Errorf("Failed to Archive, error: %s", err)
	}

	return parseCodeSigningSettingsFromXcodeOutput(xcoutput), xcoutput, nil
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
