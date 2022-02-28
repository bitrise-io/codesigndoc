package ios

import "github.com/bitrise-io/bitrise-init/models"

//------------------
// ScannerInterface
//------------------

// Scanner ...
type Scanner struct {
	DetectResult DetectResult

	ConfigDescriptors []ConfigDescriptor

	ExcludeAppIcon            bool
	SuppressPodFileParseError bool
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return string(XcodeProjectTypeIOS)
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	result, err := ParseProjects(XcodeProjectTypeIOS, searchDir, scanner.ExcludeAppIcon, scanner.SuppressPodFileParseError)
	if err != nil {
		return false, err
	}

	scanner.DetectResult = result
	detected := len(result.Projects) > 0
	return detected, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return []string{}
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	options, configDescriptors, icons, warnings, err := GenerateOptions(XcodeProjectTypeIOS, scanner.DetectResult)
	if err != nil {
		return models.OptionNode{}, warnings, nil, err
	}

	scanner.ConfigDescriptors = configDescriptors

	return options, warnings, icons, nil
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionNode {
	return GenerateDefaultOptions(XcodeProjectTypeIOS)
}

// Configs ...
func (scanner *Scanner) Configs(isPrivateRepository bool) (models.BitriseConfigMap, error) {
	return GenerateConfig(XcodeProjectTypeIOS, scanner.ConfigDescriptors, isPrivateRepository)
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return GenerateDefaultConfig(XcodeProjectTypeIOS)
}

// GetProjectType returns the project_type property used in a bitrise config
func (Scanner) GetProjectType() string {
	return string(XcodeProjectTypeIOS)
}
