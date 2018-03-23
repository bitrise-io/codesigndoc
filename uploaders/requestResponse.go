package uploaders

// Owner ...
type Owner struct {
	AccountType string `json:"account_type"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
}

// Appliocation ...
type Appliocation struct {
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
	Data []Appliocation `json:"data"`
}

// RegisterProvisioningProfileResponseData ...
type RegisterProvisioningProfileResponseData struct {
	Data           []Appliocation `json:"data"`
	UploadFileName string         `json:"upload_file_name"`
	UploadFileSize int64          `json:"upload_file_size"`
	Slug           string         `json:"slug"`
	Processed      bool           `json:"processed"`
	IsExpose       bool           `json:"is_expose"`
	IsProtected    bool           `json:"is_protected"`
	UploadURL      string         `json:"upload_url"`
}

// RegisterProvisioningProfileResponse ...
type RegisterProvisioningProfileResponse struct {
	Data RegisterProvisioningProfileResponseData `json:"data"`
}
