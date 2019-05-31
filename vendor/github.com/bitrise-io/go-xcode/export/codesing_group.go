package export

import (
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/profileutil"
)

// CodeSignGroup ...
type CodeSignGroup interface {
	Certificate() certificateutil.CertificateInfoModel
	InstallerCertificate() *certificateutil.CertificateInfoModel
	BundleIDProfileMap() map[string]profileutil.ProvisioningProfileInfoModel
}
