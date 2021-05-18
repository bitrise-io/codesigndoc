package xcodeproj

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/bitrise-io/xcode-project/serialized"
)

// TargetsToAppIconSets maps target names to an array app icon set absolute paths.
type TargetsToAppIconSets map[string][]string

// AppIconSetPaths parses an Xcode project and returns targets mapped to app icon set absolute paths.
func AppIconSetPaths(projectPath string) (TargetsToAppIconSets, error) {
	absPth, err := pathutil.AbsPath(projectPath)
	if err != nil {
		return TargetsToAppIconSets{}, err
	}

	proj, err := Open(absPth)
	if err != nil {
		return TargetsToAppIconSets{}, err
	}

	objects, err := proj.RawProj.Object("objects")
	if err != nil {
		return TargetsToAppIconSets{}, err
	}

	return appIconSetPaths(proj.Proj, projectPath, objects)
}

func appIconSetPaths(project Proj, projectPath string, objects serialized.Object) (TargetsToAppIconSets, error) {
	targetToAppIcons := map[string][]string{}
	for _, target := range project.Targets {
		appIconSetNames := getAppIconSetNames(target)
		if len(appIconSetNames) == 0 {
			continue
		}

		assetCatalogs, err := assetCatalogs(target, project.ID, objects)
		if err != nil {
			return nil, err
		} else if len(assetCatalogs) == 0 {
			continue
		}

		appIcons := []string{}
		for _, appIconSetName := range appIconSetNames {
			appIconSetPaths, err := lookupAppIconPaths(projectPath, assetCatalogs, appIconSetName, project.ID, objects)
			if err != nil {
				return nil, err
			} else if len(appIconSetPaths) == 0 {
				return nil, fmt.Errorf("not found app icon set (%s) on paths: %s", appIconSetName, assetCatalogs)
			}
			appIcons = append(appIcons, appIconSetPaths...)
		}
		targetToAppIcons[target.ID] = sliceutil.UniqueStringSlice(appIcons)
	}

	return targetToAppIcons, nil
}

func lookupAppIconPaths(projectPath string, assetCatalogs []fileReference, appIconSetName string, projectID string, objects serialized.Object) ([]string, error) {
	var icons []string
	for _, fileReference := range assetCatalogs {
		resolvedPath, err := resolveObjectAbsolutePath(fileReference.id, projectID, projectPath, objects)
		if err != nil {
			return nil, err
		} else if resolvedPath == "" {
			return nil, fmt.Errorf("could not resolve path")
		}

		re := regexp.MustCompile(`\$\{(.+)\}`)
		wildcharAppIconSetName := re.ReplaceAllString(appIconSetName, "*")

		matches, err := filepath.Glob(path.Join(regexp.QuoteMeta(resolvedPath), wildcharAppIconSetName+".appiconset"))
		if err != nil {
			return nil, err
		}

		icons = append(icons, matches...)
	}

	return icons, nil
}

func assetCatalogs(target Target, projectID string, objects serialized.Object) ([]fileReference, error) {
	if target.Type == NativeTargetType { // Ignoring PBXAggregateTarget and PBXLegacyTarget as may not contain buildPhases key
		resourcesBuildPhase, err := filterResourcesBuildPhase(target.buildPhaseIDs, objects)
		if err != nil {
			return nil, fmt.Errorf("getting resource build phases failed, error: %s", err)
		}
		assetCatalogs, err := filterAssetCatalogs(resourcesBuildPhase, projectID, objects)
		if err != nil {
			return nil, err
		}
		return assetCatalogs, nil
	}
	return nil, nil
}

func filterResourcesBuildPhase(buildPhases []string, objects serialized.Object) (resourcesBuildPhase, error) {
	for _, buildPhaseUUID := range buildPhases {
		rawBuildPhase, err := objects.Object(buildPhaseUUID)
		if err != nil {
			return resourcesBuildPhase{}, err
		}
		if isResourceBuildPhase(rawBuildPhase) {
			buildPhase, err := parseResourcesBuildPhase(buildPhaseUUID, objects)
			if err != nil {
				return resourcesBuildPhase{}, fmt.Errorf("failed to parse ResourcesBuildPhase, error: %s", err)
			}
			return buildPhase, nil
		}
	}
	return resourcesBuildPhase{}, fmt.Errorf("resource build phase not found")
}

func filterAssetCatalogs(buildPhase resourcesBuildPhase, projectID string, objects serialized.Object) ([]fileReference, error) {
	assetCatalogs := []fileReference{}
	for _, fileUUID := range buildPhase.files {
		buildFile, err := parseBuildFile(fileUUID, objects)
		if err != nil {
			// ignore:
			// D0177B971F26869C0044446D /* (null) in Resources */ = {isa = PBXBuildFile; };
			continue
		}

		// can be PBXVariantGroup or PBXFileReference
		rawElement, err := objects.Object(buildFile.fileRef)
		if err != nil {
			return nil, err
		}
		if ok, err := isFileReference(rawElement); err != nil {
			return nil, err
		} else if !ok {
			// ignore PBXVariantGroup
			continue
		}

		fileReference, err := parseFileReference(buildFile.fileRef, objects)
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(fileReference.path, ".xcassets") {
			assetCatalogs = append(assetCatalogs, fileReference)
		}
	}
	return assetCatalogs, nil
}

func getAppIconSetNames(target Target) []string {
	const appIconSetNameKey = "ASSETCATALOG_COMPILER_APPICON_NAME"

	appIconSetNames := []string{}
	for _, configuration := range target.BuildConfigurationList.BuildConfigurations {
		appIconSetName, err := configuration.BuildSettings.String(appIconSetNameKey)
		if err != nil {
			return nil
		}
		appIconSetNames = append(appIconSetNames, appIconSetName)
	}

	return appIconSetNames
}
