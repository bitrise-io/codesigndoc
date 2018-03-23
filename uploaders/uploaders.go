package uploaders

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/urlutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

const (
	baseURL                      string = "https://api.bitrise.io/v0.1/"
	myAppsEndPoint               string = "/me/apps"
	appsEndPoint                 string = "/apps"
	provisioningProfilesEndPoint string = "/provisioning-profiles"
)

// AskAccessToken ...
func AskAccessToken() (token string, err error) {
	messageToAsk := "Please copy your personal access token to Bitrise.\n" +
		"(To acquire a Personal Access Token for your user, sign in with that user on bitrise.io, go to your Account Settings page,\nand select the Security tab on the left side.)"

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

// FetchMyApps ...
func FetchMyApps(accessToken string) (apps []Appliocation, err error) {
	log.Infof("Asking your application list from Bitrise...")

	requestURL := createRequestURLForMyApps()

	headers := map[string]string{
		"Authorization": "token " + accessToken,
	}

	request, err := createRequest(requestURL, map[string]interface{}{}, Get, headers)
	if err != nil {
		failWithMessage("Failed to create request, error: %#v", err)
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
		return []Appliocation{}, err
	}

	// Success
	log.Donef("Request succeeded with status code: %d", responseStatusCode)
	logDebugPretty(appListResponse)

	return appListResponse.Data, nil
}

// RegisterProvisioningProfile ...
func RegisterProvisioningProfile(accessToken string, appSlug string, provisioningProfSize int64, profile profileutil.ProvisioningProfileInfoModel) (RegisterProvisioningProfileResponseData, error) {
	log.Infof("Register %s on Bitrise...", profile.Name)

	requestURL := createRequestURLForProvProfSlug(appSlug)

	fields := map[string]interface{}{
		"upload_file_name": profile.Name,
		"upload_file_size": provisioningProfSize,
	}

	headers := map[string]string{
		"Authorization": "token " + accessToken,
	}

	request, err := createRequest(requestURL, fields, Post, headers)
	if err != nil {
		failWithMessage("Failed to create request, error: %#v", err)
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
func UploadProvisioningProfile(responseData RegisterProvisioningProfileResponseData, outputDirPath string, exportFileName string) error {
	log.Infof("Upload %s to Bitrise...", exportFileName)

	requestURL := responseData.UploadURL

	files := map[string]string{
		responseData.UploadFileName: (outputDirPath + "/" + exportFileName),
	}

	request := createUploadRequest(requestURL, files, Put, map[string]string{})

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

func createUploadRequest(url string, files map[string]string, requestMethod httpMethod, headers map[string]string) *http.Request {
	var fContent []byte

	for _, file := range files {
		f, err := os.Open(file)

		if err != nil {
			failWithMessage("Failed to read file, err: %s", err)
		}

		fContent, err = ioutil.ReadAll(f)

		if err != nil {
			failWithMessage("Failed to read file, err: %s", err)
		}

	}

	req, err := http.NewRequest(Put.String(), url, bytes.NewReader(fContent))

	if err != nil {
		failWithMessage("Failed to create upload request, err: %s", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)

	}

	return req
}

func createRequest(url string, fields map[string]interface{}, requestMethod httpMethod, headers map[string]string) (*http.Request, error) {
	b := new(bytes.Buffer)

	if len(fields) > 0 {
		err := json.NewEncoder(b).Encode(fields)
		if err != nil {
			return nil, err
		}
	}

	by := b.Bytes()
	log.Debugf("Request body: %s", string(by))

	req, err := http.NewRequest(requestMethod.String(), url, bytes.NewReader(by))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(os.Stdout, req.Body); err != nil {
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

func createRequestURLForMyApps() string {
	requestURL, err := urlutil.Join(baseURL, myAppsEndPoint)
	if err != nil {
		failWithMessage("failed to create request URL, error: %s", err)
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	return requestURL
}

func createRequestURLForProvProfSlug(appSlug string) string {
	requestURL, err := urlutil.Join(baseURL, appsEndPoint, appSlug, provisioningProfilesEndPoint)
	if err != nil {
		failWithMessage("failed to create request URL, error: %s", err)
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	return requestURL
}

func failWithMessage(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func logDebugPretty(v interface{}) {
	indentedBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println()
	log.Debugf("Response: %+v\n", string(indentedBytes))
}
