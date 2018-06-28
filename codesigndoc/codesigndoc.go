package codesigndoc

import (
	"errors"
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/bitriseio"
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/export"
	"github.com/bitrise-tools/go-xcode/profileutil"
	"github.com/bitrise-tools/go-xcode/xcarchive"
)

const collectCodesigningFilesInfo = `To collect available code sign files, we search for installed Provisioning Profiles:"
- which has installed Codesign Identity in your Keychain"
- which can provision your application target's bundle ids"
- which has the project defined Capabilities set"
- which matches to the selected export method"
`

// ExportCodesignFiles exports the codesigning files required to create an xcode archive
// and exports the codesigning files for the specified export method
func ExportCodesignFiles(archivePath, outputDirPath string, certificatesOnly bool, askForPassword bool) (bool, bool, error) {
	// Find out the XcArchive type
	isMacOs, err := xcarchive.IsMacOS(archivePath)
	if err != nil {
		return false, false, err
	}

	// Set up the XcArchive type for certs and profiles.
	certificateType := IOSCertificate
	profileType := profileutil.ProfileTypeIos
	if isMacOs {
		certificateType = MacOSCertificate
		profileType = profileutil.ProfileTypeMacOs
	}

	// Certificates
	certificates, err := installedCertificates(certificateType)
	if err != nil {
		return false, false, fmt.Errorf("failed to list installed code signing identities, error: %s", err)
	}

	installerCertificates, err := certificateutil.InstalledInstallerCertificateInfos()
	if err != nil {
		return false, false, fmt.Errorf("failed to list installed code signing identities, error: %s", err)
	}

	log.Debugf("Installed certificates:")
	for _, installedCertificate := range certificates {
		log.Debugf(installedCertificate.String())
	}

	// Profiles
	profiles, err := profileutil.InstalledProvisioningProfileInfos(profileType)
	if err != nil {
		return false, false, fmt.Errorf("failed to list installed provisioning profiles, error: %s", err)
	}

	log.Debugf("Installed profiles:")
	for _, profileInfo := range profiles {
		log.Debugf(profileInfo.String(certificates...))
	}

	certificatesToExport, profilesToExport, err := getFilesToExport(archivePath, certificates, installerCertificates, profiles, certificatesOnly)
	if err != nil {
		return false, false, err
	}

	// export code sign settings
	fmt.Println()
	fmt.Println()
	log.Printf("🔦  Analyzing the archive, to get export code signing settings...")

	if err := collectAndExportIdentities(certificatesToExport, outputDirPath, askForPassword); err != nil {
		return false, false, err
	}

	if err := collectAndExportProvisioningProfiles(profilesToExport, outputDirPath); err != nil {
		return false, false, err
	}

	provProfilesUploaded := (len(profilesToExport) == 0)
	certsUploaded := (len(certificatesToExport) == 0)

	if len(profilesToExport) > 0 || len(certificatesToExport) > 0 {
		fmt.Println()
		shouldUpload, err := goinp.AskForBoolFromReader("Do you want to upload the provisioning profiles and certificates to Bitrise?", os.Stdin)
		if err != nil {
			return false, false, err
		}

		if shouldUpload {
			certsUploaded, provProfilesUploaded, err = bitriseio.UploadCodesigningFiles(certificatesToExport, profilesToExport, outputDirPath)
			if err != nil {
				return false, false, err
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

	return certsUploaded, provProfilesUploaded, nil
}

func getFilesToExport(archivePath string, installedCertificates []certificateutil.CertificateInfoModel, installedInstallerCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel, certificatesOnly bool) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel, error) {
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
		exportCertificate, err := collectExportCertificate(macOS, certificate, installedCertificates, installedInstallerCertificates)
		if err != nil {
			return nil, nil, err
		}

		certificatesToExport = append(certificatesToExport, certificate)
		certificatesToExport = append(certificatesToExport, exportCertificate...)
	} else {
		certificatesToExport, profilesToExport, err = collectCertificatesAndProfiles(archive, certificate, installedCertificates, installedProfiles, certificatesToExport, profilesToExport, achiveCodeSignGroup)
		if err != nil {
			return nil, nil, err
		}
	}

	return certificatesToExport, profilesToExport, nil
}

func collectCertificatesAndProfiles(archive Archive, certificate certificateutil.CertificateInfoModel,
	installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel,
	certificatesToExport []certificateutil.CertificateInfoModel, profilesToExport []profileutil.ProvisioningProfileInfoModel,
	achiveCodeSignGroup export.CodeSignGroup) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel, error) {

	_, macOS := archive.(xcarchive.MacosArchive)

	groups, err := collectExportCodeSignGroups(archive, installedCertificates, installedProfiles)
	if err != nil {
		return nil, nil, err
	}

	var exportCodeSignGroups []export.CodeSignGroup
	for _, group := range groups {
		if macOS {
			exportCodeSignGroup, ok := group.(*export.MacCodeSignGroup)
			if ok {
				exportCodeSignGroups = append(exportCodeSignGroups, exportCodeSignGroup)
			}
		} else {
			exportCodeSignGroup, ok := group.(*export.IosCodeSignGroup)
			if ok {
				exportCodeSignGroups = append(exportCodeSignGroups, exportCodeSignGroup)
			}
		}
	}

	if len(exportCodeSignGroups) == 0 {
		return nil, nil, errors.New("no export code sign groups collected")
	}

	codeSignGroups := append(exportCodeSignGroups, achiveCodeSignGroup)
	certificates, profiles := extractCertificatesAndProfiles(codeSignGroups...)
	certificatesToExport = append(certificatesToExport, certificates...)
	profilesToExport = append(profilesToExport, profiles...)

	return certificatesToExport, profilesToExport, nil
}
