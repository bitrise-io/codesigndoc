package codesign

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/bitrise-io/codesigndoc/bitriseio"
	"github.com/bitrise-io/codesigndoc/bitriseio/bitrise"
	"github.com/bitrise-io/codesigndoc/models"
	"github.com/bitrise-io/codesigndoc/osxkeychain"
	"github.com/bitrise-io/codesigndoc/utility"
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

// WriteFilesConfig controls writing artifacts as files
type WriteFilesConfig struct {
	WriteFiles       WriteFilesLevel
	AbsOutputDirPath string
}

// WriteFilesLevel describes if codesigning files should be written to the output directory
type WriteFilesLevel int

const (
	// Invalid represents an invalid value
	Invalid WriteFilesLevel = iota
	// WriteFilesAlways writes build logs and codesigning files always
	WriteFilesAlways
	// WriteFilesFallback writes artifacts when upload was not chosen or failed
	WriteFilesFallback
	// WriteFilesDisabled does not write any files
	WriteFilesDisabled
)

// ExportReport describes the output of codesigning files export
type ExportReport struct {
	CertificatesUploaded         bool
	ProvisioningProfilesUploaded bool
	CodesignFilesWritten         bool
}

// ExportCodesigningFiles exports certificates from the Keychain and provisoining profiles from their directory
func ExportCodesigningFiles(certificatesRequired []certificateutil.CertificateInfoModel, profilesRequired []profileutil.ProvisioningProfileInfoModel, askForPassword bool) (models.Certificates, []models.ProvisioningProfile, error) {
	certificates, err := exportIdentities(certificatesRequired, askForPassword)
	if err != nil {
		return models.Certificates{}, nil, err
	}

	profiles, err := exportProvisioningProfiles(profilesRequired)
	if err != nil {
		return models.Certificates{}, nil, err
	}

	return certificates, profiles, nil
}

// UploadAndWriteCodesignFiles exports then uploads codesign files to bitrise.io and saves them to output folder
func UploadAndWriteCodesignFiles(certificates models.Certificates, provisioningProfiles []models.ProvisioningProfile, writeFilesConfig WriteFilesConfig, uploadConfig UploadConfig) (ExportReport, error) {
	var client *bitrise.Client
	// both or none CLI flags are required
	if uploadConfig.PersonalAccessToken != "" && uploadConfig.AppSlug != "" {
		// Upload automatically if token is provided as CLI paramter, do not export to filesystem
		// Used to upload artifacts as part of an other CLI tool
		var err error
		client, err = bitrise.NewClient(uploadConfig.PersonalAccessToken)
		if err != nil {
			return ExportReport{}, err
		}

		client.SetSelectedAppSlug(uploadConfig.AppSlug)
	}

	if client == nil {
		uploadConfirmMsg := "Do you want to upload the provisioning profiles and certificates to Bitrise?"
		if len(provisioningProfiles) == 0 {
			uploadConfirmMsg = "Do you want to upload the certificates to Bitrise?"
		}
		fmt.Println()

		shouldUpload, err := goinp.AskForBoolFromReader(uploadConfirmMsg, os.Stdin)
		if err != nil {
			return ExportReport{}, err
		}

		if shouldUpload {
			if client, err = bitriseio.GetInteractiveConfigClient(); err != nil {
				return ExportReport{}, err
			}
		}
	}

	var filesWritten bool
	if writeFilesConfig.WriteFiles == WriteFilesAlways ||
		writeFilesConfig.WriteFiles == WriteFilesFallback && client == nil {
		if err := writeFiles(certificates, provisioningProfiles, writeFilesConfig); err != nil {
			return ExportReport{}, err
		}
		filesWritten = true
	}

	if client == nil {
		return ExportReport{
			CertificatesUploaded:         len(certificates.Info) == 0,
			ProvisioningProfilesUploaded: len(provisioningProfiles) == 0,
			CodesignFilesWritten:         filesWritten,
		}, nil
	}

	certificatesUploaded, profilesUploaded, err := bitriseio.UploadCodesigningFiles(client, certificates, provisioningProfiles)
	return ExportReport{
		CertificatesUploaded:         certificatesUploaded,
		ProvisioningProfilesUploaded: profilesUploaded,
		CodesignFilesWritten:         filesWritten,
	}, err
}

func writeFiles(identities models.Certificates, provisioningProfiles []models.ProvisioningProfile, writeFilesConfig WriteFilesConfig) error {
	if err := os.MkdirAll(writeFilesConfig.AbsOutputDirPath, 0700); err != nil {
		return fmt.Errorf("failed to create output directory for codesigning files, error: %s", err)
	}

	entries, err := ioutil.ReadDir(writeFilesConfig.AbsOutputDirPath)
	if err != nil && err != os.ErrNotExist {
		return fmt.Errorf("failed to check output directory contents, error: %s", err)
	}
	containsArtifacts := false
	for _, entry := range entries {
		if !entry.IsDir() && (path.Ext(entry.Name()) != ".log") {
			containsArtifacts = true
			break
		}
	}
	if containsArtifacts {
		fmt.Println()
		log.Warnf("Export output directory exists and is not empty.")
	}

	if err := writeIdentities(identities.Content, writeFilesConfig.AbsOutputDirPath); err != nil {
		return err
	}
	if err := writeProvisioningProfiles(provisioningProfiles, writeFilesConfig.AbsOutputDirPath); err != nil {
		return err
	}
	return nil
}

// exportIdentities exports the given certificates merged in a single .p12 file
func exportIdentities(certificates []certificateutil.CertificateInfoModel, isAskForPassword bool) (models.Certificates, error) {
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

	identities, err := osxkeychain.ExportFromKeychain(identityKechainRefs, isAskForPassword)
	if err != nil {
		return models.Certificates{}, fmt.Errorf("failed to export from Keychain: %s", err)
	}
	return models.Certificates{
		Info:    certificates,
		Content: identities,
	}, nil
}

// writeIdentities writes identities to a file path
func writeIdentities(identites []byte, absExportOutputDirPath string) error {
	return ioutil.WriteFile(filepath.Join(absExportOutputDirPath, "Identities.p12"), identites, 0600)
}

// exportProvisioningProfiles returns provisioning profies
func exportProvisioningProfiles(profiles []profileutil.ProvisioningProfileInfoModel) ([]models.ProvisioningProfile, error) {
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

// writeProvisioningProfiles writes provisioning profiles to the filesystem
func writeProvisioningProfiles(profiles []models.ProvisioningProfile, absExportOutputDirPath string) error {
	for _, profile := range profiles {
		exportFileName := utility.ProfileExportFileNameNoPath(profile.Info)
		exportPth := filepath.Join(absExportOutputDirPath, exportFileName)
		if err := ioutil.WriteFile(exportPth, profile.Content, 0600); err != nil {
			return fmt.Errorf("failed to write file, error: %s", err)
		}
	}
	return nil
}
