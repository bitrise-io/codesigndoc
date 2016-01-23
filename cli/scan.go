package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/codesigndoc/osxkeychain"
	"github.com/bitrise-tools/codesigndoc/utils"
	"github.com/bitrise-tools/codesigndoc/xcode"
	"github.com/codegangsta/cli"
)

func scan(c *cli.Context) {
	projectPath := c.String(FileParamKey)
	if projectPath == "" {
		askText := `Please drag-and-drop your Xcode Project (` + colorstring.Green(".xcodeproj") + `)
   or Workspace (` + colorstring.Green(".xcworkspace") + `) file, the one you usually open in Xcode,
   then hit Enter.

  (Note: if you have a Workspace file you should most likely use that)`
		fmt.Println()
		projpth, err := goinp.AskForString(askText)
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
	fmt.Println("========================================")

	fmt.Println()
	utils.Printlnf("=== Required Provisioning Profiles (%d) ===", len(codeSigningSettings.ProvProfiles))
	for idx, aProvProfile := range codeSigningSettings.ProvProfiles {
		utils.Printlnf(" * (%d): %s (UUID: %s)", idx+1, aProvProfile.Title, aProvProfile.UUID)
	}
	fmt.Println("======================================")

	if len(codeSigningSettings.Identities) < 1 {
		log.Fatal("No Code Signing Identity detected!")
	}
	if len(codeSigningSettings.Identities) > 1 {
		log.Warning("More than one Code Signing Identity (certificate) is required to sign your app!")
		log.Warning("You should check your settings and make sure a single Identity/Certificate can be used")
		log.Warning(" for Archiving your app!")
	}

	identityExportRefs := osxkeychain.CreateEmptyCFTypeRefSlice()
	defer osxkeychain.ReleaseRefList(identityExportRefs)

	for _, aIdentity := range codeSigningSettings.Identities {
		identityRefs, err := osxkeychain.FindIdentity(aIdentity.Title)
		if err != nil {
			log.Fatalf("Failed to Export Identity: %s", err)
		}
		log.Printf("identityRefs: %d", len(identityRefs))
		if len(identityRefs) < 1 {
			log.Fatalf("No Identity found in Keychain!")
		}
		if len(identityRefs) > 1 {
			log.Fatalf("Multiple matching Identities found in Keychain! Most likely you have duplicate identity in separate Keychains, like one in System.keychain and one in your Login.keychain")
		}
		identityExportRefs = append(identityExportRefs, identityRefs...)
	}

	if err := osxkeychain.ExportFromKeychain(identityExportRefs, "./Identities.p12"); err != nil {
		log.Fatalf("Failed to export from Keychain: %s", err)
	}
}
