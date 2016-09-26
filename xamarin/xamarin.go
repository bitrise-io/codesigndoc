package xamarin

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/maputil"
	"github.com/bitrise-io/go-utils/progress"
	"github.com/bitrise-io/go-utils/regexputil"
	"github.com/bitrise-tools/codesigndoc/common"
	"github.com/bitrise-tools/codesigndoc/provprofile"
)

// CommandModel ...
type CommandModel struct {
	SolutionFilePath  string
	ProjectName       string
	ConfigurationName string
}

// ScanCodeSigningSettings ...
func (xamarinCmd CommandModel) ScanCodeSigningSettings() (common.CodeSigningSettings, string, error) {
	cmdOut := ""
	var err error

	progress.SimpleProgress(".", 1*time.Second, func() {
		cmdOut, err = xamarinCmd.RunBuildCommand()
	})
	fmt.Println()

	if err != nil {
		return common.CodeSigningSettings{}, cmdOut, fmt.Errorf("Failed to Archive, error: %s", err)
	}

	return parseCodeSigningSettingsFromOutput(cmdOut), cmdOut, nil
}

// RunBuildCommand ...
func (xamarinCmd CommandModel) RunBuildCommand() (string, error) {
	mdtoolPth := "/Applications/Xamarin Studio.app/Contents/MacOS/mdtool"
	cmdArgs := []string{mdtoolPth, "build",
		xamarinCmd.SolutionFilePath,
		fmt.Sprintf("-c:%s", xamarinCmd.ConfigurationName),
		fmt.Sprintf("-p:%s", xamarinCmd.ProjectName),
	}
	log.Infof("$ %s", cmdex.PrintableCommandArgs(true, cmdArgs))
	fmt.Print("Running and analyzing log ...")
	cmd, err := cmdex.NewCommandFromSlice(cmdArgs)
	if err != nil {
		return "", fmt.Errorf("Failed to create Xamarin command, error: %s", err)
	}
	xcoutput, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return xcoutput, fmt.Errorf("Failed to run Xamarin command, error: %s", err)
	}

	log.Debugf("xcoutput: %s", xcoutput)
	return xcoutput, nil
}

func parseCodeSigningSettingsFromOutput(logOutput string) common.CodeSigningSettings {
	scanner := bufio.NewScanner(strings.NewReader(logOutput))

	identitiesMap := map[string]common.CodeSigningIdentityInfo{}
	provProfilesMap := map[string]provprofile.ProvisioningProfileInfo{}
	teamIDsMap := map[string]interface{}{}
	appBundleIDsMap := map[string]interface{}{}
	for scanner.Scan() {
		line := scanner.Text()

		// App ID
		if rexp := regexp.MustCompile(`^[[:space:]]*App Id: (?P<appid>.+)$`); rexp.MatchString(line) {
			results, err := regexputil.NamedFindStringSubmatch(rexp, line)
			if err != nil {
				log.Errorf("Failed to scan App Bundle ID: %s", err)
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
			appBundleIDsMap[appID] = 1
		}

		// Signing Identity
		if rexp := regexp.MustCompile(`^[[:space:]]*Code Signing Key: "(?P<title>.+)" \((?P<identityid>[a-zA-Z0-9]+)\)$`); rexp.MatchString(line) {
			results, err := regexputil.NamedFindStringSubmatch(rexp, line)
			if err != nil {
				log.Errorf("Failed to scan Signing Identity title: %s", err)
				continue
			}
			codeSigningID := common.CodeSigningIdentityInfo{Title: results["title"]}
			identitiesMap[codeSigningID.Title] = codeSigningID
		}
		// Prov. Profile - title line
		if rexp := regexp.MustCompile(`^[[:space:]]*Provisioning Profile: "(?P<title>.+)" \((?P<uuid>[a-zA-Z0-9-]+)\)$`); rexp.MatchString(line) {
			results, err := regexputil.NamedFindStringSubmatch(rexp, line)
			if err != nil {
				log.Errorf("Failed to scan Provisioning Profile: %s", err)
				continue
			}
			tmpProvProfile := provprofile.ProvisioningProfileInfo{Title: results["title"]}
			tmpProvProfile.UUID = results["uuid"]
			provProfilesMap[tmpProvProfile.Title] = tmpProvProfile
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
	appBundleIDs := maputil.KeysOfStringInterfaceMap(appBundleIDsMap)

	return common.CodeSigningSettings{
		Identities:   identities,
		ProvProfiles: provProfiles,
		TeamIDs:      teamIDs,
		AppBundleIDs: appBundleIDs,
	}
}
