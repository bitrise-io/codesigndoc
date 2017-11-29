package certificateutil

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
)

// CertificateInfoModel ...
type CertificateInfoModel struct {
	CommonName string
	TeamName   string
	TeamID     string
	EndDate    time.Time
	StartDate  time.Time

	Serial          string
	SHA1Fingerprint string

	certificate x509.Certificate
}

// String ...
func (info CertificateInfoModel) String() string {
	printable := map[string]interface{}{}
	printable["name"] = info.CommonName
	printable["serial"] = info.Serial
	printable["team"] = fmt.Sprintf("%s (%s)", info.TeamName, info.TeamID)
	printable["expire"] = info.EndDate.String()

	errors := []string{}
	if err := info.CheckValidity(); err != nil {
		errors = append(errors, err.Error())
	}
	if len(errors) > 0 {
		printable["errors"] = errors
	}

	data, err := json.MarshalIndent(printable, "", "\t")
	if err != nil {
		log.Errorf("Failed to marshal: %v, error: %s", printable, err)
		return ""
	}

	return string(data)
}

// CheckValidity ...
func CheckValidity(certificate x509.Certificate) error {
	timeNow := time.Now()
	if !timeNow.After(certificate.NotBefore) {
		return fmt.Errorf("Certificate is not yet valid - validity starts at: %s", certificate.NotBefore)
	}
	if !timeNow.Before(certificate.NotAfter) {
		return fmt.Errorf("Certificate is not valid anymore - validity ended at: %s", certificate.NotAfter)
	}
	return nil
}

// CheckValidity ...
func (info CertificateInfoModel) CheckValidity() error {
	return CheckValidity(info.certificate)
}

// NewCertificateInfo ...
func NewCertificateInfo(certificate x509.Certificate) CertificateInfoModel {
	fingerprint := sha1.Sum(certificate.Raw)
	fingerprintStr := fmt.Sprintf("%x", fingerprint)

	return CertificateInfoModel{
		CommonName:      certificate.Subject.CommonName,
		TeamName:        strings.Join(certificate.Subject.Organization, " "),
		TeamID:          strings.Join(certificate.Subject.OrganizationalUnit, " "),
		EndDate:         certificate.NotAfter,
		StartDate:       certificate.NotBefore,
		Serial:          certificate.SerialNumber.String(),
		SHA1Fingerprint: fingerprintStr,
		certificate:     certificate,
	}
}

// CertificateInfos ...
func CertificateInfos(certificates []*x509.Certificate) []CertificateInfoModel {
	infos := []CertificateInfoModel{}
	for _, certificate := range certificates {
		if certificate != nil {
			info := NewCertificateInfo(*certificate)
			infos = append(infos, info)
		}
	}

	return infos
}

// NewCertificateInfosFromPKCS12 ...
func NewCertificateInfosFromPKCS12(pkcs12Pth, password string) ([]CertificateInfoModel, error) {
	certificates, err := CertificatesFromPKCS12File(pkcs12Pth, password)
	if err != nil {
		return nil, err
	}
	return CertificateInfos(certificates), nil
}

// InstalledCodesigningCertificateInfos ...
func InstalledCodesigningCertificateInfos() ([]CertificateInfoModel, error) {
	certificates, err := InstalledCodesigningCertificates()
	if err != nil {
		return nil, err
	}
	return CertificateInfos(certificates), nil
}

// InstalledInstallerCertificateInfos ...
func InstalledInstallerCertificateInfos() ([]CertificateInfoModel, error) {
	certificates, err := InstalledMacAppStoreCertificates()
	if err != nil {
		return nil, err
	}

	installerCertificates := FilterCertificateInfoModelsByFilterFunc(CertificateInfos(certificates), func(cert CertificateInfoModel) bool {
		return strings.Contains(cert.CommonName, "Installer")
	})

	return installerCertificates, nil
}
