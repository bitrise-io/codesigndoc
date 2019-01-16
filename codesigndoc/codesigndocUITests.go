package codesigndoc

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/bitriseio"
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

const collectCodesigningFilesInfo__UITests = `To collect available code sign files, we search for installed Provisioning Profiles:"
- which has installed Codesign Identity in your Keychain"
- which can provision your application target's bundle ids"
- which has the project defined Capabilities set"
- which matches to the selected export method"
`

// ExportCodesignFiles_UITests exports the codesigning files required to create an xcode archive
// and exports the codesigning files for the specified export method
func ExportCodesignFiles_UITests(buildPath, outputDirPath string, certificatesOnly bool, askForPassword bool) (bool, bool, error) {
	// Find out the XcArchive type
	certificateType := IOSCertificate
	profileType := profileutil.ProfileTypeIos

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

	certificatesToExport, profilesToExport, err := getFilesToExport(buildPath, certificates, installerCertificates, profiles, certificatesOnly)
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
