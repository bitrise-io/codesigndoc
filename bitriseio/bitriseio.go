package bitriseio

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bitrise-io/codesigndoc/bitriseio/bitrise"
	"github.com/bitrise-io/codesigndoc/codesign"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/profileutil"
	"github.com/bitrise-io/goinp/goinp"
)

// UploadCodesigningFiles ...
func UploadCodesigningFiles(certificates []certificateutil.CertificateInfoModel, profiles []profileutil.ProvisioningProfileInfoModel, certsOnly bool, outputDirPath string) (bool, bool, error) {
	accessToken, err := askAccessToken()
	if err != nil {
		return false, false, err
	}

	bitriseClient, appList, err := bitrise.NewClient(accessToken)
	if err != nil {
		return false, false, err
	}

	selectedAppSlug, err := selectApp(appList)
	if err != nil {
		return false, false, err
	}

	bitriseClient.SetSelectedAppSlug(selectedAppSlug)

	var provProfilesUploaded bool
	if !certsOnly {
		provProfilesUploaded, err = uploadExportedProvProfiles(bitriseClient, profiles, outputDirPath)
		if err != nil {
			return false, false, err
		}
	}

	certsUploaded, err := uploadExportedIdentity(bitriseClient, certificates, outputDirPath)
	if err != nil {
		return false, false, err
	}
	return certsUploaded, provProfilesUploaded, nil
}

// SetupClient asks for access token and app, returns a bitrise client
func SetupClient() (*bitrise.Client, error) {
	accessToken, err := askAccessToken()
	if err != nil {
		return nil, err
	}

	client, appList, err := bitrise.NewClient(accessToken)
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

// UploadCodesigningFilesAsStream ...
func UploadCodesigningFilesAsStream(client *bitrise.Client, certificates codesign.Certificates, profiles []profileutil.ProvisioningProfileInfoModel, certsOnly bool) (bool, bool, error) {
	var provProfilesUploaded bool
	if !certsOnly {
		var err error
		provProfilesUploaded, err = uploadExportedProvProfilesAsStream(client, profiles)
		if err != nil {
			return false, false, err
		}
	}

	certsUploaded, err := uploadExportedIdentityAsStream(client, certificates)
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

func uploadExportedProvProfiles(bitriseClient *bitrise.Client, profilesToExport []profileutil.ProvisioningProfileInfoModel, outputDirPath string) (bool, error) {
	fmt.Println()
	log.Infof("Uploading provisioning profiles...")

	profilesToUpload, err := filterAlreadyUploadedProvProfiles(bitriseClient, profilesToExport)
	if err != nil {
		return false, err
	}

	if len(profilesToUpload) > 0 {
		if err := uploadProvisioningProfiles(bitriseClient, profilesToUpload, outputDirPath); err != nil {
			return false, err
		}
	} else {
		log.Warnf("There is no new provisioning profile to upload...")
	}

	return true, nil
}

func uploadExportedProvProfilesAsStream(bitriseClient *bitrise.Client, profilesToExport []profileutil.ProvisioningProfileInfoModel) (bool, error) {
	fmt.Println()
	log.Infof("Uploading provisioning profiles...")

	profilesToUpload, err := filterAlreadyUploadedProvProfiles(bitriseClient, profilesToExport)
	if err != nil {
		return false, err
	}

	if len(profilesToUpload) > 0 {
		if err := uploadProvisioningProfilesAsStream(bitriseClient, profilesToUpload); err != nil {
			return false, err
		}
	} else {
		log.Warnf("There is no new provisioning profile to upload...")
	}

	return true, nil
}

func filterAlreadyUploadedProvProfiles(client *bitrise.Client, localProfiles []profileutil.ProvisioningProfileInfoModel) ([]profileutil.ProvisioningProfileInfoModel, error) {
	log.Printf("Looking for provisioning profile duplicates on Bitrise...")

	uploadedProfileUUIDList := map[string]bool{}
	profilesToUpload := []profileutil.ProvisioningProfileInfoModel{}

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
		contains, _ := uploadedProfileUUIDList[localProfile.UUID]
		if contains {
			log.Warnf("Already on Bitrise: - %s - (UUID: %s) ", localProfile.Name, localProfile.UUID)
		} else {
			profilesToUpload = append(profilesToUpload, localProfile)
		}
	}

	return profilesToUpload, nil
}

func uploadProvisioningProfiles(bitriseClient *bitrise.Client, profilesToUpload []profileutil.ProvisioningProfileInfoModel, outputDirPath string) error {
	for _, profile := range profilesToUpload {
		exportFileName := provProfileExportFileName(profile, outputDirPath)

		provProfile, err := os.Open(outputDirPath + "/" + exportFileName)
		if err != nil {
			return err
		}

		defer func() {
			if err := provProfile.Close(); err != nil {
				log.Warnf("Provisioning profile close failed, err: %s", err)
			}

		}()

		info, err := provProfile.Stat()
		if err != nil {
			return err
		}

		log.Debugf("\n%s size: %d", exportFileName, info.Size())

		provProfSlugResponseData, err := bitriseClient.RegisterProvisioningProfile(info.Size(), exportFileName)
		if err != nil {
			return err
		}

		if err := bitriseClient.UploadProvisioningProfile(provProfSlugResponseData.UploadURL, provProfSlugResponseData.UploadFileName, outputDirPath, exportFileName); err != nil {
			return err
		}

		if err := bitriseClient.ConfirmProvisioningProfileUpload(provProfSlugResponseData.Slug, provProfSlugResponseData.UploadFileName); err != nil {
			return err
		}
	}

	return nil
}

func uploadProvisioningProfilesAsStream(bitriseClient *bitrise.Client, profilesToUpload []profileutil.ProvisioningProfileInfoModel) error {
	for _, profile := range profilesToUpload {
		exportFileName := provProfileExportFileNameNoPath(profile)
		exportSize := int64(len(profile.Content()))

		log.Debugf("\n%s size: %d", exportFileName, exportSize)

		provProfSlugResponseData, err := bitriseClient.RegisterProvisioningProfile(exportSize, exportFileName)
		if err != nil {
			return err
		}

		if err := bitriseClient.UploadProvisioningProfileAsStream(provProfSlugResponseData.UploadURL, provProfSlugResponseData.UploadFileName, bytes.NewReader(profile.Content())); err != nil {
			return err
		}

		if err := bitriseClient.ConfirmProvisioningProfileUpload(provProfSlugResponseData.Slug, provProfSlugResponseData.UploadFileName); err != nil {
			return err
		}
	}

	return nil
}

func provProfileExportFileName(info profileutil.ProvisioningProfileInfoModel, path string) string {
	replaceRexp, err := regexp.Compile("[^A-Za-z0-9_.-]")
	if err != nil {
		log.Warnf("Invalid regex, error: %s", err)
		return ""
	}
	safeTitle := replaceRexp.ReplaceAllString(info.Name, "")
	extension := ".mobileprovision"
	if strings.HasSuffix(path, ".provisionprofile") {
		extension = ".provisionprofile"
	}

	return info.UUID + "." + safeTitle + extension
}

func provProfileExportFileNameNoPath(info profileutil.ProvisioningProfileInfoModel) string {
	replaceRexp, err := regexp.Compile("[^A-Za-z0-9_.-]")
	if err != nil {
		log.Warnf("Invalid regex, error: %s", err)
		return ""
	}
	safeTitle := replaceRexp.ReplaceAllString(info.Name, "")
	extension := ".mobileprovision"
	if info.Type == profileutil.ProfileTypeMacOs {
		extension = ".provisionprofile"
	}

	return info.UUID + "." + safeTitle + extension
}

func uploadExportedIdentity(bitriseClient *bitrise.Client, certificatesToExport []certificateutil.CertificateInfoModel, outputDirPath string) (bool, error) {
	fmt.Println()
	log.Infof("Uploading certificate...")

	shouldUploadIdentities, err := shouldUploadCertificates(bitriseClient, certificatesToExport)
	if err != nil {
		return false, err
	}

	if shouldUploadIdentities {

		if err := uploadIdentity(bitriseClient, outputDirPath); err != nil {
			return false, err
		}
	} else {
		log.Warnf("There is no new certificate to upload...")
	}

	return true, err
}

func uploadExportedIdentityAsStream(bitriseClient *bitrise.Client, certificates codesign.Certificates) (bool, error) {
	fmt.Println()
	log.Infof("Uploading certificate...")

	shouldUploadIdentities, err := shouldUploadCertificates(bitriseClient, certificates.Certificates)
	if err != nil {
		return false, err
	}

	if shouldUploadIdentities {
		if err := uploadIdentityAsStream(bitriseClient, certificates.Contents); err != nil {
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

func uploadIdentity(bitriseClient *bitrise.Client, outputDirPath string) error {
	identities, err := os.Open(outputDirPath + "/" + "Identities.p12")
	if err != nil {
		return err
	}

	defer func() {
		if err := identities.Close(); err != nil {
			log.Warnf("Identities failed, err: %s", err)
		}

	}()

	info, err := identities.Stat()
	if err != nil {
		return err
	}

	log.Debugf("\n%s size: %d", "Identities.p12", info.Size())

	certificateResponseData, err := bitriseClient.RegisterIdentity(info.Size())
	if err != nil {
		return err
	}

	if err := bitriseClient.UploadIdentity(certificateResponseData.UploadURL, certificateResponseData.UploadFileName, outputDirPath, "Identities.p12"); err != nil {
		return err
	}

	return bitriseClient.ConfirmIdentityUpload(certificateResponseData.Slug, certificateResponseData.UploadFileName)
}

func uploadIdentityAsStream(bitriseClient *bitrise.Client, identities []byte) error {
	const identitiesName = "Identities.p12"
	identitiesSize := int64(len(identities))

	log.Debugf("\n%s size: %d", identitiesName, identitiesSize)

	certificateResponseData, err := bitriseClient.RegisterIdentity(identitiesSize)
	if err != nil {
		return err
	}

	if err := bitriseClient.UploadIdentityAsStream(certificateResponseData.UploadURL, certificateResponseData.UploadFileName, identitiesName, bytes.NewReader(identities)); err != nil {
		return err
	}

	return bitriseClient.ConfirmIdentityUpload(certificateResponseData.Slug, certificateResponseData.UploadFileName)
}
