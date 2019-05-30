package bitrise

import (
	"net/http"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/urlutil"
	"github.com/bitrise-io/go-xcode/profileutil"
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
func (client *Client) FetchProvisioningProfiles() ([]ProvisioningProfileListData, error) {
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
func (client *Client) GetUploadedProvisioningProfileUUIDby(profileSlug string) (UUID string, err error) {
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

	data, err := profileutil.NewProvisioningProfileInfo(*plistData)
	if err != nil {
		return "", err
	}

	return data.UUID, nil
}

func (client *Client) getUploadedProvisioningProfileDownloadURLBy(profileSlug string) (downloadURL string, err error) {
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

func (client *Client) downloadUploadedProvisioningProfile(downloadURL string) (content string, err error) {
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
func (client *Client) RegisterProvisioningProfile(provisioningProfSize int64, exportedProfileName string) (RegisterProvisioningProfileData, error) {
	log.Printf("Register %s on Bitrise...", exportedProfileName)

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, provisioningProfilesEndPoint)
	if err != nil {
		return RegisterProvisioningProfileData{}, err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	fields := map[string]interface{}{
		"upload_file_name": exportedProfileName,
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

// ConfirmProvisioningProfileUpload ...
func (client *Client) ConfirmProvisioningProfileUpload(profileSlug string, provUploadName string) error {
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
