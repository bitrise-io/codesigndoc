package bitrise

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

// MyAppsResponse ...
type MyAppsResponse struct {
	Data   []Application `json:"data"`
	Paging Paging        `json:"paging"`
}

// Client ...
type Client struct {
	accessToken     string
	selectedAppSlug string
	headers         map[string]string
	client          http.Client
}

// NewClient ...
func NewClient(accessToken string) (*Client, error) {
	client := &Client{accessToken, "", map[string]string{"Authorization": "token " + accessToken}, http.Client{}}
	return client, nil
}

// GetAppList returns the list of apps for the given access token
func (client *Client) GetAppList() ([]Application, error) {
	var apps []Application

	log.Infof("Fetching your application list from Bitrise...")

	requestURL, err := urlutil.Join(baseURL, appsEndPoint)
	if err != nil {
		return nil, err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	// Response struct
	var appListResponse MyAppsResponse
	stillPaging := true
	var next string

	for stillPaging {
		headers := client.headers

		request, err := createRequest(http.MethodGet, requestURL, headers, nil)
		if err != nil {
			return nil, err
		}

		if len(next) > 0 {
			quearryValues := request.URL.Query()
			quearryValues.Add("next", next)
			request.URL.RawQuery = quearryValues.Encode()
		}

		// Perform request
		response, _, err := RunRequest(client, request, &appListResponse)
		if err != nil {
			return nil, err
		}

		appListResponse = *response.(*MyAppsResponse)
		apps = append(apps, appListResponse.Data...)

		if len(appListResponse.Paging.Next) > 0 {
			next = appListResponse.Paging.Next
			appListResponse = MyAppsResponse{}
		} else {
			stillPaging = false
		}
	}

	return apps, nil
}

// SetSelectedAppSlug ...
func (client *Client) SetSelectedAppSlug(slug string) {
	client.selectedAppSlug = slug
}

// UploadArtifact ...
func (client *Client) UploadArtifact(uploadURL string, content io.Reader) error {
	request, err := http.NewRequest(http.MethodPut, uploadURL, content)
	if err != nil {
		return err
	}

	_, _, err = RunRequest(client, request, nil)
	if err != nil {
		return err
	}

	return nil
}

// RunRequest ...
func RunRequest(client *Client, req *http.Request, requestResponse interface{}) (interface{}, []byte, error) {
	var responseBody []byte

	if err := retry.Times(1).Wait(5 * time.Second).Try(func(attempt uint) error {
		body, statusCode, err := performRequest(client, req)
		if err != nil {
			log.Warnf("Attempt (%d) failed, error: %s", attempt+1, err)
			if !strings.Contains(err.Error(), "failed to perform request") {
				log.Warnf("Response status: %d", statusCode)
				log.Warnf("Body: %s", string(body))
			}
			return err
		}

		// Parse JSON body
		if requestResponse != nil {
			if err := json.Unmarshal([]byte(body), &requestResponse); err != nil {
				return fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
			}

			logDebugPretty(&requestResponse)
		}
		responseBody = body

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return requestResponse, responseBody, nil
}

func createRequest(requestMethod string, url string, headers map[string]string, fields map[string]interface{}) (*http.Request, error) {
	var b bytes.Buffer

	if len(fields) > 0 {
		err := json.NewEncoder(&b).Encode(fields)
		if err != nil {
			return nil, err
		}
	}

	log.Debugf("Request body: %s", string(b.Bytes()))

	req, err := http.NewRequest(requestMethod, url, bytes.NewReader(b.Bytes()))
	if err != nil {
		return nil, err
	}
	addHeaders(req, headers)

	return req, nil
}

func performRequest(bitriseClient *Client, request *http.Request) (body []byte, statusCode int, err error) {
	response, err := bitriseClient.client.Do(request)
	if err != nil {
		// On error, any Response can be ignored
		return nil, -1, fmt.Errorf("failed to perform request, error: %s", err)
	}

	// The client must close the response body when finished with it
	defer func() {
		if cerr := response.Body.Close(); err != nil {
			cerr = fmt.Errorf("Failed to close response body, error: %s", cerr)
		}
	}()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, response.StatusCode, fmt.Errorf("failed to read response body, error: %s", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode > http.StatusMultipleChoices {
		return body, response.StatusCode, errors.New("non success status code")
	}

	return body, response.StatusCode, nil
}

func addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

func logDebugPretty(v interface{}) {
	indentedBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	log.Debugf("Response: %+v\n", string(indentedBytes))
}
