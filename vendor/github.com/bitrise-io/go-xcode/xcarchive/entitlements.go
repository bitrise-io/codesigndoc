package xcarchive

import (
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-xcode/plistutil"
)

func executableNameFromInfoPlist(infoPlist plistutil.PlistData) string {
	if name, ok := infoPlist.GetString("CFBundleExecutable"); ok {
		return name
	}
	return ""
}

func getEntitlements(basePath, executableRelativePath string) (plistutil.PlistData, error) {
	entitlements, err := entitlementsFromExecutable(basePath, executableRelativePath)
	if err != nil {
		return plistutil.PlistData{}, err
	}

	if entitlements != nil {
		return *entitlements, nil
	}

	return plistutil.PlistData{}, nil
}

func entitlementsFromExecutable(basePath, executableRelativePath string) (*plistutil.PlistData, error) {
	cmd := command.New("codesign", "--display", "--entitlements", ":-", filepath.Join(basePath, executableRelativePath))
	entitlementsString, err := cmd.RunAndReturnTrimmedOutput()
	if err != nil {
		return nil, err
	}

	plist, err := plistutil.NewPlistDataFromContent(entitlementsString)
	if err != nil {
		return nil, err
	}

	return &plist, nil
}
