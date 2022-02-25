package ios

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
)

func getMainTarget(project xcodeproj.XcodeProj, scheme xcscheme.Scheme) (xcodeproj.Target, error) {
	entry, found := scheme.AppBuildActionEntry()
	if !found {
		return xcodeproj.Target{}, fmt.Errorf("scheme (%s) does not contain app buildable reference in project (%s)", scheme.Name, project.Path)
	}

	blueprintID := entry.BuildableReference.BlueprintIdentifier
	mainTarget, found := project.Proj.Target(blueprintID)
	if !found {
		return xcodeproj.Target{}, fmt.Errorf("no target found for blueprint ID (%s) in project (%s)", blueprintID, project.Path)
	}

	return mainTarget, nil
}

func lookupIconByScheme(project xcodeproj.XcodeProj, scheme xcscheme.Scheme, basePath string) (models.Icons, error) {
	mainTarget, err := getMainTarget(project, scheme)
	if err != nil {
		log.TDebugf("%s", err)
		return nil, nil
	}

	return lookupIconByTarget(project.Path, mainTarget, basePath)
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
