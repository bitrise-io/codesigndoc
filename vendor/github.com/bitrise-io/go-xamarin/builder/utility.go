package builder

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xamarin/analyzers/solution"
	"github.com/bitrise-io/go-xamarin/constants"
	"github.com/bitrise-io/go-xamarin/utility"
)

func validateSolutionPth(pth string) error {
	ext := filepath.Ext(pth)
	if ext != constants.SolutionExt {
		return fmt.Errorf("path is not a solution file path: %s", pth)
	}
	if exist, err := pathutil.IsPathExists(pth); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("solution not exist at: %s", pth)
	}
	return nil
}

func validateSolutionConfig(solution solution.Model, configuration, platform string) error {
	config := utility.ToConfig(configuration, platform)
	if _, ok := solution.ConfigMap[config]; !ok {
		return fmt.Errorf("invalid solution config, available: %v", solution.ConfigList())
	}
	return nil
}

func whitelistAllows(projectType constants.SDK, projectTypeWhiteList ...constants.SDK) bool {
	if len(projectTypeWhiteList) == 0 {
		return true
	}

	for _, filter := range projectTypeWhiteList {
		switch filter {
		case constants.SDKIOS:
			if projectType == constants.SDKIOS {
				return true
			}
		case constants.SDKTvOS:
			if projectType == constants.SDKTvOS {
				return true
			}
		case constants.SDKMacOS:
			if projectType == constants.SDKMacOS {
				return true
			}
		case constants.SDKAndroid:
			if projectType == constants.SDKAndroid {
				return true
			}
		}
	}

	return false
}

// IsDeviceArch based on:
// default architecture: ARMv7
// iPhone architecture: <MtouchArch>ARMv7,ARMv7s,ARM64</MtouchArch>
// iPhoneSimulator architecture: <MtouchArch>i386, x86_64</MtouchArch>
func IsDeviceArch(architectures ...string) bool {
	return len(architectures) == 0 || strings.HasPrefix(strings.ToLower(architectures[0]), "arm")
}

func isPlatformAnyCPU(platform string) bool {
	return (platform == "Any CPU" || platform == "AnyCPU")
}

func androidPackageName(manifestPth string) (string, error) {
	content, err := fileutil.ReadStringFromFile(manifestPth)
	if err != nil {
		return "", err
	}

	return androidPackageNameFromManifestContent(content)
}

func androidPackageNameFromManifestContent(manifestContent string) (string, error) {
	// package is attribute of the rott xml element
	manifestContent = "<a>" + manifestContent + "</a>"

	type Manifest struct {
		Package string `xml:"package,attr"`
	}

	type Result struct {
		Manifest Manifest `xml:"manifest"`
	}

	var result Result
	if err := xml.Unmarshal([]byte(manifestContent), &result); err != nil {
		return "", err
	}

	return result.Manifest.Package, nil
}
