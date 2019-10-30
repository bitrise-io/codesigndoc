package bitrise

import (
	"math/big"
	"net/http"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/urlutil"
	"github.com/bitrise-io/go-xcode/certificateutil"
)

// RegisterIdentityData ...
type RegisterIdentityData struct {
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
	Data RegisterIdentityData `json:"data"`
}

// ConfirmIdentityUploadData ...
type ConfirmIdentityUploadData struct {
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
	Data ConfirmIdentityUploadData `json:"data"`
}

// IdentityListData ...
type IdentityListData struct {
	UploadFileName      string `json:"upload_file_name"`
	UploadFileSize      int    `json:"upload_file_size"`
	Slug                string `json:"slug"`
	Processed           bool   `json:"processed"`
	CertificatePassword string `json:"certificate_password"`
	IsExpose            bool   `json:"is_expose"`
	IsProtected         bool   `json:"dais_protectedta"`
}

// IdentityListResponse ...
type IdentityListResponse struct {
	Data []IdentityListData `json:"data"`
}

// IdentityData ...
type IdentityData struct {
	UploadFileName      string `json:"upload_file_name"`
	UploadFileSize      int    `json:"upload_file_size"`
	Slug                string `json:"slug"`
	Processed           bool   `json:"processed"`
	CertificatePassword string `json:"certificate_password"`
	IsExpose            bool   `json:"is_expose"`
	IsProtected         bool   `json:"dais_protectedta"`
	DownloadURL         string `json:"download_url"`
}

// IdentityResponse ...
type IdentityResponse struct {
	Data IdentityData `json:"data"`
}

// FetchUploadedIdentities ...
func (client *Client) FetchUploadedIdentities() ([]IdentityListData, error) {
	log.Debugf("\nDownloading provisioning profile list from Bitrise...")

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, certificatesEndPoint)
	if err != nil {
		return []IdentityListData{}, err
	}

	request, err := createRequest(http.MethodGet, requestURL, client.headers, nil)
	if err != nil {
		return []IdentityListData{}, err
	}
	log.Debugf("\nRequest URL: %s", requestURL)

	// Response struct
	var requestResponse IdentityListResponse
	//
	// Perform request
	response, _, err := RunRequest(client, request, &requestResponse)
	if err != nil {
		return nil, err
	}

	requestResponse = *response.(*IdentityListResponse)
	return requestResponse.Data, nil
}

// GetUploadedCertificatesSerialby ...
func (client *Client) GetUploadedCertificatesSerialby(identitySlug string) (certificateSerialList []big.Int, err error) {
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
		var (
			serial big.Int
			base   = 10
		)

		if serialRef, ok := serial.SetString(certificate.Serial, base); ok {
			serialList = append(serialList, *serialRef)
		} else {
			log.Warnf("Error converting serial ID (%s) with base (%d): ", certificate.Serial, base)
		}
	}
	return serialList, nil
}

func (client *Client) getUploadedIdentityDownloadURLBy(certificateSlug string) (downloadURL string, password string, err error) {
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
	var requestResponse IdentityResponse

	//
	// Perform request
	response, _, err := RunRequest(client, request, &requestResponse)
	if err != nil {
		return "", "", err
	}

	requestResponse = *response.(*IdentityResponse)
	return requestResponse.Data.DownloadURL, requestResponse.Data.CertificatePassword, nil
}

func (client *Client) downloadUploadedIdentity(downloadURL string) (content string, err error) {
	log.Debugf("\nDownloading identities from Bitrise...")
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

// RegisterIdentity ...
func (client *Client) RegisterIdentity(certificateSize int64) (RegisterIdentityData, error) {
	log.Printf("Register %s on Bitrise...", "Identities.p12")

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, certificatesEndPoint)
	if err != nil {
		return RegisterIdentityData{}, err
	}

	log.Debugf("\nRequest URL: %s", requestURL)

	fields := map[string]interface{}{
		"upload_file_name": "Identities.p12",
		"upload_file_size": certificateSize,
	}

	request, err := createRequest(http.MethodPost, requestURL, client.headers, fields)
	if err != nil {
		return RegisterIdentityData{}, err
	}

	// Response struct
	var requestResponse RegisterIdentityResponse

	//
	// Perform request
	response, _, err := RunRequest(client, request, &requestResponse)
	if err != nil {
		return RegisterIdentityData{}, err
	}

	requestResponse = *response.(*RegisterIdentityResponse)
	return requestResponse.Data, nil
}

// ConfirmIdentityUpload ...
func (client *Client) ConfirmIdentityUpload(certificateSlug string, certificateUploadName string) error {
	log.Printf("Confirm - %s - upload to Bitrise...", certificateUploadName)

	requestURL, err := urlutil.Join(baseURL, appsEndPoint, client.selectedAppSlug, "build-certificates", certificateSlug, "uploaded")
	if err != nil {
		return err
	}

	request, err := createRequest(http.MethodPost, requestURL, client.headers, nil)
	if err != nil {
		return err
	}

	// Response struct
	requestResponse := ConfirmProvProfileUploadResponse{}

	_, _, err = RunRequest(client, request, &requestResponse)
	if err != nil {
		return err
	}
	return nil
}
