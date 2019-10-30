package certificateutil

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/pkcs12"
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

	Certificate x509.Certificate
	PrivateKey  interface{}
}

// String ...
func (info CertificateInfoModel) String() string {
	team := fmt.Sprintf("%s (%s)", info.TeamName, info.TeamID)
	certInfo := fmt.Sprintf("Serial: %s, Name: %s, Team: %s, Expiry: %s", info.Serial, info.CommonName, team, info.EndDate)

	err := info.CheckValidity()
	if err != nil {
		certInfo = certInfo + fmt.Sprintf(", error: %s", err)
	}

	return certInfo
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
	return CheckValidity(info.Certificate)
}

// EncodeToP12 encodes a CertificateInfoModel in pkcs12 (.p12) format.
func (info CertificateInfoModel) EncodeToP12(passphrase string) ([]byte, error) {
	return pkcs12.Encode(rand.Reader, info.PrivateKey, &info.Certificate, nil, passphrase)
}

// NewCertificateInfo ...
func NewCertificateInfo(certificate x509.Certificate, privateKey interface{}) CertificateInfoModel {
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

		Certificate: certificate,
		PrivateKey:  privateKey,
	}
}

// InstalledCodesigningCertificateInfos ...
func InstalledCodesigningCertificateInfos() ([]CertificateInfoModel, error) {
	certificates, err := InstalledCodesigningCertificates()
	if err != nil {
		return nil, err
	}

	infos := []CertificateInfoModel{}
	for _, certificate := range certificates {
		if certificate != nil {
			infos = append(infos, NewCertificateInfo(*certificate, nil))
		}
	}

	return infos, nil
}

// InstalledInstallerCertificateInfos ...
func InstalledInstallerCertificateInfos() ([]CertificateInfoModel, error) {
	certificates, err := InstalledMacAppStoreCertificates()
	if err != nil {
		return nil, err
	}

	infos := []CertificateInfoModel{}
	for _, certificate := range certificates {
		if certificate != nil {
			infos = append(infos, NewCertificateInfo(*certificate, nil))
		}
	}

	installerCertificates := FilterCertificateInfoModelsByFilterFunc(infos, func(cert CertificateInfoModel) bool {
		return strings.Contains(cert.CommonName, "Installer")
	})

	return installerCertificates, nil
}
