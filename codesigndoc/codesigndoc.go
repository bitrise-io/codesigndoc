package codesigndoc

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/colorstring"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/bitriseio"
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
	"github.com/bitrise-tools/go-xcode/xcarchive"
)

const collectCodesigningFilesInfo = `To collect available code sign files, we search for installed Provisioning Profiles:"
- which has installed Codesign Identity in your Keychain"
- which can provision your application target's bundle ids"
- which has the project defined Capabilities set"
- which matches to the selected ipa export method"
`

// ExportCodesignFiles exports the codesigning files required to create an xcode archive
// and exports the codesigning files for the specified export method
func ExportCodesignFiles(archivePath, outputDirPath string, certificatesOnly bool, askForPassword bool) (bool, bool, error) {
	certificates, err := installedCertificates(IOSCertificate)
	if err != nil {
		return false, false, fmt.Errorf("failed to list installed code signing identities, error: %s", err)
	}

	profiles, err := profileutil.InstalledProvisioningProfileInfos(profileutil.ProfileTypeIos)
	if err != nil {
		return false, false, fmt.Errorf("failed to list installed provisioning profiles, error: %s", err)
	}

	// archive code sign settings
	archive, err := xcarchive.NewIosArchive(archivePath)
	if err != nil {
		return false, false, fmt.Errorf("failed to analyze archive, error: %s", err)
	}

	archiveCodeSignGroup, err := analyzeArchive(archive, certificates)
	if err != nil {
		return false, false, fmt.Errorf("failed to analyze the archive, error: %s", err)
	}

	fmt.Println()
	log.Infof("Codesign settings used for archive:")
	fmt.Println()
	printCodesignGroup(archiveCodeSignGroup)

	// ipa export code sign settings
	fmt.Println()
	fmt.Println()
	log.Printf("🔦  Analyzing the archive, to get ipa export code signing settings...")

	certificatesToExport := []certificateutil.CertificateInfoModel{}
	profilesToExport := []profileutil.ProvisioningProfileInfoModel{}

	if certificatesOnly {
		ipaExportCertificate, err := collectIpaExportCertificate(archiveCodeSignGroup.Certificate, certificates)
		if err != nil {
			return false, false, err
		}

		certificatesToExport = append(certificatesToExport, archiveCodeSignGroup.Certificate, ipaExportCertificate)
	} else {
		ipaExportCodeSignGroups, err := collectIpaExportCodeSignGroups(archive, certificates, profiles)
		if err != nil {
			return false, false, err
		}

		if len(ipaExportCodeSignGroups) == 0 {
			errorString := "\n🚨  " + colorstring.Red("Failed to collect codesigning files for the selected distribution type.\n") +
				colorstring.Yellow("Export an ipa with the same export method which code signing files you want to collect (e.g app-store if you want to collect the code signing files for app-store distribution) in your local xcode and run codesigndoc again.\n") +
				colorstring.Yellow("If the tool fails please report the issue with the codesigndoc log and the local ipa exportOptions.plist")
			return false, false, fmt.Errorf(errorString)
		}

		codeSignGroups := append(ipaExportCodeSignGroups, archiveCodeSignGroup)
		certificates, profiles := extractCertificatesAndProfiles(codeSignGroups...)
		certificatesToExport = append(certificatesToExport, certificates...)
		profilesToExport = append(profilesToExport, profiles...)
	}

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
