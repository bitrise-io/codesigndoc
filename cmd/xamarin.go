package cmd

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/xamarin"
	"github.com/bitrise-tools/go-xamarin/analyzers/project"
	"github.com/bitrise-tools/go-xamarin/analyzers/solution"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
	"github.com/bitrise-tools/go-xcode/xcarchive"
	"github.com/spf13/cobra"
)

// xamarinCmd represents the xamarin command
var xamarinCmd = &cobra.Command{
	Use:   "xamarin",
	Short: "Xamarin project scanner",
	Long:  `Scan a Xamarin project`,

	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          scanXamarinProject,
}

var (
	paramXamarinSolutionFilePath  = ""
	paramXamarinProjectName       = ""
	paramXamarinConfigurationName = ""
)

func init() {
	scanCmd.AddCommand(xamarinCmd)

	xamarinCmd.Flags().StringVar(&paramXamarinSolutionFilePath, "file", "", `Xamarin Solution file path`)
	xamarinCmd.Flags().StringVar(&paramXamarinProjectName, "project", "", `Xamarin iOS Project Name (e.g.: "MyProject.iOS")`)
	xamarinCmd.Flags().StringVar(&paramXamarinConfigurationName, "config", "", `Xamarin Configuration Name (e.g.: "Release|iPhone")`)
}

func printXamarinScanFinishedWithError(format string, args ...interface{}) error {
	return printFinishedWithError("Xamarin Studio", format, args...)
}

func scanXamarinProject(cmd *cobra.Command, args []string) error {
	initLogVerbosity()
	absExportOutputDirPath, err := initExportOutputDir()
	if err != nil {
		return printXamarinScanFinishedWithError("Failed to prepare Export directory: %s", err)
	}

	xamarinCmd := xamarin.CommandModel{}
	// --- Inputs ---

	// Xamarin Solution Path
	xamarinCmd.SolutionFilePath = paramXamarinSolutionFilePath
	if xamarinCmd.SolutionFilePath == "" {
		askText := `Please drag-and-drop your Xamarin Solution (` + colorstring.Green(".sln") + `) file,
and then hit Enter`
		fmt.Println()
		projpth, err := goinp.AskForPath(askText)
		if err != nil {
			return printXamarinScanFinishedWithError("Failed to read input: %s", err)
		}
		xamarinCmd.SolutionFilePath = projpth
	}
	log.Debugf("xamSolutionPth: %s", xamarinCmd.SolutionFilePath)

	xamSln, err := solution.New(xamarinCmd.SolutionFilePath, true)
	if err != nil {
		return printXamarinScanFinishedWithError("Failed to analyze Xamarin solution: %s", err)
	}
	log.Debugf("xamSln: %#v", xamSln)

	// filter only the iOS "app"" projects
	xamarinProjectsToChooseFrom := []project.Model{}
	for _, aXamarinProject := range xamSln.ProjectMap {
		switch aXamarinProject.SDK {
		case constants.SDKIOS, constants.SDKTvOS, constants.SDKMacOS:
			if aXamarinProject.OutputType == "exe" {
				// possible project
				xamarinProjectsToChooseFrom = append(xamarinProjectsToChooseFrom, aXamarinProject)
			}
		default:
			continue
		}
	}
	log.Debugf("%d xamarinProjectsToChooseFrom: %#v", len(xamarinProjectsToChooseFrom), xamarinProjectsToChooseFrom)

	// Xamarin Project
	selectedXamarinProject := project.Model{}
	{
		if len(xamarinProjectsToChooseFrom) < 1 {
			return printXamarinScanFinishedWithError(
				"No acceptable Project found in the provided Solution, or none can be used for iOS Archive.",
			)
		}

		if paramXamarinProjectName != "" {
			// project specified via flag/param
			for _, aProj := range xamarinProjectsToChooseFrom {
				if paramXamarinProjectName == aProj.Name {
					selectedXamarinProject = aProj
					break
				}
			}
			if selectedXamarinProject.Name == "" {
				return printXamarinScanFinishedWithError(
					`Invalid Project specified (%s), either not found in the provided Solution or it can't be used for iOS "Archive for Publishing".`,
					paramXamarinProjectName)
			}
		} else {
			// no project CLI param specified
			if len(xamarinProjectsToChooseFrom) == 1 {
				selectedXamarinProject = xamarinProjectsToChooseFrom[0]
			} else {
				projectNames := []string{}
				for _, aProj := range xamarinProjectsToChooseFrom {
					projectNames = append(projectNames, aProj.Name)
				}
				sort.Strings(projectNames)

				fmt.Println()
				answerValue, err := goinp.SelectFromStringsWithDefault(`Select the Project Name you use for "Archive for Publishing" (usually ends with ".iOS", e.g.: MyProject.iOS)?`, 1, projectNames)
				if err != nil {
					return printXamarinScanFinishedWithError("Failed to select Project: %s", err)
				}
				log.Debugf("selected project: %v", answerValue)
				for _, aProj := range xamarinProjectsToChooseFrom {
					if answerValue == aProj.Name {
						selectedXamarinProject = aProj
						break
					}
				}
			}
		}
	}
	xamarinCmd.ProjectName = selectedXamarinProject.Name
	log.Debugf("xamarinCmd.ProjectName: %s", xamarinCmd.ProjectName)
	log.Debugf("selectedXamarinProject.Configs: %#v", selectedXamarinProject.Configs)

	// Xamarin Configuration Name
	selectedXamarinConfigurationName := ""
	{
		acceptableConfigs := []string{}
		for configName, aConfig := range selectedXamarinProject.Configs {
			if aConfig.Platform == "iPhone" {
				if aConfig.Configuration == "Release" {
					// ios & tvOS app
					acceptableConfigs = append(acceptableConfigs, configName)
				}
			} else if aConfig.Platform == "x86" {
				if aConfig.Configuration == "Release" || aConfig.Configuration == "Debug" {
					// MacOS app
					acceptableConfigs = append(acceptableConfigs, configName)
				}
			}
		}
		if len(acceptableConfigs) < 1 {
			return printXamarinScanFinishedWithError(
				`No acceptable Configuration found in the provided Solution and Project, or none can be used for iOS "Archive for Publishing".`,
			)
		}

		if paramXamarinConfigurationName != "" {
			// configuration specified via flag/param
			for _, aConfigName := range acceptableConfigs {
				if paramXamarinConfigurationName == aConfigName {
					selectedXamarinConfigurationName = aConfigName
					break
				}
			}
			if selectedXamarinConfigurationName == "" {
				return printXamarinScanFinishedWithError(
					"Invalid Configuration specified (%s), either not found in the provided Solution and Project or it can't be used for iOS Archive.",
					paramXamarinConfigurationName)
			}
		} else {
			// no configuration CLI param specified
			if len(acceptableConfigs) == 1 {
				selectedXamarinConfigurationName = acceptableConfigs[0]
			} else {
				fmt.Println()
				answerValue, err := goinp.SelectFromStrings(
					`Select the Configuration Name you use for "Archive for Publishing" (usually Release|iPhone)?`,
					acceptableConfigs,
				)
				if err != nil {
					return printXamarinScanFinishedWithError("Failed to select Configuration: %s", err)
				}
				log.Debugf("selected configuration: %v", answerValue)
				selectedXamarinConfigurationName = answerValue
			}
		}
	}
	if selectedXamarinConfigurationName == "" {
		return printXamarinScanFinishedWithError(
			`No acceptable Configuration found (it was empty) in the provided Solution and Project, or none can be used for iOS "Archive for Publishing".`,
		)
	}
	if err := xamarinCmd.SetConfigurationPlatformCombination(selectedXamarinConfigurationName); err != nil {
		return printXamarinScanFinishedWithError("Failed to set Configuration Platform combination for the command, error: %s", err)
	}

	fmt.Println()
	fmt.Println()
	log.Printf(`🔦  Running a Build, to get all the required code signing settings...`)
	archivePath, xamLogOut, err := xamarinCmd.GenerateArchive()
	logOutput := xamLogOut
	// save the xamarin output into a debug log file
	logOutputFilePath := filepath.Join(absExportOutputDirPath, "xamarin-build-output.log")
	{
		log.Infof("💡  "+colorstring.Yellow("Saving xamarin output into file")+": %s", logOutputFilePath)
		if logWriteErr := fileutil.WriteStringToFile(logOutputFilePath, logOutput); logWriteErr != nil {
			log.Errorf("Failed to save xamarin build output into file (%s), error: %s", logOutputFilePath, logWriteErr)
		} else if err != nil {
			log.Warnf("Please check the logfile (%s) to see what caused the error", logOutputFilePath)
			log.Warnf(`and make sure that you can "Archive for Publishing" this project from Xamarin!`)
			fmt.Println()
			log.Infof("Open the project: %s", xamarinCmd.SolutionFilePath)
			log.Infof(`And do "Archive for Publishing", after selecting the Configuration+Platform: %s|%s`, xamarinCmd.Configuration, xamarinCmd.Platform)
			fmt.Println()
		}
	}
	if err != nil {
		return printXamarinScanFinishedWithError("Failed to run xamarin build command: %s", err)
	}

	// archive code sign settings
	installedCertificates, err := certificateutil.InstalledCodesigningCertificateInfos()
	if err != nil {
		return printXamarinScanFinishedWithError("Failed to list installed code signing identities, error: %s", err)
	}
	installedCertificates = certificateutil.FilterValidCertificateInfos(installedCertificates)

	log.Debugf("Installed certificates:")
	for _, installedCertificate := range installedCertificates {
		log.Debugf(installedCertificate.String())
	}

	installedProfiles, err := profileutil.InstalledProvisioningProfileInfos(profileutil.ProfileTypeIos)
	if err != nil {
		return err
	}

	log.Debugf("Installed profiles:")
	for _, profileInfo := range installedProfiles {
		log.Debugf(profileInfo.String(installedCertificates...))
	}

	archive, err := xcarchive.NewIosArchive(archivePath)
	if err != nil {
		return printXamarinScanFinishedWithError("Failed to analyze archive, error: %s", err)
	}

	archiveCodeSignGroup, err := analyzeArchive(archive, installedCertificates)
	if err != nil {
		return printXamarinScanFinishedWithError("Failed to analyze the archive, error: %s", err)
	}

	fmt.Println()
	log.Infof("Codesign settings used for Xamarin archive:")
	fmt.Println()
	printCodesignGroup(archiveCodeSignGroup)

	// ipa export code sign settings
	fmt.Println()
	fmt.Println()
	log.Printf("🔦  Analyzing the archive, to get ipa export code signing settings...")

	certificatesToExport := []certificateutil.CertificateInfoModel{}
	profilesToExport := []profileutil.ProvisioningProfileInfoModel{}

	if certificatsOnly {
		ipaExportCertificate, err := collectIpaExportCertificate(archiveCodeSignGroup.Certificate, installedCertificates)
		if err != nil {
			return err
		}

		certificatesToExport = append(certificatesToExport, archiveCodeSignGroup.Certificate, ipaExportCertificate)
	} else {
		ipaExportCodeSignGroups, err := collectIpaExportCodeSignGroups(archive, installedCertificates, installedProfiles)
		if err != nil {
			return printXcodeScanFinishedWithError("Failed to collect ipa export code sign groups, error: %s", err)
		}

		codeSignGroups := append(ipaExportCodeSignGroups, archiveCodeSignGroup)
		certificates, profiles := extractCertificatesAndProfiles(codeSignGroups...)

		certificatesToExport = append(certificatesToExport, certificates...)
		profilesToExport = append(profilesToExport, profiles...)
	}

	if err := collectAndExportIdentities(certificatesToExport, absExportOutputDirPath); err != nil {
		return printXcodeScanFinishedWithError("Failed to export codesign identities, error: %s", err)
	}

	if err := collectAndExportProvisioningProfiles(profilesToExport, absExportOutputDirPath); err != nil {
		return printXcodeScanFinishedWithError("Failed to export provisioning profiles, error: %s", err)
	}

	printFinished()

	return nil
}
