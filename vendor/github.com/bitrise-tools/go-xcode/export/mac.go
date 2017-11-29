package export

import (
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/exportoptions"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

// MacCodeSignGroup ...
type MacCodeSignGroup struct {
	Certificate          certificateutil.CertificateInfoModel
	InstallerCertificate *certificateutil.CertificateInfoModel
	BundleIDProfileMap   map[string]profileutil.ProvisioningProfileInfoModel
}

// CreateMacCodeSignGroup ...
func CreateMacCodeSignGroup(selectableGroups []SelectableCodeSignGroup, installedInstallerCertificates []certificateutil.CertificateInfoModel, exportMethod exportoptions.Method) []MacCodeSignGroup {
	macosCodeSignGroups := []MacCodeSignGroup{}

	iosCodesignGroups := CreateIosCodeSignGroups(selectableGroups)

	for _, group := range iosCodesignGroups {
		if exportMethod == exportoptions.MethodAppStore {
			installerCertificates := []certificateutil.CertificateInfoModel{}

			for _, installerCertificate := range installedInstallerCertificates {
				if installerCertificate.TeamID == group.Certificate.TeamID {
					installerCertificates = append(installerCertificates, installerCertificate)
				}
			}

			if len(installerCertificates) > 0 {
				installerCertificate := installerCertificates[0]
				macosCodeSignGroups = append(macosCodeSignGroups, MacCodeSignGroup{
					Certificate:          group.Certificate,
					InstallerCertificate: &installerCertificate,
					BundleIDProfileMap:   group.BundleIDProfileMap,
				})
			}
		} else {
			macosCodeSignGroups = append(macosCodeSignGroups, MacCodeSignGroup{
				Certificate:        group.Certificate,
				BundleIDProfileMap: group.BundleIDProfileMap,
			})
		}
	}

	return macosCodeSignGroups
}
