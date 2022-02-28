package xcworkspace

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/xcodebuild"
	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
	"golang.org/x/text/unicode/norm"
)

const (
	// XCWorkspaceExtension ...
	XCWorkspaceExtension = ".xcworkspace"
)

// Workspace represents an Xcode workspace
type Workspace struct {
	FileRefs []FileRef `xml:"FileRef"`
	Groups   []Group   `xml:"Group"`

	Name string
	Path string
}

// Scheme returns the scheme by name and it's container's absolute path.
func (w Workspace) Scheme(name string) (*xcscheme.Scheme, string, error) {
	schemesByContainer, err := w.Schemes()
	if err != nil {
		return nil, "", err
	}

	normName := norm.NFC.String(name)
	for container, schemes := range schemesByContainer {
		for _, scheme := range schemes {
			if norm.NFC.String(scheme.Name) == normName {
				return &scheme, container, nil
			}
		}
	}

	return nil, "", xcscheme.NotFoundError{Scheme: name, Container: w.Name}
}

// SchemeBuildSettings ...
func (w Workspace) SchemeBuildSettings(scheme, configuration string, customOptions ...string) (serialized.Object, error) {
	commandModel := xcodebuild.NewShowBuildSettingsCommand(w.Path)
	commandModel.SetScheme(scheme)
	commandModel.SetConfiguration(configuration)
	commandModel.SetCustomOptions(customOptions)
	return commandModel.RunAndReturnSettings()
}

// Schemes ...
func (w Workspace) Schemes() (map[string][]xcscheme.Scheme, error) {
	schemesByContainer := map[string][]xcscheme.Scheme{}

	workspaceSchemes, err := xcscheme.FindSchemesIn(w.Path)
	if err != nil {
		return nil, err
	}

	schemesByContainer[w.Path] = workspaceSchemes

	// project schemes
	projectLocations, err := w.ProjectFileLocations()
	if err != nil {
		return nil, err
	}

	for _, projectLocation := range projectLocations {
		if exist, err := pathutil.IsPathExists(projectLocation); err != nil {
			return nil, fmt.Errorf("failed to check if project exist at: %s, error: %s", projectLocation, err)
		} else if !exist {
			// at this point we are interested the schemes visible for the workspace
			continue
		}

		project, err := xcodeproj.Open(projectLocation)
		if err != nil {
			return nil, err
		}

		projectSchemes, err := project.Schemes()
		if err != nil {
			return nil, err
		}

		schemesByContainer[project.Path] = projectSchemes
	}

	return schemesByContainer, nil
}

// FileLocations ...
func (w Workspace) FileLocations() ([]string, error) {
	var fileLocations []string

	for _, fileRef := range w.FileRefs {
		pth, err := fileRef.AbsPath(filepath.Dir(w.Path))
		if err != nil {
			return nil, err
		}

		fileLocations = append(fileLocations, pth)
	}

	for _, group := range w.Groups {
		groupFileLocations, err := group.FileLocations(filepath.Dir(w.Path))
		if err != nil {
			return nil, err
		}

		fileLocations = append(fileLocations, groupFileLocations...)
	}

	return fileLocations, nil
}

// ProjectFileLocations ...
func (w Workspace) ProjectFileLocations() ([]string, error) {
	var projectLocations []string
	fileLocations, err := w.FileLocations()
	if err != nil {
		return nil, err
	}
	for _, fileLocation := range fileLocations {
		if xcodeproj.IsXcodeProj(fileLocation) {
			projectLocations = append(projectLocations, fileLocation)
		}
	}
	return projectLocations, nil
}

// Open ...
func Open(pth string) (Workspace, error) {
	contentsPth := filepath.Join(pth, "contents.xcworkspacedata")
	b, err := fileutil.ReadBytesFromFile(contentsPth)
	if err != nil {
		return Workspace{}, err
	}

	var workspace Workspace
	if err := xml.Unmarshal(b, &workspace); err != nil {
		return Workspace{}, fmt.Errorf("failed to unmarshal workspace file: %s, error: %s", pth, err)
	}

	workspace.Name = strings.TrimSuffix(filepath.Base(pth), filepath.Ext(pth))
	workspace.Path = pth

	return workspace, nil
}

// IsWorkspace ...
func IsWorkspace(pth string) bool {
	return filepath.Ext(pth) == ".xcworkspace"
}
