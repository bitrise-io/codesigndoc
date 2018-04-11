package cmd

import (
	"fmt"
	"strings"

	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/export"
	"github.com/bitrise-tools/go-xcode/profileutil"
	"github.com/pkg/errors"
)

func extractCertificatesAndProfiles(codeSignGroups ...export.IosCodeSignGroup) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel) {
	certificateMap := map[string]certificateutil.CertificateInfoModel{}
	profilesMap := map[string]profileutil.ProvisioningProfileInfoModel{}
	for _, group := range codeSignGroups {
		certificate := group.Certificate

		certificateMap[certificate.Serial] = certificate

		for _, profile := range group.BundleIDProfileMap {
			profilesMap[profile.UUID] = profile
		}
	}

	certificates := []certificateutil.CertificateInfoModel{}
	profiles := []profileutil.ProvisioningProfileInfoModel{}
	for _, certificate := range certificateMap {
		certificates = append(certificates, certificate)
	}
	for _, profile := range profilesMap {
		profiles = append(profiles, profile)
	}
	return certificates, profiles
}

func extractMacOsCertificatesAndProfiles(codeSignGroups ...export.MacCodeSignGroup) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel) {
	certificateMap := map[string]certificateutil.CertificateInfoModel{}
	profilesMap := map[string]profileutil.ProvisioningProfileInfoModel{}
	for _, group := range codeSignGroups {
		certificate := group.Certificate

		certificateMap[certificate.Serial] = certificate

		for _, profile := range group.BundleIDProfileMap {
			profilesMap[profile.UUID] = profile
		}
	}

	certificates := []certificateutil.CertificateInfoModel{}
	profiles := []profileutil.ProvisioningProfileInfoModel{}
	for _, certificate := range certificateMap {
		certificates = append(certificates, certificate)
	}
	for _, profile := range profilesMap {
		profiles = append(profiles, profile)
	}
	return certificates, profiles
}

func exportMethod(group export.IosCodeSignGroup) string {
	for _, profile := range group.BundleIDProfileMap {
		return string(profile.ExportType)
	}
	return ""
}

func exportMacOsMethod(group export.MacCodeSignGroup) string {
	for _, profile := range group.BundleIDProfileMap {
		return string(profile.ExportType)
	}
	return ""
}

func findCertificate(nameOrSHA1Fingerprint string, certificates []certificateutil.CertificateInfoModel) (certificateutil.CertificateInfoModel, error) {
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

func isDistributionCertificate(certificate certificateutil.CertificateInfoModel) bool {
	return strings.HasPrefix(certificate.CommonName, "iPhone Distribution:")
}

func mapCertificatesByTeam(certificates []certificateutil.CertificateInfoModel) map[string][]certificateutil.CertificateInfoModel {
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
