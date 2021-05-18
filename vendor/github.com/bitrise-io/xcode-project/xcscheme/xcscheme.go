package xcscheme

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

// BuildableReference ...
type BuildableReference struct {
	BlueprintIdentifier string `xml:"BlueprintIdentifier,attr"`
	BlueprintName       string `xml:"BlueprintName,attr"`
	BuildableName       string `xml:"BuildableName,attr"`
	ReferencedContainer string `xml:"ReferencedContainer,attr"`
}

// IsAppReference ...
func (r BuildableReference) IsAppReference() bool {
	return filepath.Ext(r.BuildableName) == ".app"
}

// ReferencedContainerAbsPath ...
func (r BuildableReference) ReferencedContainerAbsPath(schemeContainerDir string) (string, error) {
	s := strings.Split(r.ReferencedContainer, ":")
	if len(s) != 2 {
		return "", fmt.Errorf("unknown referenced container (%s)", r.ReferencedContainer)
	}

	base := s[1]
	absPth := filepath.Join(schemeContainerDir, base)
	return pathutil.AbsPath(absPth)
}

// BuildActionEntry ...
type BuildActionEntry struct {
	BuildForTesting    string `xml:"buildForTesting,attr"`
	BuildForArchiving  string `xml:"buildForArchiving,attr"`
	BuildableReference BuildableReference
}

// BuildAction ...
type BuildAction struct {
	BuildActionEntries []BuildActionEntry `xml:"BuildActionEntries>BuildActionEntry"`
}

// TestableReference ...
type TestableReference struct {
	Skipped            string `xml:"skipped,attr"`
	BuildableReference BuildableReference
}

// TestAction ...
type TestAction struct {
	Testables          []TestableReference `xml:"Testables>TestableReference"`
	BuildConfiguration string              `xml:"buildConfiguration,attr"`
}

// ArchiveAction ...
type ArchiveAction struct {
	BuildConfiguration string `xml:"buildConfiguration,attr"`
}

// Scheme ...
type Scheme struct {
	BuildAction   BuildAction
	ArchiveAction ArchiveAction
	TestAction    TestAction

	Name string
	Path string
}

// Open ...
func Open(pth string) (Scheme, error) {
	b, err := fileutil.ReadBytesFromFile(pth)
	if err != nil {
		return Scheme{}, err
	}

	var scheme Scheme
	if err := xml.Unmarshal(b, &scheme); err != nil {
		return Scheme{}, fmt.Errorf("failed to unmarshal scheme file: %s, error: %s", pth, err)
	}

	scheme.Name = strings.TrimSuffix(filepath.Base(pth), filepath.Ext(pth))
	scheme.Path = pth

	return scheme, nil
}

// AppBuildActionEntry ...
func (s Scheme) AppBuildActionEntry() (BuildActionEntry, bool) {
	var entry BuildActionEntry
	for _, e := range s.BuildAction.BuildActionEntries {
		if e.BuildForArchiving != "YES" {
			continue
		}
		if !e.BuildableReference.IsAppReference() {
			continue
		}
		entry = e
		break
	}

	return entry, (entry.BuildableReference.BlueprintIdentifier != "")
}
