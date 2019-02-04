package macos

import (
	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/ios"
)

//------------------
// ScannerInterface
//------------------

// Scanner ...
type Scanner struct {
	searchDir         string
	configDescriptors []ios.ConfigDescriptor
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return string(ios.XcodeProjectTypeMacOS)
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	scanner.searchDir = searchDir

	detected, err := ios.Detect(ios.XcodeProjectTypeMacOS, searchDir)
	if err != nil {
		return false, err
	}

	return detected, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return []string{}
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, error) {
	options, configDescriptors, warnings, err := ios.GenerateOptions(ios.XcodeProjectTypeMacOS, scanner.searchDir)
	if err != nil {
		return models.OptionNode{}, warnings, err
	}

	scanner.configDescriptors = configDescriptors

	return options, warnings, nil
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionNode {
	return ios.GenerateDefaultOptions(ios.XcodeProjectTypeMacOS)
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	return ios.GenerateConfig(ios.XcodeProjectTypeMacOS, scanner.configDescriptors, true)
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return ios.GenerateDefaultConfig(ios.XcodeProjectTypeMacOS, true)
}
