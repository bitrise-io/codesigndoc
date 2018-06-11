package codesigndoc

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/export"
	"github.com/bitrise-tools/go-xcode/exportoptions"
	"github.com/bitrise-tools/go-xcode/profileutil"
	"github.com/bitrise-tools/go-xcode/xcarchive"
)

// extractCertificatesAndProfiles returns the certificates and provisioning profiles of the given codesign group
func extractCertificatesAndProfiles(codeSignGroups ...export.IosCodeSignGroup) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel) {
	certificateMap := map[string]certificateutil.CertificateInfoModel{}
	profilesMap := map[string]profileutil.ProvisioningProfileInfoModel{}
	for _, group := range codeSignGroups {
		certificate := group.Certificate

		certificateMap[certificate.Serial] = certificate

		for _, profile := range group.BundleIDProfileMap {
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
func exportMethod(group export.IosCodeSignGroup) string {
	for _, profile := range group.BundleIDProfileMap {
		return string(profile.ExportType)
	}
	return ""
}

// printCodesignGroup prints the given codesign group
func printCodesignGroup(group export.IosCodeSignGroup) {
	fmt.Printf("%s %s (%s)\n", colorstring.Green("development team:"), group.Certificate.TeamName, group.Certificate.TeamID)
	fmt.Printf("%s %s [%s]\n", colorstring.Green("codesign identity:"), group.Certificate.CommonName, group.Certificate.Serial)
	idx := -1
	for bundleID, profile := range group.BundleIDProfileMap {
		idx++
		if idx == 0 {
			fmt.Printf("%s %s -> %s\n", colorstring.Greenf("provisioning profiles:"), profile.Name, bundleID)
		} else {
			fmt.Printf("%s%s -> %s\n", strings.Repeat(" ", len("provisioning profiles: ")), profile.Name, bundleID)
		}
	}
}

// collectIpaExportCertificate returns the certificate to use for the ipa export
func collectIpaExportCertificate(archiveCertificate certificateutil.CertificateInfoModel, installedCertificates []certificateutil.CertificateInfoModel) (certificateutil.CertificateInfoModel, error) {
	fmt.Println()
	fmt.Println()

	selectedCertificate := certificateutil.CertificateInfoModel{}

	// Asking the user over and over until we find a valid certificate for the selected export method.
	for searchingValidCertificate := true; searchingValidCertificate; {

		// Export method
		exportMethods := []string{"development", "app-store", "ad-hoc", "enterprise"}

		selectedExportMethod, err := goinp.SelectFromStringsWithDefault("Select the ipa export method", 1, exportMethods)
		if err != nil {
			return certificateutil.CertificateInfoModel{}, fmt.Errorf("failed to read input: %s", err)
		}
		log.Debugf("selected export method: %v", selectedExportMethod)

		// Export method needs Development or Distribution certificate?
		exportDistCert := true
		if selectedExportMethod == "development" {
			exportDistCert = false
		}

		selectedTeam := ""
		log.Debugf("InstalledCerts: %v\n", installedCertificates)

		// Filter the installed certificates by distribution type
		var certsForSelectedExport []certificateutil.CertificateInfoModel
		if exportDistCert {
			certsForSelectedExport = certificateutil.FilterCertificateInfoModelsByFilterFunc(installedCertificates, func(certInfo certificateutil.CertificateInfoModel) bool {
				return isDistributionCertificate(certInfo)
			})
			log.Debugf("Distribution certificated: %v\n", certsForSelectedExport)
		} else {
			certsForSelectedExport = certificateutil.FilterCertificateInfoModelsByFilterFunc(installedCertificates, func(certInfo certificateutil.CertificateInfoModel) bool {
				return !isDistributionCertificate(certInfo)
			})
			log.Debugf("Developer certificated: %v\n", certsForSelectedExport)
		}

		filteredCertificatesByTeam := mapCertificatesByTeam(certsForSelectedExport)
		log.Debugf("Filtered certificates (by distribution type) by team: %v\n", filteredCertificatesByTeam)

		if len(filteredCertificatesByTeam) == 0 {
			log.Warnf("🚨  We couldn't find any certificate for the %s export method.", selectedExportMethod)
			continue
		}

		useArchiveTeam := true
		_, contains := filteredCertificatesByTeam[fmt.Sprintf("%s - %s", archiveCertificate.TeamID, archiveCertificate.TeamName)]

		// Ask the question if there is multiple valid team and the archiving team is one of them.
		// Skip it if only 1 team has certificates on the machine. Or the archiving team does'n have the desired certificate type.
		// Skip the question + set the useArchiveTeam = false, if multiple team has certificate for the export method but the archiving team is not one of them.
		if len(filteredCertificatesByTeam) > 1 && contains {
			question := fmt.Sprintf(`The archive used codesigning files of team: %s - %s
Would you like to use this team to export an ipa file?`, archiveCertificate.TeamID, archiveCertificate.TeamName)
			useArchiveTeam, err = goinp.AskForBoolWithDefault(question, true)
			if err != nil {
				return certificateutil.CertificateInfoModel{}, fmt.Errorf("failed to read input: %s", err)
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
				if hasCertificateForDistType(exportDistCert, filteredCertificatesByTeam[team]) {
					teams = append(teams, team)
				}
			}

			fmt.Println()
			selectedTeam, err = goinp.SelectFromStringsWithDefault("Select the Development team to export your app", 1, teams)
			if err != nil {
				return certificateutil.CertificateInfoModel{}, fmt.Errorf("failed to read input: %s", err)
			}
		} else {
			selectedTeam = fmt.Sprintf("%s - %s", archiveCertificate.TeamID, archiveCertificate.TeamName)
		}

		// Find the specific development certificate.
		filteredTeamCertificates := filteredCertificatesByTeam[selectedTeam]
		certificateOptions := []string{}

		for _, certInfo := range filteredTeamCertificates {
			certificateOption := fmt.Sprintf("%s [%s]", certInfo.CommonName, certInfo.Serial)
			certificateOptions = append(certificateOptions, certificateOption)
		}

		certType := "distribution"
		if !exportDistCert {
			certType = "development"
		}

		fmt.Println()
		question := fmt.Sprintf("Please select a %s certificate:", certType)
		selectedCertificateOption, err := goinp.SelectFromStringsWithDefault(question, 1, certificateOptions)
		if err != nil {
			return certificateutil.CertificateInfoModel{}, fmt.Errorf("failed to read input: %s", err)
		}

		for _, certInfo := range filteredTeamCertificates {
			certificateOption := fmt.Sprintf("%s [%s]", certInfo.CommonName, certInfo.Serial)
			if certificateOption == selectedCertificateOption {
				selectedCertificate = certInfo
				searchingValidCertificate = false
			}
		}
	}

	return selectedCertificate, nil
}

// collectIpaExportCodeSignGroups returns the codesigngroups required to export an ipa with the selected export methods
func collectIpaExportCodeSignGroups(archive xcarchive.IosArchive, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) ([]export.IosCodeSignGroup, error) {
	iosCodeSignGroups := []export.IosCodeSignGroup{}

	codeSignGroups := collectIpaExportSelectableCodeSignGroups(archive, installedCertificates, installedProfiles)
	if len(codeSignGroups) == 0 {
		return nil, errors.New("no code sign files (Codesign Identities and Provisioning Profiles) are installed to export an ipa\n" + collectCodesigningFilesInfo)
	}

	exportMethods := []string{"development", "app-store", "ad-hoc", "enterprise"}

	for true {
		fmt.Println()
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
			log.Errorf("🚨  Could not find the codesigning files for %s ipa export:", selectedExportMethod)
			log.Warnf(collectCodesigningFilesInfo)
			fmt.Println()
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

			fmt.Println()
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
			profiles = filterLatestProfiles(profiles)
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

		iosCodeSignGroup := export.IosCodeSignGroup{
			Certificate:        *selectedCertificate,
			BundleIDProfileMap: selectedBundleIDProfileMap,
		}

		fmt.Println()
		fmt.Println()
		log.Infof("Codesign settings will be used for %s ipa export:", exportMethod(iosCodeSignGroup))
		fmt.Println()
		printCodesignGroup(iosCodeSignGroup)

		iosCodeSignGroups = append(iosCodeSignGroups, iosCodeSignGroup)

		fmt.Println()
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

	return iosCodeSignGroups, nil
}

// collectIpaExportSelectableCodeSignGroups returns every possible codesigngroup which can be used to export an ipa file
func collectIpaExportSelectableCodeSignGroups(archive xcarchive.IosArchive, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) []export.SelectableCodeSignGroup {
	bundleIDEntitlemenstMap := archive.BundleIDEntitlementsMap()

	fmt.Println()
	fmt.Println()
	log.Infof("Targets to sign:")
	fmt.Println()
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
func hasCertificateForDistType(isDistCert bool, certificates []certificateutil.CertificateInfoModel) bool {
	if !isDistCert {
		developmentCertificates := certificateutil.FilterCertificateInfoModelsByFilterFunc(certificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return !isDistributionCertificate(certInfo)
		})
		return len(developmentCertificates) > 0
	}

	distributionCertificates := certificateutil.FilterCertificateInfoModelsByFilterFunc(certificates, func(certInfo certificateutil.CertificateInfoModel) bool {
		return isDistributionCertificate(certInfo)
	})
	return len(distributionCertificates) > 0
}
