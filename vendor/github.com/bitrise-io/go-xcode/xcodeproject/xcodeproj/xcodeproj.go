package xcodeproj

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	plist "github.com/bitrise-io/go-plist"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/pretty"
	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcodebuild"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
	"golang.org/x/text/unicode/norm"
)

const (
	// XcodeProjExtension ...
	XcodeProjExtension = ".xcodeproj"
)

// XcodeProj ...
type XcodeProj struct {
	Proj    Proj
	RawProj serialized.Object
	Format  int
	// Used to replace project in-place. This leaves the order of objects and comments for unchanged objects unchanged.
	// It allows better compatibility with Cordova and the Xcode agvtool
	originalContents                  []byte
	originalPbxProj, annotatedPbxProj serialized.Object

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

// ForceTargetCodeSignEntitlement updates the project descriptor behind p. It
// searches for the entitlements file for target and configuration and sets the
// entitlement key to the value provided.
// Error is returned:
// - if there's an error during reading or writing the file
// - if the file does not exist for the given target and configuration
func (p XcodeProj) ForceTargetCodeSignEntitlement(target, configuration, entitlement string, value interface{}) error {
	codeSignEntitlementsPth, err := p.TargetCodeSignEntitlementsPath(target, configuration)
	if err != nil {
		return err
	}

	codeSignEntitlements, format, err := ReadPlistFile(codeSignEntitlementsPth)
	if err != nil {
		return err
	}

	codeSignEntitlements[entitlement] = value

	return WritePlistFile(codeSignEntitlementsPth, codeSignEntitlements, format)
}

// TargetCodeSignEntitlements ...
func (p XcodeProj) TargetCodeSignEntitlements(target, configuration string) (serialized.Object, error) {
	codeSignEntitlementsPth, err := p.TargetCodeSignEntitlementsPath(target, configuration)
	if err != nil {
		return nil, err
	}

	codeSignEntitlements, _, err := ReadPlistFile(codeSignEntitlementsPth)
	if err != nil {
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

// ForceTargetBundleID updates the projects bundle ID for the specified target
// and configuration.
// An error is returned if:
// - the target or configuration is not found
// - the given target or configuration is not found
func (p XcodeProj) ForceTargetBundleID(target, configuration, bundleID string) error {
	t, targetFound := p.Proj.TargetByName(target)
	if !targetFound {
		return fmt.Errorf("could not find target (%s)", target)
	}

	var configurationFound bool
	buildConfigurations := t.BuildConfigurationList.BuildConfigurations
	for _, c := range buildConfigurations {
		if c.Name == configuration {
			configurationFound = true
			c.BuildSettings["PRODUCT_BUNDLE_IDENTIFIER"] = bundleID
		}
	}

	if !configurationFound {
		return fmt.Errorf("could not find configuration (%s) for target (%s)", configuration, target)
	}

	return p.Save()
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
		return Resolve(bundleID, buildSettings)
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

	return Resolve(bundleID, buildSettings)
}

// Resolve returns the resolved bundleID. We need this, because the bundleID is not exposed in the .pbxproj file ( raw ).
// If the raw BundleID contains an environment variable we have to replace it.
//
//**Example:**
//BundleID in the .pbxproj: Bitrise.Test.$(PRODUCT_NAME:rfc1034identifier).Suffix
//BundleID after the env is expanded: Bitrise.Test.Sample.Suffix
func Resolve(bundleID string, buildSettings serialized.Object) (string, error) {
	resolvedBundleIDs := map[string]bool{}
	resolved := bundleID
	for true {
		if !strings.Contains(resolved, "$") {
			return resolved, nil
		}

		var err error
		resolved, err = expand(resolved, buildSettings)
		if err != nil {
			return "", err
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
	r, err := regexp.Compile("[$][{(][^$]*?[)}]")
	if err != nil {
		return "", err
	}
	if r.MatchString(bundleID) {
		// envs like:  $(PRODUCT_NAME:rfc1034identifier) || $(PRODUCT_NAME) || ${PRODUCT_NAME:rfc1034identifier} || ${PRODUCT_NAME}
		return expandComplexEnv(bundleID, buildSettings)
	}
	// envs like: $PRODUCT_NAME
	return expandSimpleEnv(bundleID, buildSettings)
}

// expandComplexEnv expands the env with the "[$][{(][^$]*?[)}]" regex
// **Example:** `prefix.$(ENV_KEY:rfc1034identifier).suffix.$(ENV_KEY:rfc1034identifier)` **=>** `auto_provision.ios-simple-objc.suffix.ios-simple-objc`
func expandComplexEnv(bundleID string, buildSettings serialized.Object) (string, error) {
	r, err := regexp.Compile("[$][{(][^$]*?[)}]")
	if err != nil {
		return "", err
	}
	rawEnvKey := r.FindString(bundleID)

	replacer := strings.NewReplacer("$", "", "(", "", ")", "", "{", "", "}", "")
	envKey := strings.Split(replacer.Replace(rawEnvKey), ":")[0]

	envValue, ok := envInBuildSettings(envKey, buildSettings)
	if !ok {
		return "", fmt.Errorf("failed to find env in build settings: %s", envKey)
	}
	return strings.Replace(bundleID, rawEnvKey, envValue, -1), nil
}

// expandSimpleEnv expands the env with the "[$][^$]*" regex
// **Example:** `prefix.$ENV_KEY.suffix.$ENV_KEY` **=>** `auto_provision.ios-simple-objc.suffix.ios-simple-objc`
func expandSimpleEnv(bundleID string, buildSettings serialized.Object) (string, error) {
	r, err := regexp.Compile("[$][^$]*")
	if err != nil {
		return "", err
	}
	if !r.MatchString(bundleID) {
		return "", fmt.Errorf("failed to match regex [$][^$]* for %s", bundleID)
	}
	envKey := r.FindString(bundleID)

	var envValue string
	for len(envKey) > 1 {
		var ok bool
		envValue, ok = envInBuildSettings(strings.Replace(envKey, "$", "", 1), buildSettings)
		if ok {
			break
		}

		envKey = envKey[:len(envKey)-1]
	}
	return strings.Replace(bundleID, envKey, envValue, -1), nil

}

func envInBuildSettings(envKey string, buildSettings serialized.Object) (string, bool) {
	envValue, err := buildSettings.String(envKey)
	if err != nil {
		return "", false
	}
	return envValue, true
}

// TargetBuildSettings ...
func (p XcodeProj) TargetBuildSettings(target, configuration string, customOptions ...string) (serialized.Object, error) {
	return xcodebuild.ShowProjectBuildSettings(p.Path, target, configuration, customOptions...)
}

// Scheme returns the project's scheme by name and the project's absolute path.
func (p XcodeProj) Scheme(name string) (*xcscheme.Scheme, string, error) {
	schemes, err := p.Schemes()
	if err != nil {
		return nil, "", err
	}

	normName := norm.NFC.String(name)
	for _, scheme := range schemes {
		if norm.NFC.String(scheme.Name) == normName {
			return &scheme, p.Path, nil
		}
	}

	return nil, "", xcscheme.NotFoundError{Scheme: name, Container: p.Name}
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

	content, err := fileutil.ReadBytesFromFile(pbxProjPth)
	if err != nil {
		return XcodeProj{}, err
	}

	p, err := parsePBXProjContent(content)
	if err != nil {
		return XcodeProj{}, err
	}

	p.Path = absPth
	p.Name = strings.TrimSuffix(filepath.Base(absPth), filepath.Ext(absPth))

	return *p, nil
}

func parsePBXProjContent(content []byte) (*XcodeProj, error) {
	var rawPbxProj serialized.Object
	format, err := plist.UnmarshalWithCustomAnnotation(content, &rawPbxProj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal project.pbxproj: %s", err)
	}

	annotatedPbxProj := deepCopyObject(rawPbxProj) // Preserve annotations
	rawPbxProj = removeCustomInfoObject(rawPbxProj)
	originalPbxProj := deepCopyObject(rawPbxProj)

	objects, err := rawPbxProj.Object("objects")
	if err != nil {
		return nil, err
	}

	var projectID string
	for id := range objects {
		var object serialized.Object
		object, err = objects.Object(id)
		if err != nil {
			return nil, err
		}

		var objectISA string
		objectISA, err = object.String("isa")
		if err != nil {
			return nil, err
		}

		if objectISA == "PBXProject" {
			projectID = id
			break
		}
	}

	if projectID == "" {
		return nil, fmt.Errorf("failed to find PBXProject's id in project.pbxproj")
	}

	proj, err := parseProj(projectID, objects)
	if err != nil {
		return nil, err
	}

	return &XcodeProj{
		Proj:             proj,
		RawProj:          rawPbxProj,
		Format:           format,
		originalPbxProj:  originalPbxProj,
		annotatedPbxProj: annotatedPbxProj,
		originalContents: content,
	}, nil
}

// IsXcodeProj ...
func IsXcodeProj(pth string) bool {
	return filepath.Ext(pth) == XcodeProjExtension
}

// ForceCodeSign modifies the project's code signing settings to use manual code signing.
//
// Overrides the target's `ProvisioningStyle`, `DevelopmentTeam` and clears the `DevelopmentTeamName` in the **TargetAttributes**.
// Overrides the target's `CODE_SIGN_STYLE`, `DEVELOPMENT_TEAM`, `CODE_SIGN_IDENTITY`, `CODE_SIGN_IDENTITY[sdk=iphoneos*]` `PROVISIONING_PROFILE_SPECIFIER`,
// `PROVISIONING_PROFILE` and `PROVISIONING_PROFILE[sdk=iphoneos*]` in the **BuildSettings**.
func (p *XcodeProj) ForceCodeSign(configuration, targetName, developmentTeam, codesignIdentity, provisioningProfileUUID string) error {
	target, ok := p.Proj.TargetByName(targetName)
	if !ok {
		return fmt.Errorf("failed to find target with name: %s", targetName)
	}

	buildConfigurationList, err := p.BuildConfigurationList(target.ID)
	if err != nil {
		return fmt.Errorf("failed to get target's (%s) buildConfigurationList, error: %s", target.ID, err)
	}
	buildConfigurations, err := p.BuildConfigurations(buildConfigurationList)
	if err != nil {
		return fmt.Errorf("failed to get buildConfigurations of buildConfigurationList (%s), error: %s", pretty.Object(buildConfigurationList), err)
	}

	var buildConfiguration serialized.Object
	for _, b := range buildConfigurations {
		if b["name"] == configuration {
			buildConfiguration = b
			break
		}
	}

	if buildConfiguration == nil {
		return fmt.Errorf("failed to find buildConfiguration for configuration %s in the buildConfiguration list: %s", configuration, pretty.Object(buildConfigurations))
	}

	// Override BuildSettings
	if err = forceCodeSignOnBuildConfiguration(buildConfiguration, developmentTeam, provisioningProfileUUID, codesignIdentity); err != nil {
		return fmt.Errorf("failed to change code signing in build settings, error: %s", err)
	}

	if targetAttributes, err := p.TargetAttributes(); err == nil {
		// Override TargetAttributes
		if err = forceCodeSignOnTargetAttributes(targetAttributes, target.ID, developmentTeam); err != nil {
			return fmt.Errorf("failed to change code signing in target attributes, error: %s", err)
		}
	} else if !serialized.IsKeyNotFoundError(err) {
		return fmt.Errorf("failed to get project's target attributes, error: %s", err)
	}

	return nil
}

// forceCodeSignOnTargetAttributes sets the TargetAttributes for the provided targetID.
// **Overrides the ProvisioningStyle, developmentTeam and clears the DevelopmentTeamName in the provided `targetAttributes`!**
func forceCodeSignOnTargetAttributes(targetAttributes serialized.Object, targetID, developmentTeam string) error {
	targetAttribute, err := targetAttributes.Object(targetID)
	if err != nil {
		// Skip projects not using target attributes
		if serialized.IsKeyNotFoundError(err) {
			return nil
		}
		return fmt.Errorf("failed to get target's (%s) attributes, error: %s", targetID, err)
	}

	targetAttribute["ProvisioningStyle"] = "Manual"
	targetAttribute["DevelopmentTeam"] = developmentTeam
	targetAttribute["DevelopmentTeamName"] = ""
	return nil
}

// forceCodeSignOnBuildConfiguration sets the BuildSettings for the provided build configuration.
// **Overrides the CODE_SIGN_STYLE, DEVELOPMENT_TEAM, CODE_SIGN_IDENTITY, PROVISIONING_PROFILE
// and clears the PROVISIONING_PROFILE_SPECIFIER in the provided `buildConfiguration`,
// each modification also applies for the sdk specific settings too (CODE_SIGN_IDENTITY[sdk=iphoneos*])!**
func forceCodeSignOnBuildConfiguration(buildConfiguration serialized.Object, developmentTeam, provisioningProfileUUID, codesignIdentity string) error {
	buildSettings, err := buildConfiguration.Object("buildSettings")
	if err != nil {
		return fmt.Errorf("failed to get buildSettings of buildConfiguration (%s), error: %s", pretty.Object(buildConfiguration), err)
	}

	forceAttributes := map[string]string{
		"CODE_SIGN_STYLE":                "Manual",
		"DEVELOPMENT_TEAM":               developmentTeam,
		"CODE_SIGN_IDENTITY":             codesignIdentity,
		"PROVISIONING_PROFILE_SPECIFIER": "",
		"PROVISIONING_PROFILE":           provisioningProfileUUID,
	}
	for key, value := range forceAttributes {
		writeAttributeForAllSDKs(buildSettings, key, value)
	}

	return nil
}

func writeAttributeForAllSDKs(buildSettings serialized.Object, newKey string, newValue string) {
	buildSettings[newKey] = newValue

	// override specific build setting if any: https://stackoverflow.com/a/5382708/5842489
	// Example: CODE_SIGN_IDENTITY[sdk=iphoneos*]
	matcher := regexp.MustCompile(fmt.Sprintf(`^%s\[sdk=.*\]$`, regexp.QuoteMeta(newKey)))
	for oldKey := range buildSettings {
		if matcher.Match([]byte(oldKey)) {
			buildSettings[oldKey] = newValue
		}
	}
}

// Save the XcodeProj
//
// Overrides the project.pbxproj file of the XcodeProj with the contents of `rawProj`
func (p XcodeProj) Save() error {
	return p.savePBXProj()
}

// savePBXProj overrides the project.pbxproj file of  the XcodeProj with the contents of `rawProj`
func (p XcodeProj) savePBXProj() error {
	pth := path.Join(p.Path, "project.pbxproj")
	newContent, merr := p.perObjectModify()
	if merr == nil {
		return ioutil.WriteFile(pth, newContent, 0644)
	}
	// merr != nil
	log.Warnf("failed to modify project in-place: %v", merr)

	newContent, err := plist.MarshalIndent(p.RawProj, p.Format, "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal .pbxproj: %v", err)
	}

	return ioutil.WriteFile(pth, newContent, 0644)
}

const (
	customAnnotationKey = plist.CustomAnnotationKey
	startKey            = plist.CustomAnnotationStartKey
	endKey              = plist.CustomAnnotationEndKey
)

func removeCustomInfoObject(o serialized.Object) serialized.Object {
	for _, v := range o {
		removeCustomInfo(v)
	}

	return o
}

func removeCustomInfo(o interface{}) interface{} {
	switch container := o.(type) {
	case map[string]interface{}:
		{
			delete(container, customAnnotationKey)
			for _, val := range container {
				removeCustomInfo(val)
			}

			return container
		}
	case []interface{}:
		{
			for _, element := range container {
				removeCustomInfo(element)
			}

			return container
		}
	default:
		return o
	}
}

func deepCopyObject(object serialized.Object) serialized.Object {
	newObj := make(map[string]interface{})
	for k, v := range object {
		newObj[k] = deepCopy(v)
	}

	return newObj
}

func deepCopy(o interface{}) interface{} {
	switch container := o.(type) {
	case map[string]interface{}:
		{
			newObj := make(map[string]interface{})
			for k, v := range container {
				newObj[k] = deepCopy(v)
			}

			return newObj
		}
	case []interface{}:
		{
			destArray := make([]interface{}, len(container))
			for i, element := range container {
				destArray[i] = deepCopy(element)
			}

			return destArray
		}
	default:
		return o
	}
}

type change struct {
	start, end int
	rawObject  []byte
}

func (p XcodeProj) perObjectModify() ([]byte, error) {
	objectsMod, err := p.RawProj.Object("objects")
	if err != nil {
		return nil, fmt.Errorf("failed to parse project: %v", err)
	}
	objectsOrig, err := p.originalPbxProj.Object("objects")
	if err != nil {
		return nil, fmt.Errorf("failed to parse project: %v", err)
	}
	objectsAnnotated, err := p.annotatedPbxProj.Object("objects")
	if err != nil {
		return nil, fmt.Errorf("failed to parse project: %v", err)
	}

	var mods []change
	for keyMod := range objectsMod {
		objectMod, err := objectsMod.Object(keyMod)
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}

		objectOrig, err := objectsOrig.Object(keyMod)
		if err != nil {
			return nil, fmt.Errorf("new object added, not in original project: %v", err)
		}

		objectsAnnotated, err := objectsAnnotated.Object(keyMod)
		if err != nil {
			return nil, fmt.Errorf("new object added, not in original annotated project: %v", err)
		}

		// If object did not change do nothing
		if reflect.DeepEqual(objectOrig, objectMod) {
			continue
		}

		customPosDict, err := objectsAnnotated.Object(customAnnotationKey)
		if err != nil {
			return nil, fmt.Errorf("no raw object position available: %v", err)
		}
		startPos, err := customPosDict.Int64(startKey)
		if err != nil {
			return nil, fmt.Errorf("no raw object start position available: %v", err)
		}
		endPos, err := customPosDict.Int64(endKey)
		if err != nil {
			return nil, fmt.Errorf("no raw end position availbale: %v", err)
		}

		contentMod, err := plist.MarshalIndent(objectMod, p.Format, "\t")
		if err != nil {
			return nil, fmt.Errorf("could not marshal object (%s): %v", objectsMod, err)
		}

		mods = append(mods, change{
			start:     int(startPos),
			end:       int(endPos),
			rawObject: contentMod,
		})
	}

	if len(mods) == 0 {
		return p.originalContents, nil
	}

	sort.Slice(mods, func(i, j int) bool {
		if mods[i].start == mods[j].start {
			return mods[i].end < mods[j].end
		}
		return mods[i].start < mods[j].start
	})

	var contentsMod []byte
	previousEndPos := 0
	for i, mod := range mods {
		if i < len(mods)-1 && mod.end >= mods[i+1].start {
			return nil, fmt.Errorf("overlapping changes: %d, %d", mods[i].end, mods[i+1].start)
		}

		contentsMod = append(contentsMod, p.originalContents[previousEndPos:mod.start]...)
		contentsMod = append(contentsMod, mod.rawObject...)
		previousEndPos = mod.end
	}

	if previousEndPos <= len(p.originalContents)-1 {
		contentsMod = append(contentsMod, p.originalContents[previousEndPos:]...)
	}

	return contentsMod, nil
}
