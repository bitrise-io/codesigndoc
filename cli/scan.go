package cli

import (
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/certutil"
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
	log.Println("🔦  Running an Xcode Archive, to get all the required code signing settings...")
	codeSigningSettings, err := xcodeCmd.ScanCodeSigningSettings()
	if err != nil {
		log.Fatalf("Failed to detect code signing settings: %s", err)
	}
	log.Debugf("codeSigningSettings: %#v", codeSigningSettings)

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
		isShouldExport, err := goinp.AskForBool("Do you want to export these files?")
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
	log.Println("Exporting the required Identities (Certificates) ...")
	fmt.Println(" You'll most likely see popups (one for each Identity) from Keychain,")
	fmt.Println(" you will have to accept (Allow) those to be able to export the Identity")
	fmt.Println()

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
		log.Debugf("Export output dir already exists at path: %s", absExportOutputDirPath)
	}

	identityExportRefs := osxkeychain.CreateEmptyCFTypeRefSlice()
	defer osxkeychain.ReleaseRefList(identityExportRefs)

	fmt.Println()
	for _, aIdentity := range codeSigningSettings.Identities {
		log.Infof(" * "+colorstring.Blue("Exporting Identity")+": %s", aIdentity.Title)
		foundIdentityRefs, err := osxkeychain.FindIdentity(aIdentity.Title)
		if err != nil {
			log.Fatalf("Failed to Export Identity: %s", err)
		}
		log.Debugf("foundIdentityRefs: %d", len(foundIdentityRefs))
		if len(foundIdentityRefs) < 1 {
			log.Fatalf("No Identity found in Keychain!")
		}

		// check validity
		validIdentityRefs := osxkeychain.CreateEmptyCFTypeRefSlice()
		for _, aIdentityRef := range foundIdentityRefs {
			cert, err := osxkeychain.GetCertificateDataFromIdentityRef(aIdentityRef)
			if err != nil {
				log.Fatalf("Failed to read certificate data: %s", err)
			}

			if err := certutil.CheckCertificateValidity(cert); err != nil {
				log.Warning(colorstring.Yellowf("Certificate is not valid, skipping: %s", err))
				continue
			}

			validIdentityRefs = append(validIdentityRefs, aIdentityRef)
		}

		if len(validIdentityRefs) < 1 {
			log.Fatalf("Identity found found in Keychain, but no Valid identity found!")
		}
		if len(validIdentityRefs) > 1 {
			log.Warning(colorstring.Yellow("Multiple matching Identities found in Keychain! Most likely you have duplicated identities in separate Keychains, like one in System.keychain and one in your Login.keychain, or you have revoked versions of the Certificate."))
		}
		identityExportRefs = append(identityExportRefs, validIdentityRefs...)
	}

	if err := osxkeychain.ExportFromKeychain(identityExportRefs, path.Join(absExportOutputDirPath, "Identities.p12")); err != nil {
		log.Fatalf("Failed to export from Keychain: %s", err)
	}

	fmt.Println()
	for _, aProvProfile := range codeSigningSettings.ProvProfiles {
		log.Infof(" * "+colorstring.Blue("Exporting Provisioning Profile")+": %s (UUID: %s)", aProvProfile.Title, aProvProfile.UUID)
		filePth, err := provprofile.FindProvProfileFile(aProvProfile)
		if err != nil {
			log.Fatalf("Failed to find Provisioning Profile: %s", err)
		}
		log.Infof("   File found at: %s", filePth)

		if err := cmdex.RunCommand("cp", filePth, absExportOutputDirPath+"/"); err != nil {
			log.Fatalf("Failed to copy the Provisioning Profile into the export directory: %s", err)
		}

		// if err := provprofile.PrintFileInfo(filePth); err != nil {
		// 	log.Fatalf("Err: %s", err)
		// }
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
