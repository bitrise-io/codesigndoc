package xcodeproj

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
)

const (
	yes                         = "YES"
	no                          = "NO"
	buildableID                 = "primary"
	defaultDebugConfiguration   = "Debug"
	defaultReleaseConfiguration = "Release"
	debuggerID                  = "Xcode.DebuggerFoundation.Debugger.LLDB"
	launcherID                  = "Xcode.DebuggerFoundation.Launcher.LLDB"
)

// SaveSharedScheme saves or overwrites a shared Scheme in the Project
// The file name will be determined using the Name field of the Scheme
func (p XcodeProj) SaveSharedScheme(scheme xcscheme.Scheme) error {
	dir := filepath.Join(p.Path, "xcshareddata", "xcschemes")
	path := filepath.Join(dir, fmt.Sprintf("%s.xcscheme", scheme.Name))

	contents, err := scheme.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal Scheme: %v", err)
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := ioutil.WriteFile(path, contents, 0600); err != nil {
		return fmt.Errorf("failed to write Scheme file (%s): %v", path, err)
	}

	return nil
}

// ReCreateSchemes creates new schemes based on the available Targets
func (p XcodeProj) ReCreateSchemes() []xcscheme.Scheme {
	var schemes []xcscheme.Scheme
	for _, buildTarget := range p.Proj.Targets {
		if buildTarget.Type != NativeTargetType || buildTarget.IsTest() {
			continue
		}

		var testTargets []Target
		for _, testTarget := range p.Proj.Targets {
			if testTarget.IsTest() && testTarget.DependesOn(buildTarget.ID) {
				testTargets = append(testTargets, testTarget)
			}
		}

		schemes = append(schemes, newScheme(buildTarget, testTargets, filepath.Base(p.Path)))
	}

	return schemes
}

func newScheme(buildTarget Target, testTargets []Target, projectname string) xcscheme.Scheme {
	return xcscheme.Scheme{
		Name:               buildTarget.Name,
		LastUpgradeVersion: "1240",
		Version:            "1.3",
		BuildAction:        newBuildAction(buildTarget, projectname),
		TestAction:         newTestAction(buildTarget, testTargets, projectname),
		LaunchAction:       newLaunchAction(buildTarget, projectname),
		ProfileAction:      newProfileAction(buildTarget, projectname),
		AnalyzeAction:      newAnalyzeAction(buildTarget),
		ArchiveAction:      newArchiveAction(buildTarget),
	}
}

func newBuildableReference(target Target, projectName string) xcscheme.BuildableReference {
	return xcscheme.BuildableReference{
		BuildableIdentifier: buildableID,
		BlueprintIdentifier: target.ID,
		BuildableName:       path.Base(target.ProductReference.Path),
		BlueprintName:       target.Name,
		ReferencedContainer: fmt.Sprintf("container:%s", projectName),
	}
}

func newBuildAction(target Target, projectName string) xcscheme.BuildAction {
	return xcscheme.BuildAction{
		ParallelizeBuildables:     yes,
		BuildImplicitDependencies: yes,
		BuildActionEntries: []xcscheme.BuildActionEntry{
			{
				BuildForTesting:    yes,
				BuildForRunning:    yes,
				BuildForProfiling:  yes,
				BuildForArchiving:  yes,
				BuildForAnalyzing:  yes,
				BuildableReference: newBuildableReference(target, projectName),
			},
		},
	}
}

func newTestableReference(target Target, projectName string) xcscheme.TestableReference {
	return xcscheme.TestableReference{
		Skipped:            no,
		BuildableReference: newBuildableReference(target, projectName),
	}
}

func newTestAction(buildTarget Target, testTargets []Target, projectName string) xcscheme.TestAction {
	if len(testTargets) == 0 {
		return xcscheme.TestAction{}
	}

	testAction := xcscheme.TestAction{
		BuildConfiguration:           debugConfigurationName(testTargets[0]),
		SelectedDebuggerIdentifier:   debuggerID,
		SelectedLauncherIdentifier:   launcherID,
		ShouldUseLaunchSchemeArgsEnv: yes,
		MacroExpansion: xcscheme.MacroExpansion{
			BuildableReference: newBuildableReference(buildTarget, projectName),
		},
		Testables: []xcscheme.TestableReference{},
	}

	for _, testTarget := range testTargets {
		testAction.Testables = append(
			testAction.Testables,
			newTestableReference(testTarget, projectName),
		)
	}

	return testAction
}

func newBuildableProductRunnable(target Target, projectName string) xcscheme.BuildableProductRunnable {
	return xcscheme.BuildableProductRunnable{
		RunnableDebuggingMode: "0",
		BuildableReference:    newBuildableReference(target, projectName),
	}
}

func newLaunchAction(target Target, projectName string) xcscheme.LaunchAction {
	return xcscheme.LaunchAction{
		BuildConfiguration:             debugConfigurationName(target),
		SelectedDebuggerIdentifier:     debuggerID,
		SelectedLauncherIdentifier:     launcherID,
		LaunchStyle:                    "0",
		UseCustomWorkingDirectory:      no,
		IgnoresPersistentStateOnLaunch: no,
		DebugDocumentVersioning:        yes,
		DebugServiceExtension:          "internal",
		AllowLocationSimulation:        yes,
		BuildableProductRunnable:       newBuildableProductRunnable(target, projectName),
	}
}

func newProfileAction(target Target, projectName string) xcscheme.ProfileAction {
	return xcscheme.ProfileAction{
		BuildConfiguration:           releaseConfigurationName(target),
		ShouldUseLaunchSchemeArgsEnv: yes,
		UseCustomWorkingDirectory:    no,
		DebugDocumentVersioning:      yes,
		BuildableProductRunnable:     newBuildableProductRunnable(target, projectName),
	}
}

func newAnalyzeAction(target Target) xcscheme.AnalyzeAction {
	return xcscheme.AnalyzeAction{
		BuildConfiguration: debugConfigurationName(target),
	}
}

func newArchiveAction(target Target) xcscheme.ArchiveAction {
	return xcscheme.ArchiveAction{
		BuildConfiguration:       releaseConfigurationName(target),
		RevealArchiveInOrganizer: yes,
	}
}

func debugConfigurationName(target Target) string {
	for _, buildConfig := range target.BuildConfigurationList.BuildConfigurations {
		if buildConfig.Name == defaultDebugConfiguration {
			return defaultDebugConfiguration
		}
	}

	return target.BuildConfigurationList.DefaultConfigurationName
}

func releaseConfigurationName(target Target) string {
	for _, buildConfig := range target.BuildConfigurationList.BuildConfigurations {
		if buildConfig.Name == defaultReleaseConfiguration {
			return defaultReleaseConfiguration
		}
	}

	return target.BuildConfigurationList.DefaultConfigurationName
}
