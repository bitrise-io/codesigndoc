package xcarchive

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xcode/plistutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
	"github.com/bitrise-tools/go-xcode/utility"
)

type iosBaseApplication struct {
	Path                string
	InfoPlist           plistutil.PlistData
	Entitlements        plistutil.PlistData
	ProvisioningProfile profileutil.ProvisioningProfileInfoModel
}

// BundleIdentifier ...
func (app iosBaseApplication) BundleIdentifier() string {
	bundleID, _ := app.InfoPlist.GetString("CFBundleIdentifier")
	return bundleID
}

func newIosBaseApplication(path string) (iosBaseApplication, error) {
	infoPlist := plistutil.PlistData{}
	{
		infoPlistPath := filepath.Join(path, "Info.plist")
		if exist, err := pathutil.IsPathExists(infoPlistPath); err != nil {
			return iosBaseApplication{}, fmt.Errorf("failed to check if Info.plist exists at: %s, error: %s", infoPlistPath, err)
		} else if !exist {
			return iosBaseApplication{}, fmt.Errorf("Info.plist not exists at: %s", infoPlistPath)
		}
		plist, err := plistutil.NewPlistDataFromFile(infoPlistPath)
		if err != nil {
			return iosBaseApplication{}, err
		}
		infoPlist = plist
	}

	provisioningProfile := profileutil.ProvisioningProfileInfoModel{}
	{
		provisioningProfilePath := filepath.Join(path, "embedded.mobileprovision")
		if exist, err := pathutil.IsPathExists(provisioningProfilePath); err != nil {
			return iosBaseApplication{}, fmt.Errorf("failed to check if profile exists at: %s, error: %s", provisioningProfilePath, err)
		} else if !exist {
			return iosBaseApplication{}, fmt.Errorf("profile not exists at: %s", provisioningProfilePath)
		}

		profile, err := profileutil.NewProvisioningProfileInfoFromFile(provisioningProfilePath)
		if err != nil {
			return iosBaseApplication{}, err
		}
		provisioningProfile = profile
	}

	entitlements := plistutil.PlistData{}
	{
		entitlementsPath := filepath.Join(path, "archived-expanded-entitlements.xcent")
		if exist, err := pathutil.IsPathExists(entitlementsPath); err != nil {
			return iosBaseApplication{}, fmt.Errorf("failed to check if entitlements exists at: %s, error: %s", entitlementsPath, err)
		} else if exist {
			plist, err := plistutil.NewPlistDataFromFile(entitlementsPath)
			if err != nil {
				return iosBaseApplication{}, err
			}
			entitlements = plist
		}
	}

	return iosBaseApplication{
		Path:                path,
		InfoPlist:           infoPlist,
		Entitlements:        entitlements,
		ProvisioningProfile: provisioningProfile,
	}, nil
}

// IosExtension ...
type IosExtension struct {
	iosBaseApplication
}

// NewIosExtension ...
func NewIosExtension(path string) (IosExtension, error) {
	baseApp, err := newIosBaseApplication(path)
	if err != nil {
		return IosExtension{}, err
	}

	return IosExtension{
		baseApp,
	}, nil
}

// IosWatchApplication ...
type IosWatchApplication struct {
	iosBaseApplication
	Extensions []IosExtension
}

// NewIosWatchApplication ...
func NewIosWatchApplication(path string) (IosWatchApplication, error) {
	baseApp, err := newIosBaseApplication(path)
	if err != nil {
		return IosWatchApplication{}, err
	}

	extensions := []IosExtension{}
	pattern := filepath.Join(path, "PlugIns/*.appex")
	pths, err := filepath.Glob(pattern)
	if err != nil {
		return IosWatchApplication{}, fmt.Errorf("failed to search for watch application's extensions using pattern: %s, error: %s", pattern, err)
	}
	for _, pth := range pths {
		extension, err := NewIosExtension(pth)
		if err != nil {
			return IosWatchApplication{}, err
		}

		extensions = append(extensions, extension)
	}

	return IosWatchApplication{
		iosBaseApplication: baseApp,
		Extensions:         extensions,
	}, nil
}

// IosApplication ...
type IosApplication struct {
	iosBaseApplication
	WatchApplication *IosWatchApplication
	Extensions       []IosExtension
}

// NewIosApplication ...
func NewIosApplication(path string) (IosApplication, error) {
	baseApp, err := newIosBaseApplication(path)
	if err != nil {
		return IosApplication{}, err
	}

	var watchApp *IosWatchApplication
	{
		pattern := filepath.Join(path, "Watch/*.app")
		pths, err := filepath.Glob(pattern)
		if err != nil {
			return IosApplication{}, err
		}
		if len(pths) > 0 {
			watchPath := pths[0]
			app, err := NewIosWatchApplication(watchPath)
			if err != nil {
				return IosApplication{}, err
			}
			watchApp = &app
		}
	}

	extensions := []IosExtension{}
	{
		pattern := filepath.Join(path, "PlugIns/*.appex")
		pths, err := filepath.Glob(pattern)
		if err != nil {
			return IosApplication{}, fmt.Errorf("failed to search for watch application's extensions using pattern: %s, error: %s", pattern, err)
		}
		for _, pth := range pths {
			extension, err := NewIosExtension(pth)
			if err != nil {
				return IosApplication{}, err
			}

			extensions = append(extensions, extension)
		}
	}

	return IosApplication{
		iosBaseApplication: baseApp,
		WatchApplication:   watchApp,
		Extensions:         extensions,
	}, nil
}

// IosArchive ...
type IosArchive struct {
	Path        string
	InfoPlist   plistutil.PlistData
	Application IosApplication
}

// NewIosArchive ...
func NewIosArchive(path string) (IosArchive, error) {
	infoPlist := plistutil.PlistData{}
	{
		infoPlistPath := filepath.Join(path, "Info.plist")
		if exist, err := pathutil.IsPathExists(infoPlistPath); err != nil {
			return IosArchive{}, fmt.Errorf("failed to check if Info.plist exists at: %s, error: %s", infoPlistPath, err)
		} else if !exist {
			return IosArchive{}, fmt.Errorf("Info.plist not exists at: %s", infoPlistPath)
		}
		plist, err := plistutil.NewPlistDataFromFile(infoPlistPath)
		if err != nil {
			return IosArchive{}, err
		}
		infoPlist = plist
	}

	application := IosApplication{}
	{
		pattern := filepath.Join(path, "Products/Applications/*.app")
		pths, err := filepath.Glob(pattern)
		if err != nil {
			return IosArchive{}, err
		}

		appPath := ""
		if len(pths) > 0 {
			appPath = pths[0]
		} else {
			return IosArchive{}, fmt.Errorf("failed to find main app, using pattern: %s", pattern)
		}

		app, err := NewIosApplication(appPath)
		if err != nil {
			return IosArchive{}, err
		}
		application = app
	}

	return IosArchive{
		Path:        path,
		InfoPlist:   infoPlist,
		Application: application,
	}, nil
}

// IsXcodeManaged ...
func (archive IosArchive) IsXcodeManaged() bool {
	return archive.Application.ProvisioningProfile.IsXcodeManaged()
}

// SigningIdentity ...
func (archive IosArchive) SigningIdentity() string {
	properties, found := archive.InfoPlist.GetMapStringInterface("ApplicationProperties")
	if found {
		identity, _ := properties.GetString("SigningIdentity")
		return identity
	}
	return ""
}

// BundleIDEntitlementsMap ...
func (archive IosArchive) BundleIDEntitlementsMap() map[string]plistutil.PlistData {
	bundleIDEntitlementsMap := map[string]plistutil.PlistData{}

	bundleID := archive.Application.BundleIdentifier()
	bundleIDEntitlementsMap[bundleID] = archive.Application.Entitlements

	for _, plugin := range archive.Application.Extensions {
		bundleID := plugin.BundleIdentifier()
		bundleIDEntitlementsMap[bundleID] = plugin.Entitlements
	}

	if archive.Application.WatchApplication != nil {
		watchApplication := *archive.Application.WatchApplication

		bundleID := watchApplication.BundleIdentifier()
		bundleIDEntitlementsMap[bundleID] = watchApplication.Entitlements

		for _, plugin := range watchApplication.Extensions {
			bundleID := plugin.BundleIdentifier()
			bundleIDEntitlementsMap[bundleID] = plugin.Entitlements
		}
	}

	return bundleIDEntitlementsMap
}

// BundleIDProfileInfoMap ...
func (archive IosArchive) BundleIDProfileInfoMap() map[string]profileutil.ProvisioningProfileInfoModel {
	bundleIDProfileMap := map[string]profileutil.ProvisioningProfileInfoModel{}

	bundleID := archive.Application.BundleIdentifier()
	bundleIDProfileMap[bundleID] = archive.Application.ProvisioningProfile

	for _, plugin := range archive.Application.Extensions {
		bundleID := plugin.BundleIdentifier()
		bundleIDProfileMap[bundleID] = plugin.ProvisioningProfile
	}

	if archive.Application.WatchApplication != nil {
		watchApplication := *archive.Application.WatchApplication

		bundleID := watchApplication.BundleIdentifier()
		bundleIDProfileMap[bundleID] = watchApplication.ProvisioningProfile

		for _, plugin := range watchApplication.Extensions {
			bundleID := plugin.BundleIdentifier()
			bundleIDProfileMap[bundleID] = plugin.ProvisioningProfile
		}
	}

	return bundleIDProfileMap
}

// FindDSYMs ...
func (archive IosArchive) FindDSYMs() (string, []string, error) {
	dsymsDirPth := filepath.Join(archive.Path, "dSYMs")
	dsyms, err := utility.ListEntries(dsymsDirPth, utility.ExtensionFilter(".dsym", true))
	if err != nil {
		return "", []string{}, err
	}

	appDSYM := ""
	frameworkDSYMs := []string{}
	for _, dsym := range dsyms {
		if strings.HasSuffix(dsym, ".app.dSYM") {
			appDSYM = dsym
		} else {
			frameworkDSYMs = append(frameworkDSYMs, dsym)
		}
	}
	if appDSYM == "" && len(frameworkDSYMs) == 0 {
		return "", []string{}, fmt.Errorf("no dsym found")
	}

	return appDSYM, frameworkDSYMs, nil
}
