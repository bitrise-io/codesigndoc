package cli

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
	"github.com/codegangsta/cli"
)

const (
	confExportOutputDirPath = "./codesigndoc_exports"
)

func printFinished() {
	fmt.Println()
	fmt.Println(colorstring.Green("That's all."))
	fmt.Println("You just have to upload the found code signing files and you'll be good to go!")
}

func scan(c *cli.Context) {
	absExportOutputDirPath, err := pathutil.AbsPath(confExportOutputDirPath)
	log.Debugf("absExportOutputDirPath: %s", absExportOutputDirPath)
	if err != nil {
		log.Fatalf("Failed to determin Absolute path of export dir: %s", confExportOutputDirPath)
	}
	if exist, err := pathutil.IsDirExists(absExportOutputDirPath); err != nil {
		log.Fatalf("Failed to determin whether the export directory already exists: %s", err)
	} else if !exist {
		if err := os.Mkdir(absExportOutputDirPath, 0777); err != nil {
			log.Fatalf("Failed to create export output directory at path: %s | error: %s", absExportOutputDirPath, err)
		}
	} else {
		log.Infof("Export output dir already exists at path: %s", absExportOutputDirPath)
	}

	projectPath := c.String(FileParamKey)
	if projectPath == "" {
		askText := `Please drag-and-drop your Xcode Project (` + colorstring.Green(".xcodeproj") + `)
   or Workspace (` + colorstring.Green(".xcworkspace") + `) file, the one you usually open in Xcode,
   then hit Enter.

  (Note: if you have a Workspace file you should most likely use that)`
		fmt.Println()
		projpth, err := goinp.AskForPath(askText)
		if err != nil {
			log.Fatalf("Failed to read input: %s", err)
		}
		projectPath = projpth
	}
	log.Debugf("projectPath: %s", projectPath)
	xcodeCmd := xcode.CommandModel{
		ProjectFilePath: projectPath,
	}

	schemeToUse := c.String(SchemeParamKey)
	if schemeToUse == "" {
		log.Println("🔦  Scanning Schemes ...")
		schemes, err := xcodeCmd.ScanSchemes()
		if err != nil {
			log.Fatalf("Failed to scan Schemes: %s", err)
		}
		log.Debugf("schemes: %v", schemes)

		fmt.Println()
		selectedScheme, err := goinp.SelectFromStrings("Select the Scheme you usually use in Xcode", schemes)
		if err != nil {
			log.Fatalf("Failed to select Scheme: %s", err)
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
	{
		xcodebuildOutputFilePath := filepath.Join(absExportOutputDirPath, "xcodebuild-output.log")
		log.Infof("  💡  Saving xcodebuild output into file: %s", xcodebuildOutputFilePath)
		if err := fileutil.WriteStringToFile(xcodebuildOutputFilePath, xcodebuildOutput); err != nil {
			log.Errorf("Failed to save xcodebuild output into file (%s), error: %s", xcodebuildOutputFilePath, err)
		}
	}
	if err != nil {
		log.Fatalf("Failed to detect code signing settings: %s", err)
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
		log.Fatal("No Code Signing Identity detected!")
	}
	if len(codeSigningSettings.Identities) > 1 {
		log.Warning(colorstring.Yellow("More than one Code Signing Identity (certificate) is required to sign your app!"))
		log.Warning("You should check your settings and make sure a single Identity/Certificate can be used")
		log.Warning(" for Archiving your app!")
	}

	if len(codeSigningSettings.ProvProfiles) < 1 {
		log.Fatal("No Provisioning Profiles detected!")
	}

	//
	// --- Export
	//

	if !c.Bool(AllowExportParamKey) {
		isShouldExport, err := goinp.AskForBoolWithDefault("Do you want to export these files?", true)
		if err != nil {
			log.Fatalf("Failed to process your input: %s", err)
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
			log.Fatalf("Failed to export, error: %s", err)
		}

		if len(validIdentityRefs) < 1 {
			log.Fatalf("Identity not found in the keychain, or it was invalid (expired)!")
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
			log.Fatalf("Failed to export, error: %s", err)
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
	log.Infoln(colorstring.Blue("Exporting from Keychain") + ", " + colorstring.Yellow("using empty Passphrase") + " ...")
	log.Info(" This means that " + colorstring.Yellow("if you want to import the file the passphrase at import should be left empty") + ",")
	log.Info(" you don't have to type in anything, just leave the passphrase input empty.")
	fmt.Println()
	log.Info(colorstring.Blue("You'll most likely see popups") + " (one for each Identity) from Keychain,")
	log.Info(colorstring.Yellow(" you will have to accept (Allow)") + " those to be able to export the Identities!")
	fmt.Println()
	if err := osxkeychain.ExportFromKeychain(identityKechainRefs, filepath.Join(absExportOutputDirPath, "Identities.p12")); err != nil {
		log.Fatalf("Failed to export from Keychain: %s", err)
	}

	fmt.Println()
	log.Println(colorstring.Green("Exporting base Provisioning Profile(s)"), "...")
	fmt.Println()

	for _, aProvProfile := range codeSigningSettings.ProvProfiles {
		log.Infof(" * "+colorstring.Blue("Exporting Provisioning Profile")+": %s (UUID: %s)", aProvProfile.Title, aProvProfile.UUID)
		filePth, err := provprofile.FindProvProfileFileByUUID(aProvProfile.UUID)
		if err != nil {
			log.Fatalf("Failed to find Provisioning Profile: %s", err)
		}
		log.Infof("   File found at: %s", filePth)

		exportFileName := provProfileExportFileName(aProvProfile.UUID, aProvProfile.Title)
		exportPth := filepath.Join(absExportOutputDirPath, exportFileName)
		if err := cmdex.RunCommand("cp", filePth, exportPth); err != nil {
			log.Fatalf("Failed to copy the Provisioning Profile into the export directory: %s", err)
		}
	}

	fmt.Println()
	log.Println(colorstring.Green("Exporting additinal, Distribution Provisioning Profile(s)"), "...")
	fmt.Println()
	for _, aAppBundleID := range codeSigningSettings.AppBundleIDs {
		log.Infof(" * "+colorstring.Blue("Searching for Provisioning Profiles with Bundle ID")+": %s", aAppBundleID)
		filePths, err := provprofile.FindProvProfilesFileByAppID(aAppBundleID)
		if err != nil {
			log.Fatalf("Failed to find Provisioning Profile: %s", err)
		}
		if len(filePths) < 1 {
			log.Warn("   No Provisioning Profile found for this Bundle ID")
			continue
		}

		for _, aFilePth := range filePths {
			log.Info("   " + colorstring.Green("Exporting Provisioning Profile:") + " " + aFilePth)
			exportFileName := provProfileExportFileName(
				strings.TrimSuffix(filepath.Base(aFilePth), ".mobileprovision"),
				aAppBundleID,
			)
			exportPth := filepath.Join(absExportOutputDirPath, exportFileName)
			if err := cmdex.RunCommand("cp", aFilePth, exportPth); err != nil {
				log.Fatalf("Failed to copy the Provisioning Profile into the export directory: %s", err)
			}
		}
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

func provProfileExportFileName(provProfileUUID, title string) string {
	replaceRexp, err := regexp.Compile("[^A-Za-z0-9_.-]")
	if err != nil {
		log.Warn("Invalid regex, error: %s", err)
		return ""
	}
	safeTitle := replaceRexp.ReplaceAllString(title, "")

	return provProfileUUID + "_" + safeTitle + ".mobileprovision"
}
