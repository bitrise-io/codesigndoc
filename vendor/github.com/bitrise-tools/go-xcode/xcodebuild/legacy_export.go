package xcodebuild

import (
	"os"
	"os/exec"

	"github.com/bitrise-io/go-utils/command"
)

/*
xcodebuild -exportArchive \
	-exportFormat format \
	-archivePath xcarchivepath \
    -exportPath destinationpath \
    [-exportProvisioningProfile profilename] \
	[-exportSigningIdentity identityname] \
	[-exportInstallerIdentity identityname]
*/

// LegacyExportCommandModel ...
type LegacyExportCommandModel struct {
	exportFormat                  string
	archivePath                   string
	exportPath                    string
	exportProvisioningProfileName string
}

// NewLegacyExportCommand ...
func NewLegacyExportCommand() *LegacyExportCommandModel {
	return &LegacyExportCommandModel{}
}

// SetExportFormat ...
func (c *LegacyExportCommandModel) SetExportFormat(exportFormat string) *LegacyExportCommandModel {
	c.exportFormat = exportFormat
	return c
}

// SetArchivePath ...
func (c *LegacyExportCommandModel) SetArchivePath(archivePath string) *LegacyExportCommandModel {
	c.archivePath = archivePath
	return c
}

// SetExportPath ...
func (c *LegacyExportCommandModel) SetExportPath(exportPath string) *LegacyExportCommandModel {
	c.exportPath = exportPath
	return c
}

// SetExportProvisioningProfileName ...
func (c *LegacyExportCommandModel) SetExportProvisioningProfileName(exportProvisioningProfileName string) *LegacyExportCommandModel {
	c.exportProvisioningProfileName = exportProvisioningProfileName
	return c
}

func (c LegacyExportCommandModel) cmdSlice() []string {
	slice := []string{toolName, "-exportArchive"}
	if c.exportFormat != "" {
		slice = append(slice, "-exportFormat", c.exportFormat)
	}
	if c.archivePath != "" {
		slice = append(slice, "-archivePath", c.archivePath)
	}
	if c.exportPath != "" {
		slice = append(slice, "-exportPath", c.exportPath)
	}
	if c.exportProvisioningProfileName != "" {
		slice = append(slice, "-exportProvisioningProfile", c.exportProvisioningProfileName)
	}
	return slice
}

// PrintableCmd ...
func (c LegacyExportCommandModel) PrintableCmd() string {
	cmdSlice := c.cmdSlice()
	return command.PrintableCommandArgs(false, cmdSlice)
}

// Command ...
func (c LegacyExportCommandModel) Command() *command.Model {
	cmdSlice := c.cmdSlice()
	return command.New(cmdSlice[0], cmdSlice[1:]...)
}

// Cmd ...
func (c LegacyExportCommandModel) Cmd() *exec.Cmd {
	command := c.Command()
	return command.GetCmd()
}

// Run ...
func (c LegacyExportCommandModel) Run() error {
	command := c.Command()

	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	return command.Run()
}
