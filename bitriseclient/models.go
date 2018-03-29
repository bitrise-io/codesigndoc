package bitriseclient

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
	Data []Application `json:"data"`
}

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
