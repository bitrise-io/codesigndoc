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
)

const (
	baseURL                      = "https://api.bitrise.io/v0.1/"
	appsEndPoint                 = "/apps"
	provisioningProfilesEndPoint = "/provisioning-profiles"
	certificatesEndPoint         = "/build-certificates"
)

// Paging ...
type Paging struct {
	TotalItemCount int    `json:"total_item_count"`
	PageItemLimit  int    `json:"page_item_limit"`
	Next           string `json:"next"`
}

// Owner ...
type Owner struct {
	AccountType string `json:"account_type"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
}

// Application ...
type Application struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	ProjectType string `json:"project_type"`
	Provider    string `json:"provider"`
	RepoOwner   string `json:"repo_owner"`
	RepoURL     string `json:"repo_url"`
	RepoSlug    string `json:"repo_slug"`
	IsDisabled  bool   `json:"is_disabled"`
	Status      int    `json:"status"`
	IsPublic    bool   `json:"is_public"`
	Owner       Owner  `json:"owner"`
}

// FetchMyAppsResponse ...
type FetchMyAppsResponse struct {
	Data   []Application `json:"data"`
	Paging Paging        `json:"paging"`
}

// BitriseClient ...
type BitriseClient struct {
	accessToken     string
	selectedAppSlug string
	headers         map[string]string
}

// NewBitriseClient ...
func NewBitriseClient(accessToken string) (client *BitriseClient, apps []Application, err error) {
	client = &BitriseClient{accessToken, "", map[string]string{"Authorization": "token " + accessToken}}

	log.Infof("Asking your application list from Bitrise...")

	requestURL, cerr := urlutil.Join(baseURL, appsEndPoint)
	if cerr != nil {
		err = cerr
		return
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	// Response struct
	var appListResponse FetchMyAppsResponse
	responseStatusCode := -1

	stillPaging := true
	var next string
	for stillPaging {

		headers := client.headers

		request, cerr := createRequest(http.MethodGet, requestURL, headers, nil)
		if cerr != nil {
			err = cerr
			return
		}

		if len(next) > 0 {
			quearryValues := request.URL.Query()
			quearryValues.Add("next", next)
			request.URL.RawQuery = quearryValues.Encode()
		}

		//
		// Perform request
		if cerr := retry.Times(1).Wait(5 * time.Second).Try(func(attempt uint) error {
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

		}); cerr != nil {
			err = cerr
			return
		}

		logDebugPretty(appListResponse)

		apps = append(apps, appListResponse.Data...)

		if len(appListResponse.Paging.Next) > 0 {
			next = appListResponse.Paging.Next
			appListResponse = FetchMyAppsResponse{}
		} else {
			stillPaging = false
		}

	}

	return
}

// SetSelectedAppSlug ...
func (client *BitriseClient) SetSelectedAppSlug(slug string) {
	client.selectedAppSlug = slug
}

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
