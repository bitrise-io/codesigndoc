package models

import (
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/profileutil"
)

// Certificates contains all the certificates as io.Reader, besides an array of the certificate infos.
type Certificates struct {
	Info    []certificateutil.CertificateInfoModel
	Content []byte
}

// ProvisioningProfile contains parsed data in the provisioning profile and the original profile file contents.
type ProvisioningProfile struct {
	Info    profileutil.ProvisioningProfileInfoModel
	Content []byte
}
