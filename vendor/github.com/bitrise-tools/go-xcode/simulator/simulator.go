package simulator

import (
	"bufio"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-tools/go-xcode/models"
	version "github.com/hashicorp/go-version"
)

// InfoModel ...
type InfoModel struct {
	Name        string
	ID          string
	Status      string
	StatusOther string
}

// OsVersionSimulatorInfosMap ...
type OsVersionSimulatorInfosMap map[string][]InfoModel // Os version - []Info map

// getSimulatorInfoFromLine ...
// a simulator info line should look like this:
//  iPhone 5s (EA1C7E48-8137-428C-A0A5-B2C63FF276EB) (Shutdown)
// or
//  iPhone 4s (51B10EBD-C949-49F5-A38B-E658F41640FF) (Shutdown) (unavailable, runtime profile not found)
func getSimulatorInfoFromLine(lineStr string) (InfoModel, error) {
	baseInfosExp := regexp.MustCompile(`(?P<deviceName>[a-zA-Z].*[a-zA-Z0-9 -]*) \((?P<simulatorID>[a-zA-Z0-9-]{36})\) \((?P<status>[a-zA-Z]*)\)`)
	baseInfosRes := baseInfosExp.FindStringSubmatch(lineStr)
	if baseInfosRes == nil {
		return InfoModel{}, fmt.Errorf("No match found")
	}

	simInfo := InfoModel{
		Name:   baseInfosRes[1],
		ID:     baseInfosRes[2],
		Status: baseInfosRes[3],
	}

	// StatusOther
	restOfTheLine := lineStr[len(baseInfosRes[0]):]
	if len(restOfTheLine) > 0 {
		statusOtherExp := regexp.MustCompile(`\((?P<statusOther>[a-zA-Z ,]*)\)`)
		statusOtherRes := statusOtherExp.FindStringSubmatch(restOfTheLine)
		if statusOtherRes != nil {
			simInfo.StatusOther = statusOtherRes[1]
		}
	}
	return simInfo, nil
}

func getOsVersionSimulatorInfosMapFromSimctlList(simctlList string) (OsVersionSimulatorInfosMap, error) {
	simulatorsByIOSVersions := OsVersionSimulatorInfosMap{}
	currIOSVersion := ""

	fscanner := bufio.NewScanner(strings.NewReader(simctlList))
	isDevicesSectionFound := false
	for fscanner.Scan() {
		aLine := fscanner.Text()

		if aLine == "== Devices ==" {
			isDevicesSectionFound = true
			continue
		}

		if !isDevicesSectionFound {
			continue
		}
		if strings.HasPrefix(aLine, "==") {
			isDevicesSectionFound = false
			continue
		}
		if strings.HasPrefix(aLine, "--") {
			iosVersionSectionExp := regexp.MustCompile(`-- (?P<iosVersionSection>.*) --`)
			iosVersionSectionRes := iosVersionSectionExp.FindStringSubmatch(aLine)
			if iosVersionSectionRes != nil {
				currIOSVersion = iosVersionSectionRes[1]
			}
			continue
		}

		simInfo, err := getSimulatorInfoFromLine(aLine)
		if err != nil {
			fmt.Println(" [!] Error scanning the line for Simulator info: ", err)
		}

		currIOSVersionSimList := simulatorsByIOSVersions[currIOSVersion]
		currIOSVersionSimList = append(currIOSVersionSimList, simInfo)
		simulatorsByIOSVersions[currIOSVersion] = currIOSVersionSimList
	}

	return simulatorsByIOSVersions, nil
}

// GetOsVersionSimulatorInfosMap ...
func GetOsVersionSimulatorInfosMap() (OsVersionSimulatorInfosMap, error) {
	cmd := command.New("xcrun", "simctl", "list")
	simctlListOut, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return OsVersionSimulatorInfosMap{}, err
	}

	return getOsVersionSimulatorInfosMapFromSimctlList(simctlListOut)
}

func getSimulatorInfoFromSimctlOut(simctlListOut, osNameAndVersion, deviceName string) (InfoModel, error) {
	osVersionSimulatorInfosMap, err := getOsVersionSimulatorInfosMapFromSimctlList(simctlListOut)
	if err != nil {
		return InfoModel{}, err
	}

	infos, ok := osVersionSimulatorInfosMap[osNameAndVersion]
	if !ok {
		return InfoModel{}, fmt.Errorf("no simulators found for os version: %s", osNameAndVersion)
	}

	for _, info := range infos {
		if info.Name == deviceName {
			return info, nil
		}
	}

	return InfoModel{}, fmt.Errorf("no simulators found for os version: (%s), device name: (%s)", osNameAndVersion, deviceName)
}

// GetSimulatorInfo ...
func GetSimulatorInfo(osNameAndVersion, deviceName string) (InfoModel, error) {
	cmd := command.New("xcrun", "simctl", "list")
	simctlListOut, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return InfoModel{}, err
	}

	return getSimulatorInfoFromSimctlOut(simctlListOut, osNameAndVersion, deviceName)
}

func getLatestSimulatorInfoFromSimctlOut(simctlListOut, osName, deviceName string) (InfoModel, string, error) {
	osVersionSimulatorInfosMap, err := getOsVersionSimulatorInfosMapFromSimctlList(simctlListOut)
	if err != nil {
		return InfoModel{}, "", err
	}

	var latestVersionPtr *version.Version
	latestInfo := InfoModel{}
	for osVersion, infos := range osVersionSimulatorInfosMap {
		if !strings.HasPrefix(osVersion, osName) {
			continue
		}

		deviceInfo := InfoModel{}
		deviceFound := false
		for _, info := range infos {
			if info.Name == deviceName {
				deviceFound = true
				deviceInfo = info
				break
			}
		}
		if !deviceFound {
			continue
		}

		versionStr := strings.TrimPrefix(osVersion, osName)
		versionStr = strings.TrimSpace(versionStr)

		versionPtr, err := version.NewVersion(versionStr)
		if err != nil {
			return InfoModel{}, "", fmt.Errorf("failed to parse version (%s), error: %s", versionStr, err)
		}

		if latestVersionPtr == nil || versionPtr.GreaterThan(latestVersionPtr) {
			latestVersionPtr = versionPtr
			latestInfo = deviceInfo
		}
	}

	if latestVersionPtr == nil {
		return InfoModel{}, "", fmt.Errorf("failed to determin latest (%s) simulator version", osName)
	}

	versionSegments := latestVersionPtr.Segments()
	if len(versionSegments) < 2 {
		return InfoModel{}, "", fmt.Errorf("invalid version created: %s, segments count < 2", latestVersionPtr.String())
	}

	osVersion := fmt.Sprintf("%s %d.%d", osName, versionSegments[0], versionSegments[1])

	return latestInfo, osVersion, nil
}

// GetLatestSimulatorInfoAndVersion ...
func GetLatestSimulatorInfoAndVersion(osName, deviceName string) (InfoModel, string, error) {
	cmd := command.New("xcrun", "simctl", "list")
	simctlListOut, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return InfoModel{}, "", err
	}

	return getLatestSimulatorInfoFromSimctlOut(simctlListOut, osName, deviceName)
}

// Is64BitArchitecture ...
func Is64BitArchitecture(simulatorDevice string) (bool, error) {
	// 64 bit processor iPhones:
	// - iPhone 5S,
	// - iPhone 6, iPhone 6 Plus, iPhone 6S, iPhone 6S Plus,
	// - iPhone SE,
	// - iPhone 7, iPhone 7 Plus

	// 64 bit processor iPads:
	// - iPad Mini 2, iPad Mini 3, iPad Mini 4
	// - iPad Air, iPad Air 2
	// - iPad Pro (12.9 inch), iPad Pro (9.7 inch)
	deviceSplit := strings.Split(simulatorDevice, " ")
	if len(deviceSplit) == 1 && deviceSplit[0] == "iPad" {
		return false, nil
	}

	if len(deviceSplit) < 2 {
		return false, fmt.Errorf("Unexpected deivice name (%s)", simulatorDevice)
	}

	name := strings.TrimSpace(deviceSplit[0])
	versionSlice := deviceSplit[1:]

	if name == "iPhone" {
		if versionSlice[0] == "SE" {
			return true, nil
		}

		versionNumber := versionSlice[0]
		if versionNumber == "5S" {
			return true, nil
		}

		majorVersionStr := string(versionNumber[0])
		majorVersion, err := strconv.Atoi(majorVersionStr)
		if err != nil {
			return false, err
		}

		if majorVersion >= 6 {
			return true, nil
		}
	} else if name == "iPad" {
		subNameOrVersion := strings.TrimSpace(versionSlice[0])

		if subNameOrVersion == "Mini" {
			if len(versionSlice) == 2 {
				version, err := strconv.Atoi(versionSlice[1])
				if err != nil {
					return false, err
				}

				if version > 1 {
					return true, nil
				}
			}
		}

		if subNameOrVersion == "Air" {
			return true, nil
		}

		if subNameOrVersion == "Pro" {
			return true, nil
		}
	}

	return false, nil
}

func getXcodeDeveloperDirPath() (string, error) {
	cmd := command.New("xcode-select", "--print-path")
	return cmd.RunAndReturnTrimmedCombinedOutput()
}

// BootSimulator ...
func BootSimulator(simulator InfoModel, xcodebuildVersion models.XcodebuildVersionModel) error {
	simulatorApp := "Simulator"
	if xcodebuildVersion.MajorVersion == 6 {
		simulatorApp = "iOS Simulator"
	}
	xcodeDevDirPth, err := getXcodeDeveloperDirPath()
	if err != nil {
		return fmt.Errorf("failed to get Xcode Developer Directory - most likely Xcode.app is not installed")
	}
	simulatorAppFullPath := filepath.Join(xcodeDevDirPth, "Applications", simulatorApp+".app")

	openCmd := command.New("open", simulatorAppFullPath, "--args", "-CurrentDeviceUDID", simulator.ID)

	log.Printf("$ %s", openCmd.PrintableCommandArgs())

	outStr, err := openCmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start simulators (%s), output: %s, error: %s", simulator.ID, outStr, err)
	}

	return nil
}
