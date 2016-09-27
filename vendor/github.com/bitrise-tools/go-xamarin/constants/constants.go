package constants

import "fmt"

const (
	// MDToolPath ...
	MDToolPath = "/Applications/Xamarin Studio.app/Contents/MacOS/mdtool"
	// XbuildPath ...
	XbuildPath = "/Library/Frameworks/Mono.framework/Commands/xbuild"
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

// TestFramwork ...
type TestFramwork string

const (
	// XamarinUITest ...
	XamarinUITest TestFramwork = "Xamarin.UITest"
	// NunitTest ...
	NunitTest TestFramwork = "nunit.framework"
	// NunitLiteTest ...
	NunitLiteTest TestFramwork = "MonoTouch.NUnitLite"
)

// OutputType ...
type OutputType string

const (
	// OutputTypeUnknown ...
	OutputTypeUnknown OutputType = "unknown"
	// OutputTypeXCArchive ...
	OutputTypeXCArchive OutputType = "xcarchive"
	// OutputTypeIPA ...
	OutputTypeIPA OutputType = "ipa"
	// OutputTypeAPK ...
	OutputTypeAPK OutputType = "apk"
	// OutputTypePKG ...
	OutputTypePKG OutputType = "pkg"
	// OutputTypeAPP ...
	OutputTypeAPP OutputType = "app"
)

// ParseOutputType ...
func ParseOutputType(outputType string) (OutputType, error) {
	switch outputType {
	case "xcarchive":
		return OutputTypeXCArchive, nil
	case "ipa":
		return OutputTypeIPA, nil
	case "apk":
		return OutputTypeAPK, nil
	case "pkg":
		return OutputTypePKG, nil
	case "app":
		return OutputTypeAPP, nil
	default:
		return OutputTypeUnknown, fmt.Errorf("invalid output type: %s", outputType)
	}
}

// ProjectType ...
type ProjectType string

const (
	// ProjectTypeUnknown ...
	ProjectTypeUnknown ProjectType = "unknown"
	// ProjectTypeIos ...
	ProjectTypeIos ProjectType = "ios"
	// ProjectTypeAndroid ...
	ProjectTypeAndroid ProjectType = "android"
	// ProjectTypeMac ...
	ProjectTypeMac ProjectType = "mac"
	// ProjectTypeTVOs ...
	ProjectTypeTVOs ProjectType = "tvos"
)

// ParseProjectType ...
func ParseProjectType(projectType string) (ProjectType, error) {
	switch projectType {
	case "ios":
		return ProjectTypeIos, nil
	case "android":
		return ProjectTypeAndroid, nil
	case "mac":
		return ProjectTypeMac, nil
	case "tvos":
		return ProjectTypeTVOs, nil
	default:
		return ProjectTypeUnknown, fmt.Errorf("invalid project type: %s", projectType)
	}
}

// ParseProjectTypeGUID ...
func ParseProjectTypeGUID(guid string) (ProjectType, error) {
	fmt.Printf("guid: %s\n", guid)
	switch guid {
	case "E613F3A2-FE9C-494F-B74E-F63BCB86FEA6", // XamarinIos
		"6BC8ED88-2882-458C-8E55-DFD12B67127B",
		"F5B4F3BC-B597-4E2B-B552-EF5D8A32436F",
		"FEACFBD2-3405-455C-9665-78FE426C6842",
		"8FFB629D-F513-41CE-95D2-7ECE97B6EEEC",
		"EE2C853D-36AF-4FDB-B1AD-8E90477E2198":
		return ProjectTypeIos, nil
	case "EFBA0AD7-5A72-4C68-AF49-83D382785DCF", "10368E6C-D01B-4462-8E8B-01FC667A7035": // XamarinAndroid
		return ProjectTypeAndroid, nil
	case "1C533B1C-72DD-4CB1-9F6B-BF11D93BCFBE", "948B3504-5B70-4649-8FE4-BDE1FB46EC69", // MonoMac
		"42C0BBD9-55CE-4FC1-8D90-A7348ABAFB23", "A3F8F2AB-B479-4A4A-A458-A89E7DC349F1": // XamarinMac
		return ProjectTypeMac, nil
	case "06FA79CB-D6CD-4721-BB4B-1BD202089C55": // XamarinProjectTypeTVOS
		return ProjectTypeTVOs, nil
	default:
		return ProjectTypeUnknown, fmt.Errorf("Can not idetify guid: %s", guid)
	}
}
