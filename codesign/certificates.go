package codesign

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/pkg/errors"
)

type certificateType uint8

// CertificateTypes ...
const (
	IOSCertificate certificateType = iota
	MacOSCertificate
	MacOSInstallerCertificate
)

var (
	iOSCertificateNames = []string{
		"iPhone Developer",    // type: "iOS Development"
		"iPhone Distribution", // type: "iOS Distribution"
		"Apple Development",   // type: "Apple Development"
		"Apple Distribution",  // type: "Apple Distribution"
	}

	macOSCertificateNames = []string{
		"Mac Developer",                       // type: "Mac Development"
		"3rd Party Mac Developer Application", // type: "Mac App Distribution"
		"Developer ID Application",            // type: "Developer ID Application"
		"Apple Development",                   // type: "Apple Development"
		"Apple Distribution",                  // type: "Apple Distribution"
	}

	macOSInstallerCertificateNames = []string{
		"3rd Party Mac Developer Installer", // type: "Mac Installer Distribution"
		"Developer ID Installer",            // type: "Developer ID Installer"
	}
)

// InstalledCertificates returns the certificate installed in the keychain,
// the expired certificates are removed from the list
func InstalledCertificates(certType certificateType) ([]certificateutil.CertificateInfoModel, error) {
	var certs []certificateutil.CertificateInfoModel
	var err error

	if certType == MacOSInstallerCertificate {
		certs, err = certificateutil.InstalledInstallerCertificateInfos()
	} else {
		certs, err = certificateutil.InstalledCodesigningCertificateInfos()
		if err == nil {
			certs = certificateutil.FilterCertificateInfoModelsByFilterFunc(certs, func(cert certificateutil.CertificateInfoModel) bool {
				var certNames []string
				if certType == IOSCertificate {
					certNames = iOSCertificateNames
				} else {
					certNames = macOSCertificateNames
				}

				for _, name := range certNames {
					if strings.Contains(strings.ToLower(cert.CommonName), strings.ToLower(name)) {
						return true
					}
				}
				return false
			})
		}
	}

	return certificateutil.FilterValidCertificateInfos(certs).ValidCertificates, nil
}

// IsDistributionCertificate returns true if the given certificate
// is an iOS Distribution, Mac App Distribution or Developer ID Application certificate
func IsDistributionCertificate(cert certificateutil.CertificateInfoModel) bool {
	if strings.Contains(strings.ToLower(cert.CommonName), strings.ToLower("iPhone Distribution")) {
		return true
	}

	if strings.Contains(strings.ToLower(cert.CommonName), strings.ToLower("3rd Party Mac Developer Application")) {
		return true
	}

	return false
}

// IsInstallerCertificate returns true if the given certificate
// is an installer certificate
func IsInstallerCertificate(cert certificateutil.CertificateInfoModel) bool {
	if strings.Contains(strings.ToLower(cert.CommonName), strings.ToLower("installer")) {
		return true
	}

	return false
}

// MapCertificatesByTeam returns a certificate list mapped by the certificate's team (in teamdID - teamName format)
func MapCertificatesByTeam(certificates []certificateutil.CertificateInfoModel) map[string][]certificateutil.CertificateInfoModel {
	certificatesByTeam := map[string][]certificateutil.CertificateInfoModel{}
	for _, certificateInfo := range certificates {
		team := fmt.Sprintf("%s - %s", certificateInfo.TeamID, certificateInfo.TeamName)
		certs, ok := certificatesByTeam[team]
		if !ok {
			certs = []certificateutil.CertificateInfoModel{}
		}
		certs = append(certs, certificateInfo)
		certificatesByTeam[team] = certs
	}
	return certificatesByTeam
}

// FindCertificate returns the first certificate, which's common_name or SHA1 fingerprint matches to the given string
func FindCertificate(nameOrSHA1Fingerprint string, certificates []certificateutil.CertificateInfoModel) (certificateutil.CertificateInfoModel, error) {
	for _, certificate := range certificates {
		if certificate.CommonName == nameOrSHA1Fingerprint {
			return certificate, nil
		}
		if strings.ToLower(certificate.SHA1Fingerprint) == strings.ToLower(nameOrSHA1Fingerprint) {
			return certificate, nil
		}
	}
	return certificateutil.CertificateInfoModel{}, errors.Errorf("installed certificate not found with common name or sha1 hash: %s", nameOrSHA1Fingerprint)
}
