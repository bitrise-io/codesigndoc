package bitriseclient

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/urlutil"
	"github.com/bitrise-tools/go-xcode/certificateutil"
)

// RegisterIdentityResponseData ...
type RegisterIdentityResponseData struct {
	UploadFileName string `json:"upload_file_name"`
	UploadFileSize int64  `json:"upload_file_size"`
	Slug           string `json:"slug"`
	Processed      bool   `json:"processed"`
	IsExpose       bool   `json:"is_expose"`
	IsProtected    bool   `json:"is_protected"`
	UploadURL      string `json:"upload_url"`
}

// RegisterIdentityResponse ...
type RegisterIdentityResponse struct {
	Data RegisterIdentityResponseData `json:"data"`
}

// ConfirmIdentityUploadResponseData ...
type ConfirmIdentityUploadResponseData struct {
	UploadFileName      string `json:"upload_file_name"`
	UploadFileSize      int    `json:"upload_file_size"`
	Slug                string `json:"slug"`
	Processed           bool   `json:"processed"`
	CertificatePassword string `json:"certificate_password"`
	IsExpose            bool   `json:"is_expose"`
	IsProtected         bool   `json:"dais_protectedta"`
}

// ConfirmIdentityUploadResponse ...
type ConfirmIdentityUploadResponse struct {
	Data ConfirmIdentityUploadResponseData `json:"data"`
}

// FetchUploadedIdentityListResponseData ...
type FetchUploadedIdentityListResponseData struct {
	UploadFileName      string `json:"upload_file_name"`
	UploadFileSize      int    `json:"upload_file_size"`
	Slug                string `json:"slug"`
	Processed           bool   `json:"processed"`
	CertificatePassword string `json:"certificate_password"`
	IsExpose            bool   `json:"is_expose"`
	IsProtected         bool   `json:"dais_protectedta"`
}

// FetchUploadedIdentityListResponse ...
type FetchUploadedIdentityListResponse struct {
	Data []FetchUploadedIdentityListResponseData `json:"data"`
}

// FetchUploadedIdentityResponseData ...
type FetchUploadedIdentityResponseData struct {
	UploadFileName      string `json:"upload_file_name"`
	UploadFileSize      int    `json:"upload_file_size"`
	Slug                string `json:"slug"`
	Processed           bool   `json:"processed"`
	CertificatePassword string `json:"certificate_password"`
	IsExpose            bool   `json:"is_expose"`
	IsProtected         bool   `json:"dais_protectedta"`
	DownloadURL         string `json:"download_url"`
}

// FetchUploadedIdentityResponse ...
type FetchUploadedIdentityResponse struct {
	Data FetchUploadedIdentityResponseData `json:"data"`
}

// FetchUploadedIdentities ...
func (client *BitriseClient) FetchUploadedIdentities() ([]FetchUploadedIdentityListResponseData, error) {
	log.Debugf("\nDownloading provisioning profile list from Bitrise...")

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, certificatesEndPoint)
	if err != nil {
		return []FetchUploadedIdentityListResponseData{}, err
	}

	request, err := createRequest(http.MethodGet, requestURL, client.headers, nil)
	if err != nil {
		return []FetchUploadedIdentityListResponseData{}, err
	}
	log.Debugf("\nRequest URL: %s", requestURL)

	// Response struct
	requestResponse := FetchUploadedIdentityListResponse{}
	responseStatusCode := -1

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

		responseStatusCode = statusCode

		// Parse JSON body
		if err := json.Unmarshal([]byte(body), &requestResponse); err != nil {
			return fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
		}
		return nil

	}); err != nil {
		return []FetchUploadedIdentityListResponseData{}, err
	}

	logDebugPretty(requestResponse)

	return requestResponse.Data, nil
}

// GetUploadedCertificatesSerialby ...
func (client *BitriseClient) GetUploadedCertificatesSerialby(identitySlug string) (certificateSerialList []big.Int, err error) {
	downloadURL, certificatePassword, err := client.getUploadedIdentityDownloadURLBy(identitySlug)
	if err != nil {
		return nil, err
	}

	content, err := client.downloadUploadedIdentity(downloadURL)
	if err != nil {
		return nil, err
	}

	certificates, err := certificateutil.CertificatesFromPKCS12Content([]byte(content), certificatePassword)
	if err != nil {
		return nil, err
	}

	var serialList []big.Int

	for _, certificate := range certificates {
		serialList = append(serialList, *certificate.SerialNumber)
	}

	return serialList, nil
}

func (client *BitriseClient) getUploadedIdentityDownloadURLBy(certificateSlug string) (downloadURL string, password string, err error) {
	log.Debugf("\nGet downloadURL for certificate (slug - %s) from Bitrise...", certificateSlug)

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, certificatesEndPoint, certificateSlug)
	if err != nil {
		return "", "", err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	request, err := createRequest(http.MethodGet, requestURL, client.headers, nil)
	if err != nil {
		return "", "", err
	}

	// Response struct
	requestResponse := FetchUploadedIdentityResponse{}
	responseStatusCode := -1

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

		responseStatusCode = statusCode

		// Parse JSON body
		if err := json.Unmarshal([]byte(body), &requestResponse); err != nil {
			return fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
		}

		return nil

	}); err != nil {
		return "", "", err
	}

	logDebugPretty(requestResponse)

	return requestResponse.Data.DownloadURL, requestResponse.Data.CertificatePassword, nil
}

func (client *BitriseClient) downloadUploadedIdentity(downloadURL string) (content string, err error) {
	log.Debugf("\nDownloading identities from Bitrise...")
	log.Debugf("\nRequest URL: %s", downloadURL)

	request, err := createRequest(http.MethodGet, downloadURL, nil, nil)
	if err != nil {
		return "", err
	}

	// Response struct
	responseStatusCode := -1
	requestResponse := ""

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

		responseStatusCode = statusCode
		requestResponse = string(body)

		return nil

	}); err != nil {
		return "", err
	}

	return requestResponse, nil

}

// RegisterIdentity ...
func (client *BitriseClient) RegisterIdentity(certificateSize int64) (RegisterIdentityResponseData, error) {
	fmt.Println()
	log.Infof("Register %s on Bitrise...", "Identities.p12")

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, certificatesEndPoint)
	if err != nil {
		return RegisterIdentityResponseData{}, err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	fields := map[string]interface{}{
		"upload_file_name": "Identities.p12",
		"upload_file_size": certificateSize,
	}

	request, err := createRequest(http.MethodPost, requestURL, client.headers, fields)
	if err != nil {
		return RegisterIdentityResponseData{}, err
	}

	// Response struct
	requestResponse := RegisterIdentityResponse{}
	responseStatusCode := -1

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

		responseStatusCode = statusCode

		// Parse JSON body
		if err := json.Unmarshal([]byte(body), &requestResponse); err != nil {
			return fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
		}

		logDebugPretty(requestResponse)

		return nil

	}); err != nil {
		return RegisterIdentityResponseData{}, err
	}

	return requestResponse.Data, nil
}

// UploadIdentity ...
func (client *BitriseClient) UploadIdentity(uploadURL string, uploadFileName string, outputDirPath string, exportFileName string) error {
	fmt.Println()
	log.Infof("Upload %s to Bitrise...", exportFileName)

	filePth := filepath.Join(outputDirPath, exportFileName)

	request, err := createUploadRequest(http.MethodPut, uploadURL, nil, filePth)
	if err != nil {
		return err
	}

	// Response struct
	responseStatusCode := -1

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

		responseStatusCode = statusCode

		return nil

	}); err != nil {
		return err
	}

	return nil
}

// ConfirmIdentityUpload ...
func (client *BitriseClient) ConfirmIdentityUpload(certificateSlug string, certificateUploadName string) error {
	fmt.Println()
	log.Infof("Confirm - %s - upload to Bitrise...", certificateUploadName)

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, "build-certificates", certificateSlug, "uploaded")
	if err != nil {
		return err
	}

	request, err := createRequest(http.MethodPost, requestURL, client.headers, nil)
	if err != nil {
		return err
	}

	// Response struct
	responseStatusCode := -1
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

		responseStatusCode = statusCode

		// Parse JSON body
		if err := json.Unmarshal([]byte(body), &requestResponse); err != nil {
			return fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
		}

		logDebugPretty(requestResponse)

		return nil

	}); err != nil {
		return err
	}

	return nil
}
