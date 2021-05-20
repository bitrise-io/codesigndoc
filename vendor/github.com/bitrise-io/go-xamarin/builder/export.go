package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

// ModTimesByPath ...
type ModTimesByPath map[string]time.Time

// findModTimesByPath walks through on the given directory and returns a ModTimesByPath for each file. Boolean
// excludeDir indicates if it should check directories or not.
func findModTimesByPath(dir string, excludeDir bool) (ModTimesByPath, error) {
	modTimesByPath := ModTimesByPath{}

	if walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		log.Debugf("Walking for path: %s, modtime %v", path, info.ModTime())

		if excludeDir && info.IsDir() {
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
	log.Debugf("Start time: %v, End time: %v", startTime, endTime)
	if startTime.IsZero() || endTime.IsZero() || startTime.Equal(endTime) || startTime.After(endTime) {
		return ModTimesByPath{}
	}

	filteredModTimesByPath := ModTimesByPath{}

	for pth, modTime := range modTimesByPath {
		if isInTimeInterval(modTime, startTime, endTime) {
			filteredModTimesByPath[pth] = modTime
			log.Debugf("%s is in interval", pth)
		} else {
			log.Debugf("%s is not in interval, filtering out", pth)
		}
	}

	return filteredModTimesByPath
}

// finds the last modified file matching to most strict regexp
// order of regexps should be: most strict -> less strict
func findLastModifiedPathWithFileNameRegexps(modTimesByPath ModTimesByPath, regexps ...*regexp.Regexp) string {
	if len(modTimesByPath) == 0 {
		log.Debugf("Mod times by path is empty")
		return ""
	}

	var lastModifiedPth string
	var lastModTime time.Time

	if len(regexps) > 0 {
		for _, re := range regexps {
			log.Debugf("Checking match for regexp: %s", re)
			for pth, modTime := range modTimesByPath {
				log.Debugf("Checking file at: %s mod time %s", pth, modTime)
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
		log.Debugf("Regexps is empty")
		for pth, modTime := range modTimesByPath {
			log.Debugf("Checking mod time: %s. Last mod time: %s", modTime, lastModTime)
			if modTime.After(lastModTime) {
				lastModifiedPth = pth
				lastModTime = modTime
			}
		}
	}

	return lastModifiedPth
}

// exports the last modified file matching to most strict regexps within a time window
// order of regexps should be: most strict -> less strict. Boolean excludeDirs indicates that the function should search
// for directories as well or not. Please note, that for example a .xcarchive file qualifies as a directory, so if you
// want to find it, the boolean should be false.
func findArtifact(dir string, startTime, endTime time.Time, excludeDirs bool, patterns ...string) (string, error) {
	log.Debugf("Searching at %s", dir)
	regexps := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		regexps[i] = regexp.MustCompile(pattern)
	}

	modTimesByPath, err := findModTimesByPath(dir, excludeDirs)
	if err != nil {
		return "", err
	}

	modTimesByPathByTimeWindow := filterModTimesByPathByTimeWindow(modTimesByPath, startTime, endTime)
	return findLastModifiedPathWithFileNameRegexps(modTimesByPathByTimeWindow, regexps...), nil
}

func exportApk(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime, false,
		fmt.Sprintf(`(?i).*%s.*signed.*\.apk$`, assemblyName),
		fmt.Sprintf(`(?i).*%s.*\.apk$`, assemblyName),
		`(?i).*signed.*\.apk$`,
		`(?i).*\.apk$`,
	)
}

func exportIpa(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime, true,
		fmt.Sprintf(`(?i).*%s.*\.ipa$`, assemblyName),
		`(?i).*\.ipa$`,
	)
}

func exportXCArchive(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime, false,
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
	return findArtifact(outputDir, startTime, endTime, false,
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
	return findArtifact(outputDir, startTime, endTime, false,
		fmt.Sprintf(`(?i).*%s.*\.pkg$`, assemblyName),
		`(?i).*\.pkg$`,
	)
}

func exportApp(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime, false,
		fmt.Sprintf(`(?i).*%s.*\.app$`, assemblyName),
		`(?i).*\.app$`,
	)
}

func exportDLL(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	return findArtifact(outputDir, startTime, endTime, true,
		fmt.Sprintf(`(?i).*%s.*\.dll$`, assemblyName),
		`(?i).*\.dll$`,
	)
}
