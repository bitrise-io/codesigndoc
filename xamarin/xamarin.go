package xamarin

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/maputil"
	"github.com/bitrise-io/go-utils/progress"
	"github.com/bitrise-io/go-utils/readerutil"
	"github.com/bitrise-io/go-utils/regexputil"
	"github.com/bitrise-tools/codesigndoc/common"
	"github.com/bitrise-tools/codesigndoc/provprofile"
	"github.com/bitrise-tools/go-xamarin/constants"
)

// CommandModel ...
type CommandModel struct {
	SolutionFilePath  string
	ProjectName       string
	ConfigurationName string
}

// GenerateLog ...
func (xamarinCmd CommandModel) GenerateLog() (string, error) {
	cmdOut := ""
	var err error

	progress.SimpleProgress(".", 1*time.Second, func() {
		cmdOut, err = xamarinCmd.RunBuildCommand()
	})
	fmt.Println()

	if err != nil {
		return cmdOut, fmt.Errorf("Failed to Archive, error: %s", err)
	}

	return cmdOut, nil
}

// ScanCodeSigningSettings ...
func (xamarinCmd CommandModel) ScanCodeSigningSettings(logToScan string) (common.CodeSigningSettings, error) {
	return parseCodeSigningSettingsFromOutput(logToScan)
}

// RunBuildCommand ...
func (xamarinCmd CommandModel) RunBuildCommand() (string, error) {
	split := strings.Split(xamarinCmd.ConfigurationName, "|")
	if len(split) != 2 {
		return "", fmt.Errorf("failed to parse configuration: %s", xamarinCmd.ConfigurationName)
	}
	configuration := split[0]
	platform := split[1]
	projectName := strings.Replace(xamarinCmd.ProjectName, ".", "_", -1)

	cmdArgs := []string{constants.MsbuildPath,
		xamarinCmd.SolutionFilePath,
		fmt.Sprintf("/p:Configuration=%s", configuration),
		fmt.Sprintf("/p:Platform=%s", platform),
		fmt.Sprintf("/t:%s", projectName),
	}

	log.Infof("$ %s", command.PrintableCommandArgs(true, cmdArgs))
	fmt.Print("Running and analyzing log ...")
	cmd, err := command.NewFromSlice(cmdArgs)
	if err != nil {
		return "", fmt.Errorf("Failed to create Xamarin command, error: %s", err)
	}
	xamarinBuildOutput, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return xamarinBuildOutput, fmt.Errorf("Failed to run Xamarin command, error: %s", err)
	}

	log.Debugf("xamarinBuildOutput: %s", xamarinBuildOutput)
	return xamarinBuildOutput, nil
}

func parseCodeSigningSettingsFromOutput(logOutput string) (common.CodeSigningSettings, error) {
	logReader := bufio.NewReader(strings.NewReader(logOutput))

	identitiesMap := map[string]common.CodeSigningIdentityInfo{}
	provProfilesMap := map[string]provprofile.ProvisioningProfileInfo{}
	teamIDsMap := map[string]interface{}{}
	appIDsMap := map[string]interface{}{}

	// scan log line by line
	{
		line, readErr := readerutil.ReadLongLine(logReader)
		for ; readErr == nil; line, readErr = readerutil.ReadLongLine(logReader) {

			// App ID
			if rexp := regexp.MustCompile(`^[[:space:]]*App Id: (?P<appid>.+)$`); rexp.MatchString(line) {
				results, isFound := regexputil.NamedFindStringSubmatch(rexp, line)
				if !isFound {
					log.Error("Failed to scan App Bundle ID: not found in the logs")
					continue
				}
				appID := results["appid"]
				comps := strings.Split(appID, ".")
				if len(comps) < 2 {
					log.Errorf("Invalid App ID, does not include '.': %s", appID)
					continue
				}
				teamID := comps[0]
				if teamID == "" {
					log.Errorf("Invalid App ID, Team ID was empty: %s", appID)
					continue
				}
				teamIDsMap[teamID] = 1
				appIDsMap[appID] = 1
			}

			// Signing Identity
			if rexp := regexp.MustCompile(`^[[:space:]]*Code Signing Key: "(?P<title>.+)" \((?P<identityid>[a-zA-Z0-9]+)\)$`); rexp.MatchString(line) {
				results, isFound := regexputil.NamedFindStringSubmatch(rexp, line)
				if !isFound {
					log.Error("Failed to scan Signing Identity title: not found in the logs")
					continue
				}
				codeSigningID := common.CodeSigningIdentityInfo{Title: results["title"]}
				identitiesMap[codeSigningID.Title] = codeSigningID
			}
			// Prov. Profile - title line
			if rexp := regexp.MustCompile(`^[[:space:]]*Provisioning Profile: "(?P<title>.+)" \((?P<uuid>[a-zA-Z0-9-]+)\)$`); rexp.MatchString(line) {
				results, isFound := regexputil.NamedFindStringSubmatch(rexp, line)
				if !isFound {
					log.Error("Failed to scan Provisioning Profile: not found in the logs")
					continue
				}
				tmpProvProfile := provprofile.ProvisioningProfileInfo{Title: results["title"]}
				tmpProvProfile.UUID = results["uuid"]
				provProfilesMap[tmpProvProfile.Title] = tmpProvProfile
			}
		}
		if readErr != nil && readErr != io.EOF {
			return common.CodeSigningSettings{}, fmt.Errorf("Failed to scan log output, error: %s", readErr)
		}
	}

	identities := []common.CodeSigningIdentityInfo{}
	for _, v := range identitiesMap {
		identities = append(identities, v)
	}
	provProfiles := []provprofile.ProvisioningProfileInfo{}
	for _, v := range provProfilesMap {
		provProfiles = append(provProfiles, v)
	}
	teamIDs := maputil.KeysOfStringInterfaceMap(teamIDsMap)
	appIDs := maputil.KeysOfStringInterfaceMap(appIDsMap)

	return common.CodeSigningSettings{
		Identities:   identities,
		ProvProfiles: provProfiles,
		TeamIDs:      teamIDs,
		AppIDs:       appIDs,
	}, nil
}
