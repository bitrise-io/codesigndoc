package utility

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-xcode/models"
)

func getXcodeVersionFromXcodebuildOutput(outStr string) (models.XcodebuildVersionModel, error) {
	split := strings.Split(outStr, "\n")
	if len(split) == 0 {
		return models.XcodebuildVersionModel{}, fmt.Errorf("failed to parse xcodebuild version output (%s)", outStr)
	}

	xcodebuildVersion := split[0]
	buildVersion := split[1]

	split = strings.Split(xcodebuildVersion, " ")
	if len(split) != 2 {
		return models.XcodebuildVersionModel{}, fmt.Errorf("failed to parse xcodebuild version output (%s)", outStr)
	}

	version := split[1]

	split = strings.Split(version, ".")
	majorVersionStr := split[0]

	majorVersion, err := strconv.ParseInt(majorVersionStr, 10, 32)
	if err != nil {
		return models.XcodebuildVersionModel{}, fmt.Errorf("failed to parse xcodebuild version output (%s), error: %s", outStr, err)
	}

	return models.XcodebuildVersionModel{
		Version:      xcodebuildVersion,
		BuildVersion: buildVersion,
		MajorVersion: majorVersion,
	}, nil
}

// GetXcodeVersion ...
func GetXcodeVersion() (models.XcodebuildVersionModel, error) {
	cmd := command.New("xcodebuild", "-version")
	outStr, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return models.XcodebuildVersionModel{}, fmt.Errorf("xcodebuild -version failed, err: %s, details: %s", err, outStr)
	}
	return getXcodeVersionFromXcodebuildOutput(outStr)
}
