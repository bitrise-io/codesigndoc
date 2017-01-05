package xcode

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// XcodebuildVersionModel ...
type XcodebuildVersionModel struct {
	Version      string
	BuildVersion string
	MajorVersion int64
}

// GetXcodeVersion ...
func GetXcodeVersion() (XcodebuildVersionModel, error) {
	cmd := exec.Command("xcodebuild", "-version")
	outBytes, err := cmd.CombinedOutput()
	outStr := string(outBytes)
	if err != nil {
		return XcodebuildVersionModel{}, fmt.Errorf("xcodebuild -version failed, err: %s, details: %s", err, outStr)
	}

	split := strings.Split(outStr, "\n")
	if len(split) == 0 {
		return XcodebuildVersionModel{}, fmt.Errorf("failed to parse xcodebuild version output (%s)", outStr)
	}

	xcodebuildVersion := split[0]
	buildVersion := split[1]

	split = strings.Split(xcodebuildVersion, " ")
	if len(split) != 2 {
		return XcodebuildVersionModel{}, fmt.Errorf("failed to parse xcodebuild version output (%s)", outStr)
	}

	version := split[1]

	split = strings.Split(version, ".")
	majorVersionStr := split[0]

	majorVersion, err := strconv.ParseInt(majorVersionStr, 10, 32)
	if err != nil {
		return XcodebuildVersionModel{}, fmt.Errorf("failed to parse xcodebuild version output (%s), error: %s", outStr, err)
	}

	return XcodebuildVersionModel{
		Version:      xcodebuildVersion,
		BuildVersion: buildVersion,
		MajorVersion: majorVersion,
	}, nil
}
