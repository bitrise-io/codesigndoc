package profileutil

import (
	"strings"
	"time"

	"github.com/bitrise-tools/go-xcode/exportoptions"
	"github.com/bitrise-tools/go-xcode/plistutil"
	"howett.net/plist"
)

const (
	notValidParameterErrorMessage = "security: SecPolicySetValue: One or more parameters passed to a function were not valid."
)

// PlistData ...
type PlistData plistutil.PlistData

// NewPlistDataFromFile ...
func NewPlistDataFromFile(provisioningProfilePth string) (PlistData, error) {
	provisioningProfilePKCS7, err := ProvisioningProfileFromFile(provisioningProfilePth)
	if err != nil {
		return PlistData{}, err
	}

	var plistData plistutil.PlistData
	if _, err := plist.Unmarshal(provisioningProfilePKCS7.Content, &plistData); err != nil {
		return PlistData{}, err
	}

	return PlistData(plistData), nil
}

// GetUUID ...
func (profile PlistData) GetUUID() string {
	data := plistutil.PlistData(profile)
	uuid, _ := data.GetString("UUID")
	return uuid
}

// GetName ...
func (profile PlistData) GetName() string {
	data := plistutil.PlistData(profile)
	uuid, _ := data.GetString("Name")
	return uuid
}

// GetApplicationIdentifier ...
func (profile PlistData) GetApplicationIdentifier() string {
	data := plistutil.PlistData(profile)
	entitlements, ok := data.GetMapStringInterface("Entitlements")
	if !ok {
		return ""
	}

	applicationID, ok := entitlements.GetString("application-identifier")
	if !ok {
		applicationID, ok = entitlements.GetString("com.apple.application-identifier")
		if !ok {
			return ""
		}
	}
	return applicationID
}

// GetBundleIdentifier ...
func (profile PlistData) GetBundleIdentifier() string {
	applicationID := profile.GetApplicationIdentifier()

	plistData := plistutil.PlistData(profile)
	prefixes, found := plistData.GetStringArray("ApplicationIdentifierPrefix")
	if found {
		for _, prefix := range prefixes {
			applicationID = strings.TrimPrefix(applicationID, prefix+".")
		}
	}

	teamID := profile.GetTeamID()
	return strings.TrimPrefix(applicationID, teamID+".")
}

// GetExportMethod ...
func (profile PlistData) GetExportMethod() exportoptions.Method {
	data := plistutil.PlistData(profile)
	entitlements, _ := data.GetMapStringInterface("Entitlements")
	platform, _ := data.GetStringArray("Platform")

	if len(platform) != 0 {
		switch strings.ToLower(platform[0]) {
		case "osx":
			_, ok := data.GetStringArray("ProvisionedDevices")
			if !ok {
				if allDevices, ok := data.GetBool("ProvisionsAllDevices"); ok && allDevices {
					return exportoptions.MethodDeveloperID
				}
				return exportoptions.MethodAppStore
			}
			return exportoptions.MethodDevelopment
		case "ios", "tvos":
			_, ok := data.GetStringArray("ProvisionedDevices")
			if !ok {
				if allDevices, ok := data.GetBool("ProvisionsAllDevices"); ok && allDevices {
					return exportoptions.MethodEnterprise
				}
				return exportoptions.MethodAppStore
			}
			if allow, ok := entitlements.GetBool("get-task-allow"); ok && allow {
				return exportoptions.MethodDevelopment
			}
			return exportoptions.MethodAdHoc
		}
	}

	return exportoptions.MethodDefault
}

// GetEntitlements ...
func (profile PlistData) GetEntitlements() plistutil.PlistData {
	data := plistutil.PlistData(profile)
	entitlements, _ := data.GetMapStringInterface("Entitlements")
	return entitlements
}

// GetTeamID ...
func (profile PlistData) GetTeamID() string {
	data := plistutil.PlistData(profile)
	entitlements, ok := data.GetMapStringInterface("Entitlements")
	if ok {
		teamID, _ := entitlements.GetString("com.apple.developer.team-identifier")
		return teamID
	}
	return ""
}

// GetExpirationDate ...
func (profile PlistData) GetExpirationDate() time.Time {
	data := plistutil.PlistData(profile)
	expiry, _ := data.GetTime("ExpirationDate")
	return expiry
}

// GetProvisionedDevices ...
func (profile PlistData) GetProvisionedDevices() []string {
	data := plistutil.PlistData(profile)
	devices, _ := data.GetStringArray("ProvisionedDevices")
	return devices
}

// GetDeveloperCertificates ...
func (profile PlistData) GetDeveloperCertificates() [][]byte {
	data := plistutil.PlistData(profile)
	developerCertificates, _ := data.GetByteArrayArray("DeveloperCertificates")
	return developerCertificates
}

// GetTeamName ...
func (profile PlistData) GetTeamName() string {
	data := plistutil.PlistData(profile)
	teamName, _ := data.GetString("TeamName")
	return teamName
}

// GetCreationDate ...
func (profile PlistData) GetCreationDate() time.Time {
	data := plistutil.PlistData(profile)
	creationDate, _ := data.GetTime("CreationDate")
	return creationDate
}

// GetProvisionsAllDevices ...
func (profile PlistData) GetProvisionsAllDevices() bool {
	data := plistutil.PlistData(profile)
	provisionsAlldevices, _ := data.GetBool("ProvisionsAllDevices")
	return provisionsAlldevices
}
