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

// ArchiveCommandModel ...
type ArchiveCommandModel struct {
	projectPath   string
	isWorkspace   bool
	scheme        string
	configuration string

	// buildsetting
	forceDevelopmentTeam              string
	forceProvisioningProfileSpecifier string
	forceProvisioningProfile          string
	forceCodeSignIdentity             string

	// buildaction
	customBuildActions []string

	// Options
	archivePath   string
	customOptions []string
}

// NewArchiveCommand ...
func NewArchiveCommand(projectPath string, isWorkspace bool) *ArchiveCommandModel {
	return &ArchiveCommandModel{
		projectPath: projectPath,
		isWorkspace: isWorkspace,
	}
}

// SetScheme ...
func (c *ArchiveCommandModel) SetScheme(scheme string) *ArchiveCommandModel {
	c.scheme = scheme
	return c
}

// SetConfiguration ...
func (c *ArchiveCommandModel) SetConfiguration(configuration string) *ArchiveCommandModel {
	c.configuration = configuration
	return c
}

// SetForceDevelopmentTeam ...
func (c *ArchiveCommandModel) SetForceDevelopmentTeam(forceDevelopmentTeam string) *ArchiveCommandModel {
	c.forceDevelopmentTeam = forceDevelopmentTeam
	return c
}

// SetForceProvisioningProfileSpecifier ...
func (c *ArchiveCommandModel) SetForceProvisioningProfileSpecifier(forceProvisioningProfileSpecifier string) *ArchiveCommandModel {
	c.forceProvisioningProfileSpecifier = forceProvisioningProfileSpecifier
	return c
}

// SetForceProvisioningProfile ...
func (c *ArchiveCommandModel) SetForceProvisioningProfile(forceProvisioningProfile string) *ArchiveCommandModel {
	c.forceProvisioningProfile = forceProvisioningProfile
	return c
}

// SetForceCodeSignIdentity ...
func (c *ArchiveCommandModel) SetForceCodeSignIdentity(forceCodeSignIdentity string) *ArchiveCommandModel {
	c.forceCodeSignIdentity = forceCodeSignIdentity
	return c
}

// SetCustomBuildAction ...
func (c *ArchiveCommandModel) SetCustomBuildAction(buildAction ...string) *ArchiveCommandModel {
	c.customBuildActions = buildAction
	return c
}

// SetArchivePath ...
func (c *ArchiveCommandModel) SetArchivePath(archivePath string) *ArchiveCommandModel {
	c.archivePath = archivePath
	return c
}

// SetCustomOptions ...
func (c *ArchiveCommandModel) SetCustomOptions(customOptions []string) *ArchiveCommandModel {
	c.customOptions = customOptions
	return c
}

func (c *ArchiveCommandModel) cmdSlice() []string {
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
	}

	slice = append(slice, c.customBuildActions...)
	slice = append(slice, "archive")

	if c.archivePath != "" {
		slice = append(slice, "-archivePath", c.archivePath)
	}

	slice = append(slice, c.customOptions...)

	return slice
}

// PrintableCmd ...
func (c ArchiveCommandModel) PrintableCmd() string {
	cmdSlice := c.cmdSlice()
	return command.PrintableCommandArgs(false, cmdSlice)
}

// Command ...
func (c ArchiveCommandModel) Command() *command.Model {
	cmdSlice := c.cmdSlice()
	return command.New(cmdSlice[0], cmdSlice[1:]...)
}

// Cmd ...
func (c ArchiveCommandModel) Cmd() *exec.Cmd {
	command := c.Command()
	return command.GetCmd()
}

// Run ...
func (c ArchiveCommandModel) Run() error {
	command := c.Command()

	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	return command.Run()
}
