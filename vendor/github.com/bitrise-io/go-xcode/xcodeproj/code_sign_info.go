package xcodeproj

import (
	"bufio"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/command/rubyscript"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-xcode/plistutil"
	"github.com/pkg/errors"
)

// CodeSignInfo ...
type CodeSignInfo struct {
	CodeSignEntitlementsPath     string
	BundleIdentifier             string
	CodeSignIdentity             string
	ProvisioningProfileSpecifier string
	ProvisioningProfile          string
	DevelopmentTeam              string
}

// TargetMapping ...
type TargetMapping struct {
	Configuration  string              `json:"configuration"`
	ProjectTargets map[string][]string `json:"project_targets"`
}

func clearRubyScriptOutput(out string) string {
	reader := strings.NewReader(out)
	scanner := bufio.NewScanner(reader)

	jsonLines := []string{}
	jsonResponseStart := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if !jsonResponseStart && trimmed == "{" {
			jsonResponseStart = true
		}
		if !jsonResponseStart {
			continue
		}

		jsonLines = append(jsonLines, line)
	}

	if len(jsonLines) == 0 {
		return out
	}

	return strings.Join(jsonLines, "\n")
}

func readSchemeTargetMapping(projectPth, scheme, user string) (TargetMapping, error) {
	runner := rubyscript.New(codeSignInfoScriptContent)
	bundleInstallCmd, err := runner.BundleInstallCommand(gemfileContent, "")
	if err != nil {
		return TargetMapping{}, fmt.Errorf("failed to create bundle install command, error: %s", err)
	}

	if out, err := bundleInstallCmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return TargetMapping{}, fmt.Errorf("bundle install failed, output: %s, error: %s", out, err)
	}

	runCmd, err := runner.RunScriptCommand()
	if err != nil {
		return TargetMapping{}, fmt.Errorf("failed to create script runner command, error: %s", err)
	}

	envsToAppend := []string{
		"project=" + projectPth,
		"scheme=" + scheme,
		"user=" + user}
	envs := append(runCmd.GetCmd().Env, envsToAppend...)

	runCmd.SetEnvs(envs...)

	out, err := runCmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return TargetMapping{}, fmt.Errorf("failed to run code signing analyzer script, output: %s, error: %s", out, err)
	}

	type OutputModel struct {
		Data  TargetMapping `json:"data"`
		Error string        `json:"error"`
	}
	var output OutputModel
	if err := json.Unmarshal([]byte(out), &output); err != nil {
		out = clearRubyScriptOutput(out)
		if err := json.Unmarshal([]byte(out), &output); err != nil {
			return TargetMapping{}, fmt.Errorf("failed to unmarshal output: %s", out)
		}
	}

	if output.Error != "" {
		return TargetMapping{}, fmt.Errorf("failed to get provisioning profile - bundle id mapping, error: %s", output.Error)
	}

	return output.Data, nil
}

func parseBuildSettingsOut(out string) (map[string]string, error) {
	reader := strings.NewReader(out)
	scanner := bufio.NewScanner(reader)

	buildSettings := map[string]string{}
	isBuildSettings := false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Build settings for") {
			isBuildSettings = true
			continue
		}
		if !isBuildSettings {
			continue
		}

		split := strings.Split(line, " = ")
		if len(split) > 1 {
			key := strings.TrimSpace(split[0])
			value := strings.TrimSpace(strings.Join(split[1:], " = "))

			buildSettings[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return map[string]string{}, errors.Wrap(err, "Failed to scan build settings")
	}

	return buildSettings, nil
}

func getTargetBuildSettingsWithXcodebuild(projectPth, target, configuration string) (map[string]string, error) {
	args := []string{"-showBuildSettings", "-project", projectPth, "-target", target}
	if configuration != "" {
		args = append(args, "-configuration", configuration)
	}

	cmd := command.New("xcodebuild", args...)
	cmd.SetDir(filepath.Dir(projectPth))

	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		if errorutil.IsExitStatusError(err) {
			return map[string]string{}, errors.Wrapf(err, "%s failed with output: %s", cmd.PrintableCommandArgs(), out)
		}
		return map[string]string{}, errors.Wrapf(err, "%s failed", cmd.PrintableCommandArgs())
	}

	return parseBuildSettingsOut(out)
}

func getBundleIDWithPlistbuddy(infoPlistPth string) (string, error) {
	plistData, err := plistutil.NewPlistDataFromFile(infoPlistPth)
	if err != nil {
		return "", err
	}

	bundleID, _ := plistData.GetString("CFBundleIdentifier")
	return bundleID, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// ResolveCodeSignInfo ...
func ResolveCodeSignInfo(projectOrWorkspacePth, scheme, user string) (map[string]CodeSignInfo, error) {
	projectTargetsMapping, err := readSchemeTargetMapping(projectOrWorkspacePth, scheme, user)
	if err != nil {
		return nil, err
	}

	resolvedCodeSignInfoMap := map[string]CodeSignInfo{}
	for projectPth, targets := range projectTargetsMapping.ProjectTargets {
		for _, targetName := range targets {
			if targetName == "" {
				return nil, errors.New("target name is empty")
			}

			if projectPth == "" {
				return nil, fmt.Errorf("failed to resolve which project contains target: %s", targetName)
			}

			buildSettings, err := getTargetBuildSettingsWithXcodebuild(projectPth, targetName, projectTargetsMapping.Configuration)
			if err != nil {
				return nil, fmt.Errorf("failed to read project build settings, error: %s", err)
			}

			// resolve Info.plist path
			infoPlistPth := buildSettings["INFOPLIST_FILE"]
			if infoPlistPth != "" {
				projectDir := filepath.Dir(projectPth)
				infoPlistPth = filepath.Join(projectDir, infoPlistPth)
			}
			// ---

			// resolve bundle id
			// best case if it presents in the buildSettings, since it is expanded
			bundleID := buildSettings["PRODUCT_BUNDLE_IDENTIFIER"]
			if bundleID == "" && infoPlistPth != "" {
				// try to find the bundle id in the Info.plist file, unless it contains env var
				id, err := getBundleIDWithPlistbuddy(infoPlistPth)
				if err != nil {
					return nil, fmt.Errorf("failed to resolve bundle id, error: %s", err)
				}
				bundleID = id
			}
			if bundleID == "" {
				return nil, fmt.Errorf("failed to resolve bundle id")
			}
			// ---

			codeSignEntitlementsPth := buildSettings["CODE_SIGN_ENTITLEMENTS"]
			if codeSignEntitlementsPth != "" {
				projectDir := filepath.Dir(projectPth)
				codeSignEntitlementsPth = filepath.Join(projectDir, codeSignEntitlementsPth)
			}

			codeSignIdentity := buildSettings["CODE_SIGN_IDENTITY"]
			provisioningProfileSpecifier := buildSettings["PROVISIONING_PROFILE_SPECIFIER"]
			provisioningProfile := buildSettings["PROVISIONING_PROFILE"]
			developmentTeam := buildSettings["DEVELOPMENT_TEAM"]

			resolvedCodeSignInfo := CodeSignInfo{
				CodeSignEntitlementsPath:     codeSignEntitlementsPth,
				BundleIdentifier:             bundleID,
				CodeSignIdentity:             codeSignIdentity,
				ProvisioningProfileSpecifier: provisioningProfileSpecifier,
				ProvisioningProfile:          provisioningProfile,
				DevelopmentTeam:              developmentTeam,
			}

			resolvedCodeSignInfoMap[targetName] = resolvedCodeSignInfo
		}
	}

	return resolvedCodeSignInfoMap, nil
}
