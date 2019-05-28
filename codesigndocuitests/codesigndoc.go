package codesigndocuitests

import (
	"errors"
	"fmt"

	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/export"
	"github.com/bitrise-io/go-xcode/profileutil"
)

const collectCodesigningFilesInfo = `To collect available code sign files, we search for installed Provisioning Profiles:"
- which has installed Codesign Identity in your Keychain"
- which can provision your application target's bundle ids"
- which has the project defined Capabilities set"
`

// CollectCodesignFiles collects the codesigning files for the UITests-Runner.app
// and filters them for the specified export method
func CollectCodesignFiles(buildPath string, certificatesOnly bool) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel, error) {
	// Find out the XcArchive type
	certificateType := codesign.IOSCertificate
	profileType := profileutil.ProfileTypeIos

	// Certificates
	certificates, err := codesign.InstalledCertificates(certificateType)
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

	return getFilesToExport(buildPath, certificates, profiles, certificatesOnly)
}

func getFilesToExport(buildPath string, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel, certificatesOnly bool) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel, error) {
	certificatesToExport := []certificateutil.CertificateInfoModel{}
	profilesToExport := []profileutil.ProvisioningProfileInfoModel{}

	if certificatesOnly {
		exportCertificate, err := collectExportCertificate(installedCertificates)
		if err != nil {
			return nil, nil, err
		}

		certificatesToExport = append(certificatesToExport, exportCertificate...)
	} else {
		testRunners, err := NewIOSTestRunners(buildPath)
		if err != nil {
			return nil, nil, err
		}

		for _, testRunner := range testRunners {
			certsToExport, profsToExport, err := collectCertificatesAndProfiles(*testRunner, installedCertificates, installedProfiles)
			if err != nil {
				return nil, nil, err
			}

			certificatesToExport = append(certificatesToExport, certsToExport...)
			profilesToExport = append(profilesToExport, profsToExport...)
		}

	}

	return certificatesToExport, profilesToExport, nil
}

func collectCertificatesAndProfiles(testRunner IOSTestRunner, installedCertificates []certificateutil.CertificateInfoModel,
	installedProfiles []profileutil.ProvisioningProfileInfoModel) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel, error) {

	groups, err := collectExportCodeSignGroups(testRunner, installedCertificates, installedProfiles)
	if err != nil {
		return nil, nil, err
	}

	var exportCodeSignGroups []export.CodeSignGroup
	for _, group := range groups {
		exportCodeSignGroup, ok := group.(*export.IosCodeSignGroup)
		if ok {
			exportCodeSignGroups = append(exportCodeSignGroups, exportCodeSignGroup)
		}
	}

	if len(exportCodeSignGroups) == 0 {
		return nil, nil, errors.New("no export code sign groups collected")
	}

	certificates, profiles := extractCertificatesAndProfiles(exportCodeSignGroups...)
	return certificates, profiles, nil
}
