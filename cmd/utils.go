package cmd

import (
	"fmt"
	"strings"

	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/export"
	"github.com/bitrise-tools/go-xcode/profileutil"
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

func exportMethod(group export.IosCodeSignGroup) string {
	for _, profile := range group.BundleIDProfileMap {
		return string(profile.ExportType)
	}
	return ""
}

func findCertificate(nameOrSHA1Fingerprint string, certificates []certificateutil.CertificateInfoModel) *certificateutil.CertificateInfoModel {
	for _, certificate := range certificates {
		if certificate.CommonName == nameOrSHA1Fingerprint {
			return &certificate
		}
		if strings.ToLower(certificate.SHA1Fingerprint) == strings.ToLower(nameOrSHA1Fingerprint) {
			return &certificate
		}
	}
	return nil
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
