package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/osxkeychain"
	"github.com/bitrise-tools/codesigndoc/provprofile"
	"github.com/bitrise-tools/codesigndoc/utils"
	"github.com/bitrise-tools/codesigndoc/xcode"
	"github.com/spf13/cobra"
)

// xcodeCmd represents the xcode command
var xcodeCmd = &cobra.Command{
	Use:   "xcode",
	Short: "Xcode project scanner",
	Long:  `Scan an Xcode project`,

	Run: scanXcodeProject,
}

var (
	xcodeProjectFilePath = ""
	xcodeScheme          = ""
)

func init() {
	scanCmd.AddCommand(xcodeCmd)

	xcodeCmd.Flags().StringVar(&xcodeProjectFilePath,
		"file", "",
		"Xcode Project/Workspace file")
	xcodeCmd.Flags().StringVar(&xcodeScheme,
		"scheme", "",
		"Xcode Scheme")
}

const (
	confExportOutputDirPath = "./codesigndoc_exports"
)

func printFinished() {
	fmt.Println()
	fmt.Println(colorstring.Green("That's all."))
	fmt.Println("You just have to upload the found code signing files and you'll be good to go!")
}

func failWithError(format string, args ...interface{}) {
	log.Errorf(colorstring.Red("Error: ")+format, args...)
	fmt.Println()
	fmt.Println("------------------------------")
	fmt.Println("First of all " + colorstring.Red("please make sure that you can Archive your app from Xcode."))
	fmt.Println("codesigndoc only works if you can archive your app from Xcode.")
	fmt.Println("If you can, and you get a valid IPA file if you export from Xcode,")
	fmt.Println(colorstring.Red("please create an issue") + " on GitHub at: https://github.com/bitrise-tools/codesigndoc/issues")
	fmt.Println("with as many details & logs as you can share!")
	fmt.Println("------------------------------")
	fmt.Println()
	os.Exit(1)
}

func scanXcodeProject(cmd *cobra.Command, args []string) {
	absExportOutputDirPath, err := pathutil.AbsPath(confExportOutputDirPath)
	log.Debugf("absExportOutputDirPath: %s", absExportOutputDirPath)
	if err != nil {
		failWithError("Failed to determin Absolute path of export dir: %s", confExportOutputDirPath)
	}
	if exist, err := pathutil.IsDirExists(absExportOutputDirPath); err != nil {
		failWithError("Failed to determin whether the export directory already exists: %s", err)
	} else if !exist {
		if err := os.Mkdir(absExportOutputDirPath, 0777); err != nil {
			failWithError("Failed to create export output directory at path: %s | error: %s", absExportOutputDirPath, err)
		}
	} else {
		log.Infof("Export output dir already exists at path: %s", absExportOutputDirPath)
	}

	projectPath := xcodeProjectFilePath
	if projectPath == "" {
		askText := `Please drag-and-drop your Xcode Project (` + colorstring.Green(".xcodeproj") + `)
   or Workspace (` + colorstring.Green(".xcworkspace") + `) file, the one you usually open in Xcode,
   then hit Enter.

  (Note: if you have a Workspace file you should most likely use that)`
		fmt.Println()
		projpth, err := goinp.AskForPath(askText)
		if err != nil {
			failWithError("Failed to read input: %s", err)
		}
		projectPath = projpth
	}
	log.Debugf("projectPath: %s", projectPath)
	xcodeCmd := xcode.CommandModel{
		ProjectFilePath: projectPath,
	}

	schemeToUse := xcodeScheme
	if schemeToUse == "" {
		log.Println("🔦  Scanning Schemes ...")
		schemes, err := xcodeCmd.ScanSchemes()
		if err != nil {
			failWithError("Failed to scan Schemes: %s", err)
		}
		log.Debugf("schemes: %v", schemes)

		fmt.Println()
		selectedScheme, err := goinp.SelectFromStrings("Select the Scheme you usually use in Xcode", schemes)
		if err != nil {
			failWithError("Failed to select Scheme: %s", err)
		}
		log.Debugf("selected scheme: %v", selectedScheme)
		schemeToUse = selectedScheme
	}
	xcodeCmd.Scheme = schemeToUse

	fmt.Println()
	fmt.Println()
	log.Println("🔦  Running an Xcode Archive, to get all the required code signing settings...")
	codeSigningSettings, xcodebuildOutput, err := xcodeCmd.ScanCodeSigningSettings()
	// save the xcodebuild output into a debug log file
	xcodebuildOutputFilePath := filepath.Join(absExportOutputDirPath, "xcodebuild-output.log")
	{
		log.Infof("  💡  "+colorstring.Yellow("Saving xcodebuild output into file")+": %s", xcodebuildOutputFilePath)
		if err := fileutil.WriteStringToFile(xcodebuildOutputFilePath, xcodebuildOutput); err != nil {
			log.Errorf("Failed to save xcodebuild output into file (%s), error: %s", xcodebuildOutputFilePath, err)
		}
		log.Infoln(colorstring.Yellow("Please check the logfile (" + xcodebuildOutputFilePath + ") to see what caused the error"))
		log.Infoln(colorstring.Red("and make sure that you can Archive this project from Xcode!"))
		fmt.Println()
		log.Infoln("Open the project: ", xcodeCmd.ProjectFilePath)
		log.Infoln("And Archive, using the Scheme: ", xcodeCmd.Scheme)
		fmt.Println()
	}
	if err != nil {
		failWithError("Failed to detect code signing settings: %s", err)
	}
	log.Debugf("codeSigningSettings: %#v", codeSigningSettings)

	fmt.Println()
	fmt.Println()
	utils.Printlnf("=== Required Identities/Certificates (%d) ===", len(codeSigningSettings.Identities))
	for idx, anIdentity := range codeSigningSettings.Identities {
		utils.Printlnf(" * (%d): %s", idx+1, anIdentity.Title)
	}
	fmt.Println("============================================")

	fmt.Println()
	utils.Printlnf("=== Required Provisioning Profiles (%d) ===", len(codeSigningSettings.ProvProfiles))
	for idx, aProvProfile := range codeSigningSettings.ProvProfiles {
		utils.Printlnf(" * (%d): %s (UUID: %s)", idx+1, aProvProfile.Title, aProvProfile.UUID)
	}
	fmt.Println("==========================================")

	fmt.Println()
	utils.Printlnf("=== Team IDs (%d) ===", len(codeSigningSettings.TeamIDs))
	for idx, aTeamID := range codeSigningSettings.TeamIDs {
		utils.Printlnf(" * (%d): %s", idx+1, aTeamID)
	}
	fmt.Println("==========================================")

	fmt.Println()
	utils.Printlnf("=== App/Bundle IDs (%d) ===", len(codeSigningSettings.AppBundleIDs))
	for idx, anAppBundleID := range codeSigningSettings.AppBundleIDs {
		utils.Printlnf(" * (%d): %s", idx+1, anAppBundleID)
	}
	fmt.Println("==========================================")
	fmt.Println()

	//
	// --- Code Signing issue checks / report
	//

	if len(codeSigningSettings.Identities) < 1 {
		failWithError("No Code Signing Identity detected!")
	}
	if len(codeSigningSettings.Identities) > 1 {
		log.Warning(colorstring.Yellow("More than one Code Signing Identity (certificate) is required to sign your app!"))
		log.Warning("You should check your settings and make sure a single Identity/Certificate can be used")
		log.Warning(" for Archiving your app!")
	}

	if len(codeSigningSettings.ProvProfiles) < 1 {
		failWithError("No Provisioning Profiles detected!")
	}

	//
	// --- Export
	//

	if !isAllowExport {
		isShouldExport, err := goinp.AskForBoolWithDefault("Do you want to export these files?", true)
		if err != nil {
			failWithError("Failed to process your input: %s", err)
		}
		if !isShouldExport {
			printFinished()
			return
		}
	} else {
		log.Debug("Allow Export flag was set - doing export without asking")
	}

	fmt.Println()
	log.Println("Collecting the required Identities (Certificates) for a base Xcode Archive ...")
	fmt.Println()

	identitiesWithKeychainRefs := []osxkeychain.IdentityWithRefModel{}
	defer osxkeychain.ReleaseIdentityWithRefList(identitiesWithKeychainRefs)

	for _, aIdentity := range codeSigningSettings.Identities {
		log.Infof(" * "+colorstring.Blue("Searching for Identity")+": %s", aIdentity.Title)
		validIdentityRefs, err := osxkeychain.FindAndValidateIdentity(aIdentity.Title, true)
		if err != nil {
			failWithError("Failed to export, error: %s", err)
		}

		if len(validIdentityRefs) < 1 {
			failWithError("Identity not found in the keychain, or it was invalid (expired)!")
		}
		if len(validIdentityRefs) > 1 {
			log.Warning(colorstring.Yellow("Multiple matching Identities found in Keychain! Most likely you have duplicated identities in separate Keychains, e.g. one in System.keychain and one in your Login.keychain, or you have revoked versions of the Certificate."))
		}

		identitiesWithKeychainRefs = append(identitiesWithKeychainRefs, validIdentityRefs...)
	}

	fmt.Println()
	log.Println("Collecting additional identities, for Distribution builds ...")
	fmt.Println()

	for _, aTeamID := range codeSigningSettings.TeamIDs {
		log.Infof(" * "+colorstring.Blue("Searching for Identities with Team ID")+": %s", aTeamID)
		validIdentityRefs, err := osxkeychain.FindAndValidateIdentity(fmt.Sprintf("(%s)", aTeamID), false)
		if err != nil {
			failWithError("Failed to export, error: %s", err)
		}

		if len(validIdentityRefs) < 1 {
			log.Infoln("No valid identity found for this Team ID")
		}

		identitiesWithKeychainRefs = append(identitiesWithKeychainRefs, validIdentityRefs...)
	}

	fmt.Println()
	log.Println(colorstring.Green("Exporting the Identities") + " (Certificates):")
	fmt.Println()

	identityKechainRefs := osxkeychain.CreateEmptyCFTypeRefSlice()
	for _, aIdentityWithRefItm := range identitiesWithKeychainRefs {
		fmt.Println(" * "+colorstring.Blue("Identity")+":", aIdentityWithRefItm.Label)
		identityKechainRefs = append(identityKechainRefs, aIdentityWithRefItm.KeychainRef)
	}

	fmt.Println()
	if isAskForPassword {
		log.Infoln(colorstring.Blue("Exporting from Keychain"))
		log.Infoln(colorstring.Yellow(" You'll be asked to provide a Passphrase for the .p12 file!"))
	} else {
		log.Infoln(colorstring.Blue("Exporting from Keychain") + ", " + colorstring.Yellow("using empty Passphrase") + " ...")
		log.Info(" This means that " + colorstring.Yellow("if you want to import the file the passphrase at import should be left empty") + ",")
		log.Info(" you don't have to type in anything, just leave the passphrase input empty.")
	}
	fmt.Println()
	log.Info(colorstring.Blue("You'll most likely see popups") + " (one for each Identity) from Keychain,")
	log.Info(colorstring.Yellow(" you will have to accept (Allow)") + " those to be able to export the Identities!")
	fmt.Println()

	if err := osxkeychain.ExportFromKeychain(identityKechainRefs, filepath.Join(absExportOutputDirPath, "Identities.p12"), isAskForPassword); err != nil {
		failWithError("Failed to export from Keychain: %s", err)
	}

	fmt.Println()
	log.Println(colorstring.Green("Searching for required Provisioning Profiles"), "...")
	fmt.Println()

	provProfileUUIDLookupMap := map[string]provprofile.ProvisioningProfileFileInfoModel{}
	for _, aProvProfile := range codeSigningSettings.ProvProfiles {
		log.Infof(" * "+colorstring.Blue("Searching for required Provisioning Profile")+": %s (UUID: %s)", aProvProfile.Title, aProvProfile.UUID)
		provProfileFileInfo, err := provprofile.FindProvProfileByUUID(aProvProfile.UUID)
		if err != nil {
			failWithError("Failed to find Provisioning Profile: %s", err)
		}
		log.Infof("   File found at: %s", provProfileFileInfo.Path)

		provProfileUUIDLookupMap[provProfileFileInfo.ProvisioningProfileInfo.UUID] = provProfileFileInfo
	}

	fmt.Println()
	log.Println(colorstring.Green("Searching for additinal, Distribution Provisioning Profiles"), "...")
	fmt.Println()
	for _, aAppBundleID := range codeSigningSettings.AppBundleIDs {
		log.Infof(" * "+colorstring.Blue("Searching for Provisioning Profiles with App ID")+": %s", aAppBundleID)
		provProfileFileInfos, err := provprofile.FindProvProfilesByAppID(aAppBundleID)
		if err != nil {
			failWithError("Error during Provisioning Profile search: %s", err)
		}
		if len(provProfileFileInfos) < 1 {
			log.Warn("   No Provisioning Profile found for this Bundle ID")
			continue
		}
		log.Infof("   Found matching Provisioning Profile count: %d", len(provProfileFileInfos))

		for _, aProvProfileFileInfo := range provProfileFileInfos {
			provProfileUUIDLookupMap[aProvProfileFileInfo.ProvisioningProfileInfo.UUID] = aProvProfileFileInfo
		}
	}

	fmt.Println()
	log.Println(colorstring.Green("Exporting Provisioning Profiles"), "...")
	fmt.Println()
	provProfileFileInfos := []provprofile.ProvisioningProfileFileInfoModel{}
	for _, aProvProfFileInfo := range provProfileUUIDLookupMap {
		provProfileFileInfos = append(provProfileFileInfos, aProvProfFileInfo)
	}
	if err := exportProvisioningProfiles(provProfileFileInfos, absExportOutputDirPath); err != nil {
		failWithError("Failed to export the Provisioning Profile into the export directory: %s", err)
	}

	fmt.Println()
	fmt.Printf(colorstring.Green("Exports finished")+" you can find the exported files at: %s\n", absExportOutputDirPath)
	if err := cmdex.RunCommand("open", absExportOutputDirPath); err != nil {
		log.Errorf("Failed to open the export directory in Finder: %s", absExportOutputDirPath)
	}
	fmt.Println("Opened the directory in Finder.")
	fmt.Println()

	printFinished()
}

func exportProvisioningProfiles(provProfileFileInfos []provprofile.ProvisioningProfileFileInfoModel,
	exportTargetDirPath string) error {

	for _, aProvProfileFileInfo := range provProfileFileInfos {
		log.Infoln("   " + colorstring.Green("Exporting Provisioning Profile:") + " " + aProvProfileFileInfo.ProvisioningProfileInfo.Name)
		log.Infoln("                             UUID: " + aProvProfileFileInfo.ProvisioningProfileInfo.UUID)
		exportFileName := provProfileExportFileName(aProvProfileFileInfo)
		exportPth := filepath.Join(exportTargetDirPath, exportFileName)
		if err := cmdex.RunCommand("cp", aProvProfileFileInfo.Path, exportPth); err != nil {
			return fmt.Errorf("Failed to copy Provisioning Profile (from: %s) (to: %s), error: %s",
				aProvProfileFileInfo.Path, exportPth, err)
		}
	}
	return nil
}

func provProfileExportFileName(provProfileFileInfo provprofile.ProvisioningProfileFileInfoModel) string {
	replaceRexp, err := regexp.Compile("[^A-Za-z0-9_.-]")
	if err != nil {
		log.Warn("Invalid regex, error: %s", err)
		return ""
	}
	safeTitle := replaceRexp.ReplaceAllString(provProfileFileInfo.ProvisioningProfileInfo.Name, "")
	extension := ".mobileprovision"
	if strings.HasSuffix(provProfileFileInfo.Path, ".provisionprofile") {
		extension = ".provisionprofile"
	}

	return provProfileFileInfo.ProvisioningProfileInfo.UUID + "." + safeTitle + extension
}
