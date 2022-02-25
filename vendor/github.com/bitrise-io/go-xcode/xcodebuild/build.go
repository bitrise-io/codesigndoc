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

// const ...
const (
	ArchiveAction Action = "archiveAction"
	BuildAction   Action = "buildAction"
	AnalyzeAction Action = "analyzeAction"
)

// Action ...
type Action string

// CommandBuilder ...
type CommandBuilder struct {
	projectPath    string
	isWorkspace    bool
	scheme         string
	configuration  string
	destination    string
	xcconfigPath   string
	authentication *AuthenticationParams

	// buildsetting
	disableCodesign bool

	// buildaction
	customBuildActions []string

	// Options
	archivePath      string
	customOptions    []string
	sdk              string
	resultBundlePath string

	// Archive
	action Action
}

// NewCommandBuilder ...
func NewCommandBuilder(projectPath string, isWorkspace bool, action Action) *CommandBuilder {
	return &CommandBuilder{
		projectPath: projectPath,
		isWorkspace: isWorkspace,
		action:      action,
	}
}

// SetScheme ...
func (c *CommandBuilder) SetScheme(scheme string) *CommandBuilder {
	c.scheme = scheme
	return c
}

// SetConfiguration ...
func (c *CommandBuilder) SetConfiguration(configuration string) *CommandBuilder {
	c.configuration = configuration
	return c
}

// SetDestination ...
func (c *CommandBuilder) SetDestination(destination string) *CommandBuilder {
	c.destination = destination
	return c
}

// SetXCConfigPath ...
func (c *CommandBuilder) SetXCConfigPath(xcconfigPath string) *CommandBuilder {
	c.xcconfigPath = xcconfigPath
	return c
}

// SetAuthentication ...
func (c *CommandBuilder) SetAuthentication(authenticationParams AuthenticationParams) *CommandBuilder {
	c.authentication = &authenticationParams
	return c
}

// SetCustomBuildAction ...
func (c *CommandBuilder) SetCustomBuildAction(buildAction ...string) *CommandBuilder {
	c.customBuildActions = buildAction
	return c
}

// SetArchivePath ...
func (c *CommandBuilder) SetArchivePath(archivePath string) *CommandBuilder {
	c.archivePath = archivePath
	return c
}

// SetResultBundlePath ...
func (c *CommandBuilder) SetResultBundlePath(resultBundlePath string) *CommandBuilder {
	c.resultBundlePath = resultBundlePath
	return c
}

// SetCustomOptions ...
func (c *CommandBuilder) SetCustomOptions(customOptions []string) *CommandBuilder {
	c.customOptions = customOptions
	return c
}

// SetSDK ...
func (c *CommandBuilder) SetSDK(sdk string) *CommandBuilder {
	c.sdk = sdk
	return c
}

// SetDisableCodesign ...
func (c *CommandBuilder) SetDisableCodesign(disable bool) *CommandBuilder {
	c.disableCodesign = disable
	return c
}

func (c *CommandBuilder) cmdSlice() []string {
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

	if c.configuration != "" {
		slice = append(slice, "-configuration", c.configuration)
	}

	if c.destination != "" {
		// "-destination" "id=07933176-D03B-48D3-A853-0800707579E6" => (need the plus `"` marks between the `destination` and the `id`)
		slice = append(slice, "-destination", c.destination)
	}

	if c.xcconfigPath != "" {
		slice = append(slice, "-xcconfig", c.xcconfigPath)
	}

	if c.disableCodesign {
		slice = append(slice, "CODE_SIGNING_ALLOWED=NO")
	}

	slice = append(slice, c.customBuildActions...)

	switch c.action {
	case ArchiveAction:
		slice = append(slice, "archive")

		if c.archivePath != "" {
			slice = append(slice, "-archivePath", c.archivePath)
		}
	case BuildAction:
		slice = append(slice, "build")
	case AnalyzeAction:
		slice = append(slice, "analyze")
	}

	if c.sdk != "" {
		slice = append(slice, "-sdk", c.sdk)
	}

	if c.resultBundlePath != "" {
		slice = append(slice, "-resultBundlePath", c.resultBundlePath)
	}

	if c.authentication != nil {
		slice = append(slice, c.authentication.args()...)
	}

	slice = append(slice, c.customOptions...)

	return slice
}

// PrintableCmd ...
func (c CommandBuilder) PrintableCmd() string {
	cmdSlice := c.cmdSlice()
	return command.PrintableCommandArgs(false, cmdSlice)
}

// Command ...
func (c CommandBuilder) Command() *command.Model {
	cmdSlice := c.cmdSlice()
	return command.New(cmdSlice[0], cmdSlice[1:]...)
}

// ExecCommand ...
func (c CommandBuilder) ExecCommand() *exec.Cmd {
	command := c.Command()
	return command.GetCmd()
}

// Run ...
func (c CommandBuilder) Run() error {
	command := c.Command()

	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	return command.Run()
}
