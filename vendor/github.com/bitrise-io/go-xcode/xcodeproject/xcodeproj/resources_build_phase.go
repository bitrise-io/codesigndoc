package xcodeproj

import (
	"fmt"
	"path"

	"github.com/bitrise-io/go-utils/sliceutil"

	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
)

// resourcesBuildPhase represents a PBXResourcesBuildPhase element
type resourcesBuildPhase struct {
	ID    string
	files []string
}

func isResourceBuildPhase(raw serialized.Object) bool {
	if isa, err := raw.String("isa"); err != nil {
		return false
	} else if isa != "PBXResourcesBuildPhase" {
		return false
	}
	return true
}

func parseResourcesBuildPhase(id string, objects serialized.Object) (resourcesBuildPhase, error) {
	rawResourcesBuildPhase, err := objects.Object(id)
	if err != nil {
		return resourcesBuildPhase{}, err
	}

	if !isResourceBuildPhase(rawResourcesBuildPhase) {
		return resourcesBuildPhase{}, fmt.Errorf("not a PBXResourcesBuildPhase element")
	}

	files, err := rawResourcesBuildPhase.StringSlice("files")
	if err != nil {
		return resourcesBuildPhase{}, err
	}

	return resourcesBuildPhase{
		ID:    id,
		files: files,
	}, nil
}

// buildFile represents a PBXBuildFile element
// 47C11A4A21FF63970084FD7F /* Assets.xcassets in Resources */ = {isa = PBXBuildFile; fileRef = 47C11A4921FF63970084FD7F /* Assets.xcassets */; };
type buildFile struct {
	fileRef string
}

func parseBuildFile(id string, objects serialized.Object) (buildFile, error) {
	rawBuildFile, err := objects.Object(id)
	if err != nil {
		return buildFile{}, err
	}
	if isa, err := rawBuildFile.String("isa"); err != nil {
		return buildFile{}, err
	} else if isa != "PBXBuildFile" {
		return buildFile{}, fmt.Errorf("not a PBXBuildFile element")
	}

	fileRef, err := rawBuildFile.String("fileRef")
	if err != nil {
		return buildFile{}, err
	}

	return buildFile{
		fileRef: fileRef,
	}, nil
}

type sourceTree int

const (
	unsupportedParent sourceTree = iota
	groupParent
	absoluteParentPath
	undefinedParent
)

// PBXFileReference
// 47C11A4921FF63970084FD7F /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
type fileReference struct {
	id   string
	path string
}

const fileReferenceElementType = "PBXFileReference"

func isFileReference(raw serialized.Object) (bool, error) {
	if isa, err := raw.String("isa"); err != nil {
		return false, err
	} else if isa == fileReferenceElementType {
		return true, nil
	}
	return false, nil
}

func parseFileReference(id string, objects serialized.Object) (fileReference, error) {
	rawFileReference, err := objects.Object(id)
	if err != nil {
		return fileReference{}, err
	}

	if ok, err := isFileReference(rawFileReference); err != nil {
		return fileReference{}, err
	} else if !ok {
		return fileReference{}, fmt.Errorf("not a %s element", fileReferenceElementType)
	}

	path, err := rawFileReference.String("path")
	if err != nil {
		return fileReference{}, err
	}

	return fileReference{
		id:   id,
		path: path,
	}, nil
}

// PBXGroup example:
// 01801EC11A3360B1002B4718 /* Resources */ = {
// 	isa = PBXGroup;
// 	children = (
// 		A045E5E11EDC5C1700BC8A92 /* Localizable.strings */,
// 		01801EA51A32CA2A002B4718 /* Images.xcassets */,
// 	);
// 	name = Resources;
// 	sourceTree = "<group>";
// };

func resolveObjectAbsolutePath(targetID string, projectID string, projectPath string, objects serialized.Object) (string, error) {
	_, err := objects.Object(targetID)
	if err != nil {
		return "", err
	}
	project, err := objects.Object(projectID)
	if err != nil {
		return "", err
	}

	projectDirPath, err := project.String("projectDirPath")
	if err != nil {
		return "", fmt.Errorf("key projectDirPath not found, project: %s, error: %s", project, err)
	}
	projectRoot, err := project.String("projectRoot")
	if err != nil {
		return "", fmt.Errorf("key projectRoot not found, project: %s, error: %s", project, err)
	}
	mainGroup, err := project.String("mainGroup")
	if err != nil {
		return "", fmt.Errorf("key mainGroup not found, project: %s, error: %s", project, err)
	}

	pathInProjectTree, err := findInProjectTree(targetID, mainGroup, objects, &[]string{})
	if err != nil {
		return "", fmt.Errorf("failed to find target ID in project, error: %s", err)
	}
	pathInProjectTree = append(pathInProjectTree, projectEntry{
		path:         path.Join(projectPath, "..", projectDirPath, projectRoot),
		pathRelation: absoluteParentPath,
	})

	path, err := resolveFilePath(pathInProjectTree)
	if err != nil {
		return "", err
	}
	return path, nil
}

type projectEntry struct {
	id           string
	pathRelation sourceTree
	path         string
}

func findInProjectTree(target string, currentID string, object serialized.Object, visited *[]string) ([]projectEntry, error) {
	if sliceutil.IsStringInSlice(currentID, *visited) {
		return nil, fmt.Errorf("circular reference in project, id: %s", currentID)
	}
	*visited = append(*visited, currentID)

	entry, err := object.Object(currentID)
	if err != nil {
		return nil, fmt.Errorf("object not found, id: %s, error: %s", currentID, err)
	}

	entryPath, err := entry.String("path")
	if err != nil {
	}
	sourceTreeRaw, err := entry.String("sourceTree")
	if err != nil {
		return nil, err
	}
	var pathRelation sourceTree
	switch sourceTreeRaw {
	case "<group>":
		pathRelation = groupParent
	case "<absolute>":
		pathRelation = absoluteParentPath
	case "":
		pathRelation = undefinedParent
	default:
		pathRelation = unsupportedParent
	}

	treeNode := projectEntry{
		id:           currentID,
		path:         entryPath,
		pathRelation: pathRelation,
	}

	if currentID == target {
		return []projectEntry{treeNode}, nil
	}

	childIDs, err := entry.StringSlice("children")
	if err != nil {
		return nil, nil
	}
	for _, childID := range childIDs {
		pathInProjectTree, err := findInProjectTree(target, childID, object, visited)
		if err != nil {
			return nil, err
		} else if pathInProjectTree != nil {
			return append(pathInProjectTree, treeNode), nil
		}
	}
	return nil, nil
}

func resolveFilePath(nodes []projectEntry) (string, error) {
	var partialPath string
	for _, entry := range nodes {
		switch entry.pathRelation {
		case groupParent:
			partialPath = path.Join(entry.path, partialPath)
		case absoluteParentPath:
			return path.Join(entry.path, partialPath), nil
		case undefinedParent:
		case unsupportedParent:
			return "", fmt.Errorf("failed to resolve path, unsupported path relation")
		}
	}
	return partialPath, nil
}
