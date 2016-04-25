package xcode

import (
	"testing"

	"github.com/bitrise-io/go-utils/testutil"
	"github.com/bitrise-tools/codesigndoc/provprofile"
	"github.com/stretchr/testify/require"
)

func Test_parseSchemesFromXcodeOutput(t *testing.T) {
	xcout := `Information about project "SampleAppWithCocoapods":
    Targets:
        SampleAppWithCocoapods
        SampleAppWithCocoapodsTests

    Build Configurations:
        Debug
        Release

    If no build configuration is specified and -scheme is not passed then "Release" is used.

    Schemes:
        SampleAppWithCocoapods`
	parsedSchemes := parseSchemesFromXcodeOutput(xcout)
	require.Equal(t, []string{"SampleAppWithCocoapods"}, parsedSchemes)
}

func Test_parseCodeSigningSettingsFromXcodeOutput(t *testing.T) {
	t.Log("A single Identity & Prov Profile")
	{
		xcout := `CodeSign /Users/bitrise/Library/...
    cd /Users/...
    export CODESIGN_ALLOCATE=/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/codesign_allocate
    export PATH="/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform...

Signing Identity:     "iPhone Developer: First Last (F72Z82XD37)"
Provisioning Profile: "Prov Profile 42"
                      (87af6d83-cb65-4dbe-aee7-f97a87d6fec1)

    /usr/bin/codesign --force --sign E7D5FA3770F4ECC529CFCF683CBCDF874F7870FB --entitlements /Users/...`

		parsedCodeSigningSettings := parseCodeSigningSettingsFromXcodeOutput(xcout)
		require.Equal(t, []CodeSigningIdentityInfo{
			CodeSigningIdentityInfo{Title: "iPhone Developer: First Last (F72Z82XD37)"},
		}, parsedCodeSigningSettings.Identities)
		require.Equal(t, []provprofile.ProvisioningProfileInfo{
			provprofile.ProvisioningProfileInfo{Title: "Prov Profile 42", UUID: "87af6d83-cb65-4dbe-aee7-f97a87d6fec1"},
		}, parsedCodeSigningSettings.ProvProfiles)
	}

	t.Log("A single Identity & Prov Profile - different style, and not after each other")
	{
		xcout := `CodeSign /Users/bitrise/Library/...
    cd /Users/...
    export CODESIGN_ALLOCATE=/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/codesign_allocate

Signing Identity:     "iPhone Distribution: First Last Company (F72Z82XD37)"

    export PATH="/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform...

Provisioning Profile: "com.domain.app AdHoc"
                      (87af6d83-cb65-4dbe-aee7-f97a87d6fec1)

    /usr/bin/codesign --force --sign E7D5FA3770F4ECC529CFCF683CBCDF874F7870FB --entitlements /Users/...`

		parsedCodeSigningSettings := parseCodeSigningSettingsFromXcodeOutput(xcout)
		require.Equal(t, []CodeSigningIdentityInfo{
			CodeSigningIdentityInfo{Title: "iPhone Distribution: First Last Company (F72Z82XD37)"},
		}, parsedCodeSigningSettings.Identities)
		require.Equal(t, []provprofile.ProvisioningProfileInfo{
			provprofile.ProvisioningProfileInfo{Title: "com.domain.app AdHoc", UUID: "87af6d83-cb65-4dbe-aee7-f97a87d6fec1"},
		}, parsedCodeSigningSettings.ProvProfiles)
	}

	t.Log("A single Identity & Prov Profile - wildcard Prov Profile")
	{
		xcout := `CodeSign /Users/bitrise/Library/...
    cd /Users/...
Signing Identity:     "iPhone Developer: First Last (F72Z82XD37)"
Provisioning Profile: "iOS Team Provisioning Profile: *"
                      (87af6d83-cb65-4dbe-aee7-f97a87d6fec1)

    /usr/bin/codesign --force --sign E7D5FA3770F4ECC529CFCF683CBCDF874F7870FB --entitlements /Users/...`

		parsedCodeSigningSettings := parseCodeSigningSettingsFromXcodeOutput(xcout)
		require.Equal(t, []CodeSigningIdentityInfo{
			CodeSigningIdentityInfo{Title: "iPhone Developer: First Last (F72Z82XD37)"},
		}, parsedCodeSigningSettings.Identities)
		require.Equal(t, []provprofile.ProvisioningProfileInfo{
			provprofile.ProvisioningProfileInfo{Title: "iOS Team Provisioning Profile: *", UUID: "87af6d83-cb65-4dbe-aee7-f97a87d6fec1"},
		}, parsedCodeSigningSettings.ProvProfiles)
	}

	t.Log("Multiple Identity & Prov Profiles")
	{
		xcout := `CodeSign /Users/bitrise/Library/...
    cd /Users/...
    export CODESIGN_ALLOCATE=/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/codesign_allocate
    export PATH="/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform...

Signing Identity:     "iPhone Developer: First Last (F72Z82XD37)"
Provisioning Profile: "Prov Profile 42"
                      (87af6d83-cb65-4dbe-aee7-f97a87d6fec1)

    /usr/bin/codesign --force --sign E7D5FA3770F4ECC529CFCF683CBCDF874F7870FB --entitlements /Users/...
CodeSign /Users/bitrise/Library/...
    cd /Users/...
    export CODESIGN_ALLOCATE=/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/codesign_allocate
    export PATH="/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform...

Signing Identity:     "iPhone Distribution: BFirst BLast (B72Z82XD37)"
Provisioning Profile: "Prov Profile 43"
                      (97af6d83-cb65-4dbe-aee7-f97a87d6fec1)

    /usr/bin/codesign --force --sign E7D5FA3770F4ECC529CFCF683CBCDF874F7870FB --entitlements /Users/...`

		parsedCodeSigningSettings := parseCodeSigningSettingsFromXcodeOutput(xcout)
		testutil.EqualSlicesWithoutOrder(t, []CodeSigningIdentityInfo{
			CodeSigningIdentityInfo{Title: "iPhone Developer: First Last (F72Z82XD37)"},
			CodeSigningIdentityInfo{Title: "iPhone Distribution: BFirst BLast (B72Z82XD37)"},
		}, parsedCodeSigningSettings.Identities)
		testutil.EqualSlicesWithoutOrder(t, []provprofile.ProvisioningProfileInfo{
			provprofile.ProvisioningProfileInfo{Title: "Prov Profile 42", UUID: "87af6d83-cb65-4dbe-aee7-f97a87d6fec1"},
			provprofile.ProvisioningProfileInfo{Title: "Prov Profile 43", UUID: "97af6d83-cb65-4dbe-aee7-f97a87d6fec1"},
		}, parsedCodeSigningSettings.ProvProfiles)
	}
}
