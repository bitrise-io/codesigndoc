package bitriseclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/urlutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

// RegisterProvisioningProfileResponseData ...
type RegisterProvisioningProfileResponseData struct {
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
	Data RegisterProvisioningProfileResponseData `json:"data"`
}

// ConfirmProvProfileUploadResponseData ...
type ConfirmProvProfileUploadResponseData struct {
	UploadFileName string `json:"upload_file_name"`
	UploadFileSize int    `json:"upload_file_size"`
	Slug           string `json:"slug"`
	Processed      bool   `json:"processed"`
	IsExpose       bool   `json:"is_expose"`
	IsProtected    bool   `json:"dais_protectedta"`
}

// ConfirmProvProfileUploadResponse ...
type ConfirmProvProfileUploadResponse struct {
	Data ConfirmProvProfileUploadResponseData `json:"data"`
}

// FetchProvisioningProfileListResponseData ...
type FetchProvisioningProfileListResponseData struct {
	UploadFileName string `json:"upload_file_name"`
	UploadFileSize int    `json:"upload_file_size"`
	Slug           string `json:"slug"`
	Processed      bool   `json:"processed"`
	IsExpose       bool   `json:"is_expose"`
	IsProtected    bool   `json:"dais_protectedta"`
}

// FetchProvisioningProfileListResponse ...
type FetchProvisioningProfileListResponse struct {
	Data []FetchProvisioningProfileListResponseData `json:"data"`
}

// FetchUploadedProvisioningProfileResponseData ...
type FetchUploadedProvisioningProfileResponseData struct {
	UploadFileName string `json:"upload_file_name"`
	UploadFileSize int    `json:"upload_file_size"`
	Slug           string `json:"slug"`
	Processed      bool   `json:"processed"`
	IsExpose       bool   `json:"is_expose"`
	IsProtected    bool   `json:"dais_protectedta"`
	DownloadURL    string `json:"download_url"`
}

// FetchUploadedProvisioningProfileResponse ...
type FetchUploadedProvisioningProfileResponse struct {
	Data FetchUploadedProvisioningProfileResponseData `json:"data"`
}

// FetchProvisioningProfiles ...
func (client *BitriseClient) FetchProvisioningProfiles() ([]FetchProvisioningProfileListResponseData, error) {
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
	var requestResponse FetchProvisioningProfileListResponse

	//
	// Perform request
	if err := retry.Times(1).Wait(5 * time.Second).Try(func(attempt uint) error {
		body, statusCode, err := performRequest(request)
		if err != nil {
			log.Warnf("Attempt (%d) failed, error: %s", attempt+1, err)
			if !strings.Contains(err.Error(), "failed to perform request") {
				log.Warnf("Response status: %d", statusCode)
				log.Warnf("Body: %s", string(body))
			}
			return err
		}

		// Parse JSON body
		if err := json.Unmarshal([]byte(body), &requestResponse); err != nil {
			return fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
		}
		return nil

	}); err != nil {
		return nil, err
	}

	logDebugPretty(requestResponse)

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
	requestResponse := FetchUploadedProvisioningProfileResponse{}

	//
	// Perform request
	if err := retry.Times(1).Wait(5 * time.Second).Try(func(attempt uint) error {
		body, statusCode, err := performRequest(request)
		if err != nil {
			log.Warnf("Attempt (%d) failed, error: %s", attempt+1, err)
			if !strings.Contains(err.Error(), "failed to perform request") {
				log.Warnf("Response status: %d", statusCode)
				log.Warnf("Body: %s", string(body))
			}
			return err
		}

		// Parse JSON body
		if err := json.Unmarshal([]byte(body), &requestResponse); err != nil {
			return fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
		}

		return nil

	}); err != nil {
		return "", err
	}

	logDebugPretty(requestResponse)

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
	if err := retry.Times(1).Wait(5 * time.Second).Try(func(attempt uint) error {
		body, statusCode, err := performRequest(request)
		if err != nil {
			log.Warnf("Attempt (%d) failed, error: %s", attempt+1, err)
			if !strings.Contains(err.Error(), "failed to perform request") {
				log.Warnf("Response status: %d", statusCode)
				log.Warnf("Body: %s", string(body))
			}
			return err
		}

		requestResponse = string(body)

		return nil

	}); err != nil {
		return "", err
	}

	return requestResponse, nil

}

// RegisterProvisioningProfile ...
func (client *BitriseClient) RegisterProvisioningProfile(provisioningProfSize int64, profile profileutil.ProvisioningProfileInfoModel) (RegisterProvisioningProfileResponseData, error) {
	fmt.Println()
	log.Infof("Register %s on Bitrise...", profile.Name)

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, provisioningProfilesEndPoint)
	if err != nil {
		return RegisterProvisioningProfileResponseData{}, err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	fields := map[string]interface{}{
		"upload_file_name": profile.Name,
		"upload_file_size": provisioningProfSize,
	}

	request, err := createRequest(http.MethodPost, requestURL, client.headers, fields)
	if err != nil {
		return RegisterProvisioningProfileResponseData{}, err
	}

	// Response struct
	requestResponse := RegisterProvisioningProfileResponse{}

	//
	// Perform request
	if err := retry.Times(1).Wait(5 * time.Second).Try(func(attempt uint) error {
		body, statusCode, err := performRequest(request)
		if err != nil {
			log.Warnf("Attempt (%d) failed, error: %s", attempt+1, err)
			if !strings.Contains(err.Error(), "failed to perform request") {
				log.Warnf("Response status: %d", statusCode)
				log.Warnf("Body: %s", string(body))
			}
			return err
		}

		// Parse JSON body
		if err := json.Unmarshal([]byte(body), &requestResponse); err != nil {
			return fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
		}
		return nil

	}); err != nil {
		return RegisterProvisioningProfileResponseData{}, err
	}

	// Success
	logDebugPretty(requestResponse)

	return requestResponse.Data, nil
}

// UploadProvisioningProfile ...
func (client *BitriseClient) UploadProvisioningProfile(uploadURL string, uploadFileName string, outputDirPath string, exportFileName string) error {
	fmt.Println()
	log.Infof("Upload %s to Bitrise...", exportFileName)

	filePth := filepath.Join(outputDirPath, exportFileName)

	request, err := createUploadRequest(http.MethodPut, uploadURL, nil, filePth)
	if err != nil {
		return err
	}

	//
	// Perform request
	if err := retry.Times(1).Wait(5 * time.Second).Try(func(attempt uint) error {
		body, statusCode, err := performRequest(request)
		if err != nil {
			log.Warnf("Attempt (%d) failed, error: %s", attempt+1, err)
			if !strings.Contains(err.Error(), "failed to perform request") {
				log.Warnf("Response status: %d", statusCode)
				log.Warnf("Body: %s", string(body))
			}
			return err
		}

		return nil

	}); err != nil {
		return err
	}

	return nil
}

// ConfirmProvisioningProfileUpload ...
func (client *BitriseClient) ConfirmProvisioningProfileUpload(profileSlug string, provUploadName string) error {
	fmt.Println()
	log.Infof("Confirm - %s - upload to Bitrise...", provUploadName)

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
	if err := retry.Times(1).Wait(5 * time.Second).Try(func(attempt uint) error {
		body, statusCode, err := performRequest(request)
		if err != nil {
			log.Warnf("Attempt (%d) failed, error: %s", attempt+1, err)
			if !strings.Contains(err.Error(), "failed to perform request") {
				log.Warnf("Response status: %d", statusCode)
				log.Warnf("Body: %s", string(body))
			}
			return err
		}

		// Parse JSON body
		if err := json.Unmarshal([]byte(body), &requestResponse); err != nil {
			return fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
		}

		return nil

	}); err != nil {
		return err
	}

	logDebugPretty(requestResponse)

	return nil
}
