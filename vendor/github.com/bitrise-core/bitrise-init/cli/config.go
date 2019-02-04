package cli

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/output"
	"github.com/bitrise-core/bitrise-init/scanner"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/urfave/cli"
)

const (
	defaultScanResultDir = "_scan_result"
)

var configCommand = cli.Command{
	Name:  "config",
	Usage: "Generates a bitrise config files based on your project.",
	Action: func(c *cli.Context) error {
		if err := initConfig(c); err != nil {
			log.TErrorf(err.Error())
			os.Exit(1)
		}
		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "dir",
			Usage: "Directory to scan.",
			Value: "./",
		},
		cli.StringFlag{
			Name:  "output-dir",
			Usage: "Directory to save scan results.",
			Value: "./_scan_result",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "Output format, options [json, yaml].",
			Value: "yaml",
		},
	},
}

func writeScanResult(scanResult models.ScanResultModel, outputDir string, format output.Format) (string, error) {
	pth := path.Join(outputDir, "result")
	return output.WriteToFile(scanResult, format, pth)
}

func initConfig(c *cli.Context) error {
	// Config
	isCI := c.GlobalBool("ci")
	searchDir := c.String("dir")
	outputDir := c.String("output-dir")
	formatStr := c.String("format")

	if isCI {
		log.TInfof(colorstring.Yellow("CI mode"))
	}
	log.TInfof(colorstring.Yellowf("scan dir: %s", searchDir))
	log.TInfof(colorstring.Yellowf("output dir: %s", outputDir))
	log.TInfof(colorstring.Yellowf("output format: %s", formatStr))
	fmt.Println()

	currentDir, err := pathutil.AbsPath("./")
	if err != nil {
		return fmt.Errorf("Failed to expand path (%s), error: %s", outputDir, err)
	}

	if searchDir == "" {
		searchDir = currentDir
	}
	searchDir, err = pathutil.AbsPath(searchDir)
	if err != nil {
		return fmt.Errorf("Failed to expand path (%s), error: %s", outputDir, err)
	}

	if outputDir == "" {
		outputDir = filepath.Join(currentDir, defaultScanResultDir)
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
		return fmt.Errorf("Failed to parse format (%s), error: %s", formatStr, err)
	}
	if format != output.JSONFormat && format != output.YAMLFormat {
		return fmt.Errorf("Not allowed output format (%s), options: [%s, %s]", format.String(), output.YAMLFormat.String(), output.JSONFormat.String())
	}
	// ---

	scanResult := scanner.Config(searchDir)

	platforms := []string{}
	for platform := range scanResult.ScannerToOptionRoot {
		platforms = append(platforms, platform)
	}

	if len(platforms) == 0 {
		cmd := command.New("which", "tree")
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		if err != nil || out == "" {
			log.TErrorf("tree not installed, can not list files")
		} else {
			fmt.Println()
			cmd := command.NewWithStandardOuts("tree", ".", "-L", "3")
			log.TPrintf("$ %s", cmd.PrintableCommandArgs())
			if err := cmd.Run(); err != nil {
				log.TErrorf("Failed to list files in current directory, error: %s", err)
			}
		}

		log.TInfof("Saving outputs:")
		scanResult.AddError("general", "No known platform detected")

		outputPth, err := writeScanResult(scanResult, outputDir, format)
		if err != nil {
			return fmt.Errorf("Failed to write output, error: %s", err)
		}

		log.TPrintf("scan result: %s", outputPth)
		return fmt.Errorf("No known platform detected")
	}

	// Write output to files
	if isCI {
		log.TInfof("Saving outputs:")

		outputPth, err := writeScanResult(scanResult, outputDir, format)
		if err != nil {
			return fmt.Errorf("Failed to write output, error: %s", err)
		}

		log.TPrintf("  scan result: %s", outputPth)
		return nil
	}
	// ---

	// Select option
	log.TInfof("Collecting inputs:")

	config, err := scanner.AskForConfig(scanResult)
	if err != nil {
		return err
	}

	if exist, err := pathutil.IsDirExists(outputDir); err != nil {
		return err
	} else if !exist {
		if err := os.MkdirAll(outputDir, 0700); err != nil {
			return fmt.Errorf("Failed to create (%s), error: %s", outputDir, err)
		}
	}

	pth := path.Join(outputDir, "bitrise.yml")
	outputPth, err := output.WriteToFile(config, format, pth)
	if err != nil {
		return fmt.Errorf("Failed to print result, error: %s", err)
	}
	log.TInfof("  bitrise.yml template: %s", outputPth)
	fmt.Println()
	// ---

	return nil
}
