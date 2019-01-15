package xcodeproj

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/bitrise-tools/xcode-project/xcodebuild"
	"github.com/bitrise-tools/xcode-project/xcscheme"
	"howett.net/plist"
)

// XcodeProj ...
type XcodeProj struct {
	Proj Proj

	Name string
	Path string
}

func (p XcodeProj) buildSettingsFilePath(target, configuration, key string) (string, error) {
	buildSettings, err := p.TargetBuildSettings(target, configuration)
	if err != nil {
		return "", err
	}

	pth, err := buildSettings.String(key)
	if err != nil {
		return "", err
	}

	if pathutil.IsRelativePath(pth) {
		pth = filepath.Join(filepath.Dir(p.Path), pth)
	}

	return pth, nil
}

// TargetCodeSignEntitlementsPath ...
func (p XcodeProj) TargetCodeSignEntitlementsPath(target, configuration string) (string, error) {
	return p.buildSettingsFilePath(target, configuration, "CODE_SIGN_ENTITLEMENTS")
}

// TargetCodeSignEntitlements ...
func (p XcodeProj) TargetCodeSignEntitlements(target, configuration string) (serialized.Object, error) {
	codeSignEntitlementsPth, err := p.TargetCodeSignEntitlementsPath(target, configuration)
	if err != nil {
		return nil, err
	}

	codeSignEntitlementsContent, err := fileutil.ReadBytesFromFile(codeSignEntitlementsPth)
	if err != nil {
		return nil, err
	}

	var codeSignEntitlements serialized.Object
	if _, err := plist.Unmarshal([]byte(codeSignEntitlementsContent), &codeSignEntitlements); err != nil {
		return nil, err
	}

	return codeSignEntitlements, nil
}

// TargetInformationPropertyListPath ...
func (p XcodeProj) TargetInformationPropertyListPath(target, configuration string) (string, error) {
	return p.buildSettingsFilePath(target, configuration, "INFOPLIST_FILE")
}

// TargetInformationPropertyList ...
func (p XcodeProj) TargetInformationPropertyList(target, configuration string) (serialized.Object, error) {
	informationPropertyListPth, err := p.TargetInformationPropertyListPath(target, configuration)
	if err != nil {
		return nil, err
	}

	informationPropertyListContent, err := fileutil.ReadBytesFromFile(informationPropertyListPth)
	if err != nil {
		return nil, err
	}

	var informationPropertyList serialized.Object
	if _, err := plist.Unmarshal([]byte(informationPropertyListContent), &informationPropertyList); err != nil {
		return nil, err
	}

	return informationPropertyList, nil
}

// TargetBundleID ...
func (p XcodeProj) TargetBundleID(target, configuration string) (string, error) {
	buildSettings, err := p.TargetBuildSettings(target, configuration)
	if err != nil {
		return "", err
	}

	bundleID, err := buildSettings.String("PRODUCT_BUNDLE_IDENTIFIER")
	if err != nil && !serialized.IsKeyNotFoundError(err) {
		return "", err
	}

	if bundleID != "" {
		return resolve(bundleID, buildSettings)
	}

	informationPropertyList, err := p.TargetInformationPropertyList(target, configuration)
	if err != nil {
		return "", err
	}

	bundleID, err = informationPropertyList.String("CFBundleIdentifier")
	if err != nil {
		return "", err
	}

	if bundleID == "" {
		return "", errors.New("no PRODUCT_BUNDLE_IDENTIFIER build settings nor CFBundleIdentifier information property found")
	}

	return resolve(bundleID, buildSettings)
}

func resolve(bundleID string, buildSettings serialized.Object) (string, error) {
	resolvedBundleIDs := map[string]bool{}
	resolved := bundleID
	for true {
		var err error
		resolved, err = expand(resolved, buildSettings)
		if err != nil {
			return "", err
		}

		if !strings.Contains(resolved, "$") {
			return resolved, nil
		}

		_, ok := resolvedBundleIDs[resolved]
		if ok {
			return "", fmt.Errorf("bundle id reference cycle found")
		}
		resolvedBundleIDs[resolved] = true
	}
	return "", fmt.Errorf("failed to resolve bundle id: %s", bundleID)
}

func expand(bundleID string, buildSettings serialized.Object) (string, error) {
	if !strings.Contains(bundleID, "$") {
		return bundleID, nil
	}

	pattern := `(.*)\$\((.*)\)(.*)`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(bundleID)
	if len(match) != 4 {
		return "", fmt.Errorf("%s does not match to pattern: %s", bundleID, pattern)
	}

	prefix := match[1]
	suffix := match[3]
	envKey := match[2]

	split := strings.Split(envKey, ":")
	envKey = split[0]

	envValue, err := buildSettings.String(envKey)
	if err != nil {
		if serialized.IsKeyNotFoundError(err) {
			return "", fmt.Errorf("%s build settings not found", envKey)
		}
		return "", err
	}

	return prefix + envValue + suffix, nil
}

// TargetBuildSettings ...
func (p XcodeProj) TargetBuildSettings(target, configuration string, customOptions ...string) (serialized.Object, error) {
	return xcodebuild.ShowProjectBuildSettings(p.Path, target, configuration, customOptions...)
}

// Scheme ...
func (p XcodeProj) Scheme(name string) (xcscheme.Scheme, bool) {
	schemes, err := p.Schemes()
	if err != nil {
		return xcscheme.Scheme{}, false
	}

	for _, scheme := range schemes {
		if scheme.Name == name {
			return scheme, true
		}
	}

	return xcscheme.Scheme{}, false
}

// Schemes ...
func (p XcodeProj) Schemes() ([]xcscheme.Scheme, error) {
	return xcscheme.FindSchemesIn(p.Path)
}

// Open ...
func Open(pth string) (XcodeProj, error) {
	absPth, err := pathutil.AbsPath(pth)
	if err != nil {
		return XcodeProj{}, err
	}

	pbxProjPth := filepath.Join(absPth, "project.pbxproj")

	b, err := fileutil.ReadBytesFromFile(pbxProjPth)
	if err != nil {
		return XcodeProj{}, err
	}

	var raw serialized.Object
	if _, err := plist.Unmarshal(b, &raw); err != nil {
		return XcodeProj{}, fmt.Errorf("failed to generate json from Pbxproj - error: %s", err)
	}

	objects, err := raw.Object("objects")
	if err != nil {
		return XcodeProj{}, err
	}

	projectID := ""
	for id := range objects {
		object, err := objects.Object(id)
		if err != nil {
			return XcodeProj{}, err
		}

		objectISA, err := object.String("isa")
		if err != nil {
			return XcodeProj{}, err
		}

		if objectISA == "PBXProject" {
			projectID = id
			break
		}
	}

	p, err := parseProj(projectID, objects)
	if err != nil {
		return XcodeProj{}, err
	}

	return XcodeProj{
		Proj: p,
		Path: absPth,
		Name: strings.TrimSuffix(filepath.Base(absPth), filepath.Ext(absPth)),
	}, nil
}

// IsXcodeProj ...
func IsXcodeProj(pth string) bool {
	return filepath.Ext(pth) == ".xcodeproj"
}
