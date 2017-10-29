package cli

import (
	"fmt"

	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xamarin/builder"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/tools/buildtools"
	"github.com/urfave/cli"
)

func buildCmd(c *cli.Context) error {
	solutionPth := c.String(solutionFilePathKey)
	solutionConfiguration := c.String(solutionConfigurationKey)
	solutionPlatform := c.String(solutionPlatformKey)
	buildToolName := c.String(buildToolKey)

	fmt.Println()
	log.Infof("Config:")
	log.Printf("- solution: %s", solutionPth)
	log.Printf("- configuration: %s", solutionConfiguration)
	log.Printf("- platform: %s", solutionPlatform)
	log.Printf("- build-tool: %s", buildToolName)

	if solutionPth == "" {
		return fmt.Errorf("missing required input: %s", solutionFilePathKey)
	}
	if solutionConfiguration == "" {
		return fmt.Errorf("missing required input: %s", solutionConfigurationKey)
	}
	if solutionPlatform == "" {
		return fmt.Errorf("missing required input: %s", solutionPlatformKey)
	}

	buildTool := buildtools.Msbuild
	if buildToolName == "xbuild" {
		buildTool = buildtools.Xbuild
	}

	buildHandler, err := builder.New(solutionPth, nil, buildTool)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println()
	log.Infof("Building all projects in solution: %s", solutionPth)

	callback := func(solutionName string, projectName string, sdk constants.SDK, testFramwork constants.TestFramework, commandStr string, alreadyPerformed bool) {
		if projectName != "" {
			fmt.Println()
			log.Infof("Building project: %s", projectName)
			log.Donef("$ %s", commandStr)
			if alreadyPerformed {
				log.Warnf("build command already performed, skipping...")
			}
			fmt.Println()
		}
	}

	startTime := time.Now()

	warnings, err := buildHandler.BuildAllProjects(solutionConfiguration, solutionPlatform, nil, callback)
	for _, warning := range warnings {
		log.Warnf(warning)
	}
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	endTime := time.Now()

	fmt.Println()
	log.Infof("Collecting generated outputs")

	outputMap, err := buildHandler.CollectProjectOutputs(solutionConfiguration, solutionPlatform, startTime, endTime)
	if err != nil {
		return err
	}

	for projectName, projectOutput := range outputMap {
		fmt.Println()
		log.Infof("%s outputs:", projectName)

		for _, output := range projectOutput.Outputs {
			log.Donef("%s: %s", output.OutputType, output.Pth)
		}
	}

	return nil
}
