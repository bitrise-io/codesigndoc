package codesigndoc

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/codesigndoc/bitriseio/bitrise"

	"github.com/bitrise-io/codesigndoc/bitriseio"
	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/export"
	"github.com/bitrise-io/go-xcode/profileutil"
	"github.com/bitrise-io/go-xcode/xcarchive"
	"github.com/bitrise-io/goinp/goinp"
)

const collectCodesigningFilesInfo = `To collect available code sign files, we search for installed Provisioning Profiles:"
- which has installed Codesign Identity in your Keychain"
- which can provision your application target's bundle ids"
- which has the project defined Capabilities set"
- which matches to the selected export method"
`

// UploadConfig contains configuration to automatically upload artifacts to bitrise.io
type UploadConfig struct {
	PersonalAccessToken string
	AppSlug             string
}

func (config *UploadConfig) isValid() bool {
	return (strings.TrimSpace(config.PersonalAccessToken) != "") &&
		(strings.TrimSpace(config.AppSlug) != "")
}

// CollectCodesignFiles collects the codesigning files required to create an xcode archive
// and filers them for the specified export method
func CollectCodesignFiles(archivePath string, certificatesOnly bool) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel, error) {
	// Find out the XcArchive type
	isMacOs, err := xcarchive.IsMacOS(archivePath)
	if err != nil {
		return nil, nil, err
	}

	// Set up the XcArchive type for certs and profiles.
	certificateType := codesign.IOSCertificate
	profileType := profileutil.ProfileTypeIos
	if isMacOs {
		certificateType = codesign.MacOSCertificate
		profileType = profileutil.ProfileTypeMacOs
	}

	// Certificates
	certificates, err := codesign.InstalledCertificates(certificateType)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list installed code signing identities, error: %s", err)
	}

	installerCertificates, err := certificateutil.InstalledInstallerCertificateInfos()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list installed code signing identities, error: %s", err)
	}

	log.Debugf("Installed certificates:")
	for _, installedCertificate := range certificates {
		log.Debugf(installedCertificate.String())
	}

	// Profiles
	profiles, err := profileutil.InstalledProvisioningProfileInfos(profileType)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list installed provisioning profiles, error: %s", err)
	}

	log.Debugf("Installed profiles:")
	for _, profileInfo := range profiles {
		log.Debugf(profileInfo.String(certificates...))
	}

	// export code sign settings
	fmt.Println()
	fmt.Println()
	log.Printf("🔦  Analyzing the archive, to get export code signing settings...")
	return getFilesToExport(archivePath, certificates, installerCertificates, profiles, certificatesOnly)
}

// UploadAndWriteCodesignFiles exports then uploads codesign files to bitrise.io and saves them to output folder
func UploadAndWriteCodesignFiles(certificates []certificateutil.CertificateInfoModel, profiles []profileutil.ProvisioningProfileInfoModel, askForPassword bool, outputDirPath string, uploadConfig UploadConfig) (bool, bool, error) {
	identities, err := codesign.CollectAndExportIdentitiesAsReader(certificates, askForPassword)
	if err != nil {
		return false, false, err
	}

	provisioningProfiles, err := codesign.CollectAndExportProvisioningProfilesAsReader(profiles)
	if err != nil {
		return false, false, err
	}

	var client *bitrise.Client
	if uploadConfig.isValid() {
		// Upload automatically if token is provided as CLI paramter, do not export to filesystem
		// Used to upload artifacts as part of an other CLI tool
		client, err = bitrise.NewClientAsStream(uploadConfig.PersonalAccessToken)
		if err != nil {
			return false, false, err
		}
		client.SetSelectedAppSlug(uploadConfig.AppSlug)
	}

	if client == nil {
		uploadConfirmMsg := "Do you want to upload the provisioning profiles and certificates to Bitrise?"
		if len(provisioningProfiles) == 0 {
			uploadConfirmMsg = "Do you want to upload the certificates to Bitrise?"
		}
		fmt.Println()
		if shouldUpload, err := goinp.AskForBoolFromReader(uploadConfirmMsg, os.Stdin); err != nil {
			return false, false, err
		} else if shouldUpload {
			client, err = bitriseio.GetInteractiveConfigClient()
		}
	}

	provProfilesUploaded := (len(profiles) == 0)
	certsUploaded := (len(certificates) == 0)
	if client != nil {
		certsUploaded, provProfilesUploaded, err = bitriseio.UploadCodesigningFilesAsStream(client, identities, provisioningProfiles)
		if err != nil {
			return false, false, err
		}
	}

	if strings.TrimSpace(outputDirPath) != "" {
		if err := codesign.WriteIdentities(identities.Content, outputDirPath); err != nil {
			return false, false, err
		}
		if err := codesign.WriteProvisioningProfilesAsStream(provisioningProfiles, outputDirPath); err != nil {
			return false, false, err
		}
		fmt.Println()
		log.Successf("Exports finished you can find the exported files at: %s", outputDirPath)

		if err := command.RunCommand("open", outputDirPath); err != nil {
			log.Errorf("Failed to open the export directory in Finder: %s", outputDirPath)
		} else {
			fmt.Println("Opened the directory in Finder.")
		}
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
