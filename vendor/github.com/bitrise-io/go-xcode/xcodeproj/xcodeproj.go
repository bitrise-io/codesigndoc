package xcodeproj

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
)

// Extensions
const (
	// XCWorkspaceExt ...
	XCWorkspaceExt = ".xcworkspace"
	// XCodeProjExt ...
	XCodeProjExt = ".xcodeproj"
	// XCSchemeExt ...
	XCSchemeExt = ".xcscheme"
)

// IsXCodeProj ...
func IsXCodeProj(pth string) bool {
	return strings.HasSuffix(pth, XCodeProjExt)
}

// IsXCWorkspace ...
func IsXCWorkspace(pth string) bool {
	return strings.HasSuffix(pth, XCWorkspaceExt)
}

// GetBuildConfigSDKs ...
func GetBuildConfigSDKs(pbxprojPth string) ([]string, error) {
	content, err := fileutil.ReadStringFromFile(pbxprojPth)
	if err != nil {
		return []string{}, err
	}

	return getBuildConfigSDKsFromContent(content)
}

func getBuildConfigSDKsFromContent(pbxprojContent string) ([]string, error) {
	sdkMap := map[string]bool{}

	beginXCBuildConfigurationSection := `/* Begin XCBuildConfiguration section */`
	endXCBuildConfigurationSection := `/* End XCBuildConfiguration section */`
	isXCBuildConfigurationSection := false

	// SDKROOT = macosx;
	pattern := `SDKROOT = (?P<sdk>.*);`
	regexp := regexp.MustCompile(pattern)

	scanner := bufio.NewScanner(strings.NewReader(pbxprojContent))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == endXCBuildConfigurationSection {
			break
		}

		if strings.TrimSpace(line) == beginXCBuildConfigurationSection {
			isXCBuildConfigurationSection = true
			continue
		}

		if !isXCBuildConfigurationSection {
			continue
		}

		if match := regexp.FindStringSubmatch(line); len(match) == 2 {
			sdk := match[1]
			sdkMap[sdk] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	sdks := []string{}
	for sdk := range sdkMap {
		sdks = append(sdks, sdk)
	}

	return sdks, nil
}
