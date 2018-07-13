package xcodebuild

import (
	"fmt"
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
	projectPath   string
	isWorkspace   bool
	scheme        string
	configuration string

	// buildsetting
	forceDevelopmentTeam              string
	forceProvisioningProfileSpecifier string
	forceProvisioningProfile          string
	forceCodeSignIdentity             string
	disableCodesign                   bool

	// buildaction
	customBuildActions []string

	// Options
	archivePath   string
	customOptions []string
	sdk           string

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

// SetForceDevelopmentTeam ...
func (c *CommandBuilder) SetForceDevelopmentTeam(forceDevelopmentTeam string) *CommandBuilder {
	c.forceDevelopmentTeam = forceDevelopmentTeam
	return c
}

// SetForceProvisioningProfileSpecifier ...
func (c *CommandBuilder) SetForceProvisioningProfileSpecifier(forceProvisioningProfileSpecifier string) *CommandBuilder {
	c.forceProvisioningProfileSpecifier = forceProvisioningProfileSpecifier
	return c
}

// SetForceProvisioningProfile ...
func (c *CommandBuilder) SetForceProvisioningProfile(forceProvisioningProfile string) *CommandBuilder {
	c.forceProvisioningProfile = forceProvisioningProfile
	return c
}

// SetForceCodeSignIdentity ...
func (c *CommandBuilder) SetForceCodeSignIdentity(forceCodeSignIdentity string) *CommandBuilder {
	c.forceCodeSignIdentity = forceCodeSignIdentity
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

	if c.forceDevelopmentTeam != "" {
		slice = append(slice, fmt.Sprintf("DEVELOPMENT_TEAM=%s", c.forceDevelopmentTeam))
	}
	if c.forceProvisioningProfileSpecifier != "" {
		slice = append(slice, fmt.Sprintf("PROVISIONING_PROFILE_SPECIFIER=%s", c.forceProvisioningProfileSpecifier))
	}
	if c.forceProvisioningProfile != "" {
		slice = append(slice, fmt.Sprintf("PROVISIONING_PROFILE=%s", c.forceProvisioningProfile))
	}
	if c.forceCodeSignIdentity != "" {
		slice = append(slice, fmt.Sprintf("CODE_SIGN_IDENTITY=%s", c.forceCodeSignIdentity))
	} else if c.disableCodesign {
		slice = append(slice, "CODE_SIGN_IDENTITY=")
		slice = append(slice, "CODE_SIGNING_REQUIRED=NO")
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
