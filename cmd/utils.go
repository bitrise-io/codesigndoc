package cmd

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/export"
	"github.com/bitrise-tools/go-xcode/profileutil"
	"github.com/bitrise-tools/go-xcode/xcarchive"
	"github.com/pkg/errors"
)

func extractIOSCertificatesAndProfiles(codeSignGroups ...export.IosCodeSignGroup) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel) {
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

func extractMacOSCertificatesAndProfiles(codeSignGroups ...export.MacCodeSignGroup) ([]certificateutil.CertificateInfoModel, []profileutil.ProvisioningProfileInfoModel) {
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

func exportIOSMethod(group export.IosCodeSignGroup) string {
	for _, profile := range group.BundleIDProfileMap {
		return string(profile.ExportType)
	}
	return ""
}

func exportMacOSMethod(group export.MacCodeSignGroup) string {
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

func getIOSCodeSignGroup(archivePath string, installedCertificates []certificateutil.CertificateInfoModel) (xcarchive.IosArchive, export.IosCodeSignGroup, error) {
	archive, err := xcarchive.NewIosArchive(archivePath)
	if err != nil {
		return xcarchive.IosArchive{}, export.IosCodeSignGroup{}, fmt.Errorf("failed to analyze archive, error: %s", err)
	}

	certificate, bundleIDProfileInfoMap, err := analyzeArchive(archive.SigningIdentity(), archive.BundleIDProfileInfoMap(), installedCertificates)
	archiveCodeSignGroup := export.IosCodeSignGroup{
		Certificate:        certificate,
		BundleIDProfileMap: bundleIDProfileInfoMap,
	}
	if err != nil {
		return xcarchive.IosArchive{}, export.IosCodeSignGroup{}, fmt.Errorf("failed to analyze the archive, error: %s", err)
	}

	fmt.Println()
	log.Infof("Codesign settings used for archive:")
	fmt.Println()
	printIOSCodesignGroup(archiveCodeSignGroup)

	return archive, archiveCodeSignGroup, nil
}

func getMacOSCodeSignGroup(archivePath string, installedCertificates []certificateutil.CertificateInfoModel) (xcarchive.MacosArchive, export.MacCodeSignGroup, error) {
	archive, err := xcarchive.NewMacosArchive(archivePath)
	if err != nil {
		return xcarchive.MacosArchive{}, export.MacCodeSignGroup{}, fmt.Errorf("failed to analyze archive, error: %s", err)
	}

	certificate, bundleIDProfileInfoMap, err := analyzeArchive(archive.SigningIdentity(), archive.BundleIDProfileInfoMap(), installedCertificates)
	archiveCodeSignGroup := export.MacCodeSignGroup{
		Certificate:        certificate,
		BundleIDProfileMap: bundleIDProfileInfoMap,
	}
	if err != nil {
		return xcarchive.MacosArchive{}, export.MacCodeSignGroup{}, fmt.Errorf("failed to analyze the archive, error: %s", err)
	}

	fmt.Println()
	log.Infof("Codesign settings used for archive:")
	fmt.Println()
	printMacOsCodesignGroup(archiveCodeSignGroup)

	return archive, archiveCodeSignGroup, nil
}

func collectIpaExportSelectableCodeSignGroups(archive Archive, installedCertificates []certificateutil.CertificateInfoModel, installedProfiles []profileutil.ProvisioningProfileInfoModel) []export.SelectableCodeSignGroup {
	bundleIDEntitlemenstMap := archive.BundleIDEntitlementsMap()

	fmt.Println()
	fmt.Println()
	log.Infof("Targets to sign:")
	fmt.Println()
	for bundleID, entitlements := range bundleIDEntitlemenstMap {
		fmt.Printf("- %s with %d capabilities\n", bundleID, len(entitlements))
	}
	fmt.Println()

	bundleIDs := []string{}
	for bundleID := range bundleIDEntitlemenstMap {
		bundleIDs = append(bundleIDs, bundleID)
	}
	codeSignGroups := export.CreateSelectableCodeSignGroups(installedCertificates, installedProfiles, bundleIDs)

	log.Debugf("Codesign Groups:")
	for _, group := range codeSignGroups {
		log.Debugf(group.String())
	}

	if len(codeSignGroups) == 0 {
		return []export.SelectableCodeSignGroup{}
	}

	codeSignGroups = export.FilterSelectableCodeSignGroups(codeSignGroups,
		export.CreateEntitlementsSelectableCodeSignGroupFilter(bundleIDEntitlemenstMap),
	)

	// Handle if archive used NON xcode managed profile
	if len(codeSignGroups) > 0 && !archive.IsXcodeManaged() {
		log.Warnf("App was signed with NON xcode managed profile when archiving,")
		log.Warnf("only NOT xcode managed profiles are allowed to sign when exporting the archive.")
		log.Warnf("Removing xcode managed CodeSignInfo groups")

		codeSignGroups = export.FilterSelectableCodeSignGroups(codeSignGroups,
			export.CreateNotXcodeManagedSelectableCodeSignGroupFilter(),
		)
	}

	log.Debugf("\n")
	log.Debugf("Filtered Codesign Groups:")
	for _, group := range codeSignGroups {
		log.Debugf(group.String())
	}

	return codeSignGroups
}

func analyzeArchive(signingIdentity string, bundleIDProfileInfoMap map[string]profileutil.ProvisioningProfileInfoModel, installedCertificates []certificateutil.CertificateInfoModel) (certificateutil.CertificateInfoModel, map[string]profileutil.ProvisioningProfileInfoModel, error) {
	if signingIdentity == "" {
		return certificateutil.CertificateInfoModel{}, nil, fmt.Errorf("no signing identity found")
	}

	certificate, err := findCertificate(signingIdentity, installedCertificates)
	if err != nil {
		return certificateutil.CertificateInfoModel{}, nil, err
	}

	return certificate, bundleIDProfileInfoMap, nil
}
