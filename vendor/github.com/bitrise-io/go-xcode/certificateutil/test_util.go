package certificateutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

// GenerateTestCertificate creates a certificate (signed by a self-signed CA cert) for test purposes
func GenerateTestCertificate(serial int64, teamID, teamName, commonName string, expiry time.Time) (*x509.Certificate, *rsa.PrivateKey, error) {
	CAtemplate := &x509.Certificate{
		IsCA:                  true,
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte{1, 2, 3},
		SerialNumber:          big.NewInt(1234),
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"Pear Worldwide Developer Relations"},
			CommonName:   "Pear Worldwide Developer Relations CA",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0),
		// see http://golang.org/pkg/crypto/x509/#KeyUsage
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Self-signed certificate, parent is the template
	CAcertData, err := x509.CreateCertificate(rand.Reader, CAtemplate, CAtemplate, &privatekey.PublicKey, privatekey)
	if err != nil {
		return nil, nil, err
	}
	CAcert, err := x509.ParseCertificate(CAcertData)
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		IsCA:                  true,
		BasicConstraintsValid: true,
		SerialNumber:          big.NewInt(serial),
		Subject: pkix.Name{
			Country:            []string{"US"},
			Organization:       []string{teamName},
			OrganizationalUnit: []string{teamID},
			CommonName:         commonName,
		},
		NotBefore: time.Now(),
		NotAfter:  expiry,
		// see http://golang.org/pkg/crypto/x509/#KeyUsage
		KeyUsage: x509.KeyUsageDigitalSignature,
	}

	// generate private key
	privatekey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	certData, err := x509.CreateCertificate(rand.Reader, template, CAcert, &privatekey.PublicKey, privatekey)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, nil, err
	}

	return cert, privatekey, nil
}
