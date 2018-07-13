package xcpretty

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/command/rubycommand"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xcode/xcodebuild"
	version "github.com/hashicorp/go-version"
)

const (
	toolName = "xcpretty"
)

// CommandModel ...
type CommandModel struct {
	xcodebuildCommand xcodebuild.CommandModel

	customOptions []string
}

// New ...
func New(xcodebuildCommand xcodebuild.CommandModel) *CommandModel {
	return &CommandModel{
		xcodebuildCommand: xcodebuildCommand,
	}
}

// SetCustomOptions ...
func (c *CommandModel) SetCustomOptions(customOptions []string) *CommandModel {
	c.customOptions = customOptions
	return c
}

func (c CommandModel) cmdSlice() []string {
	slice := []string{toolName}
	slice = append(slice, c.customOptions...)
	return slice
}

// Command ...
func (c CommandModel) Command() *command.Model {
	cmdSlice := c.cmdSlice()
	return command.New(cmdSlice[0])
}

// PrintableCmd ...
func (c CommandModel) PrintableCmd() string {
	prettyCmdSlice := c.cmdSlice()
	prettyCmdStr := command.PrintableCommandArgs(false, prettyCmdSlice)

	cmdStr := c.xcodebuildCommand.PrintableCmd()

	return fmt.Sprintf("set -o pipefail && %s | %s", cmdStr, prettyCmdStr)
}

// Run ...
func (c CommandModel) Run() (string, error) {
	prettyCmd := c.Command()
	xcodebuildCmd := c.xcodebuildCommand.Command()

	// Configure cmd in- and outputs
	pipeReader, pipeWriter := io.Pipe()

	var outBuffer bytes.Buffer
	outWriter := io.MultiWriter(&outBuffer, pipeWriter)

	xcodebuildCmd.SetStdin(nil)
	xcodebuildCmd.SetStdout(outWriter)
	xcodebuildCmd.SetStderr(outWriter)

	prettyCmd.SetStdin(pipeReader)
	prettyCmd.SetStdout(os.Stdout)
	prettyCmd.SetStderr(os.Stdout)

	// Run
	if err := xcodebuildCmd.GetCmd().Start(); err != nil {
		out := outBuffer.String()
		return out, err
	}
	if err := prettyCmd.GetCmd().Start(); err != nil {
		out := outBuffer.String()
		return out, err
	}

	// Always close xcpretty outputs
	defer func() {
		if err := pipeWriter.Close(); err != nil {
			log.Warnf("Failed to close xcodebuild-xcpretty pipe, error: %s", err)
		}

		if err := prettyCmd.GetCmd().Wait(); err != nil {
			log.Warnf("xcpretty command failed, error: %s", err)
		}
	}()

	if err := xcodebuildCmd.GetCmd().Wait(); err != nil {
		out := outBuffer.String()
		return out, err
	}

	return outBuffer.String(), nil
}

// IsInstalled ...
func IsInstalled() (bool, error) {
	return rubycommand.IsGemInstalled("xcpretty", "")
}

// Install ...
func Install() error {
	cmds, err := rubycommand.GemInstall("xcpretty", "")
	if err != nil {
		return fmt.Errorf("failed to create command model, error: %s", err)
	}

	for _, cmd := range cmds {
		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create xcpretty install command, error: %s", err)
		}
	}

	return nil
}

// Version ...
func Version() (*version.Version, error) {
	cmd := command.New("xcpretty", "--version")
	versionOut, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, err
	}

	return parseVersionOut(versionOut)
}

func parseVersionOut(versionOut string) (*version.Version, error) {
	return version.NewVersion(versionOut)
}
