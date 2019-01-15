package xcodeproj

import (
	"testing"

	"github.com/bitrise-tools/xcode-project/pretty"
	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

func TestIsExecutableProduct(t *testing.T) {
	var raw serialized.Object
	_, err := plist.Unmarshal([]byte(rawNativeTarget), &raw)
	require.NoError(t, err)

	{
		target, err := parseTarget("13E76E0D1F4AC90A0028096E", raw)
		require.NoError(t, err)

		require.True(t, target.IsAppProduct())
		require.False(t, target.IsAppExtensionProduct())
		require.True(t, target.IsExecutableProduct())
	}

	{

		target, err := parseTarget("13E76E461F4AC94F0028096E", raw)
		require.NoError(t, err)

		require.False(t, target.IsAppProduct())
		require.True(t, target.IsAppExtensionProduct())
		require.True(t, target.IsExecutableProduct())
	}

}

func TestParseTarget(t *testing.T) {
	t.Log("PBXNativeTarget")
	{
		var raw serialized.Object
		_, err := plist.Unmarshal([]byte(rawNativeTarget), &raw)
		require.NoError(t, err)

		target, err := parseTarget("13E76E0D1F4AC90A0028096E", raw)
		require.NoError(t, err)
		// fmt.Printf("target:\n%s\n", pretty.Object(target))
		require.Equal(t, expectedNativeTarget, pretty.Object(target))
	}

	t.Log("PBXAggregateTarget")
	{
		var raw serialized.Object
		_, err := plist.Unmarshal([]byte(rawAggregateTarget), &raw)
		require.NoError(t, err)

		target, err := parseTarget("FD55DAD914CE0B0000F84D24", raw)
		require.NoError(t, err)
		// fmt.Printf("target:\n%s\n", pretty.Object(target))
		require.Equal(t, expectedAggregateTarget, pretty.Object(target))
	}

	t.Log("PBXLegacyTarget")
	{
		var raw serialized.Object
		_, err := plist.Unmarshal([]byte(rawLegacyTarget), &raw)
		require.NoError(t, err)

		target, err := parseTarget("407952600CEA391500E202DC", raw)
		require.NoError(t, err)
		// fmt.Printf("target:\n%s\n", pretty.Object(target))
		require.Equal(t, expectedLegacyTarget, pretty.Object(target))
	}

	t.Log("Invalid Target ID")
	{
		var raw serialized.Object
		_, err := plist.Unmarshal([]byte(rawLegacyTarget), &raw)
		require.NoError(t, err)

		target, err := parseTarget("INVALID_TARGET_ID", raw)
		require.Error(t, err)
		require.Equal(t, Target{}, target)
	}
}

const rawLegacyTarget = `{
	407952600CEA391500E202DC /* build */ = {
		isa = PBXLegacyTarget;
		buildArgumentsString = all;
		buildConfigurationList = 407952610CEA393300E202DC /* Build configuration list for PBXLegacyTarget "build" */;
		buildPhases = (
		);
		buildToolPath = /usr/bin/make;
		buildWorkingDirectory = firmware;
		dependencies = (
		);
		name = build;
		passBuildSettingsInEnvironment = 1;
		productName = "Build All";
		productType = "com.apple.product-type.application";
	};

	407952610CEA393300E202DC /* Build configuration list for PBXLegacyTarget "build" */ = {
		isa = XCConfigurationList;
		buildConfigurations = (
			407952630CEA393300E202DC /* Release */,
		);
		defaultConfigurationIsVisible = 0;
		defaultConfigurationName = Release;
	};

	407952630CEA393300E202DC /* Release */ = {
		isa = XCBuildConfiguration;
		buildSettings = {
			PATH = "$(PATH):/usr/local/CrossPack-AVR/bin";
		};
		name = Release;
	};
}`

const expectedLegacyTarget = `{
	"Type": "PBXLegacyTarget",
	"ID": "407952600CEA391500E202DC",
	"Name": "build",
	"BuildConfigurationList": {
		"ID": "407952610CEA393300E202DC",
		"DefaultConfigurationName": "Release",
		"BuildConfigurations": [
			{
				"ID": "407952630CEA393300E202DC",
				"Name": "Release",
				"BuildSettings": {
					"PATH": "$(PATH):/usr/local/CrossPack-AVR/bin"
				}
			}
		]
	},
	"Dependencies": null,
	"ProductReference": {
		"Path": ""
	},
	"ProductType": "com.apple.product-type.application"
}`

const rawAggregateTarget = `{
	FD55DAD914CE0B0000F84D24 /* rpcsvc */ = {
		isa = PBXAggregateTarget;
		buildConfigurationList = FD55DADA14CE0B0000F84D24 /* Build configuration list for PBXAggregateTarget "rpcsvc" */;
		buildPhases = (
			FD55DADC14CE0B0700F84D24 /* Run Script */,
		);
		dependencies = (
		);
		name = rpcsvc;
		productName = rpcsvc;
	};

	FD55DADA14CE0B0000F84D24 /* Build configuration list for PBXAggregateTarget "rpcsvc" */ = {
		isa = XCConfigurationList;
		buildConfigurations = (
			FD55DADB14CE0B0000F84D24 /* Release */,
		);
		defaultConfigurationIsVisible = 0;
		defaultConfigurationName = Release;
	};

	FD55DADB14CE0B0000F84D24 /* Release */ = {
		isa = XCBuildConfiguration;
		buildSettings = {
			INSTALLHDRS_SCRIPT_PHASE = YES;
			PRODUCT_NAME = "$(TARGET_NAME)";
			PRODUCT_BUNDLE_IDENTIFIER = "Bitrise.$(PRODUCT_NAME:rfc1034identifier).watch";
		};
		name = Release;
	};
}`

const expectedAggregateTarget = `{
	"Type": "PBXAggregateTarget",
	"ID": "FD55DAD914CE0B0000F84D24",
	"Name": "rpcsvc",
	"BuildConfigurationList": {
		"ID": "FD55DADA14CE0B0000F84D24",
		"DefaultConfigurationName": "Release",
		"BuildConfigurations": [
			{
				"ID": "FD55DADB14CE0B0000F84D24",
				"Name": "Release",
				"BuildSettings": {
					"INSTALLHDRS_SCRIPT_PHASE": "YES",
					"PRODUCT_BUNDLE_IDENTIFIER": "Bitrise.$(PRODUCT_NAME:rfc1034identifier).watch",
					"PRODUCT_NAME": "$(TARGET_NAME)"
				}
			}
		]
	},
	"Dependencies": null,
	"ProductReference": {
		"Path": ""
	},
	"ProductType": ""
}`

const rawNativeTarget = `{
	13E76E0D1F4AC90A0028096E /* code-sign-test */ = {
		isa = PBXNativeTarget;
		buildConfigurationList = 13E76E3A1F4AC90A0028096E /* Build configuration list for PBXNativeTarget "code-sign-test" */;
		buildPhases = (
			13E76E0A1F4AC90A0028096E /* Sources */,
			13E76E0B1F4AC90A0028096E /* Frameworks */,
			13E76E0C1F4AC90A0028096E /* Resources */,
			13E76E561F4AC94F0028096E /* Embed App Extensions */,
			13E76E811F4AC9800028096E /* Embed Watch Content */,
		);
		buildRules = (
		);
		dependencies = (
			13E76E511F4AC94F0028096E /* PBXTargetDependency */,
		);
		name = "code-sign-test";
		productName = "code-sign-test";
		productReference = 13E76E0E1F4AC90A0028096E /* code-sign-test.app */;
		productType = "com.apple.product-type.application";
	};

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

	13E76E511F4AC94F0028096E /* PBXTargetDependency */ = {
		isa = PBXTargetDependency;
		target = 13E76E461F4AC94F0028096E /* share-extension */;
		targetProxy = 13E76E501F4AC94F0028096E /* PBXContainerItemProxy */;
	};

	13E76E461F4AC94F0028096E /* share-extension */ = {
		isa = PBXNativeTarget;
		buildConfigurationList = 13E76E3A1F4AC90A0028096E /* Build configuration list for PBXNativeTarget "share-extension" */;
		buildPhases = (
			13E76E431F4AC94F0028096E /* Sources */,
			13E76E441F4AC94F0028096E /* Frameworks */,
			13E76E451F4AC94F0028096E /* Resources */,
		);
		buildRules = (
		);
		dependencies = (
		);
		name = "share-extension";
		productName = "share-extension";
		productReference = 13E76E471F4AC94F0028096E /* share-extension.appex */;
		productType = "com.apple.product-type.app-extension";
	};

	13E76E0E1F4AC90A0028096E /* code-sign-test.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = "code-sign-test.app"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E471F4AC94F0028096E /* share-extension.appex */ = {isa = PBXFileReference; explicitFileType = "wrapper.app-extension"; includeInIndex = 0; path = "share-extension.appex"; sourceTree = BUILT_PRODUCTS_DIR; };
}`

const expectedNativeTarget = `{
	"Type": "PBXNativeTarget",
	"ID": "13E76E0D1F4AC90A0028096E",
	"Name": "code-sign-test",
	"BuildConfigurationList": {
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
	},
	"Dependencies": [
		{
			"ID": "13E76E511F4AC94F0028096E",
			"Target": {
				"Type": "PBXNativeTarget",
				"ID": "13E76E461F4AC94F0028096E",
				"Name": "share-extension",
				"BuildConfigurationList": {
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
				},
				"Dependencies": null,
				"ProductReference": {
					"Path": "share-extension.appex"
				},
				"ProductType": "com.apple.product-type.app-extension"
			}
		}
	],
	"ProductReference": {
		"Path": "code-sign-test.app"
	},
	"ProductType": "com.apple.product-type.application"
}`
