package xcodeproj

const pbxprojWithouthTargetAttributes = `// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {
	};
	objectVersion = 50;
	objects = {

/* Begin PBXNativeTarget section */
		13BD62FD256BE6D000F72361 /* Target */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13BD6323256BE6D300F72361 /* Build configuration list for PBXNativeTarget "Target" */;
			buildPhases = (
				13BD62FA256BE6D000F72361 /* Sources */,
				13BD62FB256BE6D000F72361 /* Frameworks */,
				13BD62FC256BE6D000F72361 /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
			);
			name = Target;
			productName = Target;
			productReference = 13BD62FE256BE6D000F72361 /* Target.app */;
			productType = "com.apple.product-type.application";
		};
		13BD6332256BE7BF00F72361 /* TargetWithouthTargetAttributes */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13BD633A256BE7BF00F72361 /* Build configuration list for PBXNativeTarget "TargetWithouthTargetAttributes" */;
			buildPhases = (
				13BD6333256BE7BF00F72361 /* Sources */,
				13BD6336256BE7BF00F72361 /* Frameworks */,
				13BD6337256BE7BF00F72361 /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
			);
			name = TargetWithouthTargetAttributes;
			productName = Target;
			productReference = 13BD633D256BE7BF00F72361 /* TargetWithouthTargetAttributes.app */;
			productType = "com.apple.product-type.application";
		};
/* End PBXNativeTarget section */

/* Begin PBXProject section */
		13BD62F6256BE6D000F72361 /* Project object */ = {
			isa = PBXProject;
			attributes = {
				LastSwiftUpdateCheck = 1220;
				LastUpgradeCheck = 1220;
				TargetAttributes = {
					13BD62FD256BE6D000F72361 = {
						CreatedOnToolsVersion = 12.2;
					};
				};
			};
			buildConfigurationList = 13BD62F9256BE6D000F72361 /* Build configuration list for PBXProject "Target" */;
			compatibilityVersion = "Xcode 9.3";
			developmentRegion = en;
			hasScannedForEncodings = 0;
			knownRegions = (
				en,
				Base,
			);
			mainGroup = 13BD62F5256BE6D000F72361;
			productRefGroup = 13BD62FF256BE6D000F72361 /* Products */;
			projectDirPath = "";
			projectRoot = "";
			targets = (
				13BD62FD256BE6D000F72361 /* Target */,
				13BD6332256BE7BF00F72361 /* TargetWithouthTargetAttributes */,
			);
		};
/* End PBXProject section */

/* Begin XCBuildConfiguration section */
		13BD6321256BE6D300F72361 /* Debug */ = {
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
				CLANG_WARN_QUOTED_INCLUDE_IN_FRAMEWORK_HEADER = YES;
				CLANG_WARN_RANGE_LOOP_ANALYSIS = YES;
				CLANG_WARN_STRICT_PROTOTYPES = YES;
				CLANG_WARN_SUSPICIOUS_MOVE = YES;
				CLANG_WARN_UNGUARDED_AVAILABILITY = YES_AGGRESSIVE;
				CLANG_WARN_UNREACHABLE_CODE = YES;
				CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
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
				IPHONEOS_DEPLOYMENT_TARGET = 14.2;
				MTL_ENABLE_DEBUG_INFO = INCLUDE_SOURCE;
				MTL_FAST_MATH = YES;
				ONLY_ACTIVE_ARCH = YES;
				SDKROOT = iphoneos;
				SWIFT_ACTIVE_COMPILATION_CONDITIONS = DEBUG;
				SWIFT_OPTIMIZATION_LEVEL = "-Onone";
			};
			name = Debug;
		};
		13BD6324256BE6D300F72361 /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				ASSETCATALOG_COMPILER_GLOBAL_ACCENT_COLOR_NAME = AccentColor;
				CODE_SIGN_STYLE = Automatic;
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "Apple Development: Bitrise Bot (ASDF1234)";
				DEVELOPMENT_ASSET_PATHS = "\"Target/Preview Content\"";
				DEVELOPMENT_TEAM = ASDF2FASDF;
				ENABLE_PREVIEWS = YES;
				INFOPLIST_FILE = Target/Info.plist;
				IPHONEOS_DEPLOYMENT_TARGET = 14.0;
				LD_RUNPATH_SEARCH_PATHS = (
					"$(inherited)",
					"@executable_path/Frameworks",
				);
				PRODUCT_BUNDLE_IDENTIFIER = io.bitrise.target.Target;
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 5.0;
				TARGETED_DEVICE_FAMILY = "1,2";
			};
			name = Debug;
		};
		13BD633B256BE7BF00F72361 /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				ASSETCATALOG_COMPILER_GLOBAL_ACCENT_COLOR_NAME = AccentColor;
				CODE_SIGN_IDENTITY = "Apple Development: John Doe (ASDF1234)";
				CODE_SIGN_STYLE = Manual;
				DEVELOPMENT_ASSET_PATHS = "\"Target/Preview Content\"";
				DEVELOPMENT_TEAM = ASDF2FASDF;
				ENABLE_PREVIEWS = YES;
				INFOPLIST_FILE = "Target copy-Info.plist";
				IPHONEOS_DEPLOYMENT_TARGET = 14.0;
				LD_RUNPATH_SEARCH_PATHS = (
					"$(inherited)",
					"@executable_path/Frameworks",
				);
				PRODUCT_BUNDLE_IDENTIFIER = io.bitrise.target.Target;
				PRODUCT_NAME = "$(TARGET_NAME)";
				PROVISIONING_PROFILE_SPECIFIER = "Wildcard (io.bitrise.*) iOS Development";
				SWIFT_VERSION = 5.0;
				TARGETED_DEVICE_FAMILY = "1,2";
			};
			name = Debug;
		};
/* End XCBuildConfiguration section */

/* Begin XCConfigurationList section */
		13BD62F9256BE6D000F72361 /* Build configuration list for PBXProject "Target" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13BD6321256BE6D300F72361 /* Debug */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13BD6323256BE6D300F72361 /* Build configuration list for PBXNativeTarget "Target" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13BD6324256BE6D300F72361 /* Debug */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13BD633A256BE7BF00F72361 /* Build configuration list for PBXNativeTarget "TargetWithouthTargetAttributes" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13BD633B256BE7BF00F72361 /* Debug */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
/* End XCConfigurationList section */

/* Begin PBXFileReference section */
		13BD62FE256BE6D000F72361 /* Target.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = Target.app; sourceTree = BUILT_PRODUCTS_DIR; };
		13BD633D256BE7BF00F72361 /* TargetWithouthTargetAttributes.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = TargetWithouthTargetAttributes.app; sourceTree = BUILT_PRODUCTS_DIR; };
/* End PBXFileReference section */
	};
	rootObject = 13BD62F6256BE6D000F72361 /* Project object */;
}
`

// pbxprojWithouthTargetAttributes with updated settings
const pbxprojWTAafterPerObjectModify = `// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {
	};
	objectVersion = 50;
	objects = {

/* Begin PBXNativeTarget section */
		13BD62FD256BE6D000F72361 /* Target */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13BD6323256BE6D300F72361 /* Build configuration list for PBXNativeTarget "Target" */;
			buildPhases = (
				13BD62FA256BE6D000F72361 /* Sources */,
				13BD62FB256BE6D000F72361 /* Frameworks */,
				13BD62FC256BE6D000F72361 /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
			);
			name = Target;
			productName = Target;
			productReference = 13BD62FE256BE6D000F72361 /* Target.app */;
			productType = "com.apple.product-type.application";
		};
		13BD6332256BE7BF00F72361 /* TargetWithouthTargetAttributes */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13BD633A256BE7BF00F72361 /* Build configuration list for PBXNativeTarget "TargetWithouthTargetAttributes" */;
			buildPhases = (
				13BD6333256BE7BF00F72361 /* Sources */,
				13BD6336256BE7BF00F72361 /* Frameworks */,
				13BD6337256BE7BF00F72361 /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
			);
			name = TargetWithouthTargetAttributes;
			productName = Target;
			productReference = 13BD633D256BE7BF00F72361 /* TargetWithouthTargetAttributes.app */;
			productType = "com.apple.product-type.application";
		};
/* End PBXNativeTarget section */

/* Begin PBXProject section */
		13BD62F6256BE6D000F72361 /* Project object */ = {
			isa = PBXProject;
			attributes = {
				LastSwiftUpdateCheck = 1220;
				LastUpgradeCheck = 1220;
				TargetAttributes = {
					13BD62FD256BE6D000F72361 = {
						CreatedOnToolsVersion = 12.2;
					};
				};
			};
			buildConfigurationList = 13BD62F9256BE6D000F72361 /* Build configuration list for PBXProject "Target" */;
			compatibilityVersion = "Xcode 9.3";
			developmentRegion = en;
			hasScannedForEncodings = 0;
			knownRegions = (
				en,
				Base,
			);
			mainGroup = 13BD62F5256BE6D000F72361;
			productRefGroup = 13BD62FF256BE6D000F72361 /* Products */;
			projectDirPath = "";
			projectRoot = "";
			targets = (
				13BD62FD256BE6D000F72361 /* Target */,
				13BD6332256BE7BF00F72361 /* TargetWithouthTargetAttributes */,
			);
		};
/* End PBXProject section */

/* Begin XCBuildConfiguration section */
		13BD6321256BE6D300F72361 /* Debug */ = {
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
				CLANG_WARN_QUOTED_INCLUDE_IN_FRAMEWORK_HEADER = YES;
				CLANG_WARN_RANGE_LOOP_ANALYSIS = YES;
				CLANG_WARN_STRICT_PROTOTYPES = YES;
				CLANG_WARN_SUSPICIOUS_MOVE = YES;
				CLANG_WARN_UNGUARDED_AVAILABILITY = YES_AGGRESSIVE;
				CLANG_WARN_UNREACHABLE_CODE = YES;
				CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
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
				IPHONEOS_DEPLOYMENT_TARGET = 14.2;
				MTL_ENABLE_DEBUG_INFO = INCLUDE_SOURCE;
				MTL_FAST_MATH = YES;
				ONLY_ACTIVE_ARCH = YES;
				SDKROOT = iphoneos;
				SWIFT_ACTIVE_COMPILATION_CONDITIONS = DEBUG;
				SWIFT_OPTIMIZATION_LEVEL = "-Onone";
			};
			name = Debug;
		};
		13BD6324256BE6D300F72361 /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				ASSETCATALOG_COMPILER_GLOBAL_ACCENT_COLOR_NAME = AccentColor;
				CODE_SIGN_STYLE = Automatic;
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "Apple Development: Bitrise Bot (ASDF1234)";
				DEVELOPMENT_ASSET_PATHS = "\"Target/Preview Content\"";
				DEVELOPMENT_TEAM = ASDF2FASDF;
				ENABLE_PREVIEWS = YES;
				INFOPLIST_FILE = Target/Info.plist;
				IPHONEOS_DEPLOYMENT_TARGET = 14.0;
				LD_RUNPATH_SEARCH_PATHS = (
					"$(inherited)",
					"@executable_path/Frameworks",
				);
				PRODUCT_BUNDLE_IDENTIFIER = io.bitrise.target.Target;
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 5.0;
				TARGETED_DEVICE_FAMILY = "1,2";
			};
			name = Debug;
		};
		13BD633B256BE7BF00F72361 /* Debug */ = {
	buildSettings = {
		"ASSETCATALOG_COMPILER_APPICON_NAME" = AppIcon;
		"ASSETCATALOG_COMPILER_GLOBAL_ACCENT_COLOR_NAME" = AccentColor;
		"CODE_SIGN_IDENTITY" = "Apple Development: John Doe (ASDF1234)";
		"CODE_SIGN_STYLE" = Manual;
		"DEVELOPMENT_ASSET_PATHS" = "\"Target/Preview Content\"";
		"DEVELOPMENT_TEAM" = ABCD1234;
		"ENABLE_PREVIEWS" = YES;
		"INFOPLIST_FILE" = "Target copy-Info.plist";
		"IPHONEOS_DEPLOYMENT_TARGET" = "14.0";
		"LD_RUNPATH_SEARCH_PATHS" = (
			"$(inherited)",
			"@executable_path/Frameworks",
		);
		"PRODUCT_BUNDLE_IDENTIFIER" = "io.bitrise.target.Target";
		"PRODUCT_NAME" = "$(TARGET_NAME)";
		"PROVISIONING_PROFILE" = "asdf56b6-e75a-4f86-bf25-101bfc2fasdf";
		"PROVISIONING_PROFILE_SPECIFIER" = "";
		"SWIFT_VERSION" = "5.0";
		"TARGETED_DEVICE_FAMILY" = "1,2";
	};
	isa = XCBuildConfiguration;
	name = Debug;
};
/* End XCBuildConfiguration section */

/* Begin XCConfigurationList section */
		13BD62F9256BE6D000F72361 /* Build configuration list for PBXProject "Target" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13BD6321256BE6D300F72361 /* Debug */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13BD6323256BE6D300F72361 /* Build configuration list for PBXNativeTarget "Target" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13BD6324256BE6D300F72361 /* Debug */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13BD633A256BE7BF00F72361 /* Build configuration list for PBXNativeTarget "TargetWithouthTargetAttributes" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13BD633B256BE7BF00F72361 /* Debug */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
/* End XCConfigurationList section */

/* Begin PBXFileReference section */
		13BD62FE256BE6D000F72361 /* Target.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = Target.app; sourceTree = BUILT_PRODUCTS_DIR; };
		13BD633D256BE7BF00F72361 /* TargetWithouthTargetAttributes.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = TargetWithouthTargetAttributes.app; sourceTree = BUILT_PRODUCTS_DIR; };
/* End PBXFileReference section */
	};
	rootObject = 13BD62F6256BE6D000F72361 /* Project object */;
}
`
