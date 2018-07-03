package export

import (
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

// CodeSignGroup ...
type CodeSignGroup interface {
	Certificate() certificateutil.CertificateInfoModel
	InstallerCertificate() *certificateutil.CertificateInfoModel
	BundleIDProfileMap() map[string]profileutil.ProvisioningProfileInfoModel
}
