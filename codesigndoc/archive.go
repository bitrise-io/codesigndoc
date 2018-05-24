package codesigndoc

import (
	"fmt"

	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/export"
	"github.com/bitrise-tools/go-xcode/xcarchive"
)

// analyzeArchive opens the generated archive and returns a codesign group, which holds the archive signing options
func analyzeArchive(archive xcarchive.IosArchive, installedCertificates []certificateutil.CertificateInfoModel) (export.IosCodeSignGroup, error) {
	signingIdentity := archive.SigningIdentity()
	bundleIDProfileInfoMap := archive.BundleIDProfileInfoMap()

	if signingIdentity == "" {
		return export.IosCodeSignGroup{}, fmt.Errorf("no signing identity found")
	}

	certificate, err := findCertificate(signingIdentity, installedCertificates)
	if err != nil {
		return export.IosCodeSignGroup{}, err
	}

	return export.IosCodeSignGroup{
		Certificate:        certificate,
		BundleIDProfileMap: bundleIDProfileInfoMap,
	}, nil
}
