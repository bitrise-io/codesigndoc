package certificateutil

import (
	"bufio"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pkcs12"
	"github.com/pkg/errors"
)

func commandError(printableCmd string, cmdOut string, cmdErr error) error {
	return errors.Wrapf(cmdErr, "%s failed, out: %s", printableCmd, cmdOut)
}

// CertificatesFromPKCS12Content ...
func CertificatesFromPKCS12Content(content []byte, password string) ([]*x509.Certificate, error) {
	certificates, err := pkcs12.DecodeAllCerts(content, password)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}

// CertificatesFromPKCS12File ...
func CertificatesFromPKCS12File(pkcs12Pth, password string) ([]*x509.Certificate, error) {
	content, err := fileutil.ReadBytesFromFile(pkcs12Pth)
	if err != nil {
		return nil, err
	}

	return CertificatesFromPKCS12Content(content, password)
}

// CertificateFromDERContent ...
func CertificateFromDERContent(content []byte) (*x509.Certificate, error) {
	return x509.ParseCertificate(content)
}

// CeritifcateFromPemContent ...
func CeritifcateFromPemContent(content []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(content)
	if block == nil || block.Bytes == nil || len(block.Bytes) == 0 {
		return nil, fmt.Errorf("failed to parse profile from: %s", string(content))
	}
	return CertificateFromDERContent(block.Bytes)
}

func installedCodesigningCertificateNamesFromOutput(out string) ([]string, error) {
	pettern := `^[0-9]+\) (?P<hash>.*) "(?P<name>.*)"`
	re := regexp.MustCompile(pettern)

	certificateNameMap := map[string]bool{}
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := re.FindStringSubmatch(line); len(matches) == 3 {
			name := matches[2]
			certificateNameMap[name] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	names := []string{}
	for name := range certificateNameMap {
		names = append(names, name)
	}
	return names, nil
}

// InstalledCodesigningCertificateNames ...
func InstalledCodesigningCertificateNames() ([]string, error) {
	cmd := command.New("security", "find-identity", "-v", "-p", "codesigning")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, commandError(cmd.PrintableCommandArgs(), out, err)
	}
	return installedCodesigningCertificateNamesFromOutput(out)
}

// InstalledMacAppStoreCertificateNames ...
func InstalledMacAppStoreCertificateNames() ([]string, error) {
	cmd := command.New("security", "find-identity", "-v", "-p", "macappstore")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, commandError(cmd.PrintableCommandArgs(), out, err)
	}
	return installedCodesigningCertificateNamesFromOutput(out)
}

func normalizeFindCertificateOut(out string) ([]string, error) {
	certificateContents := []string{}
	pattern := `(?s)(-----BEGIN CERTIFICATE-----.*?-----END CERTIFICATE-----)`
	matches := regexp.MustCompile(pattern).FindAllString(out, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no certificates found in: %s", out)
	}

	for _, certificateContent := range matches {
		if !strings.HasPrefix(certificateContent, "\n") {
			certificateContent = "\n" + certificateContent
		}
		if !strings.HasSuffix(certificateContent, "\n") {
			certificateContent = certificateContent + "\n"
		}
		certificateContents = append(certificateContents, certificateContent)
	}

	return certificateContents, nil
}

// InstalledCodesigningCertificates ...
func InstalledCodesigningCertificates() ([]*x509.Certificate, error) {
	certificateNames, err := InstalledCodesigningCertificateNames()
	if err != nil {
		return nil, err
	}
	return getInstalledCertificatesByNameSlice(certificateNames)
}

// InstalledMacAppStoreCertificates ...
func InstalledMacAppStoreCertificates() ([]*x509.Certificate, error) {
	certificateNames, err := InstalledMacAppStoreCertificateNames()
	if err != nil {
		return nil, err
	}
	return getInstalledCertificatesByNameSlice(certificateNames)
}

func getInstalledCertificatesByNameSlice(certificateNames []string) ([]*x509.Certificate, error) {
	certificates := []*x509.Certificate{}
	for _, name := range certificateNames {
		cmd := command.New("security", "find-certificate", "-c", name, "-p", "-a")
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		if err != nil {
			return nil, commandError(cmd.PrintableCommandArgs(), out, err)
		}

		normalizedOuts, err := normalizeFindCertificateOut(out)
		if err != nil {
			return nil, err
		}

		for _, normalizedOut := range normalizedOuts {
			certificate, err := CeritifcateFromPemContent([]byte(normalizedOut))
			if err != nil {
				return nil, err
			}

			certificates = append(certificates, certificate)
		}
	}

	return certificates, nil
}
