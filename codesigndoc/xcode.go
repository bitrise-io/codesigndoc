package codesigndoc

import (
	"fmt"

	"github.com/bitrise-io/codesigndoc/xcode"

	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/codesigndoc/models"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/stringutil"
)

// GenerateXCodeArchive ...
func GenerateXCodeArchive(xcodeCmd xcode.CommandModel) (string, string, error) {
	fmt.Println()
	log.Printf("🔦  Running an Xcode Archive, to get all the required code signing settings...")
	archivePath, xcodebuildOutput, err := xcodeCmd.GenerateArchive()

	if err != nil {
		log.Warnf("Last lines of build log:")
		fmt.Println(stringutil.LastNLines(xcodebuildOutput, 15))
		fmt.Println()
		log.Printf("Open the project: %s", xcodeCmd.ProjectFilePath)
		log.Printf("and Archive, using the Scheme: %s", xcodeCmd.Scheme)
		fmt.Println()
		return "", "", err
	}

	return archivePath, xcodebuildOutput, nil
}

// CodesigningFilesForXCodeProject ...
func CodesigningFilesForXCodeProject(archivePath string, certificatesOnly bool, isAskForPassword bool) (models.Certificates, []models.ProvisioningProfile, error) {
	// If certificatesOnly is set, CollectCodesignFiles returns an empty slice for profiles
	certificatesToExport, profilesToExport, err := CollectCodesignFiles(archivePath, certificatesOnly)
	if err != nil {
		return models.Certificates{}, nil, err
	}

	return codesign.ExportCodesigningFiles(certificatesToExport, profilesToExport, isAskForPassword)
}
