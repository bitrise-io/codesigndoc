package codesigndocuitests

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xcode/plistutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

// IOSTestRunner ...
type IOSTestRunner struct {
	Path                string
	InfoPlist           plistutil.PlistData
	Entitlements        plistutil.PlistData
	ProvisioningProfile profileutil.ProvisioningProfileInfoModel
}

// NewIOSTestRunner ...
func NewIOSTestRunner(path string) (*IOSTestRunner, error) {
	runnerPattern := filepath.Join(path, "*-Runner.app")
	possibleTestRunnerPths, err := filepath.Glob(runnerPattern)
	if err != nil {
		return nil, err
	}

	if len(possibleTestRunnerPths) == 0 {
		return nil, fmt.Errorf("no Test-Runner.app found in %s", path)
	} else if len(possibleTestRunnerPths) != 1 {
		return nil, fmt.Errorf("found multiple Test-Runner.app in %s", path)
	}

	infoPlist := plistutil.PlistData{}
	{
		infoPlistPath := filepath.Join(possibleTestRunnerPths[0], "Info.plist")
		if exist, err := pathutil.IsPathExists(infoPlistPath); err != nil {
			return nil, fmt.Errorf("failed to check if Info.plist exists at: %s, error: %s", infoPlistPath, err)
		} else if !exist {
			return nil, fmt.Errorf("Info.plist not exists at: %s", infoPlistPath)
		}
		plist, err := plistutil.NewPlistDataFromFile(infoPlistPath)
		if err != nil {
			return nil, err
		}
		infoPlist = plist
	}

	provisioningProfile := profileutil.ProvisioningProfileInfoModel{}
	{
		provisioningProfilePath := filepath.Join(possibleTestRunnerPths[0], "embedded.mobileprovision")
		if exist, err := pathutil.IsPathExists(provisioningProfilePath); err != nil {
			return nil, fmt.Errorf("failed to check if profile exists at: %s, error: %s", provisioningProfilePath, err)
		} else if !exist {
			return nil, fmt.Errorf("profile not exists at: %s", provisioningProfilePath)
		}

		profile, err := profileutil.NewProvisioningProfileInfoFromFile(provisioningProfilePath)
		if err != nil {
			return nil, err
		}
		provisioningProfile = profile
	}

	entitlements := plistutil.PlistData{}
	{
		cmd := command.New("codesign", "-d", "--entitlements", "-", possibleTestRunnerPths[0])
		out, err := cmd.RunAndReturnTrimmedOutput()
		if err != nil {
			return nil, err
		}

		outSplit := strings.Split(out, "<?xml version")
		if len(outSplit) > 1 {
			out = outSplit[1]
		}

		entitlements, err = plistutil.NewPlistDataFromContent(out)
		if err != nil {
			return nil, err
		}
	}

	return &IOSTestRunner{
		Path:                path,
		InfoPlist:           infoPlist,
		Entitlements:        entitlements,
		ProvisioningProfile: provisioningProfile,
	}, nil
}

// BundleIDEntitlementsMap ...
func (runner IOSTestRunner) BundleIDEntitlementsMap() map[string]plistutil.PlistData {
	bundleIDEntitlementsMap := map[string]plistutil.PlistData{}

	bundleID := strings.TrimSuffix(runner.ProvisioningProfile.BundleID, "-Runner")
	bundleIDEntitlementsMap[bundleID] = runner.ProvisioningProfile.Entitlements

	log.Warnf("bundleIDEntitlementsMap: %+v", bundleIDEntitlementsMap)

	return bundleIDEntitlementsMap
}

// IsXcodeManaged ...
func (runner IOSTestRunner) IsXcodeManaged() bool {
	return runner.ProvisioningProfile.IsXcodeManaged()
}
