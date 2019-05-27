package codesigndocuitests

import (
	"errors"
	"fmt"
	"os"

	"github.com/bitrise-io/codesigndoc/bitriseio"
	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/export"
	"github.com/bitrise-io/go-xcode/profileutil"
	"github.com/bitrise-io/goinp/goinp"
)

const collectCodesigningFilesInfo = `To collect available code sign files, we search for installed Provisioning Profiles:"
- which has installed Codesign Identity in your Keychain"
- which can provision your application target's bundle ids"
- which has the project defined Capabilities set"
`

// ExportCodesignFiles exports the codesigning files for the UITests-Runner.app
// and exports the codesigning files for the specified export method
func ExportCodesignFiles(buildPath, outputDirPath string, certificatesOnly bool, askForPassword bool) (bool, bool, error) {
	// Find out the XcArchive type
	certificateType := codesign.IOSCertificate
	profileType := profileutil.ProfileTypeIos

	// Certificates
	certificates, err := codesign.InstalledCertificates(certificateType)
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

	certificatesToExport, profilesToExport, err := getFilesToExport(buildPath, certificates, profiles, certificatesOnly)
	if err != nil {
		return false, false, err
	}

	// export code sign settings
	fmt.Println()
	fmt.Println()
	log.Printf("🔦  Analyzing the builded artifacts, to get code signing settings...")

	if err := codesign.CollectAndExportIdentities(certificatesToExport, outputDirPath, askForPassword); err != nil {
		return false, false, err
	}

	if err := codesign.CollectAndExportProvisioningProfiles(profilesToExport, outputDirPath); err != nil {
		return false, false, err
	}

	provProfilesUploaded := (len(profilesToExport) == 0)
	certsUploaded := (len(certificatesToExport) == 0)

	var shouldUpload bool
	if !certificatesOnly {
		fmt.Println()
		shouldUpload, err = goinp.AskForBoolFromReader("Do you want to upload the provisioning profiles and certificates to Bitrise?", os.Stdin)
		if err != nil {
			return false, false, err
		}
	} else {
		shouldUpload, err = goinp.AskForBoolFromReader("Do you want to upload the certificates to Bitrise?", os.Stdin)
		if err != nil {
			return false, false, err
		}
	}

	if shouldUpload {
		certsUploaded, provProfilesUploaded, err = bitriseio.UploadCodesigningFiles(certificatesToExport, profilesToExport, certificatesOnly, outputDirPath)
		if err != nil {
			return false, false, err
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
