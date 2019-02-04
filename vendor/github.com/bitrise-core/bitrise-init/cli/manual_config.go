package cli

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/bitrise-core/bitrise-init/output"
	"github.com/bitrise-core/bitrise-init/scanner"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/urfave/cli"
)

const (
	defaultOutputDir = "_defaults"
)

var manualConfigCommand = cli.Command{
	Name:  "manual-config",
	Usage: "Generates default bitrise config files.",
	Action: func(c *cli.Context) error {
		if err := initManualConfig(c); err != nil {
			log.TErrorf(err.Error())
			os.Exit(1)
		}
		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "output-dir",
			Usage: "Directory to save scan results.",
			Value: "./_defaults",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "Output format, options [json, yaml].",
			Value: "yaml",
		},
	},
}

func initManualConfig(c *cli.Context) error {
	// Config
	isCI := c.GlobalBool("ci")
	outputDir := c.String("output-dir")
	formatStr := c.String("format")

	if isCI {
		log.TInfof(colorstring.Yellow("CI mode"))
	}
	log.TInfof(colorstring.Yellowf("output dir: %s", outputDir))
	log.TInfof(colorstring.Yellowf("output format: %s", formatStr))
	fmt.Println()

	currentDir, err := pathutil.AbsPath("./")
	if err != nil {
		return fmt.Errorf("Failed to get current directory, error: %s", err)
	}

	if outputDir == "" {
		outputDir = filepath.Join(currentDir, defaultOutputDir)
	}
	outputDir, err = pathutil.AbsPath(outputDir)
	if err != nil {
		return fmt.Errorf("Failed to expand path (%s), error: %s", outputDir, err)
	}
	if exist, err := pathutil.IsDirExists(outputDir); err != nil {
		return err
	} else if !exist {
		if err := os.MkdirAll(outputDir, 0700); err != nil {
			return fmt.Errorf("Failed to create (%s), error: %s", outputDir, err)
		}
	}

	if formatStr == "" {
		formatStr = output.YAMLFormat.String()
	}
	format, err := output.ParseFormat(formatStr)
	if err != nil {
		return fmt.Errorf("Failed to parse format, err: %s", err)
	}
	if format != output.JSONFormat && format != output.YAMLFormat {
		return fmt.Errorf("Not allowed output format (%v), options: [%s, %s]", format, output.YAMLFormat.String(), output.JSONFormat.String())
	}
	// ---

	scanResult, err := scanner.ManualConfig()
	if err != nil {
		return err
	}

	// Write output to files
	if isCI {
		log.TInfof(colorstring.Blue("Saving outputs:"))

		if err := os.MkdirAll(outputDir, 0700); err != nil {
			return fmt.Errorf("Failed to create (%s), error: %s", outputDir, err)
		}

		pth := path.Join(outputDir, "result")
		outputPth, err := output.WriteToFile(scanResult, format, pth)
		if err != nil {
			return fmt.Errorf("Failed to print result, error: %s", err)
		}
		log.TInfof("  scan result: %s", colorstring.Blue(outputPth))

		return nil
	}
	// ---

	// Select option
	log.TInfof(colorstring.Blue("Collecting inputs:"))

	config, err := scanner.AskForConfig(scanResult)
	if err != nil {
		return err
	}

	pth := path.Join(outputDir, "bitrise.yml")
	outputPth, err := output.WriteToFile(config, format, pth)
	if err != nil {
		return fmt.Errorf("Failed to print result, error: %s", err)
	}
	log.TInfof("  bitrise.yml template: %s", colorstring.Blue(outputPth))
	fmt.Println()
	// ---

	return nil
}
