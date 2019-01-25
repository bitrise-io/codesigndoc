package codesigndocuitests

import (
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/codesign"
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/export"
	"github.com/bitrise-tools/go-xcode/exportoptions"
	"github.com/bitrise-tools/go-xcode/plistutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

// extractCertificatesAndProfiles returns the certificates and provisioning profiles of the given codesign group
func extractCertificatesAndProfiles(codeSignGroups ...export.CodeSignGroup) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel) {
	certificateMap := map[string]certificateutil.CertificateInfoModel{}
	profilesMap := map[string]profileutil.ProvisioningProfileInfoModel{}
	for _, group := range codeSignGroups {
		certificate := group.Certificate()
		certificateMap[certificate.Serial] = certificate

		for _, profile := range group.BundleIDProfileMap() {
			profilesMap[profile.UUID] = profile
		}
	}

	certificates := []certificateutil.CertificateInfoModel{}
	profiles := []profileutil.ProvisioningProfileInfoModel{}
	for _, certificate := range certificateMap {
		certificates = append(certificates, certificate)
	}
	for _, profile := range profilesMap {
		profiles = append(profiles, profile)
	}
	return certificates, profiles
}

// codesignMethod returns which code sign method type is allowed by the given codesign group
func codesignMethod(group export.CodeSignGroup) string {
	for _, profile := range group.BundleIDProfileMap() {
		return string(profile.ExportType)
	}
	return ""
}

// printCodesignGroup prints the given codesign group
func printCodesignGroup(group export.CodeSignGroup) {
	fmt.Printf("%s %s (%s)\n", colorstring.Green("development team:"), group.Certificate().TeamName, group.Certificate().TeamID)
	fmt.Printf("%s %s [%s]\n", colorstring.Green("codesign identity:"), group.Certificate().CommonName, group.Certificate().Serial)

	idx := -1
	for bundleID, profile := range group.BundleIDProfileMap() {
		idx++
		if idx == 0 {
			fmt.Printf("%s %s -> %s\n", colorstring.Greenf("provisioning profiles:"), profile.Name, bundleID)
		} else {
			fmt.Printf("%s%s -> %s\n", strings.Repeat(" ", len("provisioning profiles: ")), profile.Name, bundleID)
		}
	}
}

// collectExportCertificate returns the certificate to use for the UITest-Runner
func collectExportCertificate(installedCertificates []certificateutil.CertificateInfoModel) ([]certificateutil.CertificateInfoModel, error) {
	var selectedCertificates []certificateutil.CertificateInfoModel

	// Codesign method
	codesignMethods := []string{"development", "app-store", "ad-hoc", "enterprise"}

	// Asking the user over and over until we find a valid certificate for the selected export method.
	for searchingValidCertificate := true; searchingValidCertificate; {
		fmt.Println()
		selectedCodeSignMethod, err := goinp.SelectFromStringsWithDefault("Select the code signing method", 1, codesignMethods)
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %s", err)
		}

		log.Debugf("selected export method: %v", selectedCodeSignMethod)

		selectedCertificates, err = filterCertificates(selectedCodeSignMethod, "", selectedCertificates, installedCertificates)
		if err != nil {
			return nil, err
		}

		fmt.Println()
		question := `Do you want to collect another certificate?`
		searchingValidCertificate, err = goinp.AskForBoolWithDefault(question, true)
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %s", err)
		}
	}

	return selectedCertificates, nil
}

func filterCertificates(selectedCodeSignMethod, selectedTeam string, selectedCertificates []certificateutil.CertificateInfoModel, installedCertificates []certificateutil.CertificateInfoModel) ([]certificateutil.CertificateInfoModel, error) {
	var certsForSelectedCodeSign []certificateutil.CertificateInfoModel
	var err error
	log.Debugf("InstalledCerts: %v\n", installedCertificates)

	// Filter the installed certificates by distribution type
	switch selectedCodeSignMethod {
	case "development":
		certsForSelectedCodeSign = certificateutil.FilterCertificateInfoModelsByFilterFunc(installedCertificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return !codesign.IsDistributionCertificate(certInfo)
		})

		log.Debugf("DeveloperDistribution certificates: %v\n", certsForSelectedCodeSign)
		break
	default:
		certsForSelectedCodeSign = certificateutil.FilterCertificateInfoModelsByFilterFunc(installedCertificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return codesign.IsDistributionCertificate(certInfo)
		})
		log.Debugf("Distribution certificates: %v\n", certsForSelectedCodeSign)
		break
	}

	filteredCertificatesByTeam := codesign.MapCertificatesByTeam(certsForSelectedCodeSign)
	log.Debugf("Filtered certificates (by distribution type) by team: %v\n", filteredCertificatesByTeam)

	if len(filteredCertificatesByTeam) == 0 {
		log.Warnf("🚨  We couldn't find any certificate for the %s export method.", selectedCodeSignMethod)
		return nil, nil
	}

	// If we already selected a team, we can skip it. (e.g mac app-store export)
	if selectedTeam == "" {
		// Use different team for export than archive.
		teams := []string{}
		for team := range filteredCertificatesByTeam {
			if hasCertificateForDistType(selectedCodeSignMethod, filteredCertificatesByTeam[team]) {
				teams = append(teams, team)
			}
		}

		fmt.Println()
		selectedTeam, err = goinp.SelectFromStringsWithDefault("Select the Development team to export your app", 1, teams)
		if err != nil {
			return selectedCertificates, fmt.Errorf("failed to read input: %s", err)
		}
	}

	// Find the specific development certificate.
	filteredTeamCertificates := filteredCertificatesByTeam[selectedTeam]
	certificateOptions := []string{}

	for _, certInfo := range filteredTeamCertificates {
		certificateOption := fmt.Sprintf("%s [%s]", certInfo.CommonName, certInfo.Serial)
		certificateOptions = append(certificateOptions, certificateOption)
	}

	certType := "distribution"
	if selectedCodeSignMethod == "development" {
		certType = "development"
	}

	fmt.Println()
	question := fmt.Sprintf("Please select a %s certificate:", certType)
	selectedCertificateOption, err := goinp.SelectFromStringsWithDefault(question, 1, certificateOptions)
	if err != nil {
		return selectedCertificates, fmt.Errorf("failed to read input: %s", err)
	}

	for _, certInfo := range filteredTeamCertificates {
		certificateOption := fmt.Sprintf("%s [%s]", certInfo.CommonName, certInfo.Serial)
		if certificateOption == selectedCertificateOption {
			selectedCertificates = append(selectedCertificates, certInfo)
		}
	}

	return selectedCertificates, nil
}

// collectExportCodeSignGroups returns the codesigngroups required for the UITest target with the selected code signing methods
func collectExportCodeSignGroups(testRunner IOSTestRunner, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) ([]export.CodeSignGroup, error) {
	collectedCodeSignGroups := []export.CodeSignGroup{}

	codeSignGroups := collectExportSelectableCodeSignGroups(testRunner, installedCertificates, installedProfiles)
	if len(codeSignGroups) == 0 {
		return nil, errors.New("no code sign files (Codesign Identities and Provisioning Profiles) are installed to sing the UITest target\n" + collectCodesigningFilesInfo)
	}

	testRunnerID, _ := testRunner.InfoPlist.GetString("CFBundleIdentifier")
	fmt.Println()
	log.Infof("Code signing for target with %s bundle ID", strings.TrimRight(testRunnerID, "-Runner"))

	codeSignMethods := []string{"development", "app-store", "ad-hoc", "enterprise"}
	for true {
		selectedCodeSignMethod, err := goinp.SelectFromStringsWithDefault("Select the code signing method", 1, codeSignMethods)
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %s", err)
		}
		log.Debugf("selected export method: %v", selectedCodeSignMethod)

		fmt.Println()
		filteredCodeSignGroups := export.FilterSelectableCodeSignGroups(codeSignGroups,
			export.CreateExportMethodSelectableCodeSignGroupFilter(exportoptions.Method(selectedCodeSignMethod)),
		)

		log.Debugf("\n")
		log.Debugf("Filtered Codesign Groups:")
		for _, group := range codeSignGroups {
			log.Debugf(group.String())
		}

		if len(filteredCodeSignGroups) == 0 {
			fmt.Println()
			log.Errorf(collectCodesigningFilesInfo)
			fmt.Println()
			question := fmt.Sprintf("Do you want to collect other  code sign files for (%s)", testRunnerID)
			question += "\n(select NO to finish collecting codesign files and continue)"
			anotherExport, err := goinp.AskForBoolWithDefault(question, false)
			if err != nil {
				return nil, fmt.Errorf("failed to read input: %s", err)
			}
			if !anotherExport {
				break
			}
			continue
		}

		// Select certificate
		certificates := []certificateutil.CertificateInfoModel{}
		certificateOptions := []string{}
		for _, group := range filteredCodeSignGroups {
			certificate := group.Certificate
			certificates = append(certificates, certificate)
			certificateOption := fmt.Sprintf("%s [%s] - development team: %s", certificate.CommonName, certificate.Serial, certificate.TeamName)
			certificateOptions = append(certificateOptions, certificateOption)
		}

		selectedCertificateOption := ""
		if len(certificateOptions) == 1 {
			selectedCertificateOption = certificateOptions[0]

			fmt.Printf("Codesign Indentity for %s signing: %s\n", selectedCodeSignMethod, selectedCertificateOption)
		} else {
			sort.Strings(certificateOptions)

			question := fmt.Sprintf("Select the Codesign Indentity for %s method", selectedCodeSignMethod)
			selectedCertificateOption, err = goinp.SelectFromStringsWithDefault(question, 1, certificateOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to read input: %s", err)
			}
		}

		var selectedCertificate *certificateutil.CertificateInfoModel
		for _, certificate := range certificates {
			option := fmt.Sprintf("%s [%s] - development team: %s", certificate.CommonName, certificate.Serial, certificate.TeamName)
			if option == selectedCertificateOption {
				selectedCertificate = &certificate
				break
			}
		}
		if selectedCertificate == nil {
			return nil, errors.New("failed to find selected Codesign Indentity")
		}

		// Select Profiles
		bundleIDProfilesMap := map[string][]profileutil.ProvisioningProfileInfoModel{}
		for _, group := range filteredCodeSignGroups {
			option := fmt.Sprintf("%s [%s] - development team: %s", group.Certificate.CommonName, group.Certificate.Serial, group.Certificate.TeamName)
			if option == selectedCertificateOption {
				bundleIDProfilesMap = group.BundleIDProfilesMap
				break
			}
		}
		if len(bundleIDProfilesMap) == 0 {
			return nil, errors.New("failed to find Provisioning Profiles for Code Sign Identity")
		}

		selectedBundleIDProfileMap := map[string]profileutil.ProvisioningProfileInfoModel{}
		for bundleID, profiles := range bundleIDProfilesMap {
			profiles = codesign.FilterLatestProfiles(profiles)
			profileOptions := []string{}
			for _, profile := range profiles {
				profileOption := fmt.Sprintf("%s (%s)", profile.Name, profile.UUID)
				profileOptions = append(profileOptions, profileOption)
			}

			selectedProfileOption := ""
			if len(profileOptions) == 1 {
				selectedProfileOption = profileOptions[0]

				fmt.Printf("Provisioning Profile to sign target (%s): %s\n", bundleID, selectedProfileOption)
			} else {
				sort.Strings(profileOptions)

				fmt.Println()
				question := fmt.Sprintf("Select the Provisioning Profile to sign target with bundle ID: %s", bundleID)
				selectedProfileOption, err = goinp.SelectFromStringsWithDefault(question, 1, profileOptions)
				if err != nil {
					return nil, fmt.Errorf("failed to read input: %s", err)
				}
			}

			for _, profile := range profiles {
				option := fmt.Sprintf("%s (%s)", profile.Name, profile.UUID)
				if option == selectedProfileOption {
					selectedBundleIDProfileMap[bundleID] = profile
				}
			}
		}
		if len(selectedBundleIDProfileMap) != len(bundleIDProfilesMap) {
			return nil, fmt.Errorf("failed to find Provisioning Profiles for UITest target signing")
		}

		var collectedCodeSignGroup export.CodeSignGroup

		collectedCodeSignGroup = export.NewIOSGroup(*selectedCertificate, selectedBundleIDProfileMap)

		fmt.Println()
		log.Infof("Codesign settings will be used for %s method:", codesignMethod(collectedCodeSignGroup))
		printCodesignGroup(collectedCodeSignGroup)

		collectedCodeSignGroups = append(collectedCodeSignGroups, collectedCodeSignGroup)

		fmt.Println()
		question := fmt.Sprintf("Do you want to collect other code sign files for (%s)", testRunnerID)
		question += "\n(select NO to finish collecting codesign files and continue)"
		anotherExport, err := goinp.AskForBoolWithDefault(question, false)
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %s", err)
		}
		if !anotherExport {
			break
		}
	}

	return collectedCodeSignGroups, nil
}

// collectExportSelectableCodeSignGroups returns every possible codesigngroup which can be used to sign the UITest-Runner
func collectExportSelectableCodeSignGroups(testRunner IOSTestRunner, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) []export.SelectableCodeSignGroup {
	bundleIDEntitlemenstMap := map[string]plistutil.PlistData{}
	bundleIDEntitlemenstMap = testRunner.BundleIDEntitlementsMap()

	bundleIDs := []string{}
	for bundleID := range bundleIDEntitlemenstMap {
		bundleIDs = append(bundleIDs, bundleID)
	}
	codeSignGroups := export.CreateSelectableCodeSignGroups(installedCertificates, installedProfiles, bundleIDs)

	log.Debugf("Codesign Groups:")
	for _, group := range codeSignGroups {
		log.Debugf(group.String())
	}

	if len(codeSignGroups) == 0 {
		return []export.SelectableCodeSignGroup{}
	}

	codeSignGroups = export.FilterSelectableCodeSignGroups(codeSignGroups,
		export.CreateEntitlementsSelectableCodeSignGroupFilter(bundleIDEntitlemenstMap),
	)

	fmt.Println()
	if !testRunner.IsXcodeManaged() {
		// Handle if UITest target used NON xcode managed profile
		log.Warnf("The UITest target (%s) was signed with NON xcode managed profile,", path.Base(testRunner.Path))
		log.Warnf("only NOT xcode managed profiles are allowed to sign the UITest target.")
		log.Warnf("Removing xcode managed CodeSignInfo groups")

		codeSignGroups = export.FilterSelectableCodeSignGroups(codeSignGroups,
			export.CreateNotXcodeManagedSelectableCodeSignGroupFilter(),
		)
	} else {
		// Handle if UITest target used NON xcode managed profile
		log.Warnf("The UITest target (%s) was signed with xcode managed profile,", path.Base(testRunner.Path))
		log.Warnf("only xcode managed profiles are allowed to sign the UITest target.")
		log.Warnf("Removing NON xcode managed CodeSignInfo groups")

		codeSignGroups = export.FilterSelectableCodeSignGroups(codeSignGroups,
			export.CreateXcodeManagedSelectableCodeSignGroupFilter(),
		)
	}

	log.Debugf("\n")
	log.Debugf("Filtered Codesign Groups:")
	for _, group := range codeSignGroups {
		log.Debugf(group.String())
	}

	return codeSignGroups
}

// hasCertificateForDistType returns true if the provided certificate list has certificate for the selected cert type.
// If isDistCert == true it will search for Distribution Certificates. If it's == false it will search for Developmenttion Certificates.
// If the team doesn't have any certificate for the selected cert type, it will return false.
func hasCertificateForDistType(exportMethod string, certificates []certificateutil.CertificateInfoModel) bool {
	if exportMethod == "development" {
		developmentCertificates := certificateutil.FilterCertificateInfoModelsByFilterFunc(certificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return !codesign.IsDistributionCertificate(certInfo)
		})
		return len(developmentCertificates) > 0
	}

	distributionCertificates := certificateutil.FilterCertificateInfoModelsByFilterFunc(certificates, func(certInfo certificateutil.CertificateInfoModel) bool {
		return codesign.IsDistributionCertificate(certInfo)
	})
	return len(distributionCertificates) > 0

}
