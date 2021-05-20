package xcodebuild

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
)

func parseShowBuildSettingsOutput(out string) serialized.Object {
	settings := serialized.Object{}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		split := strings.Split(line, " = ")

		if len(split) < 2 {
			continue
		}

		key := strings.TrimSpace(split[0])
		if key == "" {
			continue
		}

		value := strings.TrimSpace(strings.Join(split[1:], " = "))

		settings[key] = value
	}

	return settings
}

// ShowProjectBuildSettings ...
func ShowProjectBuildSettings(project, target, configuration string, customOptions ...string) (serialized.Object, error) {
	args := []string{"-project", project, "-target", target, "-configuration", configuration}
	args = append(args, "-showBuildSettings")
	args = append(args, customOptions...)

	cmd := command.New("xcodebuild", args...)

	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		if errorutil.IsExitStatusError(err) {
			return nil, fmt.Errorf("%s command failed: output: %s", cmd.PrintableCommandArgs(), out)
		}

		return nil, fmt.Errorf("failed to run command %s: %s", cmd.PrintableCommandArgs(), err)
	}

	return parseShowBuildSettingsOutput(out), nil
}

// ShowWorkspaceBuildSettings ...
func ShowWorkspaceBuildSettings(workspace, scheme, configuration string, customOptions ...string) (serialized.Object, error) {
	args := []string{"-workspace", workspace, "-scheme", scheme, "-configuration", configuration}
	args = append(args, "-showBuildSettings")
	args = append(args, customOptions...)

	cmd := command.New("xcodebuild", args...)

	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		if errorutil.IsExitStatusError(err) {
			return nil, fmt.Errorf("%s command failed: output: %s", cmd.PrintableCommandArgs(), out)
		}

		return nil, fmt.Errorf("failed to run command %s: %s", cmd.PrintableCommandArgs(), err)
	}

	return parseShowBuildSettingsOutput(out), nil
}
