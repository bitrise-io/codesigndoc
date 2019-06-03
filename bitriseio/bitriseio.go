package bitriseio

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/bitrise-io/codesigndoc/bitriseio/bitrise"
	"github.com/bitrise-io/codesigndoc/models"
	"github.com/bitrise-io/codesigndoc/utility"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/goinp/goinp"
)

// GetInteractiveConfigClient asks for access token and app, returns a bitrise client
func GetInteractiveConfigClient() (*bitrise.Client, error) {
	accessToken, err := askAccessToken()
	if err != nil {
		return nil, err
	}

	client, err := bitrise.NewClient(accessToken)
	if err != nil {
		return nil, err
	}

	appList, err := client.GetAppList()
	if err != nil {
		return nil, err
	}

	selectedAppSlug, err := selectApp(appList)
	if err != nil {
		return nil, err
	}
	client.SetSelectedAppSlug(selectedAppSlug)

	return client, nil
}

// UploadCodesigningFiles ...
func UploadCodesigningFiles(client *bitrise.Client, certificates models.Certificates, profiles []models.ProvisioningProfile) (bool, bool, error) {
	var provProfilesUploaded bool
	if len(profiles) != 0 {
		var err error
		provProfilesUploaded, err = uploadExportedProvProfiles(client, profiles)
		if err != nil {
			return false, false, err
		}
	}

	certsUploaded, err := uploadExportedIdentity(client, certificates)
	if err != nil {
		return false, false, err
	}
	return certsUploaded, provProfilesUploaded, nil
}

func askAccessToken() (token string, err error) {
	messageToAsk := `Please copy your personal access token to Bitrise.
(To acquire a Personal Access Token for your user, sign in with that user on bitrise.io, go to your Account Settings page,
and select the Security tab on the left side.)`
	fmt.Println()

	accesToken, err := goinp.AskForStringFromReader(messageToAsk, os.Stdin)
	if err != nil {
		return accesToken, err
	}

	fmt.Println()
	log.Infof("%s %s", colorstring.Green("Given accesToken:"), accesToken)
	fmt.Println()

	return accesToken, nil
}

func selectApp(appList []bitrise.Application) (seledtedAppSlug string, err error) {
	var selectionList []string

	for _, app := range appList {
		selectionList = append(selectionList, app.Title+" ("+app.RepoURL+")")
	}
	userSelection, err := goinp.SelectFromStringsWithDefault("Select the app which you want to upload the privisioning profiles", 1, selectionList)

	if err != nil {
		return "", fmt.Errorf("failed to read input: %s", err)

	}

	log.Debugf("selected app: %v", userSelection)

	for index, selected := range selectionList {
		if selected == userSelection {
			return appList[index].Slug, nil
		}
	}

	return "", errors.New("failed to find selected app in appList")
}

func uploadExportedProvProfiles(bitriseClient *bitrise.Client, profilesToExport []models.ProvisioningProfile) (bool, error) {
	fmt.Println()
	log.Infof("Uploading provisioning profiles...")

	profilesToUpload, err := filterAlreadyUploadedProvProfiles(bitriseClient, profilesToExport)
	if err != nil {
		return false, err
	}

	if len(profilesToUpload) > 0 {
		if err := uploadProvisioningProfiles(bitriseClient, profilesToUpload); err != nil {
			return false, err
		}
	} else {
		log.Warnf("There is no new provisioning profile to upload...")
	}

	return true, nil
}

func filterAlreadyUploadedProvProfiles(client *bitrise.Client, localProfiles []models.ProvisioningProfile) ([]models.ProvisioningProfile, error) {
	log.Printf("Looking for provisioning profile duplicates on Bitrise...")

	uploadedProfileUUIDList := map[string]bool{}
	var profilesToUpload []models.ProvisioningProfile

	uploadedProfInfoList, err := client.FetchProvisioningProfiles()
	if err != nil {
		return nil, err
	}

	for _, uploadedProfileInfo := range uploadedProfInfoList {
		uploadedProfileUUID, err := client.GetUploadedProvisioningProfileUUIDby(uploadedProfileInfo.Slug)
		if err != nil {
			return nil, err
		}

		uploadedProfileUUIDList[uploadedProfileUUID] = true
	}

	for _, localProfile := range localProfiles {
		contains, _ := uploadedProfileUUIDList[localProfile.Info.UUID]
		if contains {
			log.Warnf("Already on Bitrise: - %s - (UUID: %s) ", localProfile.Info.Name, localProfile.Info.UUID)
		} else {
			profilesToUpload = append(profilesToUpload, localProfile)
		}
	}

	return profilesToUpload, nil
}

func uploadProvisioningProfiles(bitriseClient *bitrise.Client, profilesToUpload []models.ProvisioningProfile) error {
	for _, profile := range profilesToUpload {
		exportFileName := utility.ProfileExportFileNameNoPath(profile.Info)
		exportSize := int64(len(profile.Content))

		log.Debugf("\n%s size: %d", exportFileName, exportSize)

		provProfSlugResponseData, err := bitriseClient.RegisterProvisioningProfile(exportSize, exportFileName)
		if err != nil {
			return err
		}

		log.Printf("Uploading %s to Bitrise...", provProfSlugResponseData.UploadFileName)
		if err := bitriseClient.UploadArtifact(provProfSlugResponseData.UploadURL, bytes.NewReader(profile.Content)); err != nil {
			return err
		}

		if err := bitriseClient.ConfirmProvisioningProfileUpload(provProfSlugResponseData.Slug, provProfSlugResponseData.UploadFileName); err != nil {
			return err
		}
	}

	return nil
}

func uploadExportedIdentity(bitriseClient *bitrise.Client, certificates models.Certificates) (bool, error) {
	fmt.Println()
	log.Infof("Uploading certificate...")

	shouldUploadIdentities, err := shouldUploadCertificates(bitriseClient, certificates.Info)
	if err != nil {
		return false, err
	}

	if shouldUploadIdentities {
		if err := uploadIdentity(bitriseClient, certificates.Content); err != nil {
			return false, err
		}
	} else {
		log.Warnf("There is no new certificate to upload...")
	}

	return true, err
}

func shouldUploadCertificates(client *bitrise.Client, certificatesToExport []certificateutil.CertificateInfoModel) (bool, error) {
	log.Printf("Looking for certificate duplicates on Bitrise...")

	var uploadedCertificatesSerialList []string
	localCertificatesSerialList := []string{}

	uploadedItentityList, err := client.FetchUploadedIdentities()
	if err != nil {
		return false, err
	}

	// Get uploaded certificates' serials
	for _, uploadedIdentity := range uploadedItentityList {
		var serialListAsString []string

		serialList, err := client.GetUploadedCertificatesSerialby(uploadedIdentity.Slug)
		if err != nil {
			return false, err
		}

		for _, serial := range serialList {
			serialListAsString = append(serialListAsString, serial.String())
		}
		uploadedCertificatesSerialList = append(uploadedCertificatesSerialList, serialListAsString...)
	}

	for _, certificateToExport := range certificatesToExport {
		localCertificatesSerialList = append(localCertificatesSerialList, certificateToExport.Serial)
	}

	log.Debugf("Uploaded certificates' serial list: \n\t%v", uploadedCertificatesSerialList)
	log.Debugf("Local certificates' serial list: \n\t%v", localCertificatesSerialList)

	// Search for a new certificate
	for _, localCertificateSerial := range localCertificatesSerialList {
		if !sliceutil.IsStringInSlice(localCertificateSerial, uploadedCertificatesSerialList) {
			return true, nil
		}
	}

	return false, nil
}

func uploadIdentity(bitriseClient *bitrise.Client, identities []byte) error {
	identitiesSize := int64(len(identities))
	log.Debugf("\nIdentities size: %d", identitiesSize)

	certificateResponseData, err := bitriseClient.RegisterIdentity(identitiesSize)
	if err != nil {
		return err
	}

	log.Printf("Uploading %s to Bitrise...", certificateResponseData.UploadFileName)
	if err := bitriseClient.UploadArtifact(certificateResponseData.UploadURL, bytes.NewReader(identities)); err != nil {
		return err
	}

	return bitriseClient.ConfirmIdentityUpload(certificateResponseData.Slug, certificateResponseData.UploadFileName)
}
