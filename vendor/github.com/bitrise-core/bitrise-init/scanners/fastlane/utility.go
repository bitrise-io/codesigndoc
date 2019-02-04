package fastlane

import (
	"bufio"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/fileutil"
)

const (
	fastfileBasePath = "Fastfile"
)

// FilterFastfiles ...
func FilterFastfiles(fileList []string) ([]string, error) {
	allowFastfileBaseFilter := utility.BaseFilter(fastfileBasePath, true)
	fastfiles, err := utility.FilterPaths(fileList, allowFastfileBaseFilter)
	if err != nil {
		return []string{}, err
	}

	return utility.SortPathsByComponents(fastfiles)
}

func inspectFastfileContent(content string) ([]string, error) {
	commonLanes := []string{}
	laneMap := map[string][]string{}

	// platform :ios do ...
	platformSectionStartRegexp := regexp.MustCompile(`platform\s+:(?P<platform>.*)\s+do`)
	platformSectionEndPattern := "end"
	platform := ""

	// lane :test_and_snapshot do
	laneRegexp := regexp.MustCompile(`^[\s]*lane\s+:(?P<lane>.*)\s+do`)

	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " ")

		if platform != "" && line == platformSectionEndPattern {
			platform = ""
			continue
		}

		if platform == "" {
			if match := platformSectionStartRegexp.FindStringSubmatch(line); len(match) == 2 {
				platform = match[1]
				continue
			}
		}

		if match := laneRegexp.FindStringSubmatch(line); len(match) == 2 {
			lane := match[1]

			if platform != "" {
				lanes, found := laneMap[platform]
				if !found {
					lanes = []string{}
				}
				lanes = append(lanes, lane)
				laneMap[platform] = lanes
			} else {
				commonLanes = append(commonLanes, lane)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	lanes := commonLanes
	for platform, platformLanes := range laneMap {
		for _, lane := range platformLanes {
			lanes = append(lanes, platform+" "+lane)
		}
	}

	return lanes, nil
}

// InspectFastfile ...
func InspectFastfile(fastFile string) ([]string, error) {
	content, err := fileutil.ReadStringFromFile(fastFile)
	if err != nil {
		return []string{}, err
	}

	return inspectFastfileContent(content)
}

// WorkDir ...
func WorkDir(fastfilePth string) string {
	dirPth := filepath.Dir(fastfilePth)
	dirName := filepath.Base(dirPth)
	if dirName == "fastlane" {
		return filepath.Dir(dirPth)
	}
	return dirPth
}
