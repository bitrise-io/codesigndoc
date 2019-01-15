package xcodeproj

import (
	"testing"

	"github.com/bitrise-tools/xcode-project/pretty"
	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

func TestParseConfigurationList(t *testing.T) {
	var raw serialized.Object
	_, err := plist.Unmarshal([]byte(rawConfigurationList), &raw)
	require.NoError(t, err)

	configurationList, err := parseConfigurationList("13E76E3A1F4AC90A0028096E", raw)
	require.NoError(t, err)
	// fmt.Printf("configurationList:\n%s\n", pretty.Object(configurationList))
	require.Equal(t, expectedConfigurationList, pretty.Object(configurationList))
}

const rawConfigurationList = `
{
	13E76E3A1F4AC90A0028096E /* Build configuration list for PBXNativeTarget "code-sign-test" */ = {
		isa = XCConfigurationList;
		buildConfigurations = (
			13E76E3B1F4AC90A0028096E /* Debug */,
			13E76E3C1F4AC90A0028096E /* Release */,
		);
		defaultConfigurationIsVisible = 0;
		defaultConfigurationName = Release;
	};

	13E76E3B1F4AC90A0028096E /* Debug */ = {
		isa = XCBuildConfiguration;
		buildSettings = {
			ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
			"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
			CODE_SIGN_STYLE = Automatic;
			DEVELOPMENT_TEAM = 72SA8V3WYL;
			INFOPLIST_FILE = "code-sign-test/Info.plist";
			LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks";
			PRODUCT_BUNDLE_IDENTIFIER = "com.bitrise.code-sign-test";
			PRODUCT_NAME = "$(TARGET_NAME)";
			PROVISIONING_PROFILE = "";
			PROVISIONING_PROFILE_SPECIFIER = "";
			TARGETED_DEVICE_FAMILY = "1,2";
		};
		name = Debug;
	};

	13E76E3C1F4AC90A0028096E /* Release */ = {
		isa = XCBuildConfiguration;
		buildSettings = {
			ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
			"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
			CODE_SIGN_STYLE = Automatic;
			DEVELOPMENT_TEAM = 72SA8V3WYL;
			INFOPLIST_FILE = "code-sign-test/Info.plist";
			LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks";
			PRODUCT_BUNDLE_IDENTIFIER = "com.bitrise.code-sign-test";
			PRODUCT_NAME = "$(TARGET_NAME)";
			PROVISIONING_PROFILE = "";
			PROVISIONING_PROFILE_SPECIFIER = "";
			TARGETED_DEVICE_FAMILY = "1,2";
		};
		name = Release;
	};
}`

const expectedConfigurationList = `{
	"ID": "13E76E3A1F4AC90A0028096E",
	"DefaultConfigurationName": "Release",
	"BuildConfigurations": [
		{
			"ID": "13E76E3B1F4AC90A0028096E",
			"Name": "Debug",
			"BuildSettings": {
				"ASSETCATALOG_COMPILER_APPICON_NAME": "AppIcon",
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]": "iPhone Developer",
				"CODE_SIGN_STYLE": "Automatic",
				"DEVELOPMENT_TEAM": "72SA8V3WYL",
				"INFOPLIST_FILE": "code-sign-test/Info.plist",
				"LD_RUNPATH_SEARCH_PATHS": "$(inherited) @executable_path/Frameworks",
				"PRODUCT_BUNDLE_IDENTIFIER": "com.bitrise.code-sign-test",
				"PRODUCT_NAME": "$(TARGET_NAME)",
				"PROVISIONING_PROFILE": "",
				"PROVISIONING_PROFILE_SPECIFIER": "",
				"TARGETED_DEVICE_FAMILY": "1,2"
			}
		},
		{
			"ID": "13E76E3C1F4AC90A0028096E",
			"Name": "Release",
			"BuildSettings": {
				"ASSETCATALOG_COMPILER_APPICON_NAME": "AppIcon",
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]": "iPhone Developer",
				"CODE_SIGN_STYLE": "Automatic",
				"DEVELOPMENT_TEAM": "72SA8V3WYL",
				"INFOPLIST_FILE": "code-sign-test/Info.plist",
				"LD_RUNPATH_SEARCH_PATHS": "$(inherited) @executable_path/Frameworks",
				"PRODUCT_BUNDLE_IDENTIFIER": "com.bitrise.code-sign-test",
				"PRODUCT_NAME": "$(TARGET_NAME)",
				"PROVISIONING_PROFILE": "",
				"PROVISIONING_PROFILE_SPECIFIER": "",
				"TARGETED_DEVICE_FAMILY": "1,2"
			}
		}
	]
}`
