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

// collectIpaExportCertificate returns the certificate to use for the ipa export
func collectIpaExportCertificate(archiveCertificate certificateutil.CertificateInfoModel, installedCertificates []certificateutil.CertificateInfoModel) (certificateutil.CertificateInfoModel, error) {
	fmt.Println()
	fmt.Println()
	question := fmt.Sprintf(`The archive used codesigning files of team: %s - %s
Would you like to use this team to sign your project?`, archiveCertificate.TeamID, archiveCertificate.TeamName)
	useArchiveTeam, err := goinp.AskForBoolWithDefault(question, true)
	if err != nil {
		return certificateutil.CertificateInfoModel{}, fmt.Errorf("failed to read input: %s", err)
	}

	selectedTeam := ""
	certificatesByTeam := mapCertificatesByTeam(installedCertificates)

	if !useArchiveTeam {
		teams := []string{}
		for team := range certificatesByTeam {
			teams = append(teams, team)
		}

		fmt.Println()
		selectedTeam, err = goinp.SelectFromStringsWithDefault("Select the Development team to sign your project", 1, teams)
		if err != nil {
			return certificateutil.CertificateInfoModel{}, fmt.Errorf("failed to read input: %s", err)
		}
	} else {
		selectedTeam = fmt.Sprintf("%s - %s", archiveCertificate.TeamID, archiveCertificate.TeamName)
	}

	selectedCertificate := certificateutil.CertificateInfoModel{}

	if isDistributionCertificate(archiveCertificate) {
		certificates := certificatesByTeam[selectedTeam]
		developmentCertificates := certificateutil.FilterCertificateInfoModelsByFilterFunc(certificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return !isDistributionCertificate(certInfo)
		})

		certificateOptions := []string{}
		for _, certInfo := range developmentCertificates {
			certificateOption := fmt.Sprintf("%s [%s]", certInfo.CommonName, certInfo.Serial)
			certificateOptions = append(certificateOptions, certificateOption)
		}

		fmt.Println()
		question := fmt.Sprintf(`The Xcode archive used distribution certificate: %s [%s].
Please select a development certificate:`, archiveCertificate.CommonName, archiveCertificate.Serial)
		selectedCertificateOption, err := goinp.SelectFromStringsWithDefault(question, 1, certificateOptions)
		if err != nil {
			return certificateutil.CertificateInfoModel{}, fmt.Errorf("failed to read input: %s", err)
		}

		for _, certInfo := range developmentCertificates {
			certificateOption := fmt.Sprintf("%s [%s]", certInfo.CommonName, certInfo.Serial)
			if certificateOption == selectedCertificateOption {
				selectedCertificate = certInfo
			}
		}
	} else {
		certificates := certificatesByTeam[selectedTeam]
		distributionCertificates := certificateutil.FilterCertificateInfoModelsByFilterFunc(certificates, func(certInfo certificateutil.CertificateInfoModel) bool {
			return isDistributionCertificate(certInfo)
		})

		certificateOptions := []string{}
		for _, certInfo := range distributionCertificates {
			certificateOption := fmt.Sprintf("%s [%s]", certInfo.CommonName, certInfo.Serial)
			certificateOptions = append(certificateOptions, certificateOption)
		}

		fmt.Println()
		question := fmt.Sprintf(`The Xcode archive used development certificate: %s [%s].
Please select a distribution certificate:`, archiveCertificate.CommonName, archiveCertificate.Serial)
		selectedCertificateOption, err := goinp.SelectFromStringsWithDefault(question, 1, certificateOptions)
		if err != nil {
			return certificateutil.CertificateInfoModel{}, fmt.Errorf("failed to read input: %s", err)
		}

		for _, certInfo := range distributionCertificates {
			certificateOption := fmt.Sprintf("%s [%s]", certInfo.CommonName, certInfo.Serial)
			if certificateOption == selectedCertificateOption {
				selectedCertificate = certInfo
			}
		}
	}

	return selectedCertificate, nil
}

// collectIpaExportCodeSignGroups returns the codesigngroups required to export an ipa with the selected export methods
func collectIpaExportCodeSignGroups(archive Archive, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) ([]export.CodeSignGroup, error) {
	collectedCodeSignGroups := []export.CodeSignGroup{}
	_, isMacArchive := archive.(xcarchive.MacosArchive)

	codeSignGroups := collectIpaExportSelectableCodeSignGroups(archive, installedCertificates, installedProfiles)
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
			log.Errorf(collectCodesigningFilesInfo)
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
		fmt.Println()
		log.Infof("Codesign settings will be used for %s .ipa/.app export:", exportMethod(collectedCodeSignGroup))
		fmt.Println()
		printCodesignGroup(collectedCodeSignGroup)

		collectedCodeSignGroups = append(collectedCodeSignGroups, collectedCodeSignGroup)

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

	return collectedCodeSignGroups, nil
}

// collectIpaExportSelectableCodeSignGroups returns every possible codesigngroup which can be used to export an ipa file
func collectIpaExportSelectableCodeSignGroups(archive Archive, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) []export.SelectableCodeSignGroup {
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

	certificate, err := findCertificate(archive.SigningIdentity(), installedCertificates)
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
	fmt.Println()
	printCodesignGroup(archiveCodeSignGroup)

	return archiveCodeSignGroup, nil
}
