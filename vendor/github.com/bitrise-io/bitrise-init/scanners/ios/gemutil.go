package ios

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
)

// GemVersionFromGemfileLockContent ...
func GemVersionFromGemfileLockContent(gem, content string) string {
	relevantLines := []string{}
	lines := strings.Split(content, "\n")

	specsStart := false
	for _, line := range lines {
		if strings.Contains(line, "specs:") {
			specsStart = true
		}

		trimmed := strings.Trim(line, " ")
		if trimmed == "" {
			specsStart = false
		}

		if specsStart {
			relevantLines = append(relevantLines, line)
		}
	}

	exp := regexp.MustCompile(fmt.Sprintf(`%s \((.+)\)`, gem))
	for _, line := range relevantLines {
		match := exp.FindStringSubmatch(line)
		if len(match) == 2 {
			return match[1]
		}
	}

	return ""
}

// GemVersionFromGemfileLock ...
func GemVersionFromGemfileLock(gem, gemfileLockPth string) (string, error) {
	content, err := fileutil.ReadStringFromFile(gemfileLockPth)
	if err != nil {
		return "", err
	}
	return GemVersionFromGemfileLockContent(gem, content), nil
}
