package codesign

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/bitrise-io/codesigndoc/osxkeychain"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/profileutil"
)

// CollectAndExportIdentities exports the given certificates into the given directory as a single .p12 file
func CollectAndExportIdentities(certificates []certificateutil.CertificateInfoModel, absExportOutputDirPath string, isAskForPassword bool) error {
	if len(certificates) == 0 {
		return nil
	}

	fmt.Println()
	fmt.Println()
	log.Infof("Required Identities/Certificates (%d)", len(certificates))
	for _, certificate := range certificates {
		log.Printf("- %s", certificate.CommonName)
	}

	fmt.Println()
	log.Infof("Exporting the Identities (Certificates):")

	identitiesWithKeychainRefs := []osxkeychain.IdentityWithRefModel{}
	defer osxkeychain.ReleaseIdentityWithRefList(identitiesWithKeychainRefs)

	for _, certificate := range certificates {
		log.Printf("searching for Identity: %s", certificate.CommonName)
		identityRef, err := osxkeychain.FindAndValidateIdentity(certificate.CommonName)
		if err != nil {
			return fmt.Errorf("failed to export, error: %s", err)
		}

		if identityRef == nil {
			return errors.New("identity not found in the keychain, or it was invalid (expired)")
		}

		identitiesWithKeychainRefs = append(identitiesWithKeychainRefs, *identityRef)
	}

	identityKechainRefs := osxkeychain.CreateEmptyCFTypeRefSlice()
	for _, aIdentityWithRefItm := range identitiesWithKeychainRefs {
		fmt.Println("exporting Identity:", aIdentityWithRefItm.Label)
		identityKechainRefs = append(identityKechainRefs, aIdentityWithRefItm.KeychainRef)
	}

	fmt.Println()
	if isAskForPassword {
		log.Infof("Exporting from Keychain")
		log.Warnf(" You'll be asked to provide a Passphrase for the .p12 file!")
	} else {
		log.Warnf("Exporting from Keychain using empty Passphrase...")
		log.Printf("This means that if you want to import the file the passphrase at import should be left empty,")
		log.Printf("you don't have to type in anything, just leave the passphrase input empty.")
	}
	fmt.Println()
	log.Warnf("You'll most likely see popups one for each Identity from Keychain,")
	log.Warnf("you will have to accept (Allow) those to be able to export the Identities!")
	fmt.Println()

	if err := osxkeychain.ExportFromKeychain(identityKechainRefs, filepath.Join(absExportOutputDirPath, "Identities.p12"), isAskForPassword); err != nil {
		return fmt.Errorf("failed to export from Keychain: %s", err)
	}

	return nil
}

// CollectAndExportIdentitiesAsReader exports the given certificates merged in a single .p12 file, as an io.Reader
func CollectAndExportIdentitiesAsReader(certificates []certificateutil.CertificateInfoModel, isAskForPassword bool) (Certificates, error) {
	if len(certificates) == 0 {
		return Certificates{}, nil
	}

	fmt.Println()
	fmt.Println()
	log.Infof("Required Identities/Certificates (%d)", len(certificates))
	for _, certificate := range certificates {
		log.Printf("- %s", certificate.CommonName)
	}

	fmt.Println()
	log.Infof("Exporting the Identities (Certificates):")

	identitiesWithKeychainRefs := []osxkeychain.IdentityWithRefModel{}
	defer osxkeychain.ReleaseIdentityWithRefList(identitiesWithKeychainRefs)

	for _, certificate := range certificates {
		log.Printf("searching for Identity: %s", certificate.CommonName)
		identityRef, err := osxkeychain.FindAndValidateIdentity(certificate.CommonName)
		if err != nil {
			return Certificates{}, fmt.Errorf("failed to export, error: %s", err)
		}

		if identityRef == nil {
			return Certificates{}, errors.New("identity not found in the keychain, or it was invalid (expired)")
		}

		identitiesWithKeychainRefs = append(identitiesWithKeychainRefs, *identityRef)
	}

	identityKechainRefs := osxkeychain.CreateEmptyCFTypeRefSlice()
	for _, aIdentityWithRefItm := range identitiesWithKeychainRefs {
		fmt.Println("exporting Identity:", aIdentityWithRefItm.Label)
		identityKechainRefs = append(identityKechainRefs, aIdentityWithRefItm.KeychainRef)
	}

	fmt.Println()
	if isAskForPassword {
		log.Infof("Exporting from Keychain")
		log.Warnf(" You'll be asked to provide a Passphrase for the .p12 file!")
	} else {
		log.Warnf("Exporting from Keychain using empty Passphrase...")
		log.Printf("This means that if you want to import the file the passphrase at import should be left empty,")
		log.Printf("you don't have to type in anything, just leave the passphrase input empty.")
	}
	fmt.Println()
	log.Warnf("You'll most likely see popups one for each Identity from Keychain,")
	log.Warnf("you will have to accept (Allow) those to be able to export the Identities!")
	fmt.Println()

	identities, err := osxkeychain.ExportFromKeychainToBuffer(identityKechainRefs, isAskForPassword)
	if err != nil {
		return Certificates{}, fmt.Errorf("failed to export from Keychain: %s", err)
	}
	return Certificates{
		Certificates: certificates,
		Contents:     identities,
	}, nil
}

// WriteIdentities writes identities to a file path
func WriteIdentities(identites []byte, absExportOutputDirPath string) error {
	return ioutil.WriteFile(filepath.Join(absExportOutputDirPath, "Identities.p12"), identites, 0666)
}

// CollectAndExportProvisioningProfiles copies the give profiles into the given directory
func CollectAndExportProvisioningProfiles(profiles []profileutil.ProvisioningProfileInfoModel, absExportOutputDirPath string) error {
	if len(profiles) == 0 {
		return nil
	}

	log.Infof("Required Provisioning Profiles (%d)", len(profiles))
	for _, profile := range profiles {
		log.Printf("- %s (UUID: %s)", profile.Name, profile.UUID)
	}

	fmt.Println()
	log.Infof("Exporting Provisioning Profiles...")

	for _, profile := range profiles {
		log.Printf("searching for required Provisioning Profile: %s (UUID: %s)", profile.Name, profile.UUID)
		_, pth, err := profileutil.FindProvisioningProfileInfo(profile.UUID)
		if err != nil {
			return fmt.Errorf("failed to find Provisioning Profile: %s", err)
		}

		log.Printf("file found at: %s", pth)

		exportFileName := profileExportFileName(profile, pth)
		exportPth := filepath.Join(absExportOutputDirPath, exportFileName)
		if err := command.RunCommand("cp", pth, exportPth); err != nil {
			return fmt.Errorf("Failed to copy Provisioning Profile (from: %s) (to: %s), error: %s", pth, exportPth, err)
		}
	}

	return nil
}

// CollectAndExportProvisioningProfilesAsReader returns provisioning profies as an io.Reader array
func CollectAndExportProvisioningProfilesAsReader(profiles []profileutil.ProvisioningProfileInfoModel) ([]profileutil.ProvisioningProfileInfoModel, error) {
	if len(profiles) == 0 {
		return nil, nil
	}

	log.Infof("Required Provisioning Profiles (%d)", len(profiles))
	for _, profile := range profiles {
		log.Printf("- %s (UUID: %s)", profile.Name, profile.UUID)
	}

	fmt.Println()
	log.Infof("Exporting Provisioning Profiles...")

	var exportedProfiles []profileutil.ProvisioningProfileInfoModel
	for _, profile := range profiles {
		log.Printf("searching for required Provisioning Profile: %s (UUID: %s)", profile.Name, profile.UUID)
		provisioningProfile, pth, err := profileutil.FindProvisioningProfile(profile.UUID)
		if err != nil {
			return nil, fmt.Errorf("failed to find Provisioning Profile: %s", err)
		}
		log.Printf("file found at: %s", pth)
		exportedProfile, err := profileutil.NewProvisioningProfileInfo(*provisioningProfile)
		if err != nil {
			return nil, fmt.Errorf("failed to parse exported profile, error: %s", err)
		}
		exportedProfiles = append(exportedProfiles, exportedProfile)
	}

	return exportedProfiles, nil
}

// WriteProvisioningProfiles writes provisioning profiles to the filesystem
func WriteProvisioningProfiles(profiles []profileutil.ProvisioningProfileInfoModel, absExportOutputDirPath string) error {
	fmt.Println()
	log.Infof("Exporting Provisioning Profiles...")

	for _, profile := range profiles {
		log.Printf("searching for required Provisioning Profile: %s (UUID: %s)", profile.Name, profile.UUID)
		_, pth, err := profileutil.FindProvisioningProfileInfo(profile.UUID)
		if err != nil {
			return fmt.Errorf("failed to find Provisioning Profile: %s", err)
		}

		log.Printf("file found at: %s", pth)

		exportFileName := profileExportFileName(profile, pth)
		exportPth := filepath.Join(absExportOutputDirPath, exportFileName)
		if err := command.RunCommand("cp", pth, exportPth); err != nil {
			return fmt.Errorf("Failed to copy Provisioning Profile (from: %s) (to: %s), error: %s", pth, exportPth, err)
		}
	}
	return nil
}
