package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/bitriseclient"
	"github.com/bitrise-tools/codesigndoc/osxkeychain"
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/export"
	"github.com/bitrise-tools/go-xcode/exportoptions"
	"github.com/bitrise-tools/go-xcode/profileutil"
	"github.com/bitrise-tools/go-xcode/xcarchive"
	"github.com/pkg/errors"
)

const collectCodesigningFilesInfo = `To collect available code sign files, we search for installed Provisioning Profiles:"
- which has installed Codesign Identity in your Keychain"
- which can provision your application target's bundle ids"
- which has the project defined Capabilities set"
- which matches to the selected ipa export method"
`

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

func filterLatestProfiles(profiles []profileutil.ProvisioningProfileInfoModel) []profileutil.ProvisioningProfileInfoModel {
	profilesByBundleIDAndName := map[string][]profileutil.ProvisioningProfileInfoModel{}
	for _, profile := range profiles {
		bundleID := profile.BundleID
		name := profile.Name
		bundleIDAndName := bundleID + name
		profs, ok := profilesByBundleIDAndName[bundleIDAndName]
		if !ok {
			profs = []profileutil.ProvisioningProfileInfoModel{}
		}
		profs = append(profs, profile)
		profilesByBundleIDAndName[bundleIDAndName] = profs
	}

	filteredProfiles := []profileutil.ProvisioningProfileInfoModel{}
	for _, profiles := range profilesByBundleIDAndName {
		var latestProfile *profileutil.ProvisioningProfileInfoModel
		for _, profile := range profiles {
			if latestProfile == nil || profile.ExpirationDate.After(latestProfile.ExpirationDate) {
				latestProfile = &profile
			}
		}
		filteredProfiles = append(filteredProfiles, *latestProfile)
	}
	return filteredProfiles
}

func collectIpaExportCodeSignGroups(tool Tool, archive Archive, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) ([]export.CodeSignGroup, error) {
	collectedCodeSignGroups := []export.CodeSignGroup{}
	_, isMacArchive := archive.(xcarchive.MacosArchive)

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
			collectedCodeSignGroup = export.NewMacGroup(*selectedCertificate, nil, selectedBundleIDProfileMap)
		} else {
			collectedCodeSignGroup = export.NewIOSGroup(*selectedCertificate, selectedBundleIDProfileMap)
		}

		fmt.Println()
		fmt.Println()
		log.Infof("Codesign settings will be used for %s ipa export:", exportMethod(collectedCodeSignGroup))
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

func collectIpaExportCertificate(tool Tool, archiveCertificate certificateutil.CertificateInfoModel, installedCertificates []certificateutil.CertificateInfoModel) (certificateutil.CertificateInfoModel, error) {
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
			return errors.Wrap(err, "failed to find Provisioning Profile")
		}
		profilePathInfoMap[pth] = profile
		log.Printf("file found at: %s", pth)
	}

	if err := exportProvisioningProfiles(profilePathInfoMap, absExportOutputDirPath); err != nil {
		return fmt.Errorf("failed to export the Provisioning Profile into the export directory: %s", err)
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
			return fmt.Errorf("failed to export, error: %s", err)
		}

		if identityRef == nil {
			return errors.New("identity not found in the keychain, or it was invalid (expired)")
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
		return fmt.Errorf("failed to export from Keychain: %s", err)
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

func exportCodesignFiles(tool Tool, archivePath, outputDirPath string) error {
	// archive code sign settings
	installedCertificates, err := certificateutil.InstalledCodesigningCertificateInfos()
	if err != nil {
		return fmt.Errorf("failed to list installed code signing identities, error: %s", err)
	}
	installedCertificates = certificateutil.FilterValidCertificateInfos(installedCertificates)

	log.Debugf("Installed certificates:")
	for _, installedCertificate := range installedCertificates {
		log.Debugf(installedCertificate.String())
	}

	installedProfiles, err := profileutil.InstalledProvisioningProfileInfos(profileutil.ProfileTypeIos)
	if err != nil {
		return fmt.Errorf("failed to list installed provisioning profiles, error: %s", err)
	}

	log.Debugf("Installed profiles:")
	for _, profileInfo := range installedProfiles {
		log.Debugf(profileInfo.String(installedCertificates...))
	}

	certificatesToExport, profilesToExport, err := getFilesToExport(archivePath, installedCertificates, installedProfiles, tool)
	if err != nil {
		return err
	}

	// ipa export code sign settings
	fmt.Println()
	fmt.Println()
	log.Printf("🔦  Analyzing the archive, to get ipa export code signing settings...")

	if err := collectAndExportIdentities(certificatesToExport, outputDirPath); err != nil {
		return err
	}

	if err := collectAndExportProvisioningProfiles(profilesToExport, outputDirPath); err != nil {
		return err
	}

	provProfilesUploaded := (len(profilesToExport) == 0)
	certsUploaded := (len(certificatesToExport) == 0)

	if len(profilesToExport) > 0 || len(certificatesToExport) > 0 {
		fmt.Println()
		shouldUpload, err := askUploadGeneratedFiles()
		if err != nil {
			return err
		}

		if shouldUpload {
			accessToken, err := getAccessToken()
			if err != nil {
				return err
			}

			bitriseClient, appList, err := bitriseclient.NewBitriseClient(accessToken)
			if err != nil {
				return err
			}

			selectedAppSlug, err := selectApp(appList)
			if err != nil {
				return err
			}

			bitriseClient.SetSelectedAppSlug(selectedAppSlug)

			provProfilesUploaded, err = uploadExportedProvProfiles(bitriseClient, profilesToExport, outputDirPath)
			if err != nil {
				return err
			}

			certsUploaded, err = uploadExportedIdentity(bitriseClient, certificatesToExport, outputDirPath)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println()
	log.Successf("Exports finished you can find the exported files at: %s", outputDirPath)

	if err := command.RunCommand("open", outputDirPath); err != nil {
		log.Errorf("Failed to open the export directory in Finder: %s", outputDirPath)
	} else {
		fmt.Println("Opened the directory in Finder.")
	}
	printFinished(provProfilesUploaded, certsUploaded)

	return nil
}

func getFilesToExport(archivePath string, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel, tool Tool) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel, error) {
	macOS, err := xcarchive.IsMacOS(archivePath)
	if err != nil {
		return nil, nil, err
	}

	var certificate certificateutil.CertificateInfoModel
	var archive Archive
	var achiveCodeSignGroup export.CodeSignGroup

	if macOS {
		archive, achiveCodeSignGroup, err = getMacOSCodeSignGroup(archivePath, installedCertificates)
		if err != nil {
			return nil, nil, err
		}
		certificate = achiveCodeSignGroup.Certificate()
	} else {
		archive, achiveCodeSignGroup, err = getIOSCodeSignGroup(archivePath, installedCertificates)
		if err != nil {
			return nil, nil, err
		}
		certificate = achiveCodeSignGroup.Certificate()
	}

	certificatesToExport := []certificateutil.CertificateInfoModel{}
	profilesToExport := []profileutil.ProvisioningProfileInfoModel{}

	if certificatesOnly {
		certificatesToExport, err = exportCertificatesOnly(tool, certificate, installedCertificates, certificatesToExport)
		if err != nil {
			return nil, nil, err
		}
	} else {
		certificatesToExport, profilesToExport, err = exportCertificatesAndProfiles(macOS, archive, tool, certificate, installedCertificates, installedProfiles, certificatesToExport, profilesToExport, achiveCodeSignGroup)
		if err != nil {
			return nil, nil, err
		}
	}

	return certificatesToExport, profilesToExport, nil
}

func exportCertificatesAndProfiles(macOS bool, archive Archive, tool Tool, certificate certificateutil.CertificateInfoModel,
	installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel,
	certificatesToExport []certificateutil.CertificateInfoModel, profilesToExport []profileutil.ProvisioningProfileInfoModel,
	achiveCodeSignGroup export.CodeSignGroup) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel, error) {

	groups, err := collectIpaExportCodeSignGroups(tool, archive, installedCertificates, installedProfiles)
	if err != nil {
		return nil, nil, err
	}

	var ipaExportCodeSignGroups []export.CodeSignGroup
	for _, group := range groups {
		if macOS {
			ipaExportCodeSignGroup, ok := group.(*export.MacCodeSignGroup)
			if ok {
				ipaExportCodeSignGroups = append(ipaExportCodeSignGroups, ipaExportCodeSignGroup)
			}
		} else {
			ipaExportCodeSignGroup, ok := group.(*export.IosCodeSignGroup)
			if ok {
				ipaExportCodeSignGroups = append(ipaExportCodeSignGroups, ipaExportCodeSignGroup)
			}
		}
	}

	if len(ipaExportCodeSignGroups) == 0 {
		return nil, nil, errors.New("no ipa export code sign groups collected")
	}

	codeSignGroups := append(ipaExportCodeSignGroups, achiveCodeSignGroup)
	certificates, profiles := extractCertificatesAndProfiles(codeSignGroups...)
	certificatesToExport = append(certificatesToExport, certificates...)
	profilesToExport = append(profilesToExport, profiles...)

	return certificatesToExport, profilesToExport, nil
}

func exportCertificatesOnly(tool Tool, certificate certificateutil.CertificateInfoModel, installedCertificates []certificateutil.CertificateInfoModel, certificatesToExport []certificateutil.CertificateInfoModel) ([]certificateutil.CertificateInfoModel, error) {
	ipaExportCertificate, err := collectIpaExportCertificate(tool, certificate, installedCertificates)
	if err != nil {
		return nil, err
	}

	certificatesToExport = append(certificatesToExport, certificate, ipaExportCertificate)
	return certificatesToExport, nil
}

func getAccessToken() (string, error) {
	accessToken, err := askAccessToken()
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func uploadExportedProvProfiles(bitriseClient *bitriseclient.BitriseClient, profilesToExport []profileutil.ProvisioningProfileInfoModel, outputDirPath string) (bool, error) {
	fmt.Println()
	log.Infof("Uploading provisioning profiles...")

	profilesToUpload, err := filterAlreadyUploadedProvProfiles(bitriseClient, profilesToExport)
	if err != nil {
		return false, err
	}

	if len(profilesToUpload) > 0 {
		if err := uploadProvisioningProfiles(bitriseClient, profilesToUpload, outputDirPath); err != nil {
			return false, err
		}
	} else {
		log.Warnf("There is no new provisioning profile to upload...")
	}

	return true, nil
}

func uploadExportedIdentity(bitriseClient *bitriseclient.BitriseClient, certificatesToExport []certificateutil.CertificateInfoModel, outputDirPath string) (bool, error) {
	fmt.Println()
	log.Infof("Uploading certificate...")

	shouldUploadIdentities, err := shouldUploadCertificates(bitriseClient, certificatesToExport)
	if err != nil {
		return false, err
	}

	if shouldUploadIdentities {

		if err := UploadIdentity(bitriseClient, outputDirPath); err != nil {
			return false, err
		}
	} else {
		log.Warnf("There is no new certificate to upload...")
	}

	return true, err
}

func askUploadGeneratedFiles() (bool, error) {
	messageToAsk := "Do you want to upload the provisioning profiles and certificates to Bitrise?"
	return goinp.AskForBoolFromReader(messageToAsk, os.Stdin)
}

func askUploadIdentities() (bool, error) {
	messageToAsk := "Do you want to upload the certificates to Bitrise?"
	return goinp.AskForBoolFromReader(messageToAsk, os.Stdin)
}

func filterAlreadyUploadedProvProfiles(client *bitriseclient.BitriseClient, localProfiles []profileutil.ProvisioningProfileInfoModel) ([]profileutil.ProvisioningProfileInfoModel, error) {
	log.Printf("Looking for provisioning profile duplicates on Bitrise...")

	uploadedProfileUUIDList := map[string]bool{}
	profilesToUpload := []profileutil.ProvisioningProfileInfoModel{}

	uploadedProfInfoList, err := client.FetchProvisioningProfiles()
	if err != nil {
		return nil, err
	}

	for _, uploadedProfileInfo := range uploadedProfInfoList {
		uploadedProfileUUID, err := client.GetUploadedProvisioningProfileUUIDby(uploadedProfileInfo.Slug)
		if err != nil {
			return nil, err
		}

		uploadedProfileUUIDList[uploadedProfileUUID] = true
	}

	for _, localProfile := range localProfiles {
		contains, _ := uploadedProfileUUIDList[localProfile.UUID]
		if contains {
			log.Warnf("Already on Bitrise: - %s - (UUID: %s) ", localProfile.Name, localProfile.UUID)
		} else {
			profilesToUpload = append(profilesToUpload, localProfile)
		}
	}

	return profilesToUpload, nil
}

func shouldUploadCertificates(client *bitriseclient.BitriseClient, certificatesToExport []certificateutil.CertificateInfoModel) (bool, error) {
	log.Printf("Looking for certificate duplicates on Bitrise...")

	var uploadedCertificatesSerialList []string
	localCertificatesSerialList := []string{}

	uploadedItentityList, err := client.FetchUploadedIdentities()
	if err != nil {
		return false, err
	}

	// Get uploaded certificates' serials
	for _, uploadedIdentity := range uploadedItentityList {
		var serialListAsString []string

		serialList, err := client.GetUploadedCertificatesSerialby(uploadedIdentity.Slug)
		if err != nil {
			return false, err
		}

		for _, serial := range serialList {
			serialListAsString = append(serialListAsString, serial.String())
		}
		uploadedCertificatesSerialList = append(uploadedCertificatesSerialList, serialListAsString...)
	}

	for _, certificateToExport := range certificatesToExport {
		localCertificatesSerialList = append(localCertificatesSerialList, certificateToExport.Serial)
	}

	log.Debugf("Uploaded certificates' serial list: \n\t%v", uploadedCertificatesSerialList)
	log.Debugf("Local certificates' serial list: \n\t%v", localCertificatesSerialList)

	// Search for a new certificate
	for _, localCertificateSerial := range localCertificatesSerialList {
		if !sliceutil.IsStringInSlice(localCertificateSerial, uploadedCertificatesSerialList) {
			return true, nil
		}
	}

	return false, nil
}

// ----------------------------------------------------------------
// --- Upload methods
func uploadProvisioningProfiles(bitriseClient *bitriseclient.BitriseClient, profilesToUpload []profileutil.ProvisioningProfileInfoModel, outputDirPath string) error {
	for _, profile := range profilesToUpload {
		exportFileName := provProfileExportFileName(profile, outputDirPath)

		provProfile, err := os.Open(outputDirPath + "/" + exportFileName)
		if err != nil {
			return err
		}

		defer func() {
			if err := provProfile.Close(); err != nil {
				log.Warnf("Provisioning profile close failed, err: %s", err)
			}

		}()

		info, err := provProfile.Stat()
		if err != nil {
			return err
		}

		log.Debugf("\n%s size: %d", exportFileName, info.Size())

		provProfSlugResponseData, err := bitriseClient.RegisterProvisioningProfile(info.Size(), profile)
		if err != nil {
			return err
		}

		if err := bitriseClient.UploadProvisioningProfile(provProfSlugResponseData.UploadURL, provProfSlugResponseData.UploadFileName, outputDirPath, exportFileName); err != nil {
			return err
		}

		if err := bitriseClient.ConfirmProvisioningProfileUpload(provProfSlugResponseData.Slug, provProfSlugResponseData.UploadFileName); err != nil {
			return err
		}
	}

	return nil
}

// UploadIdentity ...
func UploadIdentity(bitriseClient *bitriseclient.BitriseClient, outputDirPath string) error {
	identities, err := os.Open(outputDirPath + "/" + "Identities.p12")
	if err != nil {
		return err
	}

	defer func() {
		if err := identities.Close(); err != nil {
			log.Warnf("Identities failed, err: %s", err)
		}

	}()

	info, err := identities.Stat()
	if err != nil {
		return err
	}

	log.Debugf("\n%s size: %d", "Identities.p12", info.Size())

	certificateResponseData, err := bitriseClient.RegisterIdentity(info.Size())
	if err != nil {
		return err
	}

	if err := bitriseClient.UploadIdentity(certificateResponseData.UploadURL, certificateResponseData.UploadFileName, outputDirPath, "Identities.p12"); err != nil {
		return err
	}

	return bitriseClient.ConfirmIdentityUpload(certificateResponseData.Slug, certificateResponseData.UploadFileName)
}

func askAccessToken() (token string, err error) {
	messageToAsk := `Please copy your personal access token to Bitrise.
(To acquire a Personal Access Token for your user, sign in with that user on bitrise.io, go to your Account Settings page,
and select the Security tab on the left side.)`
	fmt.Println()

	accesToken, err := goinp.AskForStringFromReader(messageToAsk, os.Stdin)
	if err != nil {
		return accesToken, err
	}

	fmt.Println()
	log.Infof("%s %s", colorstring.Green("Given accesToken:"), accesToken)
	fmt.Println()

	return accesToken, nil
}

func selectApp(appList []bitriseclient.Application) (seledtedAppSlug string, err error) {
	var selectionList []string

	for _, app := range appList {
		selectionList = append(selectionList, app.Title+" ("+app.RepoURL+")")
	}
	userSelection, err := goinp.SelectFromStringsWithDefault("Select the app which you want to upload the privisioning profiles", 1, selectionList)

	if err != nil {
		return "", fmt.Errorf("failed to read input: %s", err)

	}

	log.Debugf("selected app: %v", userSelection)

	for index, selected := range selectionList {
		if selected == userSelection {
			return appList[index].Slug, nil
		}
	}

	return "", &appSelectionError{"failed to find selected app in appList"}
}

type appSelectionError struct {
	s string
}

func (e *appSelectionError) Error() string {
	return e.s
}
