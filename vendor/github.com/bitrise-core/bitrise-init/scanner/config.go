package scanner

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
)

const otherProjectType = "other"

type status int

const (
	// in case DetectPlatform() returned error, or false
	notDetected status = iota
	// in case DetectPlatform() returned true, but Options() or Config() returned an error
	detectedWithErrors
	// in case DetectPlatform() returned true, Options() and Config() returned no error
	detected
)

type scannerOutput struct {
	status status

	// can always be set
	// warnings returned by DetectPlatform(), Options()
	warnings models.Warnings

	// set if scanResultStatus is scanResultDetectedWithErrors
	// errors returned by Config()
	errors models.Errors

	// set if scanResultStatus is scanResultDetected
	options          models.OptionNode
	configs          models.BitriseConfigMap
	excludedScanners []string
}

// Config ...
func Config(searchDir string) models.ScanResultModel {
	result := models.ScanResultModel{}

	//
	// Setup
	currentDir, err := os.Getwd()
	if err != nil {
		result.AddError("general", fmt.Sprintf("Failed to expand current directory path, error: %s", err))
		return result
	}

	if searchDir == "" {
		searchDir = currentDir
	} else {
		absScerach, err := pathutil.AbsPath(searchDir)
		if err != nil {
			result.AddError("general", fmt.Sprintf("Failed to expand path (%s), error: %s", searchDir, err))
			return result
		}
		searchDir = absScerach
	}

	if searchDir != currentDir {
		if err := os.Chdir(searchDir); err != nil {
			result.AddError("general", fmt.Sprintf("Failed to change dir, to (%s), error: %s", searchDir, err))
			return result
		}
		defer func() {
			if err := os.Chdir(currentDir); err != nil {
				log.TWarnf("Failed to change dir, to (%s), error: %s", searchDir, err)
			}
		}()
	}
	// ---

	//
	// Scan
	log.TInfof(colorstring.Blue("Running scanners:"))
	fmt.Println()

	// Collect scanner outputs, by scanner name
	scannerToOutput := map[string]scannerOutput{}
	{
		projectScannerToOutputs := runScanners(scanners.ProjectScanners, searchDir)
		detectedProjectTypes := getDetectedScannerNames(projectScannerToOutputs)
		log.Printf("Detected project types: %s", detectedProjectTypes)
		fmt.Println()

		// Project types are needed by tool scanners, to create decision tree on which project type
		// to actually use in bitrise.yml
		if len(detectedProjectTypes) == 0 {
			detectedProjectTypes = []string{otherProjectType}
		}
		for _, toolScanner := range scanners.AutomationToolScanners {
			toolScanner.(scanners.AutomationToolScanner).SetDetectedProjectTypes(detectedProjectTypes)
		}

		toolScannerToOutputs := runScanners(scanners.AutomationToolScanners, searchDir)
		detectedAutomationToolScanners := getDetectedScannerNames(toolScannerToOutputs)
		log.Printf("Detected automation tools: %s", detectedAutomationToolScanners)
		fmt.Println()

		// Merge project and tool scanner outputs
		scannerToOutput = toolScannerToOutputs
		for scanner, scannerOutput := range projectScannerToOutputs {
			scannerToOutput[scanner] = scannerOutput
		}
	}

	scannerToWarnings := map[string]models.Warnings{}
	scannerToErrors := map[string]models.Errors{}
	scannerToOptions := map[string]models.OptionNode{}
	scannerToConfigMap := map[string]models.BitriseConfigMap{}
	for scanner, scannerOutput := range scannerToOutput {
		// Currently the tests except an empty warning list if no warnings
		// are created in the not detect case.
		if scannerOutput.status == notDetected && len(scannerOutput.warnings) > 0 ||
			scannerOutput.status != notDetected {
			scannerToWarnings[scanner] = scannerOutput.warnings
		}
		if len(scannerOutput.errors) > 0 &&
			(scannerOutput.status == detected || scannerOutput.status == detectedWithErrors) {
			scannerToErrors[scanner] = scannerOutput.errors
		}
		if len(scannerOutput.configs) > 0 && scannerOutput.status == detected {
			scannerToOptions[scanner] = scannerOutput.options
			scannerToConfigMap[scanner] = scannerOutput.configs
		}
	}
	return models.ScanResultModel{
		ScannerToOptionRoot:       scannerToOptions,
		ScannerToBitriseConfigMap: scannerToConfigMap,
		ScannerToWarnings:         scannerToWarnings,
		ScannerToErrors:           scannerToErrors,
	}
}

func runScanners(scannerList []scanners.ScannerInterface, searchDir string) map[string]scannerOutput {
	scannerOutputs := map[string]scannerOutput{}
	var excludedScannerNames []string
	for _, scanner := range scannerList {
		log.TInfof("Scanner: %s", colorstring.Blue(scanner.Name()))
		if sliceutil.IsStringInSlice(scanner.Name(), excludedScannerNames) {
			log.TWarnf("scanner is marked as excluded, skipping...")
			fmt.Println()
			continue
		}

		log.TPrintf("+------------------------------------------------------------------------------+")
		log.TPrintf("|                                                                              |")
		scannerOutput := runScanner(scanner, searchDir)
		log.TPrintf("|                                                                              |")
		log.TPrintf("+------------------------------------------------------------------------------+")
		fmt.Println()

		scannerOutputs[scanner.Name()] = scannerOutput
		excludedScannerNames = append(excludedScannerNames, scannerOutput.excludedScanners...)
	}
	return scannerOutputs
}

// Collect output of a specific scanner
func runScanner(detector scanners.ScannerInterface, searchDir string) scannerOutput {
	var detectorWarnings models.Warnings
	var detectorErrors []string

	if isDetect, err := detector.DetectPlatform(searchDir); err != nil {
		log.TErrorf("Scanner failed, error: %s", err)
		return scannerOutput{
			status:   notDetected,
			warnings: models.Warnings{err.Error()},
		}
	} else if !isDetect {
		return scannerOutput{
			status: notDetected,
		}
	}

	options, projectWarnings, err := detector.Options()
	detectorWarnings = append(detectorWarnings, projectWarnings...)

	if err != nil {
		log.TErrorf("Analyzer failed, error: %s", err)
		// Error returned as a warning
		detectorWarnings = append(detectorWarnings, err.Error())
		return scannerOutput{
			status:   detectedWithErrors,
			warnings: detectorWarnings,
		}
	}

	// Generate configs
	configs, err := detector.Configs()
	if err != nil {
		log.TErrorf("Failed to generate config, error: %s", err)
		detectorErrors = append(detectorErrors, err.Error())
		return scannerOutput{
			status:   detectedWithErrors,
			warnings: detectorWarnings,
			errors:   detectorErrors,
		}
	}

	scannerExcludedScanners := detector.ExcludedScannerNames()
	if len(scannerExcludedScanners) > 0 {
		log.TWarnf("Scanner will exclude scanners: %v", scannerExcludedScanners)
	}

	return scannerOutput{
		status:           detected,
		warnings:         detectorWarnings,
		errors:           detectorErrors,
		options:          options,
		configs:          configs,
		excludedScanners: scannerExcludedScanners,
	}
}

func getDetectedScannerNames(scannerOutputs map[string]scannerOutput) (names []string) {
	for scanner, scannerOutput := range scannerOutputs {
		if scannerOutput.status == detected {
			names = append(names, scanner)
		}
	}
	return
}
