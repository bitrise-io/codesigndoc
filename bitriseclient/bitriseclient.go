package bitriseclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/urlutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

const (
	baseURL                      string = "https://api.bitrise.io/v0.1/"
	myAppsEndPoint               string = "/me/apps"
	appsEndPoint                 string = "/apps"
	provisioningProfilesEndPoint string = "/provisioning-profiles"
)

type BitriseClient struct {
	accessToken string
}

// New returns all the application of the user on Bitrise
func (client *BitriseClient) New(accessToken string) (apps []Application, err error) {
	client.accessToken = accessToken

	log.Infof("Asking your application list from Bitrise...")

	requestURL, err := urlutil.Join(baseURL, myAppsEndPoint)
	if err != nil {
		return []Application{}, err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	headers := map[string]string{
		"Authorization": "token " + client.accessToken,
	}

	request, err := createRequest("GET", requestURL, headers, map[string]interface{}{})
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
		log.Errorf("Fetching list failed %s", err)
		return []Application{}, err
	}

	// Success
	log.Donef("Request succeeded with status code: %d", responseStatusCode)
	logDebugPretty(appListResponse)

	return appListResponse.Data, nil
}

// RegisterProvisioningProfile ...
func (client *BitriseClient) RegisterProvisioningProfile(appSlug string, provisioningProfSize int64, profile profileutil.ProvisioningProfileInfoModel) (RegisterProvisioningProfileResponseData, error) {
	log.Infof("Register %s on Bitrise...", profile.Name)

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, appSlug, provisioningProfilesEndPoint)
	if err != nil {
		return RegisterProvisioningProfileResponseData{}, err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	fields := map[string]interface{}{
		"upload_file_name": profile.Name,
		"upload_file_size": provisioningProfSize,
	}

	headers := map[string]string{
		"Authorization": "token " + client.accessToken,
	}

	request, err := createRequest(http.MethodPost, requestURL, headers, fields)
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
		log.Errorf("Fetching list failed %s", err)
		return RegisterProvisioningProfileResponseData{}, err
	}

	// Success
	log.Donef("Request succeeded with status code: %d", responseStatusCode)
	logDebugPretty(requestResponse)

	return requestResponse.Data, nil
}

// UploadProvisioningProfile ...
func (client *BitriseClient) UploadProvisioningProfile(uploadURL string, uploadFileName string, outputDirPath string, exportFileName string) error {
	log.Infof("Upload %s to Bitrise...", exportFileName)

	requestURL := uploadURL

	files := map[string]string{
		uploadFileName: (outputDirPath + "/" + exportFileName),
	}

	request, err := createUploadRequest("PUT", requestURL, map[string]string{}, files)
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

func createUploadRequest(requestMethod string, url string, headers map[string]string, files map[string]string) (*http.Request, error) {
	var fContent []byte

	for _, file := range files {
		f, err := os.Open(file)

		if err != nil {
			return nil, err
		}

		fContent, err = ioutil.ReadAll(f)

		if err != nil {
			return nil, err
		}

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

	by := b.Bytes()
	log.Debugf("Request body: %s", string(by))

	req, err := http.NewRequest(requestMethod, url, bytes.NewReader(by))
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

// func failWithMessage(format string, v ...interface{}) {
// 	log.Errorf(format, v...)
// 	os.Exit(1)
// }

func logDebugPretty(v interface{}) {
	indentedBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println()
	log.Debugf("Response: %+v\n", string(indentedBytes))
}
