package xcodebuild

import (
	"os"
	"os/exec"

	"github.com/bitrise-io/go-utils/command"
)

/*
xcodebuild [-project <projectname>] \
	-scheme <schemeName> \
	[-destination <destinationspecifier>]... \
	[-configuration <configurationname>] \
	[-arch <architecture>]... \
	[-sdk [<sdkname>|<sdkpath>]] \
	[-showBuildSettings] \
	[<buildsetting>=<value>]... \
	[<buildaction>]...

xcodebuild -workspace <workspacename> \
	-scheme <schemeName> \
	[-destination <destinationspecifier>]... \
	[-configuration <configurationname>] \
	[-arch <architecture>]... \
	[-sdk [<sdkname>|<sdkpath>]] \
	[-showBuildSettings] \
	[<buildsetting>=<value>]... \
	[<buildaction>]...
*/

// TestCommandModel ...
type TestCommandModel struct {
	projectPath string
	isWorkspace bool
	scheme      string
	destination string

	// buildsetting
	generateCodeCoverage      bool
	disableIndexWhileBuilding bool

	// buildaction
	customBuildActions []string // clean, build

	// Options
	customOptions []string
}

// NewTestCommand ...
func NewTestCommand(projectPath string, isWorkspace bool) *TestCommandModel {
	return &TestCommandModel{
		projectPath: projectPath,
		isWorkspace: isWorkspace,
	}
}

// SetScheme ...
func (c *TestCommandModel) SetScheme(scheme string) *TestCommandModel {
	c.scheme = scheme
	return c
}

// SetDestination ...
func (c *TestCommandModel) SetDestination(destination string) *TestCommandModel {
	c.destination = destination
	return c
}

// SetGenerateCodeCoverage ...
func (c *TestCommandModel) SetGenerateCodeCoverage(generateCodeCoverage bool) *TestCommandModel {
	c.generateCodeCoverage = generateCodeCoverage
	return c
}

// SetCustomBuildAction ...
func (c *TestCommandModel) SetCustomBuildAction(buildAction ...string) *TestCommandModel {
	c.customBuildActions = buildAction
	return c
}

// SetCustomOptions ...
func (c *TestCommandModel) SetCustomOptions(customOptions []string) *TestCommandModel {
	c.customOptions = customOptions
	return c
}

// SetDisableIndexWhileBuilding ...
func (c *TestCommandModel) SetDisableIndexWhileBuilding(disable bool) *TestCommandModel {
	c.disableIndexWhileBuilding = disable
	return c
}

func (c *TestCommandModel) cmdSlice() []string {
	slice := []string{toolName}

	if c.projectPath != "" {
		if c.isWorkspace {
			slice = append(slice, "-workspace", c.projectPath)
		} else {
			slice = append(slice, "-project", c.projectPath)
		}
	}

	if c.scheme != "" {
		slice = append(slice, "-scheme", c.scheme)
	}

	if c.generateCodeCoverage {
		slice = append(slice, "GCC_INSTRUMENT_PROGRAM_FLOW_ARCS=YES", "GCC_GENERATE_TEST_COVERAGE_FILES=YES")
	}

	slice = append(slice, c.customBuildActions...)
	slice = append(slice, "test")
	if c.destination != "" {
		slice = append(slice, "-destination", c.destination)
	}

	if c.disableIndexWhileBuilding {
		slice = append(slice, "COMPILER_INDEX_STORE_ENABLE=NO")
	}

	slice = append(slice, c.customOptions...)

	return slice
}

// PrintableCmd ...
func (c TestCommandModel) PrintableCmd() string {
	cmdSlice := c.cmdSlice()
	return command.PrintableCommandArgs(false, cmdSlice)
}

// Command ...
func (c TestCommandModel) Command() *command.Model {
	cmdSlice := c.cmdSlice()
	return command.New(cmdSlice[0], cmdSlice[1:]...)
}

// Cmd ...
func (c TestCommandModel) Cmd() *exec.Cmd {
	command := c.Command()
	return command.GetCmd()
}

// Run ...
func (c TestCommandModel) Run() error {
	command := c.Command()

	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	return command.Run()
}
