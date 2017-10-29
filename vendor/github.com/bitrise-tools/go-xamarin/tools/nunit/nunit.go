package nunit

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/constants"
)

const (
	nunit3Console = "nunit3-console.exe"
)

// Model ...
type Model struct {
	nunitConsolePth string

	projectPth string
	config     string

	dllPth string
	test   string

	resultLogPth string

	customOptions []string
}

// SystemNunit3ConsolePath ...
func SystemNunit3ConsolePath() (string, error) {
	nunitDir := os.Getenv("NUNIT_PATH")
	if nunitDir == "" {
		return "", fmt.Errorf("NUNIT_PATH environment is not set, failed to determin nunit console path")
	}

	nunitConsolePth := filepath.Join(nunitDir, nunit3Console)
	if exist, err := pathutil.IsPathExists(nunitConsolePth); err != nil {
		return "", fmt.Errorf("Failed to check if nunit console exist at (%s), error: %s", nunitConsolePth, err)
	} else if !exist {
		return "", fmt.Errorf("nunit console not exist at: %s", nunitConsolePth)
	}

	return nunitConsolePth, nil
}

// New ...
func New(nunitConsolePth string) (*Model, error) {
	absNunitConsolePth, err := pathutil.AbsPath(nunitConsolePth)
	if err != nil {
		return nil, fmt.Errorf("Failed to expand path (%s), error: %s", nunitConsolePth, err)
	}

	return &Model{nunitConsolePth: absNunitConsolePth}, nil
}

// SetProjectPth ...
func (nunitConsole *Model) SetProjectPth(projectPth string) *Model {
	nunitConsole.projectPth = projectPth
	return nunitConsole
}

// SetConfig ...
func (nunitConsole *Model) SetConfig(config string) *Model {
	nunitConsole.config = config
	return nunitConsole
}

// SetDLLPth ...
func (nunitConsole *Model) SetDLLPth(dllPth string) *Model {
	nunitConsole.dllPth = dllPth
	return nunitConsole
}

// SetTestToRun ...
func (nunitConsole *Model) SetTestToRun(test string) *Model {
	nunitConsole.test = test
	return nunitConsole
}

// SetResultLogPth ...
func (nunitConsole *Model) SetResultLogPth(resultLogPth string) *Model {
	nunitConsole.resultLogPth = resultLogPth
	return nunitConsole
}

// SetCustomOptions ...
func (nunitConsole *Model) SetCustomOptions(options ...string) {
	nunitConsole.customOptions = options
}

func (nunitConsole *Model) commandSlice() []string {
	cmdSlice := []string{constants.MonoPath}
	cmdSlice = append(cmdSlice, nunitConsole.nunitConsolePth)

	if nunitConsole.projectPth != "" {
		cmdSlice = append(cmdSlice, nunitConsole.projectPth)
	}
	if nunitConsole.config != "" {
		cmdSlice = append(cmdSlice, fmt.Sprintf("/config:%s", nunitConsole.config))
	}

	if nunitConsole.dllPth != "" {
		cmdSlice = append(cmdSlice, nunitConsole.dllPth)
	}
	if nunitConsole.test != "" {
		cmdSlice = append(cmdSlice, "--test", nunitConsole.test)
	}

	if nunitConsole.resultLogPth != "" {
		cmdSlice = append(cmdSlice, "--result", nunitConsole.resultLogPth)
	}

	cmdSlice = append(cmdSlice, nunitConsole.customOptions...)
	return cmdSlice
}

// PrintableCommand ...
func (nunitConsole Model) PrintableCommand() string {
	cmdSlice := nunitConsole.commandSlice()

	return command.PrintableCommandArgs(true, cmdSlice)
}

// Run ...
func (nunitConsole Model) Run() error {
	cmdSlice := nunitConsole.commandSlice()

	command, err := command.NewFromSlice(cmdSlice)
	if err != nil {
		return err
	}

	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	return command.Run()
}
