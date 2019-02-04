package models

// BitriseConfigMap ...
type BitriseConfigMap map[string]string

// Warnings ...
type Warnings []string

// Errors ...
type Errors []string

// ScanResultModel ...
type ScanResultModel struct {
	ScannerToOptionRoot       map[string]OptionNode       `json:"options,omitempty" yaml:"options,omitempty"`
	ScannerToBitriseConfigMap map[string]BitriseConfigMap `json:"configs,omitempty" yaml:"configs,omitempty"`
	ScannerToWarnings         map[string]Warnings         `json:"warnings,omitempty" yaml:"warnings,omitempty"`
	ScannerToErrors           map[string]Errors           `json:"errors,omitempty" yaml:"errors,omitempty"`
}

// AddError ...
func (result *ScanResultModel) AddError(platform string, errorMessage string) {
	if result.ScannerToErrors == nil {
		result.ScannerToErrors = map[string]Errors{}
	}
	if result.ScannerToErrors[platform] == nil {
		result.ScannerToErrors[platform] = []string{}
	}
	result.ScannerToErrors[platform] = append(result.ScannerToErrors[platform], errorMessage)
}
