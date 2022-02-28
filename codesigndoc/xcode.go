package codesigndoc

import (
	"fmt"

	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/codesigndoc/models"
	"github.com/bitrise-io/codesigndoc/xcode"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/stringutil"
	"github.com/bitrise-io/go-xcode/utility"
)

// BuildXcodeArchive builds an Xcode archive.
func BuildXcodeArchive(xcodeCmd xcode.CommandModel, handleBuildLog func(string) error) (archivePath string, err error) {
	// Output tools versions
	xcodebuildVersion, err := utility.GetXcodeVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get Xcode (xcodebuild) version, error: %s", err)
	}
	fmt.Println()
	log.Infof("%s: %s (%s)", colorstring.Green("Xcode (xcodebuild) version"), xcodebuildVersion.Version, xcodebuildVersion.BuildVersion)

	fmt.Println()
	fmt.Println()
	log.Printf("🔦  Running an Xcode Archive, to get all the required code signing settings...")

	archivePath, xcodebuildOutput, err := xcodeCmd.GenerateArchive()

	defer func() {
		if handleBuildLog != nil {
			if derr := handleBuildLog(xcodebuildOutput); derr != nil {
				if err != nil {
					err = derr
				}
			}
		}
	}()

	if err != nil {
		log.Warnf("Last lines of the build log:")
		fmt.Println(stringutil.LastNLines(xcodebuildOutput, 15))

		log.Infof(colorstring.Yellow("Please check the build log to see what caused the error."))
		fmt.Println()

		log.Errorf("Xcode Archive failed.")
		log.Infof(colorstring.Yellow("Open the project: ")+"%s", xcodeCmd.ProjectFilePath)
		log.Infof(colorstring.Yellow("and make sure that you can build an Archive, with the scheme: ")+"%s", xcodeCmd.Scheme)
		fmt.Println()

		return "", err
	}

	return archivePath, nil
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
