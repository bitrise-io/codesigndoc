package xcodebuild

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
)

// ShowBuildSettingsCommandModel ...
type ShowBuildSettingsCommandModel struct {
	projectPath string

	target        string
	scheme        string
	configuration string
	customOptions []string
}

// NewShowBuildSettingsCommand ...
func NewShowBuildSettingsCommand(projectPath string) *ShowBuildSettingsCommandModel {
	return &ShowBuildSettingsCommandModel{
		projectPath: projectPath,
	}
}

// SetTarget ...
func (c *ShowBuildSettingsCommandModel) SetTarget(target string) *ShowBuildSettingsCommandModel {
	c.target = target
	return c
}

func (c *ShowBuildSettingsCommandModel) cmdSlice() []string {
	slice := []string{toolName}

	if c.projectPath != "" {
		if filepath.Ext(c.projectPath) == ".xcworkspace" {
			slice = append(slice, "-workspace", c.projectPath)
		} else {
			slice = append(slice, "-project", c.projectPath)
		}
	}

	if c.target != "" {
		slice = append(slice, "-target", c.target)
	}

	if c.scheme != "" {
		slice = append(slice, "-scheme", c.scheme)
	}

	if c.configuration != "" {
		slice = append(slice, "-configuration", c.configuration)
	}

	slice = append(slice, "-showBuildSettings")
	slice = append(slice, c.customOptions...)

	return slice
}

// SetScheme ...
func (c *ShowBuildSettingsCommandModel) SetScheme(scheme string) *ShowBuildSettingsCommandModel {
	c.scheme = scheme
	return c
}

// SetConfiguration ...
func (c *ShowBuildSettingsCommandModel) SetConfiguration(configuration string) *ShowBuildSettingsCommandModel {
	c.configuration = configuration
	return c
}

// SetCustomOptions ...
func (c *ShowBuildSettingsCommandModel) SetCustomOptions(customOptions []string) *ShowBuildSettingsCommandModel {
	c.customOptions = customOptions
	return c
}

// Command ...
func (c ShowBuildSettingsCommandModel) Command() *command.Model {
	cmdSlice := c.cmdSlice()
	return command.New(cmdSlice[0], cmdSlice[1:]...)
}

// PrintableCmd ...
func (c ShowBuildSettingsCommandModel) PrintableCmd() string {
	cmdSlice := c.cmdSlice()
	return command.PrintableCommandArgs(false, cmdSlice)
}

func parseBuildSettings(out string) (serialized.Object, error) {
	settings := serialized.Object{}

	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if split := strings.Split(line, "="); len(split) > 1 {
			key := strings.TrimSpace(split[0])
			value := strings.TrimSpace(strings.Join(split[1:], "="))
			value = strings.Trim(value, `"`)

			settings[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return settings, nil
}

// RunAndReturnSettings ...
func (c ShowBuildSettingsCommandModel) RunAndReturnSettings() (serialized.Object, error) {
	cmd := c.Command()
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		if errorutil.IsExitStatusError(err) {
			return nil, fmt.Errorf("%s command failed: output: %s", cmd.PrintableCommandArgs(), out)
		}
		return nil, fmt.Errorf("failed to run command %s: %s", cmd.PrintableCommandArgs(), err)
	}

	return parseBuildSettings(out)
}
