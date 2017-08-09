package constants

import "fmt"

const (
	// MsbuildPath ...
	MsbuildPath = "/Library/Frameworks/Mono.framework/Versions/Current/Commands/msbuild"
	// XbuildPath ...
	XbuildPath = "/Library/Frameworks/Mono.framework/Versions/Current/Commands/xbuild"

	// MonoPath ...
	MonoPath = "/Library/Frameworks/Mono.framework/Versions/Current/Commands/mono"
)

const (
	// SolutionExt ...
	SolutionExt = ".sln"
	// CSProjExt ...
	CSProjExt = ".csproj"
	// FSProjExt ...
	FSProjExt = ".fsproj"
	// SHProjExt ...
	SHProjExt = ".shproj"
)

// SDK ...
type SDK string

const (
	// SDKUnknown ...
	SDKUnknown SDK = "unkown"
	// SDKAndroid ...
	SDKAndroid SDK = "android"
	// SDKIOS ...
	SDKIOS SDK = "ios"
	// SDKTvOS ...
	SDKTvOS SDK = "tvos"
	// SDKMacOS ...
	SDKMacOS SDK = "macos"
)

// ParseSDK ...
func ParseSDK(sdk string) (SDK, error) {
	switch sdk {
	case "android":
		return SDKAndroid, nil
	case "ios":
		return SDKIOS, nil
	case "tvos":
		return SDKTvOS, nil
	case "macos":
		return SDKMacOS, nil
	default:
		return SDKUnknown, fmt.Errorf("invalid sdk: %s", sdk)
	}
}

// TestFramework ...
type TestFramework string

const (
	// TestFrameworkUnknown ...
	TestFrameworkUnknown TestFramework = "unkown"
	// TestFrameworkXamarinUITest ...
	TestFrameworkXamarinUITest TestFramework = "xamarin-uitest"
	// TestFrameworkNunitTest ...
	TestFrameworkNunitTest TestFramework = "nunit-test"
	// TestFrameworkNunitLiteTest ...
	TestFrameworkNunitLiteTest TestFramework = "nunit-lite-test"
)

// ParseTestFramwork ...
func ParseTestFramwork(testFramwork string) (TestFramework, error) {
	switch testFramwork {
	case "xamarin-uitest":
		return TestFrameworkXamarinUITest, nil
	case "nunit-test":
		return TestFrameworkNunitTest, nil
	case "nunit-lite-test":
		return TestFrameworkNunitLiteTest, nil
	default:
		return TestFrameworkUnknown, fmt.Errorf("invalid test framwork: %s", testFramwork)
	}
}

// ParseProjectTypeGUID ...
func ParseProjectTypeGUID(guid string) (SDK, error) {
	switch guid {
	case "EFBA0AD7-5A72-4C68-AF49-83D382785DCF",
		"10368E6C-D01B-4462-8E8B-01FC667A7035": // XamarinAndroid
		return SDKAndroid, nil
	case "E613F3A2-FE9C-494F-B74E-F63BCB86FEA6", // XamarinIOS
		"6BC8ED88-2882-458C-8E55-DFD12B67127B",
		"F5B4F3BC-B597-4E2B-B552-EF5D8A32436F",
		"FEACFBD2-3405-455C-9665-78FE426C6842",
		"8FFB629D-F513-41CE-95D2-7ECE97B6EEEC",
		"EE2C853D-36AF-4FDB-B1AD-8E90477E2198":
		return SDKIOS, nil
	case "06FA79CB-D6CD-4721-BB4B-1BD202089C55": // XamarinProjectTypeTvOS
		return SDKTvOS, nil
	case "1C533B1C-72DD-4CB1-9F6B-BF11D93BCFBE", // MonoMac
		"948B3504-5B70-4649-8FE4-BDE1FB46EC69",
		"42C0BBD9-55CE-4FC1-8D90-A7348ABAFB23", // XamarinMac
		"A3F8F2AB-B479-4A4A-A458-A89E7DC349F1":
		return SDKMacOS, nil
	default:
		return SDKUnknown, fmt.Errorf("Can not identify guid: %s", guid)
	}
}

// OutputType ...
type OutputType string

const (
	// OutputTypeUnknown ...
	OutputTypeUnknown OutputType = "unknown"
	// OutputTypeAPK ...
	OutputTypeAPK OutputType = "apk"
	// OutputTypeXCArchive ...
	OutputTypeXCArchive OutputType = "xcarchive"
	// OutputTypeIPA ...
	OutputTypeIPA OutputType = "ipa"
	// OutputTypeDSYM ...
	OutputTypeDSYM OutputType = "dsym"
	// OutputTypePKG ...
	OutputTypePKG OutputType = "pkg"
	// OutputTypeAPP ...
	OutputTypeAPP OutputType = "app"
	// OutputTypeDLL ...
	OutputTypeDLL OutputType = "dll"
)

// ParseOutputType ...
func ParseOutputType(outputType string) (OutputType, error) {
	switch outputType {
	case "apk":
		return OutputTypeAPK, nil
	case "xcarchive":
		return OutputTypeXCArchive, nil
	case "ipa":
		return OutputTypeIPA, nil
	case "dsym":
		return OutputTypeDSYM, nil
	case "pkg":
		return OutputTypePKG, nil
	case "app":
		return OutputTypeAPP, nil
	case "dll":
		return OutputTypeDLL, nil
	default:
		return OutputTypeUnknown, fmt.Errorf("invalid output type: %s", outputType)
	}
}
