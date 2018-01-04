package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/bitrise-io/go-utils/pathutil"
)

// ModTimesByPath ...
type ModTimesByPath map[string]time.Time

func findModTimesByPath(dir string) (ModTimesByPath, error) {
	modTimesByPath := ModTimesByPath{}

	if walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		modTimesByPath[path] = info.ModTime()

		return nil
	}); walkErr != nil {
		return ModTimesByPath{}, walkErr
	}

	return modTimesByPath, nil
}

func isInTimeInterval(t, startTime, endTime time.Time) bool {
	if startTime.After(endTime) {
		return false
	}
	return (t.After(startTime) || t.Equal(startTime)) && (t.Before(endTime) || t.Equal(endTime))
}

func filterModTimesByPathByTimeWindow(modTimesByPath ModTimesByPath, startTime, endTime time.Time) ModTimesByPath {
	if startTime.IsZero() || endTime.IsZero() || startTime.Equal(endTime) || startTime.After(endTime) {
		return ModTimesByPath{}
	}

	filteredModTimesByPath := ModTimesByPath{}

	for pth, modTime := range modTimesByPath {
		if isInTimeInterval(modTime, startTime, endTime) {
			filteredModTimesByPath[pth] = modTime
		}
	}

	return filteredModTimesByPath
}

// finds the last modified file matching to most strict regexp
// order of regexps should be: most strict -> less strict
func findLastModifiedPathWithFileNameRegexps(modTimesByPath ModTimesByPath, regexps ...*regexp.Regexp) string {
	if len(modTimesByPath) == 0 {
		return ""
	}

	var lastModifiedPth string
	var lastModTime time.Time

	if len(regexps) > 0 {
		for _, re := range regexps {
			for pth, modTime := range modTimesByPath {
				fileName := filepath.Base(pth)
				if re.MatchString(fileName) {
					if modTime.After(lastModTime) {
						lastModifiedPth = pth
						lastModTime = modTime
					}
				}
			}

			// return with the most strict match
			if len(lastModifiedPth) > 0 {
				return lastModifiedPth
			}
		}
	} else {
		for pth, modTime := range modTimesByPath {
			if modTime.After(lastModTime) {
				lastModifiedPth = pth
				lastModTime = modTime
			}
		}
	}

	return lastModifiedPth
}

// exports the last modified file matching to most strict regexps within a time window
// order of regexps should be: most strict -> less strict
func findArtifact(dir string, startTime, endTime time.Time, patterns ...string) (string, error) {
	regexps := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		regexps[i] = regexp.MustCompile(pattern)
	}

	modTimesByPath, err := findModTimesByPath(dir)
	if err != nil {
		return "", err
	}

	modTimesByPathByTimeWindow := filterModTimesByPathByTimeWindow(modTimesByPath, startTime, endTime)
	return findLastModifiedPathWithFileNameRegexps(modTimesByPathByTimeWindow, regexps...), nil
}

func exportApk(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*signed.*\.apk$`, assemblyName),
		fmt.Sprintf(`(?i).*%s.*\.apk$`, assemblyName),
		`(?i).*signed.*\.apk$`,
		`(?i).*\.apk$`,
	)
}

func exportIpa(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.ipa$`, assemblyName),
		`(?i).*\.ipa$`,
	)
}

func exportXCArchive(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.xcarchive$`, assemblyName),
		fmt.Sprintf(`(?i).*\.xcarchive$`),
	)
}

func exportLatestXCArchiveFromXcodeArchives(assemblyName string, startTime, endTime time.Time) (string, error) {
	userHomeDir, ok := os.LookupEnv("HOME")
	if !ok {
		return "", fmt.Errorf("failed to get user home dir")
	}
	xcodeArchivesDir := filepath.Join(userHomeDir, "Library/Developer/Xcode/Archives")
	if exist, err := pathutil.IsDirExists(xcodeArchivesDir); err != nil {
		return "", err
	} else if !exist {
		return "", fmt.Errorf("no default Xcode archive path found at: %s", xcodeArchivesDir)
	}

	return exportXCArchive(xcodeArchivesDir, assemblyName, startTime, endTime)
}

func exportAppDSYM(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.app\.dSYM$`, assemblyName),
		`(?i).*\.app\.dSYM$`,
	)
}

func exportFrameworkDSYMs(outputDir string) ([]string, error) {
	// Multiplatform/iOS/bin/iPhone/Release/TTTAttributedLabel.framework.dSYM
	pattern := filepath.Join(outputDir, "*.framework.dSYM")
	return filepath.Glob(pattern)
}

func exportPKG(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.pkg$`, assemblyName),
		`(?i).*\.pkg$`,
	)
}

func exportApp(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.app$`, assemblyName),
		`(?i).*\.app$`,
	)
}

func exportDLL(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.dll$`, assemblyName),
		`(?i).*\.dll$`,
	)
}
