package xcodeproj

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/xcode-project/pretty"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/xcode-project/serialized"
	"github.com/bitrise-io/xcode-project/xcodebuild"
	"github.com/bitrise-io/xcode-project/xcscheme"
	"golang.org/x/text/unicode/norm"
	"howett.net/plist"
)

// XcodeProj ...
type XcodeProj struct {
	Proj    Proj
	RawProj serialized.Object
	Format  int

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

	format, raw, objects, projectID, err := open(pth)
	if err != nil {
		return XcodeProj{}, err
	}

	p, err := parseProj(projectID, objects)
	if err != nil {
		return XcodeProj{}, err
	}

	return XcodeProj{
		Proj:    p,
		RawProj: raw,
		Format:  format,
		Path:    absPth,
		Name:    strings.TrimSuffix(filepath.Base(absPth), filepath.Ext(absPth)),
	}, nil
}

// open parse the provided .pbxprog file.
// Returns the `raw` contents as a serialized.Object, the `objects` as serialized.Object and the PBXProject's `projectID` as string
// If there was an error during the parsing it returns an error
func open(absPth string) (format int, rawPbxProj serialized.Object, objects serialized.Object, projectID string, err error) {
	pbxProjPth := filepath.Join(absPth, "project.pbxproj")

	var b []byte
	b, err = fileutil.ReadBytesFromFile(pbxProjPth)
	if err != nil {
		return
	}

	if format, err = plist.Unmarshal(b, &rawPbxProj); err != nil {
		err = fmt.Errorf("failed to generate json from Pbxproj - error: %s", err)
		return
	}

	objects, err = rawPbxProj.Object("objects")
	if err != nil {
		return
	}

	for id := range objects {
		var object serialized.Object
		object, err = objects.Object(id)
		if err != nil {
			return
		}

		var objectISA string
		objectISA, err = object.String("isa")
		if err != nil {
			return
		}

		if objectISA == "PBXProject" {
			projectID = id
			break
		}
	}
	return
}

// IsXcodeProj ...
func IsXcodeProj(pth string) bool {
	return filepath.Ext(pth) == ".xcodeproj"
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

	targetAttributes, err := p.TargetAttributes()
	if err != nil {
		return fmt.Errorf("failed to get project's target attributes, error: %s", err)
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

	// Override TargetAttributes
	if err = foreceCodeSignOnTargetAttributes(targetAttributes, target.ID, developmentTeam); err != nil {
		return fmt.Errorf("failed to change code signing in target attributes, error: %s", err)
	}

	// Override BuildSettings
	if err = foreceCodeSignOnBuildConfiguration(buildConfiguration, target.ID, developmentTeam, provisioningProfileUUID, codesignIdentity); err != nil {
		return fmt.Errorf("failed to change code signing in build settings, error: %s", err)
	}
	return nil
}

// foreceCodeSignOnTargetAttributes sets the TargetAttributes for the provided targetID.
// **Overrides the ProvisioningStyle, developmentTeam and clears the DevelopmentTeamName in the provided `targetAttributes`!**
func foreceCodeSignOnTargetAttributes(targetAttributes serialized.Object, targetID, developmentTeam string) error {
	targetAttribute, err := targetAttributes.Object(targetID)
	if err != nil {
		return fmt.Errorf("failed to get traget's (%s) attributes, error: %s", targetID, err)
	}

	targetAttribute["ProvisioningStyle"] = "Manual"
	targetAttribute["DevelopmentTeam"] = developmentTeam
	targetAttribute["DevelopmentTeamName"] = ""
	return nil
}

// foreceCodeSignOnBuildConfiguration sets the BuildSettings for the provided targetID.
// **Overrides the CODE_SIGN_STYLE, DEVELOPMENT_TEAM, CODE_SIGN_IDENTITY, CODE_SIGN_IDENTITY[sdk=iphoneos\*], PROVISIONING_PROFILE, PROVISIONING_PROFILE[sdk=iphoneos\*] and clears the PROVISIONING_PROFILE_SPECIFIER in the provided `buildConfiguration`!**
func foreceCodeSignOnBuildConfiguration(buildConfiguration serialized.Object, targetID, developmentTeam, provisioningProfileUUID, codesignIdentity string) error {
	buildSettings, err := buildConfiguration.Object("buildSettings")
	if err != nil {
		return fmt.Errorf("failed to get buildSettings of buildConfiguration (%s), error: %s", pretty.Object(buildConfiguration), err)
	}

	buildSettings["CODE_SIGN_STYLE"] = "Manual"
	buildSettings["DEVELOPMENT_TEAM"] = developmentTeam
	buildSettings["CODE_SIGN_IDENTITY"] = codesignIdentity
	buildSettings["CODE_SIGN_IDENTITY[sdk=iphoneos*]"] = codesignIdentity
	buildSettings["PROVISIONING_PROFILE_SPECIFIER"] = ""
	buildSettings["PROVISIONING_PROFILE"] = provisioningProfileUUID
	buildSettings["PROVISIONING_PROFILE[sdk=iphoneos*]"] = provisioningProfileUUID

	return nil
}

// Save the XcodeProj
//
// Overrides the project.pbxproj file of the XcodeProj with the contents of `rawProj`
func (p XcodeProj) Save() error {
	return p.savePBXProj()
}

// savePBXProj overrides the project.pbxproj file of the XcodeProj with the contents of `rawProj`
func (p XcodeProj) savePBXProj() error {
	b, err := plist.Marshal(p.RawProj, p.Format)
	if err != nil {
		return fmt.Errorf("failed to marshal .pbxproj")
	}

	pth := path.Join(p.Path, "project.pbxproj")
	return ioutil.WriteFile(pth, b, 0644)
}
