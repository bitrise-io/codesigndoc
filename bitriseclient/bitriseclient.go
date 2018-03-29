package bitriseclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/urlutil"
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

const (
	baseURL                      = "https://api.bitrise.io/v0.1/"
	appsEndPoint                 = "/apps"
	provisioningProfilesEndPoint = "/provisioning-profiles"
	certificatesEndPoint         = "/build-certificates"
)

// Protocol ...
type Protocol interface {
	New(accessToken string) (apps []Application, err error)

	SetSelectedAppSlug(slug string)

	FetchProvisioningProfiles() ([]FetchProvisioningProfileListResponseData, error)

	getUploadedProvisioningProfileDownloadURLBy(profileSlug string) (downloadURL string, err error)

	downloadUploadedProvisioningProfile(downloadURL string) (content string, err error)

	GetUploadedProvisioningProfileUUIDby(profileSlug string) (UUID string, err error)

	RegisterProvisioningProfile(provisioningProfSize int64, profile profileutil.ProvisioningProfileInfoModel) (RegisterProvisioningProfileResponseData, error)

	UploadProvisioningProfile(uploadURL string, uploadFileName string, outputDirPath string, exportFileName string) error

	ConfirmProvisioningProfileUpload(profileSlug string, provUploadName string) error

	FetchUploadedIdentities() ([]FetchUploadedIdentityListResponseData, error)

	getUploadedIdentityDownloadURLBy(certificateSlug string) (downloadURL string, password string, err error)

	downloadUploadedIdentity(downloadURL string) (content string, err error)

	GetUploadedCertificatesSerialby(identitySlug string) (certificateSerialList []big.Int, err error)

	RegisterIdentity(certificateSize int64) (RegisterIdentityResponseData, error)

	UploadIdentity(uploadURL string, uploadFileName string, outputDirPath string, exportFileName string) error

	ConfirmIdentityUpload(certificateSlug string, certificateUploadName string) error
}

// BitriseClient ...
type BitriseClient struct {
	accessToken     string
	selectedAppSlug string
	headers         map[string]string
}

// New returns all the application of the user on Bitrise
func (client *BitriseClient) New(accessToken string) (apps []Application, err error) {
	client.accessToken = accessToken
	client.headers = map[string]string{"Authorization": "token " + client.accessToken}

	log.Infof("Asking your application list from Bitrise...")

	requestURL, err := urlutil.Join(baseURL, appsEndPoint)
	if err != nil {
		return []Application{}, err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	request, err := createRequest(http.MethodGet, requestURL, client.headers, nil)
	if err != nil {
		return []Application{}, err
	}

	// Response struct
	var appListResponse FetchMyAppsResponse
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
		if err := json.Unmarshal([]byte(body), &appListResponse); err != nil {
			return fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
		}
		return nil

	}); err != nil {
		return []Application{}, err
	}

	// Success
	log.Donef("Request succeeded with status code: %d", responseStatusCode)
	logDebugPretty(appListResponse)

	return appListResponse.Data, nil
}

// SetSelectedAppSlug ...
func (client *BitriseClient) SetSelectedAppSlug(slug string) {
	client.selectedAppSlug = slug
}

// -------------------------------------------------
// -- Provisioning profiles

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
		return nil, err
	}

	// Success
	log.Debugf("Request succeeded with status code: %d", responseStatusCode)
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
		return "", err
	}

	// Success
	log.Debugf("Request succeeded with status code: %d", responseStatusCode)
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

	// Success
	log.Debugf("Request succeeded with status code: %d", responseStatusCode)

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
		return RegisterProvisioningProfileResponseData{}, err
	}

	// Success
	log.Donef("Request succeeded with status code: %d", responseStatusCode)
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

	// Success
	log.Donef("Request succeeded with status code: %d", responseStatusCode)

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

		return nil

	}); err != nil {
		return err
	}

	// Success
	log.Donef("Request succeeded with status code: %d", responseStatusCode)
	logDebugPretty(requestResponse)

	return nil
}

// -------------------------------------------------
// -- Certificates

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

	// Success
	log.Debugf("Request succeeded with status code: %d", responseStatusCode)
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

	// Success
	log.Debugf("Request succeeded with status code: %d", responseStatusCode)
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

	// Success
	log.Debugf("Request succeeded with status code: %d", responseStatusCode)

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

	// Success
	log.Donef("Request succeeded with status code: %d", responseStatusCode)

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

	// Success
	log.Donef("Request succeeded with status code: %d", responseStatusCode)

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

	request, err := createRequest("POST", requestURL, client.headers, nil)
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

	// Success
	log.Donef("Request succeeded with status code: %d", responseStatusCode)

	return nil
}

// -------------------------------------------------
// -- Commons

func createUploadRequest(requestMethod string, url string, headers map[string]string, filePth string) (*http.Request, error) {
	var fContent []byte

	f, err := os.Open(filePth)
	if err != nil {
		return nil, err

	}

	fContent, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(fContent))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)

	}

	return req, nil
}

func createRequest(requestMethod string, url string, headers map[string]string, fields map[string]interface{}) (*http.Request, error) {
	b := new(bytes.Buffer)

	if len(fields) > 0 {
		err := json.NewEncoder(b).Encode(fields)
		if err != nil {
			return nil, err
		}
	}

	log.Debugf("Request body: %s", string(b.Bytes()))

	req, err := http.NewRequest(requestMethod, url, bytes.NewReader(b.Bytes()))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)

	}

	return req, nil
}

func performRequest(request *http.Request) ([]byte, int, error) {
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		// On error, any Response can be ignored
		return []byte{}, -1, fmt.Errorf("failed to perform request, error: %s", err)
	}

	// The client must close the response body when finished with it
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Errorf("Failed to close response body, error: %s", err)
		}
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, response.StatusCode, fmt.Errorf("failed to read response body, error: %s", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode > http.StatusMultipleChoices {
		return body, response.StatusCode, errors.New("non success status code")
	}

	return body, response.StatusCode, nil
}

func logDebugPretty(v interface{}) {
	indentedBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	log.Debugf("Response: %+v\n", string(indentedBytes))
}
