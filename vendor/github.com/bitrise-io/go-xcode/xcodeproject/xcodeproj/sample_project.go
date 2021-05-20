package xcodeproj

// https://github.com/bitrise-io/CatalystSample pbxproj -> objects
const rawCatalystProj = `{

/* Begin PBXBuildFile section */
		13917C16243F43D00087912B /* AppDelegate.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13917C15243F43D00087912B /* AppDelegate.swift */; };
		13917C18243F43D00087912B /* SceneDelegate.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13917C17243F43D00087912B /* SceneDelegate.swift */; };
		13917C1A243F43D00087912B /* ContentView.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13917C19243F43D00087912B /* ContentView.swift */; };
		13917C1C243F43D10087912B /* Assets.xcassets in Resources */ = {isa = PBXBuildFile; fileRef = 13917C1B243F43D10087912B /* Assets.xcassets */; };
		13917C1F243F43D10087912B /* Preview Assets.xcassets in Resources */ = {isa = PBXBuildFile; fileRef = 13917C1E243F43D10087912B /* Preview Assets.xcassets */; };
		13917C22243F43D10087912B /* LaunchScreen.storyboard in Resources */ = {isa = PBXBuildFile; fileRef = 13917C20243F43D10087912B /* LaunchScreen.storyboard */; };
		13917C2D243F43D10087912B /* Catalyst_SampleTests.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13917C2C243F43D10087912B /* Catalyst_SampleTests.swift */; };
		13917C38243F43D10087912B /* Catalyst_SampleUITests.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13917C37243F43D10087912B /* Catalyst_SampleUITests.swift */; };
/* End PBXBuildFile section */

/* Begin PBXContainerItemProxy section */
		13917C29243F43D10087912B /* PBXContainerItemProxy */ = {
			isa = PBXContainerItemProxy;
			containerPortal = 13917C0A243F43D00087912B /* Project object */;
			proxyType = 1;
			remoteGlobalIDString = 13917C11243F43D00087912B;
			remoteInfo = "Catalyst Sample";
		};
		13917C34243F43D10087912B /* PBXContainerItemProxy */ = {
			isa = PBXContainerItemProxy;
			containerPortal = 13917C0A243F43D00087912B /* Project object */;
			proxyType = 1;
			remoteGlobalIDString = 13917C11243F43D00087912B;
			remoteInfo = "Catalyst Sample";
		};
/* End PBXContainerItemProxy section */

/* Begin PBXFileReference section */
		13917C12243F43D00087912B /* Catalyst Sample.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = "Catalyst Sample.app"; sourceTree = BUILT_PRODUCTS_DIR; };
		13917C15243F43D00087912B /* AppDelegate.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = AppDelegate.swift; sourceTree = "<group>"; };
		13917C17243F43D00087912B /* SceneDelegate.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = SceneDelegate.swift; sourceTree = "<group>"; };
		13917C19243F43D00087912B /* ContentView.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = ContentView.swift; sourceTree = "<group>"; };
		13917C1B243F43D10087912B /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
		13917C1E243F43D10087912B /* Preview Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = "Preview Assets.xcassets"; sourceTree = "<group>"; };
		13917C21243F43D10087912B /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/LaunchScreen.storyboard; sourceTree = "<group>"; };
		13917C23243F43D10087912B /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
		13917C28243F43D10087912B /* Catalyst SampleTests.xctest */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = "Catalyst SampleTests.xctest"; sourceTree = BUILT_PRODUCTS_DIR; };
		13917C2C243F43D10087912B /* Catalyst_SampleTests.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = Catalyst_SampleTests.swift; sourceTree = "<group>"; };
		13917C2E243F43D10087912B /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
		13917C33243F43D10087912B /* Catalyst SampleUITests.xctest */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = "Catalyst SampleUITests.xctest"; sourceTree = BUILT_PRODUCTS_DIR; };
		13917C37243F43D10087912B /* Catalyst_SampleUITests.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = Catalyst_SampleUITests.swift; sourceTree = "<group>"; };
		13917C39243F43D10087912B /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
		13917C45243F44BA0087912B /* Catalyst Sample.entitlements */ = {isa = PBXFileReference; lastKnownFileType = text.plist.entitlements; path = "Catalyst Sample.entitlements"; sourceTree = "<group>"; };
/* End PBXFileReference section */

/* Begin PBXFrameworksBuildPhase section */
		13917C0F243F43D00087912B /* Frameworks */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13917C25243F43D10087912B /* Frameworks */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13917C30243F43D10087912B /* Frameworks */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXFrameworksBuildPhase section */

/* Begin PBXGroup section */
		13917C09243F43D00087912B = {
			isa = PBXGroup;
			children = (
				13917C14243F43D00087912B /* Catalyst Sample */,
				13917C2B243F43D10087912B /* Catalyst SampleTests */,
				13917C36243F43D10087912B /* Catalyst SampleUITests */,
				13917C13243F43D00087912B /* Products */,
			);
			sourceTree = "<group>";
		};
		13917C13243F43D00087912B /* Products */ = {
			isa = PBXGroup;
			children = (
				13917C12243F43D00087912B /* Catalyst Sample.app */,
				13917C28243F43D10087912B /* Catalyst SampleTests.xctest */,
				13917C33243F43D10087912B /* Catalyst SampleUITests.xctest */,
			);
			name = Products;
			sourceTree = "<group>";
		};
		13917C14243F43D00087912B /* Catalyst Sample */ = {
			isa = PBXGroup;
			children = (
				13917C45243F44BA0087912B /* Catalyst Sample.entitlements */,
				13917C15243F43D00087912B /* AppDelegate.swift */,
				13917C17243F43D00087912B /* SceneDelegate.swift */,
				13917C19243F43D00087912B /* ContentView.swift */,
				13917C1B243F43D10087912B /* Assets.xcassets */,
				13917C20243F43D10087912B /* LaunchScreen.storyboard */,
				13917C23243F43D10087912B /* Info.plist */,
				13917C1D243F43D10087912B /* Preview Content */,
			);
			path = "Catalyst Sample";
			sourceTree = "<group>";
		};
		13917C1D243F43D10087912B /* Preview Content */ = {
			isa = PBXGroup;
			children = (
				13917C1E243F43D10087912B /* Preview Assets.xcassets */,
			);
			path = "Preview Content";
			sourceTree = "<group>";
		};
		13917C2B243F43D10087912B /* Catalyst SampleTests */ = {
			isa = PBXGroup;
			children = (
				13917C2C243F43D10087912B /* Catalyst_SampleTests.swift */,
				13917C2E243F43D10087912B /* Info.plist */,
			);
			path = "Catalyst SampleTests";
			sourceTree = "<group>";
		};
		13917C36243F43D10087912B /* Catalyst SampleUITests */ = {
			isa = PBXGroup;
			children = (
				13917C37243F43D10087912B /* Catalyst_SampleUITests.swift */,
				13917C39243F43D10087912B /* Info.plist */,
			);
			path = "Catalyst SampleUITests";
			sourceTree = "<group>";
		};
/* End PBXGroup section */

/* Begin PBXNativeTarget section */
		13917C11243F43D00087912B /* Catalyst Sample */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13917C3C243F43D10087912B /* Build configuration list for PBXNativeTarget "Catalyst Sample" */;
			buildPhases = (
				13917C0E243F43D00087912B /* Sources */,
				13917C0F243F43D00087912B /* Frameworks */,
				13917C10243F43D00087912B /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
			);
			name = "Catalyst Sample";
			productName = "Catalyst Sample";
			productReference = 13917C12243F43D00087912B /* Catalyst Sample.app */;
			productType = "com.apple.product-type.application";
		};
		13917C27243F43D10087912B /* Catalyst SampleTests */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13917C3F243F43D10087912B /* Build configuration list for PBXNativeTarget "Catalyst SampleTests" */;
			buildPhases = (
				13917C24243F43D10087912B /* Sources */,
				13917C25243F43D10087912B /* Frameworks */,
				13917C26243F43D10087912B /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
				13917C2A243F43D10087912B /* PBXTargetDependency */,
			);
			name = "Catalyst SampleTests";
			productName = "Catalyst SampleTests";
			productReference = 13917C28243F43D10087912B /* Catalyst SampleTests.xctest */;
			productType = "com.apple.product-type.bundle.unit-test";
		};
		13917C32243F43D10087912B /* Catalyst SampleUITests */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13917C42243F43D10087912B /* Build configuration list for PBXNativeTarget "Catalyst SampleUITests" */;
			buildPhases = (
				13917C2F243F43D10087912B /* Sources */,
				13917C30243F43D10087912B /* Frameworks */,
				13917C31243F43D10087912B /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
				13917C35243F43D10087912B /* PBXTargetDependency */,
			);
			name = "Catalyst SampleUITests";
			productName = "Catalyst SampleUITests";
			productReference = 13917C33243F43D10087912B /* Catalyst SampleUITests.xctest */;
			productType = "com.apple.product-type.bundle.ui-testing";
		};
/* End PBXNativeTarget section */

/* Begin PBXProject section */
		13917C0A243F43D00087912B /* Project object */ = {
			isa = PBXProject;
			attributes = {
				LastSwiftUpdateCheck = 1140;
				LastUpgradeCheck = 1140;
				ORGANIZATIONNAME = "Krisztián Gödrei";
				TargetAttributes = {
					13917C11243F43D00087912B = {
						CreatedOnToolsVersion = 11.4;
					};
					13917C27243F43D10087912B = {
						CreatedOnToolsVersion = 11.4;
						TestTargetID = 13917C11243F43D00087912B;
					};
					13917C32243F43D10087912B = {
						CreatedOnToolsVersion = 11.4;
						TestTargetID = 13917C11243F43D00087912B;
					};
				};
			};
			buildConfigurationList = 13917C0D243F43D00087912B /* Build configuration list for PBXProject "Catalyst Sample" */;
			compatibilityVersion = "Xcode 9.3";
			developmentRegion = en;
			hasScannedForEncodings = 0;
			knownRegions = (
				en,
				Base,
			);
			mainGroup = 13917C09243F43D00087912B;
			productRefGroup = 13917C13243F43D00087912B /* Products */;
			projectDirPath = "";
			projectRoot = "";
			targets = (
				13917C11243F43D00087912B /* Catalyst Sample */,
				13917C27243F43D10087912B /* Catalyst SampleTests */,
				13917C32243F43D10087912B /* Catalyst SampleUITests */,
			);
		};
/* End PBXProject section */

/* Begin PBXResourcesBuildPhase section */
		13917C10243F43D00087912B /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13917C22243F43D10087912B /* LaunchScreen.storyboard in Resources */,
				13917C1F243F43D10087912B /* Preview Assets.xcassets in Resources */,
				13917C1C243F43D10087912B /* Assets.xcassets in Resources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13917C26243F43D10087912B /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13917C31243F43D10087912B /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXResourcesBuildPhase section */

/* Begin PBXSourcesBuildPhase section */
		13917C0E243F43D00087912B /* Sources */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13917C16243F43D00087912B /* AppDelegate.swift in Sources */,
				13917C18243F43D00087912B /* SceneDelegate.swift in Sources */,
				13917C1A243F43D00087912B /* ContentView.swift in Sources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13917C24243F43D10087912B /* Sources */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13917C2D243F43D10087912B /* Catalyst_SampleTests.swift in Sources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13917C2F243F43D10087912B /* Sources */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13917C38243F43D10087912B /* Catalyst_SampleUITests.swift in Sources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXSourcesBuildPhase section */

/* Begin PBXTargetDependency section */
		13917C2A243F43D10087912B /* PBXTargetDependency */ = {
			isa = PBXTargetDependency;
			target = 13917C11243F43D00087912B /* Catalyst Sample */;
			targetProxy = 13917C29243F43D10087912B /* PBXContainerItemProxy */;
		};
		13917C35243F43D10087912B /* PBXTargetDependency */ = {
			isa = PBXTargetDependency;
			target = 13917C11243F43D00087912B /* Catalyst Sample */;
			targetProxy = 13917C34243F43D10087912B /* PBXContainerItemProxy */;
		};
/* End PBXTargetDependency section */

/* Begin PBXVariantGroup section */
		13917C20243F43D10087912B /* LaunchScreen.storyboard */ = {
			isa = PBXVariantGroup;
			children = (
				13917C21243F43D10087912B /* Base */,
			);
			name = LaunchScreen.storyboard;
			sourceTree = "<group>";
		};
/* End PBXVariantGroup section */

/* Begin XCBuildConfiguration section */
		13917C3A243F43D10087912B /* Debug */ = {
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
				IPHONEOS_DEPLOYMENT_TARGET = 13.4;
				MTL_ENABLE_DEBUG_INFO = INCLUDE_SOURCE;
				MTL_FAST_MATH = YES;
				ONLY_ACTIVE_ARCH = YES;
				SDKROOT = iphoneos;
				SWIFT_ACTIVE_COMPILATION_CONDITIONS = DEBUG;
				SWIFT_OPTIMIZATION_LEVEL = "-Onone";
			};
			name = Debug;
		};
		13917C3B243F43D10087912B /* Release */ = {
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
				COPY_PHASE_STRIP = NO;
				DEBUG_INFORMATION_FORMAT = "dwarf-with-dsym";
				ENABLE_NS_ASSERTIONS = NO;
				ENABLE_STRICT_OBJC_MSGSEND = YES;
				GCC_C_LANGUAGE_STANDARD = gnu11;
				GCC_NO_COMMON_BLOCKS = YES;
				GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
				GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
				GCC_WARN_UNDECLARED_SELECTOR = YES;
				GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
				GCC_WARN_UNUSED_FUNCTION = YES;
				GCC_WARN_UNUSED_VARIABLE = YES;
				IPHONEOS_DEPLOYMENT_TARGET = 13.4;
				MTL_ENABLE_DEBUG_INFO = NO;
				MTL_FAST_MATH = YES;
				SDKROOT = iphoneos;
				SWIFT_COMPILATION_MODE = wholemodule;
				SWIFT_OPTIMIZATION_LEVEL = "-O";
				VALIDATE_PRODUCT = YES;
			};
			name = Release;
		};
		13917C3D243F43D10087912B /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				CODE_SIGN_ENTITLEMENTS = "Catalyst Sample/Catalyst Sample.entitlements";
				CODE_SIGN_IDENTITY = "iPhone Developer: Dev Portal Bot Bitrise (E89JV3W9K4)";
				"CODE_SIGN_IDENTITY[sdk=macosx*]" = "Mac Developer: Dev Portal Bot Bitrise (E89JV3W9K4)";
				CODE_SIGN_STYLE = Manual;
				DEVELOPMENT_ASSET_PATHS = "\"Catalyst Sample/Preview Content\"";
				DEVELOPMENT_TEAM = 72SA8V3WYL;
				ENABLE_PREVIEWS = YES;
				INFOPLIST_FILE = "Catalyst Sample/Info.plist";
				LD_RUNPATH_SEARCH_PATHS = (
					"$(inherited)",
					"@executable_path/Frameworks",
				);
				PRODUCT_BUNDLE_IDENTIFIER = "io.bitrise.Catalyst-Sample";
				PRODUCT_NAME = "$(TARGET_NAME)";
				PROVISIONING_PROFILE_SPECIFIER = "development-io-bitrise-ios";
				"PROVISIONING_PROFILE_SPECIFIER[sdk=macosx*]" = "development-io-bitrise-macos";
				SUPPORTS_MACCATALYST = YES;
				SWIFT_VERSION = 5.0;
				TARGETED_DEVICE_FAMILY = "1,2";
			};
			name = Debug;
		};
		13917C3E243F43D10087912B /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				CODE_SIGN_ENTITLEMENTS = "Catalyst Sample/Catalyst Sample.entitlements";
				CODE_SIGN_IDENTITY = "iPhone Developer: Dev Portal Bot Bitrise (E89JV3W9K4)";
				"CODE_SIGN_IDENTITY[sdk=macosx*]" = "Mac Developer: Dev Portal Bot Bitrise (E89JV3W9K4)";
				CODE_SIGN_STYLE = Manual;
				DEVELOPMENT_ASSET_PATHS = "\"Catalyst Sample/Preview Content\"";
				DEVELOPMENT_TEAM = 72SA8V3WYL;
				ENABLE_PREVIEWS = YES;
				INFOPLIST_FILE = "Catalyst Sample/Info.plist";
				LD_RUNPATH_SEARCH_PATHS = (
					"$(inherited)",
					"@executable_path/Frameworks",
				);
				PRODUCT_BUNDLE_IDENTIFIER = "io.bitrise.Catalyst-Sample";
				PRODUCT_NAME = "$(TARGET_NAME)";
				PROVISIONING_PROFILE_SPECIFIER = "development-io-bitrise-ios";
				"PROVISIONING_PROFILE_SPECIFIER[sdk=macosx*]" = "development-io-bitrise-macos";
				SUPPORTS_MACCATALYST = YES;
				SWIFT_VERSION = 5.0;
				TARGETED_DEVICE_FAMILY = "1,2";
			};
			name = Release;
		};
		13917C40243F43D10087912B /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				BUNDLE_LOADER = "$(TEST_HOST)";
				CODE_SIGN_STYLE = Automatic;
				DEVELOPMENT_TEAM = DT2C2FZ7U2;
				INFOPLIST_FILE = "Catalyst SampleTests/Info.plist";
				IPHONEOS_DEPLOYMENT_TARGET = 13.4;
				LD_RUNPATH_SEARCH_PATHS = (
					"$(inherited)",
					"@executable_path/Frameworks",
					"@loader_path/Frameworks",
				);
				PRODUCT_BUNDLE_IDENTIFIER = "io.bitrise.Catalyst-SampleTests";
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 5.0;
				TARGETED_DEVICE_FAMILY = "1,2";
				TEST_HOST = "$(BUILT_PRODUCTS_DIR)/Catalyst Sample.app/Catalyst Sample";
			};
			name = Debug;
		};
		13917C41243F43D10087912B /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				BUNDLE_LOADER = "$(TEST_HOST)";
				CODE_SIGN_STYLE = Automatic;
				DEVELOPMENT_TEAM = DT2C2FZ7U2;
				INFOPLIST_FILE = "Catalyst SampleTests/Info.plist";
				IPHONEOS_DEPLOYMENT_TARGET = 13.4;
				LD_RUNPATH_SEARCH_PATHS = (
					"$(inherited)",
					"@executable_path/Frameworks",
					"@loader_path/Frameworks",
				);
				PRODUCT_BUNDLE_IDENTIFIER = "io.bitrise.Catalyst-SampleTests";
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 5.0;
				TARGETED_DEVICE_FAMILY = "1,2";
				TEST_HOST = "$(BUILT_PRODUCTS_DIR)/Catalyst Sample.app/Catalyst Sample";
			};
			name = Release;
		};
		13917C43243F43D10087912B /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				CODE_SIGN_STYLE = Automatic;
				DEVELOPMENT_TEAM = DT2C2FZ7U2;
				INFOPLIST_FILE = "Catalyst SampleUITests/Info.plist";
				LD_RUNPATH_SEARCH_PATHS = (
					"$(inherited)",
					"@executable_path/Frameworks",
					"@loader_path/Frameworks",
				);
				PRODUCT_BUNDLE_IDENTIFIER = "io.bitrise.Catalyst-SampleUITests";
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 5.0;
				TARGETED_DEVICE_FAMILY = "1,2";
				TEST_TARGET_NAME = "Catalyst Sample";
			};
			name = Debug;
		};
		13917C44243F43D10087912B /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				CODE_SIGN_STYLE = Automatic;
				DEVELOPMENT_TEAM = DT2C2FZ7U2;
				INFOPLIST_FILE = "Catalyst SampleUITests/Info.plist";
				LD_RUNPATH_SEARCH_PATHS = (
					"$(inherited)",
					"@executable_path/Frameworks",
					"@loader_path/Frameworks",
				);
				PRODUCT_BUNDLE_IDENTIFIER = "io.bitrise.Catalyst-SampleUITests";
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 5.0;
				TARGETED_DEVICE_FAMILY = "1,2";
				TEST_TARGET_NAME = "Catalyst Sample";
			};
			name = Release;
		};
/* End XCBuildConfiguration section */

/* Begin XCConfigurationList section */
		13917C0D243F43D00087912B /* Build configuration list for PBXProject "Catalyst Sample" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13917C3A243F43D10087912B /* Debug */,
				13917C3B243F43D10087912B /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13917C3C243F43D10087912B /* Build configuration list for PBXNativeTarget "Catalyst Sample" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13917C3D243F43D10087912B /* Debug */,
				13917C3E243F43D10087912B /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13917C3F243F43D10087912B /* Build configuration list for PBXNativeTarget "Catalyst SampleTests" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13917C40243F43D10087912B /* Debug */,
				13917C41243F43D10087912B /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13917C42243F43D10087912B /* Build configuration list for PBXNativeTarget "Catalyst SampleUITests" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13917C43243F43D10087912B /* Debug */,
				13917C44243F43D10087912B /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
}
`
