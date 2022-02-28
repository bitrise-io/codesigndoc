package utility

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/profileutil"
	"github.com/bitrise-io/go-xcode/xcodeproject/schemeint"
	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
)

// ProfileExportFileNameNoPath creates a file name for the given profile with pattern: uuid.escaped_profile_name.[mobileprovision|provisionprofile].
func ProfileExportFileNameNoPath(info profileutil.ProvisioningProfileInfoModel) string {
	replaceRexp, err := regexp.Compile("[^A-Za-z0-9_.-]")
	if err != nil {
		log.Warnf("Invalid regex, error: %s", err)
		return ""
	}
	safeTitle := replaceRexp.ReplaceAllString(info.Name, "")
	extension := ".mobileprovision"
	if info.Type == profileutil.ProfileTypeMacOs {
		extension = ".provisionprofile"
	}

	return info.UUID + "." + safeTitle + extension
}

// Platform ...
type Platform string

const (
	iOS     Platform = "iOS"
	osX     Platform = "OS X"
	tvOS    Platform = "tvOS"
	watchOS Platform = "watchOS"
)

// TargetBuildSettingsProvider ...
type TargetBuildSettingsProvider interface {
	TargetBuildSettings(xcodeProj *xcodeproj.XcodeProj, target, configuration string, customOptions ...string) (serialized.Object, error)
}

// XcodeBuild ...
type XcodeBuild struct {
}

// TargetBuildSettings ...
func (x XcodeBuild) TargetBuildSettings(xcodeProj *xcodeproj.XcodeProj, target, configuration string, customOptions ...string) (serialized.Object, error) {
	return xcodeProj.TargetBuildSettings(target, configuration, customOptions...)
}

// BuildableTargetPlatform ...
func BuildableTargetPlatform(
	xcodeProj *xcodeproj.XcodeProj,
	scheme *xcscheme.Scheme,
	configurationName string,
	provider TargetBuildSettingsProvider,
) (Platform, error) {
	archiveEntry, ok := scheme.AppBuildActionEntry()
	if !ok {
		return "", fmt.Errorf("archivable entry not found in project: %s, scheme: %s", xcodeProj.Path, scheme.Name)
	}

	mainTarget, ok := xcodeProj.Proj.Target(archiveEntry.BuildableReference.BlueprintIdentifier)
	if !ok {
		return "", fmt.Errorf("target not found: %s", archiveEntry.BuildableReference.BlueprintIdentifier)
	}

	settings, err := provider.TargetBuildSettings(xcodeProj, mainTarget.Name, configurationName)
	if err != nil {
		return "", fmt.Errorf("failed to get target (%s) build settings: %s", mainTarget.Name, err)
	}

	return getPlatform(settings)
}

func getPlatform(buildSettings serialized.Object) (Platform, error) {
	/*
		Xcode help:
		Base SDK (SDKROOT)
		The name or path of the base SDK being used during the build.
		The product will be built against the headers and libraries located inside the indicated SDK.
		This path will be prepended to all search paths, and will be passed through the environment to the compiler and linker.
		Additional SDKs can be specified in the Additional SDKs (ADDITIONAL_SDKS) setting.

		Examples:
		- /Applications/Xcode.app/Contents/Developer/Platforms/AppleTVOS.platform/Developer/SDKs/AppleTVOS.sdk
		- /Applications/Xcode.app/Contents/Developer/Platforms/AppleTVSimulator.platform/Developer/SDKs/AppleTVSimulator13.4.sdk
		- /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS13.4.sdk
		- /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/SDKs/iPhoneSimulator.sdk
		- /Applications/Xcode.app/Contents/Developer/Platforms/MacOSX.platform/Developer/SDKs/MacOSX10.15.sdk
		- /Applications/Xcode.app/Contents/Developer/Platforms/WatchOS.platform/Developer/SDKs/WatchOS.sdk
		- /Applications/Xcode.app/Contents/Developer/Platforms/WatchSimulator.platform/Developer/SDKs/WatchSimulator.sdk
		- iphoneos
		- macosx
		- appletvos
		- watchos
	*/
	sdk, err := buildSettings.String("SDKROOT")
	if err != nil {
		return "", fmt.Errorf("failed to get SDKROOT: %s", err)
	}

	sdk = strings.ToLower(sdk)
	if filepath.Ext(sdk) == ".sdk" {
		sdk = filepath.Base(sdk)
	}

	switch {
	case strings.HasPrefix(sdk, "iphoneos"):
		return iOS, nil
	case strings.HasPrefix(sdk, "macosx"):
		return osX, nil
	case strings.HasPrefix(sdk, "appletvos"):
		return tvOS, nil
	case strings.HasPrefix(sdk, "watchos"):
		return watchOS, nil
	default:
		return "", fmt.Errorf("unknown SDKROOT: %s", sdk)
	}
}

// OpenArchivableProject ...
func OpenArchivableProject(pth, schemeName, configurationName string) (*xcodeproj.XcodeProj, *xcscheme.Scheme, string, error) {
	scheme, schemeContainerDir, err := schemeint.Scheme(pth, schemeName)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not get scheme (%s) from path (%s): %s", schemeName, pth, err)
	}
	if configurationName == "" {
		configurationName = scheme.ArchiveAction.BuildConfiguration
	}

	if configurationName == "" {
		return nil, nil, "", fmt.Errorf("no configuration provided nor default defined for the scheme's (%s) archive action", schemeName)
	}

	archiveEntry, ok := scheme.AppBuildActionEntry()
	if !ok {
		return nil, nil, "", fmt.Errorf("archivable entry not found")
	}

	projectPth, err := archiveEntry.BuildableReference.ReferencedContainerAbsPath(filepath.Dir(schemeContainerDir))
	if err != nil {
		return nil, nil, "", err
	}

	xcodeProj, err := xcodeproj.Open(projectPth)
	if err != nil {
		return nil, nil, "", err
	}
	return &xcodeProj, scheme, configurationName, nil
}
