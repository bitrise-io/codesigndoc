package xcodeproj

import (
	"testing"

	"github.com/bitrise-tools/xcode-project/pretty"
	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

func TestParseProj(t *testing.T) {
	var raw serialized.Object
	_, err := plist.Unmarshal([]byte(rawProj), &raw)
	require.NoError(t, err)

	{
		proj, err := parseProj("13E76E061F4AC90A0028096E", raw)
		require.NoError(t, err)
		// fmt.Printf("proj:\n%s\n", pretty.Object(proj))
		require.Equal(t, expectedProj, pretty.Object(proj))
	}

	{
		proj, err := parseProj("INVALID_TARGET_ID", raw)
		require.Error(t, err)
		// fmt.Printf("proj:\n%s\n", pretty.Object(proj))
		require.Equal(t, Proj{}, proj)
	}
}

const rawProj = `
{
	13E76E061F4AC90A0028096E /* Project object */ = {
		isa = PBXProject;
		attributes = {
			LastUpgradeCheck = 0900;
			ORGANIZATIONNAME = "Gödrei Krisztián";
			TargetAttributes = {
				13E76E0D1F4AC90A0028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
				};
				13E76E251F4AC90A0028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
					TestTargetID = 13E76E0D1F4AC90A0028096E;
				};
				13E76E301F4AC90A0028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
					TestTargetID = 13E76E0D1F4AC90A0028096E;
				};
				13E76E461F4AC94F0028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
				};
				13E76E581F4AC9800028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
				};
				13E76E641F4AC9800028096E = {
					CreatedOnToolsVersion = 9.0;
					ProvisioningStyle = Automatic;
				};
			};
		};
		buildConfigurationList = 7D5B35F720E28EE80022BAE6 /* Build configuration list for PBXProject "code-sign-test" */;
		compatibilityVersion = "Xcode 8.0";
		developmentRegion = en;
		hasScannedForEncodings = 0;
		knownRegions = (
			en,
			Base,
		);
		mainGroup = 13E76E051F4AC90A0028096E;
		productRefGroup = 13E76E0F1F4AC90A0028096E /* Products */;
		projectDirPath = "";
		projectRoot = "";
		targets = (
			13E76E301F4AC90A0028096E /* code-sign-testUITests */,
		);
	};

	13E76E301F4AC90A0028096E /* code-sign-testUITests */ = {
		isa = PBXNativeTarget;
		buildConfigurationList = 13E76E401F4AC90A0028096E /* Build configuration list for PBXNativeTarget "code-sign-testUITests" */;
		buildPhases = (
			13E76E2D1F4AC90A0028096E /* Sources */,
			13E76E2E1F4AC90A0028096E /* Frameworks */,
			13E76E2F1F4AC90A0028096E /* Resources */,
		);
		buildRules = (
		);
		dependencies = (
			E62D1D9F1F22422D009F8898 /* PBXTargetDependency */,
			E67597E61F138112008C1262 /* PBXTargetDependency */,
		);
		name = "code-sign-testUITests";
		productName = "code-sign-testUITests";
		productReference = 13E76E311F4AC90A0028096E /* code-sign-testUITests.xctest */;
		productType = "com.apple.product-type.bundle.ui-testing";
	};

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

	7D5B35F720E28EE80022BAE6 /* Build configuration list for PBXProject "XcodeProj" */ = {
		isa = XCConfigurationList;
		buildConfigurations = (
			7D5B360C20E28EEA0022BAE6 /* Debug */,
		);
		defaultConfigurationIsVisible = 0;
		defaultConfigurationName = Release;
	};

	13E76E401F4AC90A0028096E /* Build configuration list for PBXNativeTarget "code-sign-testUITests" */ = {
		isa = XCConfigurationList;
		buildConfigurations = (
			13E76E411F4AC90A0028096E /* Debug */,
		);
		defaultConfigurationIsVisible = 0;
		defaultConfigurationName = Release;
	};

	7D5B360C20E28EEA0022BAE6 /* Debug */ = {
		isa = XCBuildConfiguration;
		buildSettings = {
			ALWAYS_SEARCH_USER_PATHS = NO;
			CLANG_ANALYZER_NONNULL = YES;
			CLANG_ANALYZER_NUMBER_OBJECT_CONVERSION = YES_AGGRESSIVE;
			CLANG_CXX_LANGUAGE_STANDARD = "gnu++14";
			CLANG_CXX_LIBRARY = "libc++";
			CLANG_ENABLE_MODULES = YES;
			CLANG_ENABLE_OBJC_ARC = YES;
			CLANG_ENABLE_OBJC_WEAK = YES;
			CLANG_WARN_BLOCK_CAPTURE_AUTORELEASING = YES;
			CLANG_WARN_BOOL_CONVERSION = YES;
			CLANG_WARN_COMMA = YES;
			CLANG_WARN_CONSTANT_CONVERSION = YES;
			CLANG_WARN_DEPRECATED_OBJC_IMPLEMENTATIONS = YES;
			CLANG_WARN_DIRECT_OBJC_ISA_USAGE = YES_ERROR;
			CLANG_WARN_DOCUMENTATION_COMMENTS = YES;
			CLANG_WARN_EMPTY_BODY = YES;
			CLANG_WARN_ENUM_CONVERSION = YES;
			CLANG_WARN_INFINITE_RECURSION = YES;
			CLANG_WARN_INT_CONVERSION = YES;
			CLANG_WARN_NON_LITERAL_NULL_CONVERSION = YES;
			CLANG_WARN_OBJC_IMPLICIT_RETAIN_SELF = YES;
			CLANG_WARN_OBJC_LITERAL_CONVERSION = YES;
			CLANG_WARN_OBJC_ROOT_CLASS = YES_ERROR;
			CLANG_WARN_RANGE_LOOP_ANALYSIS = YES;
			CLANG_WARN_STRICT_PROTOTYPES = YES;
			CLANG_WARN_SUSPICIOUS_MOVE = YES;
			CLANG_WARN_UNGUARDED_AVAILABILITY = YES_AGGRESSIVE;
			CLANG_WARN_UNREACHABLE_CODE = YES;
			CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
			CODE_SIGN_IDENTITY = "iPhone Developer";
			COPY_PHASE_STRIP = NO;
			DEBUG_INFORMATION_FORMAT = dwarf;
			ENABLE_STRICT_OBJC_MSGSEND = YES;
			ENABLE_TESTABILITY = YES;
			GCC_C_LANGUAGE_STANDARD = gnu11;
			GCC_DYNAMIC_NO_PIC = NO;
			GCC_NO_COMMON_BLOCKS = YES;
			GCC_OPTIMIZATION_LEVEL = 0;
			GCC_PREPROCESSOR_DEFINITIONS = (
				"DEBUG=1",
				"$(inherited)",
			);
			GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
			GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
			GCC_WARN_UNDECLARED_SELECTOR = YES;
			GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
			GCC_WARN_UNUSED_FUNCTION = YES;
			GCC_WARN_UNUSED_VARIABLE = YES;
			IPHONEOS_DEPLOYMENT_TARGET = 11.4;
			MTL_ENABLE_DEBUG_INFO = YES;
			ONLY_ACTIVE_ARCH = YES;
			SDKROOT = iphoneos;
			SWIFT_ACTIVE_COMPILATION_CONDITIONS = DEBUG;
			SWIFT_OPTIMIZATION_LEVEL = "-Onone";
		};
		name = Debug;
	};

	13E76E411F4AC90A0028096E /* Debug */ = {
		isa = XCBuildConfiguration;
		buildSettings = {
			"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
			CODE_SIGN_STYLE = Automatic;
			DEVELOPMENT_TEAM = 72SA8V3WYL;
			INFOPLIST_FILE = "code-sign-testUITests/Info.plist";
			LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks @loader_path/Frameworks";
			PRODUCT_BUNDLE_IDENTIFIER = "com.bitrise.code-sign-testUITests";
			PRODUCT_NAME = "$(TARGET_NAME)";
			PROVISIONING_PROFILE = "";
			PROVISIONING_PROFILE_SPECIFIER = "";
			TARGETED_DEVICE_FAMILY = "1,2";
			TEST_TARGET_NAME = "code-sign-test";
		};
		name = Debug;
	};

	13E76E0E1F4AC90A0028096E /* code-sign-test.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = "code-sign-test.app"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E471F4AC94F0028096E /* share-extension.appex */ = {isa = PBXFileReference; explicitFileType = "wrapper.app-extension"; includeInIndex = 0; path = "share-extension.appex"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E311F4AC90A0028096E /* code-sign-testUITests.xctest */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = "code-sign-testUITests.xctest"; sourceTree = BUILT_PRODUCTS_DIR; };

	E62D1D9F1F22422D009F8898 /* PBXTargetDependency */ = {
		isa = PBXTargetDependency;
		target = FD55DAD914CE0B0000F84D24;
		targetProxy = E62D1D9E1F22422D009F8898;
	};

	E67597E61F138112008C1262 /* PBXTargetDependency */ = {
		isa = PBXTargetDependency;
		name = MHIdentityKit;
		targetProxy = E62D1D9E1F22422D009F8898;
	};
}`

const expectedProj = `{
	"ID": "13E76E061F4AC90A0028096E",
	"BuildConfigurationList": {
		"ID": "7D5B35F720E28EE80022BAE6",
		"DefaultConfigurationName": "Release",
		"BuildConfigurations": [
			{
				"ID": "7D5B360C20E28EEA0022BAE6",
				"Name": "Debug",
				"BuildSettings": {
					"ALWAYS_SEARCH_USER_PATHS": "NO",
					"CLANG_ANALYZER_NONNULL": "YES",
					"CLANG_ANALYZER_NUMBER_OBJECT_CONVERSION": "YES_AGGRESSIVE",
					"CLANG_CXX_LANGUAGE_STANDARD": "gnu++14",
					"CLANG_CXX_LIBRARY": "libc++",
					"CLANG_ENABLE_MODULES": "YES",
					"CLANG_ENABLE_OBJC_ARC": "YES",
					"CLANG_ENABLE_OBJC_WEAK": "YES",
					"CLANG_WARN_BLOCK_CAPTURE_AUTORELEASING": "YES",
					"CLANG_WARN_BOOL_CONVERSION": "YES",
					"CLANG_WARN_COMMA": "YES",
					"CLANG_WARN_CONSTANT_CONVERSION": "YES",
					"CLANG_WARN_DEPRECATED_OBJC_IMPLEMENTATIONS": "YES",
					"CLANG_WARN_DIRECT_OBJC_ISA_USAGE": "YES_ERROR",
					"CLANG_WARN_DOCUMENTATION_COMMENTS": "YES",
					"CLANG_WARN_EMPTY_BODY": "YES",
					"CLANG_WARN_ENUM_CONVERSION": "YES",
					"CLANG_WARN_INFINITE_RECURSION": "YES",
					"CLANG_WARN_INT_CONVERSION": "YES",
					"CLANG_WARN_NON_LITERAL_NULL_CONVERSION": "YES",
					"CLANG_WARN_OBJC_IMPLICIT_RETAIN_SELF": "YES",
					"CLANG_WARN_OBJC_LITERAL_CONVERSION": "YES",
					"CLANG_WARN_OBJC_ROOT_CLASS": "YES_ERROR",
					"CLANG_WARN_RANGE_LOOP_ANALYSIS": "YES",
					"CLANG_WARN_STRICT_PROTOTYPES": "YES",
					"CLANG_WARN_SUSPICIOUS_MOVE": "YES",
					"CLANG_WARN_UNGUARDED_AVAILABILITY": "YES_AGGRESSIVE",
					"CLANG_WARN_UNREACHABLE_CODE": "YES",
					"CLANG_WARN__DUPLICATE_METHOD_MATCH": "YES",
					"CODE_SIGN_IDENTITY": "iPhone Developer",
					"COPY_PHASE_STRIP": "NO",
					"DEBUG_INFORMATION_FORMAT": "dwarf",
					"ENABLE_STRICT_OBJC_MSGSEND": "YES",
					"ENABLE_TESTABILITY": "YES",
					"GCC_C_LANGUAGE_STANDARD": "gnu11",
					"GCC_DYNAMIC_NO_PIC": "NO",
					"GCC_NO_COMMON_BLOCKS": "YES",
					"GCC_OPTIMIZATION_LEVEL": "0",
					"GCC_PREPROCESSOR_DEFINITIONS": [
						"DEBUG=1",
						"$(inherited)"
					],
					"GCC_WARN_64_TO_32_BIT_CONVERSION": "YES",
					"GCC_WARN_ABOUT_RETURN_TYPE": "YES_ERROR",
					"GCC_WARN_UNDECLARED_SELECTOR": "YES",
					"GCC_WARN_UNINITIALIZED_AUTOS": "YES_AGGRESSIVE",
					"GCC_WARN_UNUSED_FUNCTION": "YES",
					"GCC_WARN_UNUSED_VARIABLE": "YES",
					"IPHONEOS_DEPLOYMENT_TARGET": "11.4",
					"MTL_ENABLE_DEBUG_INFO": "YES",
					"ONLY_ACTIVE_ARCH": "YES",
					"SDKROOT": "iphoneos",
					"SWIFT_ACTIVE_COMPILATION_CONDITIONS": "DEBUG",
					"SWIFT_OPTIMIZATION_LEVEL": "-Onone"
				}
			}
		]
	},
	"Targets": [
		{
			"Type": "PBXNativeTarget",
			"ID": "13E76E301F4AC90A0028096E",
			"Name": "code-sign-testUITests",
			"BuildConfigurationList": {
				"ID": "13E76E401F4AC90A0028096E",
				"DefaultConfigurationName": "Release",
				"BuildConfigurations": [
					{
						"ID": "13E76E411F4AC90A0028096E",
						"Name": "Debug",
						"BuildSettings": {
							"CODE_SIGN_IDENTITY[sdk=iphoneos*]": "iPhone Developer",
							"CODE_SIGN_STYLE": "Automatic",
							"DEVELOPMENT_TEAM": "72SA8V3WYL",
							"INFOPLIST_FILE": "code-sign-testUITests/Info.plist",
							"LD_RUNPATH_SEARCH_PATHS": "$(inherited) @executable_path/Frameworks @loader_path/Frameworks",
							"PRODUCT_BUNDLE_IDENTIFIER": "com.bitrise.code-sign-testUITests",
							"PRODUCT_NAME": "$(TARGET_NAME)",
							"PROVISIONING_PROFILE": "",
							"PROVISIONING_PROFILE_SPECIFIER": "",
							"TARGETED_DEVICE_FAMILY": "1,2",
							"TEST_TARGET_NAME": "code-sign-test"
						}
					}
				]
			},
			"Dependencies": null,
			"ProductReference": {
				"Path": "code-sign-testUITests.xctest"
			},
			"ProductType": "com.apple.product-type.bundle.ui-testing"
		}
	]
}`
