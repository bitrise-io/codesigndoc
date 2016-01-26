package certutil

import (
	"crypto/x509"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	certrevoke "github.com/cloudflare/cfssl/revoke"
)

// CheckCertificateValidity ...
func CheckCertificateValidity(cert *x509.Certificate) error {
	log.Debugf("cert - valid - NotBefore: %#v / %s", cert.NotBefore, cert.NotBefore)
	log.Debugf("cert - valid - NotAfter: %#v / %s", cert.NotAfter, cert.NotAfter)

	timeNow := time.Now()
	log.Debugf("cert - time now: %#v / %s", timeNow, timeNow)
	if !timeNow.After(cert.NotBefore) {
		return fmt.Errorf("Certificate is not yet valid - validity starts at: %s", cert.NotBefore)
	}
	if !timeNow.Before(cert.NotAfter) {
		return fmt.Errorf("Certificate is not valid anymore - validity ended at: %s", cert.NotAfter)
	}
	log.Debugf("Certificate is Valid, based on it's validity date-times")

	revoked, ok := certrevoke.VerifyCertificate(cert)
	log.Debugf("revoked: %#v, ok: %#v", revoked, ok)

	return nil
}
