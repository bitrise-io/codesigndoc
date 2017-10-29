package xbuild

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/constants"
)

// Model ...
type Model struct {
	BuildTool string

	SolutionPth string
	ProjectPth  string

	target        string
	configuration string
	platform      string

	buildIpa       bool
	archiveOnBuild bool

	customOptions []string
}

// New ...
func New(solutionPth, projectPth string) (*Model, error) {
	absSolutionPth, err := pathutil.AbsPath(solutionPth)
	if err != nil {
		return nil, fmt.Errorf("Failed to expand path (%s), error: %s", solutionPth, err)
	}

	absProjectPth := ""
	if projectPth != "" {
		absPth, err := pathutil.AbsPath(projectPth)
		if err != nil {
			return nil, fmt.Errorf("Failed to expand path (%s), error: %s", projectPth, err)
		}
		absProjectPth = absPth
	}

	return &Model{SolutionPth: absSolutionPth, ProjectPth: absProjectPth, BuildTool: constants.XbuildPath}, nil
}

// SetTarget ...
func (xbuild *Model) SetTarget(target string) *Model {
	xbuild.target = target
	return xbuild
}

// SetConfiguration ...
func (xbuild *Model) SetConfiguration(configuration string) *Model {
	xbuild.configuration = configuration
	return xbuild
}

// SetPlatform ...
func (xbuild *Model) SetPlatform(platform string) *Model {
	xbuild.platform = platform
	return xbuild
}

// SetBuildIpa ...
func (xbuild *Model) SetBuildIpa(buildIpa bool) *Model {
	xbuild.buildIpa = buildIpa
	return xbuild
}

// SetArchiveOnBuild ...
func (xbuild *Model) SetArchiveOnBuild(archive bool) *Model {
	xbuild.archiveOnBuild = archive
	return xbuild
}

// SetCustomOptions ...
func (xbuild *Model) SetCustomOptions(options ...string) {
	xbuild.customOptions = options
}

func (xbuild *Model) buildCommandSlice() []string {
	cmdSlice := []string{xbuild.BuildTool}

	if xbuild.ProjectPth != "" {
		cmdSlice = append(cmdSlice, xbuild.ProjectPth)
	} else {
		cmdSlice = append(cmdSlice, xbuild.SolutionPth)
	}

	if xbuild.target != "" {
		cmdSlice = append(cmdSlice, fmt.Sprintf("/target:%s", xbuild.target))
	}

	cmdSlice = append(cmdSlice, fmt.Sprintf("/p:SolutionDir=%s", filepath.Dir(xbuild.SolutionPth)))

	if xbuild.configuration != "" {
		cmdSlice = append(cmdSlice, fmt.Sprintf("/p:Configuration=%s", xbuild.configuration))
	}

	if xbuild.platform != "" {
		cmdSlice = append(cmdSlice, fmt.Sprintf("/p:Platform=%s", xbuild.platform))
	}

	if xbuild.archiveOnBuild {
		cmdSlice = append(cmdSlice, "/p:ArchiveOnBuild=true")
	}

	if xbuild.buildIpa {
		cmdSlice = append(cmdSlice, "/p:BuildIpa=true")
	}

	cmdSlice = append(cmdSlice, xbuild.customOptions...)

	//cmdSlice = append(cmdSlice, "/verbosity:minimal", "/nologo")

	return cmdSlice
}

// PrintableCommand ...
func (xbuild *Model) PrintableCommand() string {
	cmdSlice := xbuild.buildCommandSlice()

	return command.PrintableCommandArgs(true, cmdSlice)
}

// Run ...
func (xbuild *Model) Run() error {
	cmdSlice := xbuild.buildCommandSlice()

	command, err := command.NewFromSlice(cmdSlice)
	if err != nil {
		return err
	}

	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	return command.Run()
}
