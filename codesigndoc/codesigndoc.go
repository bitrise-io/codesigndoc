package codesigndoc

import (
	"errors"
	"fmt"

	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/export"
	"github.com/bitrise-io/go-xcode/profileutil"
	"github.com/bitrise-io/go-xcode/xcarchive"
)

const collectCodesigningFilesInfo = `To collect available code sign files, we search for installed Provisioning Profiles:"
- which has installed Codesign Identity in your Keychain"
- which can provision your application target's bundle ids"
- which has the project defined Capabilities set"
- which matches to the selected export method"
`

// CollectCodesignFiles collects the codesigning files required to create an xcode archive
// and filers them for the specified export method.
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

func getFilesToExport(archivePath string, installedCertificates []certificateutil.CertificateInfoModel, installedInstallerCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel, certificatesOnly bool) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel, error) {
	macOS, err := xcarchive.IsMacOS(archivePath)
	if err != nil {
		return nil, nil, err
	}

	var certificate certificateutil.CertificateInfoModel
	var archive Archive
	var archiveCodeSignGroup export.CodeSignGroup

	if macOS {
		archive, archiveCodeSignGroup, err = getMacOSCodeSignGroup(archivePath, installedCertificates)
		if err != nil {
			return nil, nil, err
		}
		certificate = archiveCodeSignGroup.Certificate()
	} else {
		archive, archiveCodeSignGroup, err = getIOSCodeSignGroup(archivePath, installedCertificates)
		if err != nil {
			return nil, nil, err
		}
		certificate = archiveCodeSignGroup.Certificate()
	}

	var certificatesToExport []certificateutil.CertificateInfoModel
	var profilesToExport []profileutil.ProvisioningProfileInfoModel

	if certificatesOnly {
		exportCertificate, err := collectExportCertificate(macOS, certificate, installedCertificates, installedInstallerCertificates)
		if err != nil {
			return nil, nil, err
		}

		certificatesToExport = append(certificatesToExport, certificate)
		certificatesToExport = append(certificatesToExport, exportCertificate...)
	} else {
		certificatesToExport, profilesToExport, err = collectCertificatesAndProfiles(archive, installedCertificates, installedProfiles, certificatesToExport, profilesToExport, archiveCodeSignGroup)
		if err != nil {
			return nil, nil, err
		}
	}

	return certificatesToExport, profilesToExport, nil
}

func collectCertificatesAndProfiles(archive Archive,
	installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel,
	certificatesToExport []certificateutil.CertificateInfoModel, profilesToExport []profileutil.ProvisioningProfileInfoModel,
	archiveCodeSignGroup export.CodeSignGroup) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel, error) {

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

	codeSignGroups := append(exportCodeSignGroups, archiveCodeSignGroup)
	certificates, profiles := extractCertificatesAndProfiles(codeSignGroups...)
	certificatesToExport = append(certificatesToExport, certificates...)
	profilesToExport = append(profilesToExport, profiles...)

	return certificatesToExport, profilesToExport, nil
}
