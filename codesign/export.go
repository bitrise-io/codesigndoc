package codesign

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/codesigndoc/bitriseio"
	"github.com/bitrise-io/codesigndoc/bitriseio/bitrise"
	"github.com/bitrise-io/codesigndoc/models"
	"github.com/bitrise-io/codesigndoc/osxkeychain"
	"github.com/bitrise-io/codesigndoc/utility"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/profileutil"
	"github.com/bitrise-io/goinp/goinp"
)

// UploadConfig contains configuration to automatically upload artifacts to bitrise.io
type UploadConfig struct {
	PersonalAccessToken string
	AppSlug             string
}

func (config *UploadConfig) isValid() bool {
	return (strings.TrimSpace(config.PersonalAccessToken) != "") &&
		(strings.TrimSpace(config.AppSlug) != "")
}

// UploadAndWriteCodesignFiles exports then uploads codesign files to bitrise.io and saves them to output folder
func UploadAndWriteCodesignFiles(certificates []certificateutil.CertificateInfoModel, profiles []profileutil.ProvisioningProfileInfoModel, askForPassword bool, outputDirPath string, uploadConfig UploadConfig) (bool, bool, error) {
	identities, err := CollectAndExportIdentitiesAsReader(certificates, askForPassword)
	if err != nil {
		return false, false, err
	}

	provisioningProfiles, err := CollectAndExportProvisioningProfilesAsReader(profiles)
	if err != nil {
		return false, false, err
	}

	var client *bitrise.Client
	if uploadConfig.isValid() {
		// Upload automatically if token is provided as CLI paramter, do not export to filesystem
		// Used to upload artifacts as part of an other CLI tool
		client, err = bitrise.NewClientAsStream(uploadConfig.PersonalAccessToken)
		if err != nil {
			return false, false, err
		}
		client.SetSelectedAppSlug(uploadConfig.AppSlug)
	}

	if client == nil {
		uploadConfirmMsg := "Do you want to upload the provisioning profiles and certificates to Bitrise?"
		if len(provisioningProfiles) == 0 {
			uploadConfirmMsg = "Do you want to upload the certificates to Bitrise?"
		}
		fmt.Println()
		if shouldUpload, err := goinp.AskForBoolFromReader(uploadConfirmMsg, os.Stdin); err != nil {
			return false, false, err
		} else if shouldUpload {
			client, err = bitriseio.GetInteractiveConfigClient()
		}
	}

	provProfilesUploaded := (len(profiles) == 0)
	certsUploaded := (len(certificates) == 0)
	if client != nil {
		certsUploaded, provProfilesUploaded, err = bitriseio.UploadCodesigningFilesAsStream(client, identities, provisioningProfiles)
		if err != nil {
			return false, false, err
		}
	}

	if strings.TrimSpace(outputDirPath) != "" {
		if err := WriteIdentities(identities.Content, outputDirPath); err != nil {
			return false, false, err
		}
		if err := WriteProvisioningProfilesAsStream(provisioningProfiles, outputDirPath); err != nil {
			return false, false, err
		}
		fmt.Println()
		log.Successf("Exports finished you can find the exported files at: %s", outputDirPath)

		if err := command.RunCommand("open", outputDirPath); err != nil {
			log.Errorf("Failed to open the export directory in Finder: %s", outputDirPath)
		} else {
			fmt.Println("Opened the directory in Finder.")
		}
	}

	return certsUploaded, provProfilesUploaded, nil
}

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
func CollectAndExportIdentitiesAsReader(certificates []certificateutil.CertificateInfoModel, isAskForPassword bool) (models.Certificates, error) {
	if len(certificates) == 0 {
		return models.Certificates{}, nil
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
			return models.Certificates{}, fmt.Errorf("failed to export, error: %s", err)
		}

		if identityRef == nil {
			return models.Certificates{}, errors.New("identity not found in the keychain, or it was invalid (expired)")
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
		return models.Certificates{}, fmt.Errorf("failed to export from Keychain: %s", err)
	}
	return models.Certificates{
		Info:    certificates,
		Content: identities,
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
func CollectAndExportProvisioningProfilesAsReader(profiles []profileutil.ProvisioningProfileInfoModel) ([]models.ProvisioningProfile, error) {
	if len(profiles) == 0 {
		return nil, nil
	}

	log.Infof("Required Provisioning Profiles (%d)", len(profiles))
	for _, profile := range profiles {
		log.Printf("- %s (UUID: %s)", profile.Name, profile.UUID)
	}

	fmt.Println()
	log.Infof("Exporting Provisioning Profiles...")

	var exportedProfiles []models.ProvisioningProfile
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
		if bytes.Compare(profile.Content(), exportedProfile.Content()) != 0 {
			return nil, fmt.Errorf("Profile found in the archive does not match found profile")
		}

		contents, err := ioutil.ReadFile(pth)
		if err != nil {
			return nil, fmt.Errorf("could not read provisioning profile file, error: %s", err)
		}

		exportedProfiles = append(exportedProfiles, models.ProvisioningProfile{
			Info:    exportedProfile,
			Content: contents,
		})
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

// WriteProvisioningProfilesAsStream writes provisioning profiles to the filesystem
func WriteProvisioningProfilesAsStream(profiles []models.ProvisioningProfile, absExportOutputDirPath string) error {
	fmt.Println()
	log.Infof("Exporting Provisioning Profiles...")

	for _, profile := range profiles {
		exportFileName := utility.ProfileExportFileNameNoPath(profile.Info)
		exportPth := filepath.Join(absExportOutputDirPath, exportFileName)
		if err := ioutil.WriteFile(exportPth, profile.Content, 0600); err != nil {
			return fmt.Errorf("failed to write file, error: %s", err)
		}
	}
	return nil
}
