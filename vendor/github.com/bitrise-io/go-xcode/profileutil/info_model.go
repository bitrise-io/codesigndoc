package profileutil

import (
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/exportoptions"
	"github.com/bitrise-io/go-xcode/plistutil"
	"github.com/fullsailor/pkcs7"
	"howett.net/plist"
)

// ProvisioningProfileInfoModel ...
type ProvisioningProfileInfoModel struct {
	UUID                  string
	Name                  string
	TeamName              string
	TeamID                string
	BundleID              string
	ExportType            exportoptions.Method
	ProvisionedDevices    []string
	DeveloperCertificates []certificateutil.CertificateInfoModel
	CreationDate          time.Time
	ExpirationDate        time.Time
	Entitlements          plistutil.PlistData
	ProvisionsAllDevices  bool
	Type                  ProfileType
}

// PrintableProvisioningProfileInfo ...
func (info ProvisioningProfileInfoModel) String(installedCertificates ...certificateutil.CertificateInfoModel) string {
	printable := map[string]interface{}{}
	printable["name"] = fmt.Sprintf("%s (%s)", info.Name, info.UUID)
	printable["export_type"] = string(info.ExportType)
	printable["team"] = fmt.Sprintf("%s (%s)", info.TeamName, info.TeamID)
	printable["bundle_id"] = info.BundleID
	printable["expire"] = info.ExpirationDate.String()
	printable["is_xcode_managed"] = info.IsXcodeManaged()
	if info.ProvisionedDevices != nil {
		printable["devices"] = info.ProvisionedDevices
	}

	certificates := []map[string]interface{}{}
	for _, certificateInfo := range info.DeveloperCertificates {
		certificate := map[string]interface{}{}
		certificate["name"] = certificateInfo.CommonName
		certificate["serial"] = certificateInfo.Serial
		certificate["team_id"] = certificateInfo.TeamID
		certificates = append(certificates, certificate)
	}
	printable["certificates"] = certificates

	errors := []string{}
	if installedCertificates != nil && !info.HasInstalledCertificate(installedCertificates) {
		errors = append(errors, "none of the profile's certificates are installed")
	}
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

// IsXcodeManaged ...
func IsXcodeManaged(profileName string) bool {
	if strings.HasPrefix(profileName, "XC") {
		return true
	}
	if strings.HasPrefix(profileName, "iOS Team") && strings.Contains(profileName, "Provisioning Profile") {
		return true
	}
	if strings.HasPrefix(profileName, "tvOS Team") && strings.Contains(profileName, "Provisioning Profile") {
		return true
	}
	if strings.HasPrefix(profileName, "Mac Team") && strings.Contains(profileName, "Provisioning Profile") {
		return true
	}
	return false
}

// IsXcodeManaged ...
func (info ProvisioningProfileInfoModel) IsXcodeManaged() bool {
	return IsXcodeManaged(info.Name)
}

// CheckValidity ...
func (info ProvisioningProfileInfoModel) CheckValidity() error {
	timeNow := time.Now()
	if !timeNow.Before(info.ExpirationDate) {
		return fmt.Errorf("Provisioning Profile is not valid anymore - validity ended at: %s", info.ExpirationDate)
	}
	return nil
}

// HasInstalledCertificate ...
func (info ProvisioningProfileInfoModel) HasInstalledCertificate(installedCertificates []certificateutil.CertificateInfoModel) bool {
	has := false
	for _, certificate := range info.DeveloperCertificates {
		for _, installedCertificate := range installedCertificates {
			if certificate.Serial == installedCertificate.Serial {
				has = true
				break
			}
		}
	}
	return has
}

// NewProvisioningProfileInfo ...
func NewProvisioningProfileInfo(provisioningProfile pkcs7.PKCS7) (ProvisioningProfileInfoModel, error) {
	var data plistutil.PlistData
	if _, err := plist.Unmarshal(provisioningProfile.Content, &data); err != nil {
		return ProvisioningProfileInfoModel{}, err
	}

	platform, _ := data.GetStringArray("Platform")
	profileType := ProfileTypeMacOs
	if len(platform) != 0 {
		if strings.ToLower(platform[0]) == string(ProfileTypeIos) {
			profileType = ProfileTypeIos
		}
	}

	profile := PlistData(data)
	info := ProvisioningProfileInfoModel{
		UUID:                 profile.GetUUID(),
		Name:                 profile.GetName(),
		TeamName:             profile.GetTeamName(),
		TeamID:               profile.GetTeamID(),
		BundleID:             profile.GetBundleIdentifier(),
		CreationDate:         profile.GetCreationDate(),
		ExpirationDate:       profile.GetExpirationDate(),
		ProvisionsAllDevices: profile.GetProvisionsAllDevices(),
		Type:                 profileType,
	}

	info.ExportType = profile.GetExportMethod()

	if devicesList := profile.GetProvisionedDevices(); devicesList != nil {
		info.ProvisionedDevices = devicesList
	}

	developerCertificates, found := data.GetByteArrayArray("DeveloperCertificates")
	if found {
		certificates := []*x509.Certificate{}
		for _, certificateBytes := range developerCertificates {
			certificate, err := certificateutil.CertificateFromDERContent(certificateBytes)
			if err == nil && certificate != nil {
				certificates = append(certificates, certificate)
			}
		}
		info.DeveloperCertificates = certificateutil.CertificateInfos(certificates)
	}

	info.Entitlements = profile.GetEntitlements()

	return info, nil
}

// NewProvisioningProfileInfoFromFile ...
func NewProvisioningProfileInfoFromFile(pth string) (ProvisioningProfileInfoModel, error) {
	provisioningProfile, err := ProvisioningProfileFromFile(pth)
	if err != nil {
		return ProvisioningProfileInfoModel{}, err
	}
	if provisioningProfile != nil {
		return NewProvisioningProfileInfo(*provisioningProfile)
	}
	return ProvisioningProfileInfoModel{}, errors.New("failed to parse provisioning profile infos")
}

// InstalledProvisioningProfileInfos ...
func InstalledProvisioningProfileInfos(profileType ProfileType) ([]ProvisioningProfileInfoModel, error) {
	provisioningProfiles, err := InstalledProvisioningProfiles(profileType)
	if err != nil {
		return nil, err
	}

	infos := []ProvisioningProfileInfoModel{}
	for _, provisioningProfile := range provisioningProfiles {
		if provisioningProfile != nil {
			info, err := NewProvisioningProfileInfo(*provisioningProfile)
			if err != nil {
				return nil, err
			}
			infos = append(infos, info)
		}
	}
	return infos, nil
}

// FindProvisioningProfileInfo ...
func FindProvisioningProfileInfo(uuid string) (ProvisioningProfileInfoModel, string, error) {
	profile, pth, err := FindProvisioningProfile(uuid)
	if err != nil {
		return ProvisioningProfileInfoModel{}, "", err
	}
	if pth == "" || profile == nil {
		return ProvisioningProfileInfoModel{}, "", nil
	}

	info, err := NewProvisioningProfileInfo(*profile)
	if err != nil {
		return ProvisioningProfileInfoModel{}, "", err
	}
	return info, pth, nil
}
