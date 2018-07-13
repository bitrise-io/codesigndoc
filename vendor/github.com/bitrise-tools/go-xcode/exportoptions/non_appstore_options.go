package exportoptions

import (
	"fmt"

	"howett.net/plist"
)

// NonAppStoreOptionsModel ...
type NonAppStoreOptionsModel struct {
	Method                             Method
	TeamID                             string
	BundleIDProvisioningProfileMapping map[string]string
	SigningCertificate                 string
	SigningStyle                       string
	ICloudContainerEnvironment         ICloudContainerEnvironment

	// for non app-store exports
	CompileBitcode                           bool
	EmbedOnDemandResourcesAssetPacksInBundle bool
	Manifest                                 Manifest
	OnDemandResourcesAssetPacksBaseURL       string
	Thinning                                 string
}

// NewNonAppStoreOptions ...
func NewNonAppStoreOptions(method Method) NonAppStoreOptionsModel {
	return NonAppStoreOptionsModel{
		Method:                                   method,
		CompileBitcode:                           CompileBitcodeDefault,
		EmbedOnDemandResourcesAssetPacksInBundle: EmbedOnDemandResourcesAssetPacksInBundleDefault,
		Thinning: ThinningDefault,
	}
}

// Hash ...
func (options NonAppStoreOptionsModel) Hash() map[string]interface{} {
	hash := map[string]interface{}{}
	if options.Method != "" {
		hash[MethodKey] = options.Method
	}
	if options.TeamID != "" {
		hash[TeamIDKey] = options.TeamID
	}
	if options.CompileBitcode != CompileBitcodeDefault {
		hash[CompileBitcodeKey] = options.CompileBitcode
	}
	if options.EmbedOnDemandResourcesAssetPacksInBundle != EmbedOnDemandResourcesAssetPacksInBundleDefault {
		hash[EmbedOnDemandResourcesAssetPacksInBundleKey] = options.EmbedOnDemandResourcesAssetPacksInBundle
	}
	if options.ICloudContainerEnvironment != "" {
		hash[ICloudContainerEnvironmentKey] = options.ICloudContainerEnvironment
	}
	if !options.Manifest.IsEmpty() {
		hash[ManifestKey] = options.Manifest.ToHash()
	}
	if options.OnDemandResourcesAssetPacksBaseURL != "" {
		hash[OnDemandResourcesAssetPacksBaseURLKey] = options.OnDemandResourcesAssetPacksBaseURL
	}
	if options.Thinning != ThinningDefault {
		hash[ThinningKey] = options.Thinning
	}
	if len(options.BundleIDProvisioningProfileMapping) > 0 {
		hash[ProvisioningProfilesKey] = options.BundleIDProvisioningProfileMapping
	}
	if options.SigningCertificate != "" {
		hash[SigningCertificateKey] = options.SigningCertificate
	}
	if options.SigningStyle != "" {
		hash[SigningStyleKey] = options.SigningStyle
	}
	return hash
}

// String ...
func (options NonAppStoreOptionsModel) String() (string, error) {
	hash := options.Hash()
	plistBytes, err := plist.MarshalIndent(hash, plist.XMLFormat, "\t")
	if err != nil {
		return "", fmt.Errorf("failed to marshal export options model, error: %s", err)
	}
	return string(plistBytes), err
}

// WriteToFile ...
func (options NonAppStoreOptionsModel) WriteToFile(pth string) error {
	return WritePlistToFile(options.Hash(), pth)
}

// WriteToTmpFile ...
func (options NonAppStoreOptionsModel) WriteToTmpFile() (string, error) {
	return WritePlistToTmpFile(options.Hash())
}
