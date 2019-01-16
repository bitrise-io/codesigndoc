package codesigndocuitests

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/pathutil"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/bitriseio"
	"github.com/bitrise-tools/codesigndoc/codesign"
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
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
	_, err := getEmbedProfile(buildPath)
	if err != nil {
		return nil, nil, err
	}

	certificatesToExport := []certificateutil.CertificateInfoModel{}
	profilesToExport := []profileutil.ProvisioningProfileInfoModel{}

	if certificatesOnly {
		exportCertificate, err := collectExportCertificate(installedCertificates)
		if err != nil {
			return nil, nil, err
		}

		certificatesToExport = append(certificatesToExport, exportCertificate...)
	} else {
		// certificatesToExport, profilesToExport, err = collectCertificatesAndProfiles(archive, certificate, installedCertificates, installedProfiles, certificatesToExport, profilesToExport, achiveCodeSignGroup)
		// if err != nil {
		// 	return nil, nil, err
		// }
	}

	return certificatesToExport, profilesToExport, nil
}

func getEmbedProfile(buildPath string) (profileutil.ProvisioningProfileInfoModel, error) {
	runnerPattern := filepath.Join(buildPath, "*-Runner.app")
	possibleTestRunnerPths, err := filepath.Glob(runnerPattern)
	if err != nil {
		return profileutil.ProvisioningProfileInfoModel{}, err
	}

	if len(possibleTestRunnerPths) == 0 {
		return profileutil.ProvisioningProfileInfoModel{}, fmt.Errorf("no Test-Runner.app found in %s", buildPath)
	} else if len(possibleTestRunnerPths) != 1 {
		return profileutil.ProvisioningProfileInfoModel{}, fmt.Errorf("found multiple Test-Runner.app in %s", buildPath)
	}

	embedProfilePath := filepath.Join(possibleTestRunnerPths[0], "embedded.mobileprovision")
	if exist, err := pathutil.IsPathExists(embedProfilePath); err != nil {
		return profileutil.ProvisioningProfileInfoModel{}, fmt.Errorf("failed to find embedded.mobileprovision in %s, error: %s", embedProfilePath, err)
	} else if !exist {
		return profileutil.ProvisioningProfileInfoModel{}, fmt.Errorf("failed to find embedded.mobileprovision in %s", embedProfilePath)
	}

	return profileutil.NewProvisioningProfileInfoFromFile(embedProfilePath)
}
