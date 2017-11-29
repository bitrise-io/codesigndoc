package xcodebuild

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/command"
)

// ShowBuildSettingsCommandModel ...
type ShowBuildSettingsCommandModel struct {
	projectPath string
	isWorkspace bool
}

// NewShowBuildSettingsCommand ...
func NewShowBuildSettingsCommand(projectPath string, isWorkspace bool) *ShowBuildSettingsCommandModel {
	return &ShowBuildSettingsCommandModel{
		projectPath: projectPath,
		isWorkspace: isWorkspace,
	}
}

func (c *ShowBuildSettingsCommandModel) cmdSlice() []string {
	slice := []string{toolName}

	if c.projectPath != "" {
		if c.isWorkspace {
			slice = append(slice, "-workspace", c.projectPath)
		} else {
			slice = append(slice, "-project", c.projectPath)
		}
	}

	return slice
}

// PrintableCmd ...
func (c ShowBuildSettingsCommandModel) PrintableCmd() string {
	cmdSlice := c.cmdSlice()
	return command.PrintableCommandArgs(false, cmdSlice)
}

// Command ...
func (c ShowBuildSettingsCommandModel) Command() *command.Model {
	cmdSlice := c.cmdSlice()
	return command.New(cmdSlice[0], cmdSlice[1:]...)
}

// Cmd ...
func (c ShowBuildSettingsCommandModel) Cmd() *exec.Cmd {
	command := c.Command()
	return command.GetCmd()
}

func parseBuildSettings(out string) (map[string]string, error) {
	settings := map[string]string{}

	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if split := strings.Split(line, "="); len(split) == 2 {
			key := strings.TrimSpace(split[0])
			value := strings.TrimSpace(split[1])
			value = strings.Trim(value, `"`)

			settings[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return map[string]string{}, err
	}

	return settings, nil
}

// RunAndReturnSettings ...
func (c ShowBuildSettingsCommandModel) RunAndReturnSettings() (map[string]string, error) {
	command := c.Command()
	out, err := command.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return map[string]string{}, err
	}

	return parseBuildSettings(out)
}
