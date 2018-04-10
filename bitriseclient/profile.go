package bitriseclient

import (
	"net/http"
	"path/filepath"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/urlutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

// RegisterProvisioningProfileData ...
type RegisterProvisioningProfileData struct {
	UploadFileName string `json:"upload_file_name"`
	UploadFileSize int64  `json:"upload_file_size"`
	Slug           string `json:"slug"`
	Processed      bool   `json:"processed"`
	IsExpose       bool   `json:"is_expose"`
	IsProtected    bool   `json:"is_protected"`
	UploadURL      string `json:"upload_url"`
}

// RegisterProvisioningProfileResponse ...
type RegisterProvisioningProfileResponse struct {
	Data RegisterProvisioningProfileData `json:"data"`
}

// ConfirmProvProfileUploadData ...
type ConfirmProvProfileUploadData struct {
	UploadFileName string `json:"upload_file_name"`
	UploadFileSize int    `json:"upload_file_size"`
	Slug           string `json:"slug"`
	Processed      bool   `json:"processed"`
	IsExpose       bool   `json:"is_expose"`
	IsProtected    bool   `json:"dais_protectedta"`
}

// ConfirmProvProfileUploadResponse ...
type ConfirmProvProfileUploadResponse struct {
	Data ConfirmProvProfileUploadData `json:"data"`
}

// ProvisioningProfileListData ...
type ProvisioningProfileListData struct {
	UploadFileName string `json:"upload_file_name"`
	UploadFileSize int    `json:"upload_file_size"`
	Slug           string `json:"slug"`
	Processed      bool   `json:"processed"`
	IsExpose       bool   `json:"is_expose"`
	IsProtected    bool   `json:"dais_protectedta"`
}

// ProvisioningProfileListResponse ...
type ProvisioningProfileListResponse struct {
	Data []ProvisioningProfileListData `json:"data"`
}

// UploadedProvisioningProfileData ...
type UploadedProvisioningProfileData struct {
	UploadFileName string `json:"upload_file_name"`
	UploadFileSize int    `json:"upload_file_size"`
	Slug           string `json:"slug"`
	Processed      bool   `json:"processed"`
	IsExpose       bool   `json:"is_expose"`
	IsProtected    bool   `json:"dais_protectedta"`
	DownloadURL    string `json:"download_url"`
}

// UploadedProvisioningProfileResponse ...
type UploadedProvisioningProfileResponse struct {
	Data UploadedProvisioningProfileData `json:"data"`
}

// FetchProvisioningProfiles ...
func (client *BitriseClient) FetchProvisioningProfiles() ([]ProvisioningProfileListData, error) {
	log.Debugf("\nDownloading provisioning profile list from Bitrise...")

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, provisioningProfilesEndPoint)
	if err != nil {
		return nil, err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	request, err := createRequest(http.MethodGet, requestURL, client.headers, nil)
	if err != nil {
		return nil, err
	}

	// Response struct
	var requestResponse ProvisioningProfileListResponse

	//
	// Perform request
	response, _, err := RunRequest(client, request, &requestResponse)
	if err != nil {
		return nil, err
	}

	requestResponse = *response.(*ProvisioningProfileListResponse)
	return requestResponse.Data, nil
}

// GetUploadedProvisioningProfileUUIDby ...
func (client *BitriseClient) GetUploadedProvisioningProfileUUIDby(profileSlug string) (UUID string, err error) {
	downloadURL, err := client.getUploadedProvisioningProfileDownloadURLBy(profileSlug)
	if err != nil {
		return "", err
	}

	content, err := client.downloadUploadedProvisioningProfile(downloadURL)
	if err != nil {
		return "", err
	}

	plistData, err := profileutil.ProvisioningProfileFromContent([]byte(content))
	if err != nil {
		return "", err
	}

	data, err := profileutil.NewProvisioningProfileInfo(*plistData, profileutil.ProfileTypeIos)
	if err != nil {
		return "", err
	}

	return data.UUID, nil
}

func (client *BitriseClient) getUploadedProvisioningProfileDownloadURLBy(profileSlug string) (downloadURL string, err error) {
	log.Debugf("\nGet downloadURL for provisioning profile (slug - %s) from Bitrise...", profileSlug)

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, provisioningProfilesEndPoint, profileSlug)
	if err != nil {
		return "", err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	request, err := createRequest(http.MethodGet, requestURL, client.headers, nil)
	if err != nil {
		return "", err
	}

	// Response struct
	requestResponse := UploadedProvisioningProfileResponse{}

	//
	// Perform request
	response, _, err := RunRequest(client, request, &requestResponse)
	if err != nil {
		return "", err
	}

	requestResponse = *response.(*UploadedProvisioningProfileResponse)
	return requestResponse.Data.DownloadURL, nil
}

func (client *BitriseClient) downloadUploadedProvisioningProfile(downloadURL string) (content string, err error) {
	log.Debugf("\nDownloading provisioning profile from Bitrise...")
	log.Debugf("\nRequest URL: %s", downloadURL)

	request, err := createRequest(http.MethodGet, downloadURL, nil, nil)
	if err != nil {
		return "", err
	}

	// Response struct
	var requestResponse string

	//
	// Perform request
	_, body, err := RunRequest(client, request, nil)
	if err != nil {
		return "", err
	}

	requestResponse = string(body)
	return requestResponse, nil

}

// RegisterProvisioningProfile ...
func (client *BitriseClient) RegisterProvisioningProfile(provisioningProfSize int64, profile profileutil.ProvisioningProfileInfoModel) (RegisterProvisioningProfileData, error) {
	log.Printf("Register %s on Bitrise...", profile.Name)

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, provisioningProfilesEndPoint)
	if err != nil {
		return RegisterProvisioningProfileData{}, err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	fields := map[string]interface{}{
		"upload_file_name": profile.Name,
		"upload_file_size": provisioningProfSize,
	}

	request, err := createRequest(http.MethodPost, requestURL, client.headers, fields)
	if err != nil {
		return RegisterProvisioningProfileData{}, err
	}

	// Response struct
	requestResponse := RegisterProvisioningProfileResponse{}

	//
	// Perform request
	response, _, err := RunRequest(client, request, &requestResponse)
	if err != nil {
		return RegisterProvisioningProfileData{}, err
	}

	requestResponse = *response.(*RegisterProvisioningProfileResponse)
	return requestResponse.Data, nil
}

// UploadProvisioningProfile ...
func (client *BitriseClient) UploadProvisioningProfile(uploadURL string, uploadFileName string, outputDirPath string, exportFileName string) error {
	log.Printf("Upload %s to Bitrise...", exportFileName)

	filePth := filepath.Join(outputDirPath, exportFileName)

	request, err := createUploadRequest(http.MethodPut, uploadURL, nil, filePth)
	if err != nil {
		return err
	}

	//
	// Perform request
	_, _, err = RunRequest(client, request, nil)
	if err != nil {
		return err
	}

	return nil
}

// ConfirmProvisioningProfileUpload ...
func (client *BitriseClient) ConfirmProvisioningProfileUpload(profileSlug string, provUploadName string) error {
	log.Printf("Confirm - %s - upload to Bitrise...", provUploadName)

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, provisioningProfilesEndPoint, profileSlug, "uploaded")
	if err != nil {
		return err
	}

	request, err := createRequest("POST", requestURL, client.headers, nil)
	if err != nil {
		return err
	}

	// Response struct
	requestResponse := ConfirmProvProfileUploadResponse{}

	//
	// Perform request
	response, _, err := RunRequest(client, request, &requestResponse)
	if err != nil {
		return err
	}

	requestResponse = *response.(*ConfirmProvProfileUploadResponse)
	return nil
}
