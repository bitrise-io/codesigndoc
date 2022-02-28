package xcscheme

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

// BuildableReference ...
type BuildableReference struct {
	BuildableIdentifier string `xml:"BuildableIdentifier,attr"`
	BlueprintIdentifier string `xml:"BlueprintIdentifier,attr"`
	BuildableName       string `xml:"BuildableName,attr"`
	BlueprintName       string `xml:"BlueprintName,attr"`
	ReferencedContainer string `xml:"ReferencedContainer,attr"`
}

// IsAppReference ...
func (r BuildableReference) IsAppReference() bool {
	return filepath.Ext(r.BuildableName) == ".app"
}

func (r BuildableReference) isTestProduct() bool {
	return filepath.Ext(r.BuildableName) == ".xctest"
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
	BuildForTesting   string `xml:"buildForTesting,attr"`
	BuildForRunning   string `xml:"buildForRunning,attr"`
	BuildForProfiling string `xml:"buildForProfiling,attr"`
	BuildForArchiving string `xml:"buildForArchiving,attr"`
	BuildForAnalyzing string `xml:"buildForAnalyzing,attr"`

	BuildableReference BuildableReference
}

// BuildAction ...
type BuildAction struct {
	ParallelizeBuildables     string             `xml:"parallelizeBuildables,attr"`
	BuildImplicitDependencies string             `xml:"buildImplicitDependencies,attr"`
	BuildActionEntries        []BuildActionEntry `xml:"BuildActionEntries>BuildActionEntry"`
}

// TestableReference ...
type TestableReference struct {
	Skipped            string `xml:"skipped,attr"`
	BuildableReference BuildableReference
}

func (r TestableReference) isTestable() bool {
	return r.Skipped == "NO" && r.BuildableReference.isTestProduct()
}

// MacroExpansion ...
type MacroExpansion struct {
	BuildableReference BuildableReference
}

// AdditionalOptions ...
type AdditionalOptions struct {
}

// TestAction ...
type TestAction struct {
	BuildConfiguration           string `xml:"buildConfiguration,attr"`
	SelectedDebuggerIdentifier   string `xml:"selectedDebuggerIdentifier,attr"`
	SelectedLauncherIdentifier   string `xml:"selectedLauncherIdentifier,attr"`
	ShouldUseLaunchSchemeArgsEnv string `xml:"shouldUseLaunchSchemeArgsEnv,attr"`

	Testables         []TestableReference `xml:"Testables>TestableReference"`
	MacroExpansion    MacroExpansion
	AdditionalOptions AdditionalOptions
}

// BuildableProductRunnable ...
type BuildableProductRunnable struct {
	RunnableDebuggingMode string `xml:"runnableDebuggingMode,attr"`
	BuildableReference    BuildableReference
}

// LaunchAction ...
type LaunchAction struct {
	BuildConfiguration             string `xml:"buildConfiguration,attr"`
	SelectedDebuggerIdentifier     string `xml:"selectedDebuggerIdentifier,attr"`
	SelectedLauncherIdentifier     string `xml:"selectedLauncherIdentifier,attr"`
	LaunchStyle                    string `xml:"launchStyle,attr"`
	UseCustomWorkingDirectory      string `xml:"useCustomWorkingDirectory,attr"`
	IgnoresPersistentStateOnLaunch string `xml:"ignoresPersistentStateOnLaunch,attr"`
	DebugDocumentVersioning        string `xml:"debugDocumentVersioning,attr"`
	DebugServiceExtension          string `xml:"debugServiceExtension,attr"`
	AllowLocationSimulation        string `xml:"allowLocationSimulation,attr"`
	BuildableProductRunnable       BuildableProductRunnable
	AdditionalOptions              AdditionalOptions
}

// ProfileAction ...
type ProfileAction struct {
	BuildConfiguration           string `xml:"buildConfiguration,attr"`
	ShouldUseLaunchSchemeArgsEnv string `xml:"shouldUseLaunchSchemeArgsEnv,attr"`
	SavedToolIdentifier          string `xml:"savedToolIdentifier,attr"`
	UseCustomWorkingDirectory    string `xml:"useCustomWorkingDirectory,attr"`
	DebugDocumentVersioning      string `xml:"debugDocumentVersioning,attr"`
	BuildableProductRunnable     BuildableProductRunnable
}

// AnalyzeAction ...
type AnalyzeAction struct {
	BuildConfiguration string `xml:"buildConfiguration,attr"`
}

// ArchiveAction ...
type ArchiveAction struct {
	BuildConfiguration       string `xml:"buildConfiguration,attr"`
	RevealArchiveInOrganizer string `xml:"revealArchiveInOrganizer,attr"`
}

// Scheme ...
type Scheme struct {
	// The last known Xcode version.
	LastUpgradeVersion string `xml:"LastUpgradeVersion,attr"`
	// The version of `.xcscheme` files supported.
	Version string `xml:"version,attr"`

	BuildAction   BuildAction
	TestAction    TestAction
	LaunchAction  LaunchAction
	ProfileAction ProfileAction
	AnalyzeAction AnalyzeAction
	ArchiveAction ArchiveAction

	Name     string `xml:"-"`
	Path     string `xml:"-"`
	IsShared bool   `xml:"-"`
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

// XMLToken ...
type XMLToken int

const (
	invalid XMLToken = iota
	// XMLStart ...
	XMLStart
	// XMLEnd ...
	XMLEnd
	// XMLAttribute ...
	XMLAttribute
)

// Marshal ...
func (s Scheme) Marshal() ([]byte, error) {
	contents, err := xml.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Scheme: %v", err)
	}

	contentsNewline := strings.ReplaceAll(string(contents), "><", ">\n<")

	// Place XML Attributes on separate lines
	re := regexp.MustCompile(`\s([^=<>/]*)\s?=\s?"([^=<>/]*)"`)
	contentsNewline = re.ReplaceAllString(contentsNewline, "\n$1 = \"$2\"")

	var contentsIndented string

	indent := 0
	for _, line := range strings.Split(contentsNewline, "\n") {
		currentLine := XMLAttribute
		if strings.HasPrefix(line, "</") {
			currentLine = XMLEnd
		} else if strings.HasPrefix(line, "<") {
			currentLine = XMLStart
		}

		if currentLine == XMLEnd && indent != 0 {
			indent--
		}

		contentsIndented += strings.Repeat("   ", indent)
		contentsIndented += line + "\n"

		if currentLine == XMLStart {
			indent++
		}
	}

	return []byte(xml.Header + contentsIndented), nil
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

// IsTestable returns true if Test is a valid action
func (s Scheme) IsTestable() bool {
	for _, testEntry := range s.TestAction.Testables {
		if testEntry.isTestable() {
			return true
		}
	}

	return false
}
