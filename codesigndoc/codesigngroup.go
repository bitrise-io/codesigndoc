package codesigndoc

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/export"
	"github.com/bitrise-io/go-xcode/exportoptions"
	"github.com/bitrise-io/go-xcode/profileutil"
	"github.com/bitrise-io/go-xcode/xcarchive"
)

// extractCertificatesAndProfiles returns the certificates and provisioning profiles of the given codesign group
func extractCertificatesAndProfiles(codeSignGroups ...export.CodeSignGroup) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel) {
	certificateMap := map[string]certificateutil.CertificateInfoModel{}
	profilesMap := map[string]profileutil.ProvisioningProfileInfoModel{}
	for _, group := range codeSignGroups {
		certificate := group.Certificate()
		certificateMap[certificate.Serial] = certificate

		installerCertificate := group.InstallerCertificate()
		if installerCertificate != nil && installerCertificate.Serial != "" {
			certificateMap[installerCertificate.Serial] = *installerCertificate
		}

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

// exportMethod returns which ipa/pkg/app export type is allowed by the given codesign group
func exportMethod(group export.CodeSignGroup) string {
	for _, profile := range group.BundleIDProfileMap() {
		return string(profile.ExportType)
	}
	return ""
}

// printCodesignGroup prints the given codesign group
func printCodesignGroup(group export.CodeSignGroup) {
	fmt.Printf("%s %s (%s)\n", colorstring.Green("development team:"), group.Certificate().TeamName, group.Certificate().TeamID)
	fmt.Printf("%s %s [%s]\n", colorstring.Green("codesign identity:"), group.Certificate().CommonName, group.Certificate().Serial)

	if group.InstallerCertificate() != nil && group.InstallerCertificate().Serial != "" {
		fmt.Printf("%s %s [%s]\n", colorstring.Green("installer codesign identity:"), group.InstallerCertificate().CommonName, group.InstallerCertificate().Serial)
	}

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

// collectExportCertificate returns the certificate to use for the ipa export
func collectExportCertificate(isMacArchive bool, archiveCertificate certificateutil.CertificateInfoModel, installedCertificates []certificateutil.CertificateInfoModel, installedInstallerCertificates []certificateutil.CertificateInfoModel) ([]certificateutil.CertificateInfoModel, error) {
	var selectedCertificates []certificateutil.CertificateInfoModel

	// Export method
	exportMethods := []string{"development", "app-store"}

	if isMacArchive {
		exportMethods = append(exportMethods, "developer-id")
	} else {
		exportMethods = append(exportMethods, "ad-hoc", "enterprise")
	}

	// Asking the user over and over until we find a valid certificate for the selected export method.
	for searchingValidCertificate := true; searchingValidCertificate; {
		fmt.Println()
		selectedExportMethod, err := goinp.SelectFromStringsWithDefault("Select the ipa export method", 1, exportMethods)
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %s", err)
		}

		log.Debugf("selected export method: %v", selectedExportMethod)

		selectedCertificates, err = filterCertificates(isMacArchive, selectedExportMethod, "", selectedCertificates, archiveCertificate, installedCertificates, installedInstallerCertificates)
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

func filterCertificates(isMacArchive bool, selectedExportMethod, selectedTeam string, selectedCertificates []certificateutil.CertificateInfoModel, archiveCertificate certificateutil.CertificateInfoModel, installedCertificates, installedInstallerCertificates []certificateutil.CertificateInfoModel) ([]certificateutil.CertificateInfoModel, error) {
	var certsForSelectedExport []certificateutil.CertificateInfoModel
	var err error
	log.Debugf("InstalledCerts: %v\n", installedCertificates)

	// Filter the installed certificates by distribution type
	switch selectedExportMethod {
	case "development":
		certsForSelectedExport = certificateutil.FilterCertificateInfoModelsByFilterFunc(installedCertificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return !codesign.IsDistributionCertificate(certInfo)
		})

		log.Debugf("DeveloperDistribution certificates: %v\n", certsForSelectedExport)
		break
	case "installer":
		certsForSelectedExport = certificateutil.FilterCertificateInfoModelsByFilterFunc(installedInstallerCertificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return codesign.IsInstallerCertificate(certInfo)
		})

		log.Debugf("Installer certificates: %v\n", certsForSelectedExport)
		break
	default:
		certsForSelectedExport = certificateutil.FilterCertificateInfoModelsByFilterFunc(installedCertificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return codesign.IsDistributionCertificate(certInfo)
		})
		log.Debugf("Distribution certificates: %v\n", certsForSelectedExport)
		break
	}

	filteredCertificatesByTeam := codesign.MapCertificatesByTeam(certsForSelectedExport)
	log.Debugf("Filtered certificates (by distribution type) by team: %v\n", filteredCertificatesByTeam)

	if len(filteredCertificatesByTeam) == 0 {
		log.Warnf("🚨  We couldn't find any certificate for the %s export method.", selectedExportMethod)
		return nil, nil
	}

	// If we already selected a team, we can skip it. (e.g mac app-store export)
	if selectedTeam == "" {
		useArchiveTeam := true
		_, contains := filteredCertificatesByTeam[fmt.Sprintf("%s - %s", archiveCertificate.TeamID, archiveCertificate.TeamName)]

		// Ask the question if there is multiple valid team and the archiving team is one of them.
		// Skip it if only 1 team has certificates on the machine. Or the archiving team does'n have the desired certificate type.
		// Skip the question + set the useArchiveTeam = false, if multiple team has certificate for the export method but the archiving team is not one of them.
		if len(filteredCertificatesByTeam) > 1 && contains {
			fmt.Println()

			question := fmt.Sprintf(`The archive used codesigning files of team: %s - %s
Would you like to use this team to export an ipa file?`, archiveCertificate.TeamID, archiveCertificate.TeamName)
			useArchiveTeam, err = goinp.AskForBoolWithDefault(question, true)
			if err != nil {
				return selectedCertificates, fmt.Errorf("failed to read input: %s", err)
			}
			// If multiple team has certificate for the export method but the archiving team is not one of them.
		} else if !contains {
			archiveTeam := fmt.Sprintf("%s - %s", archiveCertificate.TeamName, archiveCertificate.TeamID)

			fmt.Println()
			log.Warnf("🚨   The archiving team (%s) doesn't have certificate for the %s export method", archiveTeam, selectedExportMethod)
			useArchiveTeam = false
		} else {
			archiveTeam := fmt.Sprintf("%s - %s", archiveCertificate.TeamName, archiveCertificate.TeamID)

			fmt.Println()
			log.Printf("Only the archiving team (%s) has certificate for the %s export method", archiveTeam, selectedExportMethod)
		}

		// Use different team for export than archive.
		if !useArchiveTeam {
			teams := []string{}
			for team := range filteredCertificatesByTeam {
				if hasCertificateForDistType(selectedExportMethod, filteredCertificatesByTeam[team]) {
					teams = append(teams, team)
				}
			}

			fmt.Println()
			selectedTeam, err = goinp.SelectFromStringsWithDefault("Select the Development team to export your app", 1, teams)
			if err != nil {
				return selectedCertificates, fmt.Errorf("failed to read input: %s", err)
			}
		} else {
			selectedTeam = fmt.Sprintf("%s - %s", archiveCertificate.TeamID, archiveCertificate.TeamName)
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
	switch selectedExportMethod {
	case "development":
		certType = "development"
		break
	case "installer":
		certType = "installer"
		break
	default:
		break
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

	// Collect installer cert for MacOS app-store export.
	if selectedExportMethod == "app-store" && isMacArchive {
		fmt.Println()
		question := `Do you want to collect installer certificate for the app-store export? [yes,no]`
		collectInstallerCert, err := goinp.AskForBoolWithDefault(question, true)
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %s", err)
		}

		filteredInstallerCertificatesByTeam := codesign.MapCertificatesByTeam(installedInstallerCertificates)
		if !hasCertificateForDistType("installer", filteredInstallerCertificatesByTeam[selectedTeam]) {
			log.Warnf("🚨   The selected team (%s) doesn't have installer certificate for MacOS app-store export", selectedTeam)
			return selectedCertificates, nil
		}

		if collectInstallerCert {
			selectedCertificates, err = filterCertificates(isMacArchive, "installer", selectedTeam, selectedCertificates, archiveCertificate, installedCertificates, installedInstallerCertificates)
			if err != nil {
				return nil, err
			}
		}
	}

	return selectedCertificates, nil
}

// collectExportCodeSignGroups returns the codesigngroups required to export an ipa/.app with the selected export methods
func collectExportCodeSignGroups(archive Archive, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) ([]export.CodeSignGroup, error) {
	collectedCodeSignGroups := []export.CodeSignGroup{}
	_, isMacArchive := archive.(xcarchive.MacosArchive)

	codeSignGroups := collectExportSelectableCodeSignGroups(archive, installedCertificates, installedProfiles)
	if len(codeSignGroups) == 0 {
		return nil, errors.New("no code sign files (Codesign Identities and Provisioning Profiles) are installed to export an ipa\n" + collectCodesigningFilesInfo)
	}

	exportMethods := []string{"development", "app-store"}

	if isMacArchive {
		exportMethods = append(exportMethods, "developer-id")
	} else {
		exportMethods = append(exportMethods, "ad-hoc", "enterprise")
	}

	for true {
		selectedExportMethod, err := goinp.SelectFromStringsWithDefault("Select the ipa export method", 1, exportMethods)
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %s", err)
		}
		log.Debugf("selected export method: %v", selectedExportMethod)

		fmt.Println()
		filteredCodeSignGroups := export.FilterSelectableCodeSignGroups(codeSignGroups,
			export.CreateExportMethodSelectableCodeSignGroupFilter(exportoptions.Method(selectedExportMethod)),
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
			question := "Do you want to collect another ipa export code sign files"
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

			fmt.Printf("Codesign Indentity for %s ipa export: %s\n", selectedExportMethod, selectedCertificateOption)
		} else {
			sort.Strings(certificateOptions)

			question := fmt.Sprintf("Select the Codesign Indentity for %s ipa export", selectedExportMethod)
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
			return nil, fmt.Errorf("failed to find Provisioning Profiles for ipa export")
		}

		var collectedCodeSignGroup export.CodeSignGroup
		if isMacArchive {
			installedInstallerCertificates := []certificateutil.CertificateInfoModel{}

			var selectedInstallerCetrificate certificateutil.CertificateInfoModel
			if selectedExportMethod == string(exportoptions.MethodAppStore) {
				installedInstallerCertificates, err = certificateutil.InstalledInstallerCertificateInfos()
				if err != nil {
					log.Errorf("Failed to read installed Installer certificates, error: %s", err)
				}

				installedInstallerCertificates = certificateutil.FilterValidCertificateInfos(installedInstallerCertificates)

				log.Debugf("\n")
				log.Debugf("Installed installer certificates:")
				for _, certInfo := range installedInstallerCertificates {
					log.Debugf(certInfo.String())
				}

				for _, installerCetrificate := range installedInstallerCertificates {
					if installerCetrificate.TeamID == selectedCertificate.TeamID {
						selectedInstallerCetrificate = installerCetrificate
						break
					}
				}
			}

			collectedCodeSignGroup = export.NewMacGroup(*selectedCertificate, &selectedInstallerCetrificate, selectedBundleIDProfileMap)
		} else {
			collectedCodeSignGroup = export.NewIOSGroup(*selectedCertificate, selectedBundleIDProfileMap)
		}

		fmt.Println()
		log.Infof("Codesign settings will be used for %s .ipa/.app export:", exportMethod(collectedCodeSignGroup))
		printCodesignGroup(collectedCodeSignGroup)

		collectedCodeSignGroups = append(collectedCodeSignGroups, collectedCodeSignGroup)

		fmt.Println()
		question := "Do you want to collect another ipa export code sign files"
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

// collectExportSelectableCodeSignGroups returns every possible codesigngroup which can be used to export an ipa file
func collectExportSelectableCodeSignGroups(archive Archive, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) []export.SelectableCodeSignGroup {
	bundleIDEntitlemenstMap := archive.BundleIDEntitlementsMap()

	fmt.Println()
	log.Infof("Targets to sign:")
	for bundleID, entitlements := range bundleIDEntitlemenstMap {
		fmt.Printf("- %s with %d capabilities\n", bundleID, len(entitlements))
	}
	fmt.Println()

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

	// Handle if archive used NON xcode managed profile
	if len(codeSignGroups) > 0 && !archive.IsXcodeManaged() {
		log.Warnf("App was signed with NON xcode managed profile when archiving,")
		log.Warnf("only NOT xcode managed profiles are allowed to sign when exporting the archive.")
		log.Warnf("Removing xcode managed CodeSignInfo groups")

		codeSignGroups = export.FilterSelectableCodeSignGroups(codeSignGroups,
			export.CreateNotXcodeManagedSelectableCodeSignGroupFilter(),
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
	switch exportMethod {
	case "development":
		developmentCertificates := certificateutil.FilterCertificateInfoModelsByFilterFunc(certificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return !codesign.IsDistributionCertificate(certInfo)
		})
		return len(developmentCertificates) > 0
	case "installer":
		installerCertificates := certificateutil.FilterCertificateInfoModelsByFilterFunc(certificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return codesign.IsInstallerCertificate(certInfo)
		})
		return len(installerCertificates) > 0
	default:
		distributionCertificates := certificateutil.FilterCertificateInfoModelsByFilterFunc(certificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return codesign.IsDistributionCertificate(certInfo)
		})
		return len(distributionCertificates) > 0
	}
}

func getIOSCodeSignGroup(archivePath string, installedCertificates []certificateutil.CertificateInfoModel) (xcarchive.IosArchive, *export.IosCodeSignGroup, error) {
	archive, err := xcarchive.NewIosArchive(archivePath)
	if err != nil {
		return xcarchive.IosArchive{}, &export.IosCodeSignGroup{}, fmt.Errorf("failed to analyze archive, error: %s", err)
	}

	codeSignGroup, err := getCodeSignGroup(archive, installedCertificates, false)
	if err != nil {
		return xcarchive.IosArchive{}, &export.IosCodeSignGroup{}, fmt.Errorf("failed to analyze archive, error: %s", err)
	}

	archiveCodeSignGroup, ok := codeSignGroup.(*export.IosCodeSignGroup)
	if !ok {
		return xcarchive.IosArchive{}, &export.IosCodeSignGroup{}, fmt.Errorf("failed to analyze archive, error: %s", err)
	}

	return archive, archiveCodeSignGroup, nil
}

func getMacOSCodeSignGroup(archivePath string, installedCertificates []certificateutil.CertificateInfoModel) (xcarchive.MacosArchive, *export.MacCodeSignGroup, error) {
	archive, err := xcarchive.NewMacosArchive(archivePath)
	if err != nil {
		return xcarchive.MacosArchive{}, &export.MacCodeSignGroup{}, fmt.Errorf("failed to analyze archive, error: %s", err)
	}

	codeSignGroup, err := getCodeSignGroup(archive, installedCertificates, true)
	if err != nil {
		return xcarchive.MacosArchive{}, &export.MacCodeSignGroup{}, fmt.Errorf("failed to analyze archive, error: %s", err)
	}

	archiveCodeSignGroup, ok := codeSignGroup.(*export.MacCodeSignGroup)
	if !ok {
		return xcarchive.MacosArchive{}, &export.MacCodeSignGroup{}, fmt.Errorf("failed to analyze archive, error: %s", err)
	}

	return archive, archiveCodeSignGroup, nil
}

func getCodeSignGroup(archive Archive, installedCertificates []certificateutil.CertificateInfoModel, isMacArchive bool) (export.CodeSignGroup, error) {
	if archive.SigningIdentity() == "" {
		return nil, fmt.Errorf("no signing identity found")
	}

	certificate, err := codesign.FindCertificate(archive.SigningIdentity(), installedCertificates)
	if err != nil {
		return nil, err
	}

	var archiveCodeSignGroup export.CodeSignGroup
	if isMacArchive {
		archiveCodeSignGroup = export.NewMacGroup(certificate, nil, archive.BundleIDProfileInfoMap())
		if err != nil {
			return &export.MacCodeSignGroup{}, fmt.Errorf("failed to analyze the archive, error: %s", err)
		}

	} else {
		archiveCodeSignGroup = export.NewIOSGroup(certificate, archive.BundleIDProfileInfoMap())
		if err != nil {
			return &export.IosCodeSignGroup{}, fmt.Errorf("failed to analyze the archive, error: %s", err)
		}
	}

	fmt.Println()
	log.Infof("Codesign settings used for archive:")
	printCodesignGroup(archiveCodeSignGroup)

	return archiveCodeSignGroup, nil
}
