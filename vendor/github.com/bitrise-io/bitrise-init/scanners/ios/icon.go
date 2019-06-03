package ios

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/xcode-project/xcodeproj"
	"github.com/bitrise-io/xcode-project/xcscheme"
)

// lookupIconBySchemeName returns possible ios app icons for a scheme.
func lookupIconBySchemeName(projectPath string, schemeName string, basepath string) (models.Icons, error) {
	project, err := xcodeproj.Open(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open project file: %s, error: %s", projectPath, err)
	}

	scheme, found := project.Scheme(schemeName)
	if !found {
		return nil, fmt.Errorf("failed to find scheme (%s) in project (%s)", schemeName, project.Path)
	}

	blueprintID := getBlueprintID(scheme)
	if blueprintID == "" {
		log.TDebugf("scheme (%s) does not contain app buildable reference in project (%s)", scheme.Name, project.Path)
		return nil, nil
	}

	// Search for the main target
	mainTarget, found := targetByBlueprintID(project.Proj.Targets, blueprintID)
	if !found {
		return nil, fmt.Errorf("no target found for blueprint ID (%s) project (%s)", blueprintID, project.Path)
	}

	return lookupIconByTarget(projectPath, mainTarget, basepath)
}

// lookupIconByTargetName returns possible ios app icons for a target.
func lookupIconByTargetName(projectPath string, targetName string, basepath string) (models.Icons, error) {
	target, err := nameToTarget(projectPath, targetName)
	if err != nil {
		return nil, err
	}

	return lookupIconByTarget(projectPath, target, basepath)
}

func nameToTarget(projectPath string, targetName string) (xcodeproj.Target, error) {
	project, err := xcodeproj.Open(projectPath)
	if err != nil {
		return xcodeproj.Target{}, fmt.Errorf("failed to open project file: %s, error: %s", projectPath, err)
	}

	target, found := targetByName(project, targetName)
	if !found {
		return xcodeproj.Target{}, fmt.Errorf("not found target: %s, in project: %s", targetName, projectPath)
	}
	return target, nil
}

func lookupIconByTarget(projectPath string, target xcodeproj.Target, basepath string) (models.Icons, error) {
	targetToAppIconSetPaths, err := xcodeproj.AppIconSetPaths(projectPath)
	if err != nil {
		return nil, err
	}
	appIconSetPaths, ok := targetToAppIconSetPaths[target.ID]
	log.TDebugf("Appiconsets for target (%s): %s", target.Name, appIconSetPaths)
	if !ok {
		return nil, nil
	}

	iconPaths := []string{}
	for _, appIconSetPath := range appIconSetPaths {
		icon, found, err := parseResourceSet(appIconSetPath)
		if err != nil {
			return nil, fmt.Errorf("could not get icon, error: %s", err)
		} else if !found {
			log.TDebugf("No icon found at %s", appIconSetPath)
			return nil, nil
		}
		log.TDebugf("App icons: %+v", icon)

		iconPath := filepath.Join(appIconSetPath, icon.Filename)
		if _, err := os.Stat(iconPath); err != nil && os.IsNotExist(err) {
			return nil, fmt.Errorf("icon file does not exist: %s, error: %err", iconPath, err)
		}
		iconPaths = append(iconPaths, iconPath)
	}

	icons, err := utility.CreateIconDescriptors(iconPaths, basepath)
	if err != nil {
		return nil, err
	}
	return icons, nil
}

func getBlueprintID(scheme xcscheme.Scheme) string {
	var blueprintID string
	for _, entry := range scheme.BuildAction.BuildActionEntries {
		if entry.BuildableReference.IsAppReference() {
			blueprintID = entry.BuildableReference.BlueprintIdentifier
			break
		}
	}
	return blueprintID
}

func targetByBlueprintID(targets []xcodeproj.Target, blueprintID string) (xcodeproj.Target, bool) {
	for _, target := range targets {
		if target.ID == blueprintID {
			return target, true
		}
	}
	return xcodeproj.Target{}, false
}

func targetByName(proj xcodeproj.XcodeProj, target string) (xcodeproj.Target, bool) {
	projTargets := proj.Proj.Targets
	for _, t := range projTargets {
		if t.Name == target {
			return t, true
		}
	}
	return xcodeproj.Target{}, false
}
