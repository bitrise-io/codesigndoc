package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/osxkeychain"
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/export"
	"github.com/bitrise-tools/go-xcode/exportoptions"
	"github.com/bitrise-tools/go-xcode/profileutil"
	"github.com/bitrise-tools/go-xcode/xcarchive"
	"github.com/pkg/errors"
)

func initExportOutputDir() (string, error) {
	confExportOutputDirPath := "./codesigndoc_exports"
	absExportOutputDirPath, err := pathutil.AbsPath(confExportOutputDirPath)
	log.Debugf("absExportOutputDirPath: %s", absExportOutputDirPath)
	if err != nil {
		return absExportOutputDirPath, fmt.Errorf("Failed to determin Absolute path of export dir: %s", confExportOutputDirPath)
	}
	if exist, err := pathutil.IsDirExists(absExportOutputDirPath); err != nil {
		return absExportOutputDirPath, fmt.Errorf("Failed to determin whether the export directory already exists: %s", err)
	} else if !exist {
		if err := os.Mkdir(absExportOutputDirPath, 0777); err != nil {
			return absExportOutputDirPath, fmt.Errorf("Failed to create export output directory at path: %s | error: %s", absExportOutputDirPath, err)
		}
	} else {
		log.Warnf("Export output dir already exists at path: %s", absExportOutputDirPath)
	}
	return absExportOutputDirPath, nil
}

func analyzeArchive(archive xcarchive.IosArchive, installedCertificates []certificateutil.CertificateInfoModel) (export.IosCodeSignGroup, error) {
	signingIdentity := archive.SigningIdentity()
	bundleIDProfileInfoMap := archive.BundleIDProfileInfoMap()

	if signingIdentity == "" {
		return export.IosCodeSignGroup{}, fmt.Errorf("no signing identity found")
	}

	certificate, err := findCertificate(signingIdentity, installedCertificates)
	if err != nil {
		return export.IosCodeSignGroup{}, err
	}

	return export.IosCodeSignGroup{
		Certificate:        certificate,
		BundleIDProfileMap: bundleIDProfileInfoMap,
	}, nil
}

func collectIpaExportSelectableCodeSignGroups(archive xcarchive.IosArchive, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) ([]export.SelectableCodeSignGroup, error) {
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
		log.Errorf("Failed to create codesigning groups for the project")
		return []export.SelectableCodeSignGroup{}, nil
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

	if len(codeSignGroups) == 0 {
		log.Errorf("Failed to create codesigning groups for the project")
	}

	return codeSignGroups, nil
}

func collectIpaExportCodeSignGroups(archive xcarchive.IosArchive, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) ([]export.IosCodeSignGroup, error) {
	iosCodeSignGroups := []export.IosCodeSignGroup{}

	codeSignGroups, err := collectIpaExportSelectableCodeSignGroups(archive, installedCertificates, installedProfiles)
	if err != nil {
		return nil, printXcodeScanFinishedWithError("Failed to collect valid code sign settings: %s", err)
	}

	if len(codeSignGroups) == 0 {
		fmt.Println()
		log.Errorf("No code sign files (Codesign Identities and Provisioning Profiles) are installed to export an ipa")
		log.Errorf("To collect available code sign files, we search for intsalled Provisioning Profiles:")
		log.Errorf("- which has installed Codesign Identity in your Keychain")
		log.Errorf("- which can provision your application target's bundle ids")
		log.Errorf("- which has the project defined Capabilities set")
		return nil, printXcodeScanFinishedWithError("Failed to find code sign files")
	}

	exportMethods := []string{"development", "app-store", "ad-hoc", "enterprise"}

	for true {
		fmt.Println()
		selectedExportMethod, err := goinp.SelectFromStringsWithDefault("Select the ipa export method", 1, exportMethods)
		if err != nil {
			return nil, printXcodeScanFinishedWithError("Failed to select ipa export method: %s", err)
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
			log.Errorf("No code sign files (Codesign Identities and Provisioning Profiles) are installed for %s ipa export", selectedExportMethod)
			log.Errorf("To collect available code sign files, we search for intsalled Provisioning Profiles:")
			log.Errorf("- which has installed Codesign Identity in your Keychain")
			log.Errorf("- which can provision your application target's bundle ids")
			log.Errorf("- which has the project defined Capabilities set")
			log.Errorf("- which matches to the selected ipa export method")

			fmt.Println()
			fmt.Println()
			question := "Do you want to collect another ipa export code sign files"
			question += "\n(select NO to finish collecting codesign files and continue)"
			anotherExport, err := goinp.AskForBoolWithDefault(question, false)
			if err != nil {
				return nil, fmt.Errorf("failed to ask: %s", err)
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
				return nil, printXcodeScanFinishedWithError("Failed to select Codesign Indentity: %s", err)
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
			return nil, printXcodeScanFinishedWithError("Failed to find selected Codesign Indentity")
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
			return nil, printXcodeScanFinishedWithError("Failed to find Provisioning Profiles for Code Sign Identity")
		}

		selectedBundleIDProfileMap := map[string]profileutil.ProvisioningProfileInfoModel{}
		for bundleID, profiles := range bundleIDProfilesMap {
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
					return nil, printXcodeScanFinishedWithError("Failed to select Provisioning Profile: %s", err)
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
			return nil, printXcodeScanFinishedWithError("Failed to find Provisioning Profiles for ipa export")
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
			return nil, fmt.Errorf("failed to ask: %s", err)
		}
		if !anotherExport {
			break
		}
	}

	return iosCodeSignGroups, nil
}

func collectIpaExportCertificate(archiveCertificate certificateutil.CertificateInfoModel, installedCertificates []certificateutil.CertificateInfoModel) (certificateutil.CertificateInfoModel, error) {
	fmt.Println()
	fmt.Println()
	question := fmt.Sprintf(`The Xcode archive used codesigning files of team: %s - %s
Would you like to use this team to sign your project?`, archiveCertificate.TeamID, archiveCertificate.TeamName)
	useArchiveTeam, err := goinp.AskForBoolWithDefault(question, true)
	if err != nil {
		return certificateutil.CertificateInfoModel{}, fmt.Errorf("failed to ask: %s", err)
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
			return certificateutil.CertificateInfoModel{}, printXcodeScanFinishedWithError("Failed to select Codesign Indentity: %s", err)
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
			return certificateutil.CertificateInfoModel{}, printXcodeScanFinishedWithError("Failed to select Codesign Indentity: %s", err)
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
			return certificateutil.CertificateInfoModel{}, printXcodeScanFinishedWithError("Failed to select Codesign Indentity: %s", err)
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

func collectAndExportProvisioningProfiles(profiles []profileutil.ProvisioningProfileInfoModel, absExportOutputDirPath string) error {
	if len(profiles) == 0 {
		return nil
	}

	fmt.Println()
	log.Infof("Required Provisioning Profiles (%d)", len(profiles))
	fmt.Println()
	for _, profile := range profiles {
		log.Printf("- %s (UUID: %s)", profile.Name, profile.UUID)
	}

	profilePathInfoMap := map[string]profileutil.ProvisioningProfileInfoModel{}

	fmt.Println()
	log.Infof("Exporting Provisioning Profiles...")
	fmt.Println()

	for _, profile := range profiles {
		log.Printf("searching for required Provisioning Profile: %s (UUID: %s)", profile.Name, profile.UUID)
		_, pth, err := profileutil.FindProvisioningProfileInfo(profile.UUID)
		if err != nil {
			return errors.Wrap(err, "Failed to find Provisioning Profile")
		}
		profilePathInfoMap[pth] = profile
		log.Printf("file found at: %s", pth)
	}

	if err := exportProvisioningProfiles(profilePathInfoMap, absExportOutputDirPath); err != nil {
		return fmt.Errorf("Failed to export the Provisioning Profile into the export directory: %s", err)
	}

	return nil
}

func collectAndExportIdentities(certificates []certificateutil.CertificateInfoModel, absExportOutputDirPath string) error {
	if len(certificates) == 0 {
		return nil
	}

	fmt.Println()
	fmt.Println()
	log.Infof("Required Identities/Certificates (%d)", len(certificates))
	fmt.Println()
	for _, certificate := range certificates {
		log.Printf("- %s", certificate.CommonName)
	}

	fmt.Println()
	log.Infof("Exporting the Identities (Certificates):")
	fmt.Println()

	identitiesWithKeychainRefs := []osxkeychain.IdentityWithRefModel{}
	defer osxkeychain.ReleaseIdentityWithRefList(identitiesWithKeychainRefs)

	for _, certificate := range certificates {
		log.Printf("searching for Identity: %s", certificate.CommonName)
		identityRef, err := osxkeychain.FindAndValidateIdentity(certificate.CommonName)
		if err != nil {
			return fmt.Errorf("Failed to export, error: %s", err)
		}

		if identityRef == nil {
			return errors.New("Identity not found in the keychain, or it was invalid (expired)")
		}

		identitiesWithKeychainRefs = append(identitiesWithKeychainRefs, *identityRef)
	}

	identityKechainRefs := osxkeychain.CreateEmptyCFTypeRefSlice()
	for _, aIdentityWithRefItm := range identitiesWithKeychainRefs {
		fmt.Println("exporting Identity:", aIdentityWithRefItm.Label)
		identityKechainRefs = append(identityKechainRefs, aIdentityWithRefItm.KeychainRef)
	}

	fmt.Println()
	if isAskForPassword {
		log.Infof("Exporting from Keychain")
		log.Warnf(" You'll be asked to provide a Passphrase for the .p12 file!")
	} else {
		log.Warnf("Exporting from Keychain using empty Passphrase...")
		log.Printf("This means that if you want to import the file the passphrase at import should be left empty,")
		log.Printf("you don't have to type in anything, just leave the passphrase input empty.")
	}
	fmt.Println()
	log.Warnf("You'll most likely see popups one for each Identity from Keychain,")
	log.Warnf("you will have to accept (Allow) those to be able to export the Identities!")
	fmt.Println()

	if err := osxkeychain.ExportFromKeychain(identityKechainRefs, filepath.Join(absExportOutputDirPath, "Identities.p12"), isAskForPassword); err != nil {
		return fmt.Errorf("Failed to export from Keychain: %s", err)
	}

	return nil
}

func exportProvisioningProfiles(profilePathInfoMap map[string]profileutil.ProvisioningProfileInfoModel, exportTargetDirPath string) error {
	idx := -1
	for path, profileInfo := range profilePathInfoMap {
		idx++

		if idx != 0 {
			fmt.Println()
		}

		log.Printf("exporting Provisioning Profile: %s (%s)", profileInfo.Name, profileInfo.UUID)

		exportFileName := provProfileExportFileName(profileInfo, path)
		exportPth := filepath.Join(exportTargetDirPath, exportFileName)
		if err := command.RunCommand("cp", path, exportPth); err != nil {
			return fmt.Errorf("Failed to copy Provisioning Profile (from: %s) (to: %s), error: %s", path, exportPth, err)
		}
	}
	return nil
}

func provProfileExportFileName(info profileutil.ProvisioningProfileInfoModel, path string) string {
	replaceRexp, err := regexp.Compile("[^A-Za-z0-9_.-]")
	if err != nil {
		log.Warnf("Invalid regex, error: %s", err)
		return ""
	}
	safeTitle := replaceRexp.ReplaceAllString(info.Name, "")
	extension := ".mobileprovision"
	if strings.HasSuffix(path, ".provisionprofile") {
		extension = ".provisionprofile"
	}

	return info.UUID + "." + safeTitle + extension
}

func exportCodesignFiles(toolName, archivePath, utputDirPath string) error {
	// archive code sign settings
	installedCertificates, err := certificateutil.InstalledCodesigningCertificateInfos()
	if err != nil {
		return printFinishedWithError(toolName, "Failed to list installed code signing identities, error: %s", err)
	}
	installedCertificates = certificateutil.FilterValidCertificateInfos(installedCertificates)

	log.Debugf("Installed certificates:")
	for _, installedCertificate := range installedCertificates {
		log.Debugf(installedCertificate.String())
	}

	installedProfiles, err := profileutil.InstalledProvisioningProfileInfos(profileutil.ProfileTypeIos)
	if err != nil {
		return err
	}

	log.Debugf("Installed profiles:")
	for _, profileInfo := range installedProfiles {
		log.Debugf(profileInfo.String(installedCertificates...))
	}

	archive, err := xcarchive.NewIosArchive(archivePath)
	if err != nil {
		return printFinishedWithError(toolName, "Failed to analyze archive, error: %s", err)
	}

	archiveCodeSignGroup, err := analyzeArchive(archive, installedCertificates)
	if err != nil {
		return printFinishedWithError(toolName, "Failed to analyze the archive, error: %s", err)
	}

	fmt.Println()
	log.Infof("Codesign settings used for Xamarin archive:")
	fmt.Println()
	printCodesignGroup(archiveCodeSignGroup)

	// ipa export code sign settings
	fmt.Println()
	fmt.Println()
	log.Printf("🔦  Analyzing the archive, to get ipa export code signing settings...")

	certificatesToExport := []certificateutil.CertificateInfoModel{}
	profilesToExport := []profileutil.ProvisioningProfileInfoModel{}

	if certificatesOnly {
		ipaExportCertificate, err := collectIpaExportCertificate(archiveCodeSignGroup.Certificate, installedCertificates)
		if err != nil {
			return err
		}

		certificatesToExport = append(certificatesToExport, archiveCodeSignGroup.Certificate, ipaExportCertificate)
	} else {
		ipaExportCodeSignGroups, err := collectIpaExportCodeSignGroups(archive, installedCertificates, installedProfiles)
		if err != nil {
			return printFinishedWithError(toolName, "Failed to collect ipa export code sign groups, error: %s", err)
		}

		codeSignGroups := append(ipaExportCodeSignGroups, archiveCodeSignGroup)
		certificates, profiles := extractCertificatesAndProfiles(codeSignGroups...)

		certificatesToExport = append(certificatesToExport, certificates...)
		profilesToExport = append(profilesToExport, profiles...)
	}

	if err := collectAndExportIdentities(certificatesToExport, utputDirPath); err != nil {
		return printFinishedWithError(toolName, "Failed to export codesign identities, error: %s", err)
	}

	if err := collectAndExportProvisioningProfiles(profilesToExport, utputDirPath); err != nil {
		return printFinishedWithError(toolName, "Failed to export provisioning profiles, error: %s", err)
	}

	printFinished()

	return nil
}
