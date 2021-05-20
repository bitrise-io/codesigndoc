package xamarin

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

const (
	solutionExtension = ".sln"
	componentsDirName = "Components"
	// NodeModulesDirName ...
	NodeModulesDirName = "node_modules"

	solutionConfigurationStart = "GlobalSection(SolutionConfigurationPlatforms) = preSolution"
	solutionConfigurationEnd   = "EndGlobalSection"
)

var allowSolutionExtensionFilter = pathutil.ExtensionFilter(solutionExtension, true)
var forbidComponentsSolutionFilter = pathutil.ComponentFilter(componentsDirName, false)
var forbidNodeModulesDirComponentFilter = pathutil.ComponentFilter(NodeModulesDirName, false)

// FilterSolutionFiles ...
func FilterSolutionFiles(fileList []string) ([]string, error) {
	files, err := pathutil.FilterPaths(fileList,
		allowSolutionExtensionFilter,
		forbidComponentsSolutionFilter,
		forbidNodeModulesDirComponentFilter)
	if err != nil {
		return []string{}, err
	}

	return files, nil
}

// GetSolutionConfigs ...
func GetSolutionConfigs(solutionFile string) (map[string][]string, error) {
	content, err := fileutil.ReadStringFromFile(solutionFile)
	if err != nil {
		return map[string][]string{}, err
	}

	configMap := map[string][]string{}
	isNextLineScheme := false

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, solutionConfigurationStart) {
			isNextLineScheme = true
			continue
		}

		if strings.Contains(line, solutionConfigurationEnd) {
			isNextLineScheme = false
			continue
		}

		if isNextLineScheme {
			split := strings.Split(line, "=")
			if len(split) == 2 {
				configCompositStr := strings.TrimSpace(split[1])
				configSplit := strings.Split(configCompositStr, "|")

				if len(configSplit) == 2 {
					config := configSplit[0]
					platform := configSplit[1]

					platforms := configMap[config]
					platforms = append(platforms, platform)

					configMap[config] = platforms
				}
			} else {
				return map[string][]string{}, fmt.Errorf("failed to parse config line (%s)", line)
			}
		}
	}

	return configMap, nil
}
