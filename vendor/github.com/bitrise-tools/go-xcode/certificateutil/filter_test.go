package certificateutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterCertificateInfoModelsByFilterFunc(t *testing.T) {
	filterableCerts := []CertificateInfoModel{
		CertificateInfoModel{TeamID: "my-team-id"},
		CertificateInfoModel{TeamID: "find-this-team-id"},
		CertificateInfoModel{TeamID: "my--another-team-id"},
		CertificateInfoModel{TeamID: "test-team-id", CommonName: "test common name"},
		CertificateInfoModel{TeamID: "test-team-id2", CommonName: "find this common name"},
	}
	expectedCertsByTeamID := []CertificateInfoModel{
		CertificateInfoModel{TeamID: "find-this-team-id"},
	}

	foundCerts := FilterCertificateInfoModelsByFilterFunc(filterableCerts, func(cert CertificateInfoModel) bool { return cert.TeamID == "find-this-team-id" })
	require.Equal(t, expectedCertsByTeamID, foundCerts)

	expectedCertsByCommonNameExact := []CertificateInfoModel{
		CertificateInfoModel{TeamID: "test-team-id2", CommonName: "find this common name"},
	}

	foundCerts = FilterCertificateInfoModelsByFilterFunc(filterableCerts, func(cert CertificateInfoModel) bool { return cert.CommonName == "find this common name" })
	require.Equal(t, expectedCertsByCommonNameExact, foundCerts)

	expectedCertsByCommonNameMatch := []CertificateInfoModel{
		CertificateInfoModel{TeamID: "test-team-id", CommonName: "test common name"},
		CertificateInfoModel{TeamID: "test-team-id2", CommonName: "find this common name"},
	}

	foundCerts = FilterCertificateInfoModelsByFilterFunc(filterableCerts, func(cert CertificateInfoModel) bool { return strings.Contains(cert.CommonName, "common name") })
	require.Equal(t, expectedCertsByCommonNameMatch, foundCerts)
}
