package testcloud

import (
	"bufio"
	"fmt"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/constants"
)

// Parallelization ...
type Parallelization string

const (
	// ParallelizationUnknown ...
	ParallelizationUnknown Parallelization = "unkown"
	// ParallelizationByTestFixture ...
	ParallelizationByTestFixture Parallelization = "by_test_fixture"
	// ParallelizationByTestChunk ...
	ParallelizationByTestChunk Parallelization = "by_test_chunk"
)

// ParseParallelization ...
func ParseParallelization(parallelization string) (Parallelization, error) {
	switch parallelization {
	case "by_test_fixture":
		return ParallelizationByTestFixture, nil
	case "by_test_chunk":
		return ParallelizationByTestChunk, nil
	default:
		return ParallelizationUnknown, fmt.Errorf("Unkown parallelization (%s)", parallelization)
	}
}

// Model ...
type Model struct {
	testCloudExePth string

	apkPth  string
	ipaPth  string
	dsymPth string

	apiKey          string
	user            string
	assemblyDir     string
	devices         string
	isAsyncJSON     bool
	series          string
	nunitXMLPth     string
	parallelization Parallelization

	signOptions   []string
	customOptions []string
}

// NewModel ...
func NewModel(testCloudExexPth string) (*Model, error) {
	absTestCloudExexPth, err := pathutil.AbsPath(testCloudExexPth)
	if err != nil {
		return nil, fmt.Errorf("Failed to expand path (%s), error: %s", testCloudExexPth, err)
	}

	return &Model{testCloudExePth: absTestCloudExexPth}, nil
}

// SetAPKPth ...
func (testCloud *Model) SetAPKPth(apkPth string) *Model {
	testCloud.apkPth = apkPth
	return testCloud
}

// SetIPAPth ...
func (testCloud *Model) SetIPAPth(ipaPth string) *Model {
	testCloud.ipaPth = ipaPth
	return testCloud
}

// SetDSYMPth ...
func (testCloud *Model) SetDSYMPth(dsymPth string) *Model {
	testCloud.dsymPth = dsymPth
	return testCloud
}

// SetAPIKey ...
func (testCloud *Model) SetAPIKey(apiKey string) *Model {
	testCloud.apiKey = apiKey
	return testCloud
}

// SetUser ...
func (testCloud *Model) SetUser(user string) *Model {
	testCloud.user = user
	return testCloud
}

// SetAssemblyDir ...
func (testCloud *Model) SetAssemblyDir(assemblyDir string) *Model {
	testCloud.assemblyDir = assemblyDir
	return testCloud
}

// SetDevices ...
func (testCloud *Model) SetDevices(devices string) *Model {
	testCloud.devices = devices
	return testCloud
}

// SetIsAsyncJSON ...
func (testCloud *Model) SetIsAsyncJSON(isAsyncJSON bool) *Model {
	testCloud.isAsyncJSON = isAsyncJSON
	return testCloud
}

// SetSeries ...
func (testCloud *Model) SetSeries(series string) *Model {
	testCloud.series = series
	return testCloud
}

// SetNunitXMLPth ...
func (testCloud *Model) SetNunitXMLPth(nunitXMLPth string) *Model {
	testCloud.nunitXMLPth = nunitXMLPth
	return testCloud
}

// SetParallelization ...
func (testCloud *Model) SetParallelization(parallelization Parallelization) *Model {
	testCloud.parallelization = parallelization
	return testCloud
}

// SetSignOptions ...
func (testCloud *Model) SetSignOptions(options ...string) *Model {
	testCloud.signOptions = options
	return testCloud
}

// SetCustomOptions ...
func (testCloud *Model) SetCustomOptions(options ...string) *Model {
	testCloud.customOptions = options
	return testCloud
}

func (testCloud *Model) submitCommandSlice() []string {
	cmdSlice := []string{constants.MonoPath}
	cmdSlice = append(cmdSlice, testCloud.testCloudExePth)
	cmdSlice = append(cmdSlice, "submit")

	if testCloud.apkPth != "" {
		cmdSlice = append(cmdSlice, testCloud.apkPth)
	}

	if testCloud.ipaPth != "" {
		cmdSlice = append(cmdSlice, testCloud.ipaPth)
	}
	if testCloud.dsymPth != "" {
		cmdSlice = append(cmdSlice, "--dsym", testCloud.dsymPth)
	}

	cmdSlice = append(cmdSlice, testCloud.apiKey)

	for _, option := range testCloud.signOptions {
		cmdSlice = append(cmdSlice, option)
	}

	cmdSlice = append(cmdSlice, "--user", testCloud.user)
	cmdSlice = append(cmdSlice, "--assembly-dir", testCloud.assemblyDir)
	cmdSlice = append(cmdSlice, "--devices", testCloud.devices)

	if testCloud.isAsyncJSON {
		cmdSlice = append(cmdSlice, "--async-json")
	}

	cmdSlice = append(cmdSlice, "--series", testCloud.series)

	if testCloud.nunitXMLPth != "" {
		cmdSlice = append(cmdSlice, "--nunit-xml", testCloud.nunitXMLPth)
	}

	if testCloud.parallelization == ParallelizationByTestChunk {
		cmdSlice = append(cmdSlice, "--test-chunk")
	} else if testCloud.parallelization == ParallelizationByTestFixture {
		cmdSlice = append(cmdSlice, "--fixture-chunk")
	}

	cmdSlice = append(cmdSlice, testCloud.customOptions...)

	return cmdSlice
}

// String ...
func (testCloud Model) String() string {
	cmdSlice := testCloud.submitCommandSlice()
	return command.PrintableCommandArgs(true, cmdSlice)
}

// CaptureLineCallback ...
type CaptureLineCallback func(line string)

// Submit ...
func (testCloud Model) Submit(callback CaptureLineCallback) error {
	cmdSlice := testCloud.submitCommandSlice()

	command, err := command.NewFromSlice(cmdSlice)
	if err != nil {
		return err
	}

	cmd := *command.GetCmd()

	// Redirect output
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdoutReader)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if callback != nil {
				callback(line)
			}
		}
	}()
	if err := scanner.Err(); err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
