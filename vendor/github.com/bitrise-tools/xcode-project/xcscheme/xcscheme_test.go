package xcscheme

import (
	"encoding/xml"
	"testing"

	"github.com/bitrise-tools/xcode-project/testhelper"
	"github.com/stretchr/testify/require"
)

func TestOpenScheme(t *testing.T) {
	pth := testhelper.CreateTmpFile(t, "ios-simple-objc.xcscheme", schemeContent)
	scheme, err := Open(pth)
	require.NoError(t, err)

	require.Equal(t, "ios-simple-objc", scheme.Name)
	require.Equal(t, pth, scheme.Path)

	require.Equal(t, "Release", scheme.ArchiveAction.BuildConfiguration)
	require.Equal(t, 2, len(scheme.BuildAction.BuildActionEntries))

	{
		entry := scheme.BuildAction.BuildActionEntries[0]
		require.Equal(t, "YES", entry.BuildForArchiving)
		require.Equal(t, "YES", entry.BuildForTesting)
		require.Equal(t, "BA3CBE7419F7A93800CED4D5", entry.BuildableReference.BlueprintIdentifier)

		pth, err := entry.BuildableReference.ReferencedContainerAbsPath("/project.xcodeproj")
		require.NoError(t, err)
		require.Equal(t, "/project.xcodeproj/ios-simple-objc.xcodeproj", pth)
	}

	{
		entry := scheme.BuildAction.BuildActionEntries[1]
		require.Equal(t, "NO", entry.BuildForArchiving)
		require.Equal(t, "YES", entry.BuildForTesting)
		require.Equal(t, "BA3CBE9019F7A93900CED4D5", entry.BuildableReference.BlueprintIdentifier)
	}
}

func TestAppBuildActionEntry(t *testing.T) {
	var scheme Scheme
	require.NoError(t, xml.Unmarshal([]byte(schemeContent), &scheme))

	entry, ok := scheme.AppBuildActionEntry()
	require.True(t, ok)

	require.Equal(t, "YES", entry.BuildForArchiving)
	require.Equal(t, "YES", entry.BuildForTesting)
	require.Equal(t, "BA3CBE7419F7A93800CED4D5", entry.BuildableReference.BlueprintIdentifier)
	require.Equal(t, "ios-simple-objc.app", entry.BuildableReference.BuildableName)
	require.Equal(t, "ios-simple-objc", entry.BuildableReference.BlueprintName)
	require.Equal(t, "container:ios-simple-objc.xcodeproj", entry.BuildableReference.ReferencedContainer)

	require.True(t, entry.BuildableReference.IsAppReference())
}

func TestAppTestActionEntry(t *testing.T) {
	var scheme Scheme
	require.NoError(t, xml.Unmarshal([]byte(schemeContent), &scheme))

	require.Equal(t, 2, len(scheme.TestAction.Testables))
	require.Equal(t, "NO", scheme.TestAction.Testables[0].Skipped)
	require.Equal(t, "YES", scheme.TestAction.Testables[1].Skipped)
	require.Equal(t, "BA3CBE9019F7A93900CED4D5", scheme.TestAction.Testables[0].BuildableReference.BlueprintIdentifier)

	require.False(t, scheme.TestAction.Testables[0].BuildableReference.IsAppReference())
	require.False(t, scheme.TestAction.Testables[1].BuildableReference.IsAppReference())
}

const schemeContent = `<?xml version="1.0" encoding="UTF-8"?>
<Scheme
   LastUpgradeVersion = "0800"
   version = "1.3">
   <BuildAction
      parallelizeBuildables = "YES"
      buildImplicitDependencies = "YES">
      <BuildActionEntries>
         <BuildActionEntry
            buildForTesting = "YES"
            buildForRunning = "YES"
            buildForProfiling = "YES"
            buildForArchiving = "YES"
            buildForAnalyzing = "YES">
            <BuildableReference
               BuildableIdentifier = "primary"
               BlueprintIdentifier = "BA3CBE7419F7A93800CED4D5"
               BuildableName = "ios-simple-objc.app"
               BlueprintName = "ios-simple-objc"
               ReferencedContainer = "container:ios-simple-objc.xcodeproj">
            </BuildableReference>
         </BuildActionEntry>
         <BuildActionEntry
            buildForTesting = "YES"
            buildForRunning = "YES"
            buildForProfiling = "NO"
            buildForArchiving = "NO"
            buildForAnalyzing = "YES">
            <BuildableReference
               BuildableIdentifier = "primary"
               BlueprintIdentifier = "BA3CBE9019F7A93900CED4D5"
               BuildableName = "ios-simple-objcTests.xctest"
               BlueprintName = "ios-simple-objcTests"
               ReferencedContainer = "container:ios-simple-objc.xcodeproj">
            </BuildableReference>
         </BuildActionEntry>
      </BuildActionEntries>
   </BuildAction>
   <TestAction
      buildConfiguration = "Debug"
      selectedDebuggerIdentifier = "Xcode.DebuggerFoundation.Debugger.LLDB"
      selectedLauncherIdentifier = "Xcode.DebuggerFoundation.Launcher.LLDB"
      shouldUseLaunchSchemeArgsEnv = "YES">
      <Testables>
         <TestableReference
            skipped = "NO">
            <BuildableReference
               BuildableIdentifier = "primary"
               BlueprintIdentifier = "BA3CBE9019F7A93900CED4D5"
               BuildableName = "ios-simple-objcTests.xctest"
               BlueprintName = "ios-simple-objcTests"
               ReferencedContainer = "container:ios-simple-objc.xcodeproj">
            </BuildableReference>
         </TestableReference>
         <TestableReference
            skipped = "YES">
            <BuildableReference
               BuildableIdentifier = "primary"
               BlueprintIdentifier = "BA4CBE9019F7A93900CED4D5"
               BuildableName = "ios-simple-objcTests2.xctest"
               BlueprintName = "ios-simple-objcTests2"
               ReferencedContainer = "container:ios-simple-objc.xcodeproj">
            </BuildableReference>
         </TestableReference>
      </Testables>
      <MacroExpansion>
         <BuildableReference
            BuildableIdentifier = "primary"
            BlueprintIdentifier = "BA3CBE7419F7A93800CED4D5"
            BuildableName = "ios-simple-objc.app"
            BlueprintName = "ios-simple-objc"
            ReferencedContainer = "container:ios-simple-objc.xcodeproj">
         </BuildableReference>
      </MacroExpansion>
      <AdditionalOptions>
      </AdditionalOptions>
   </TestAction>
   <LaunchAction
      buildConfiguration = "Debug"
      selectedDebuggerIdentifier = "Xcode.DebuggerFoundation.Debugger.LLDB"
      selectedLauncherIdentifier = "Xcode.DebuggerFoundation.Launcher.LLDB"
      launchStyle = "0"
      useCustomWorkingDirectory = "NO"
      ignoresPersistentStateOnLaunch = "NO"
      debugDocumentVersioning = "YES"
      debugServiceExtension = "internal"
      allowLocationSimulation = "YES">
      <BuildableProductRunnable
         runnableDebuggingMode = "0">
         <BuildableReference
            BuildableIdentifier = "primary"
            BlueprintIdentifier = "BA3CBE7419F7A93800CED4D5"
            BuildableName = "ios-simple-objc.app"
            BlueprintName = "ios-simple-objc"
            ReferencedContainer = "container:ios-simple-objc.xcodeproj">
         </BuildableReference>
      </BuildableProductRunnable>
      <AdditionalOptions>
      </AdditionalOptions>
   </LaunchAction>
   <ProfileAction
      buildConfiguration = "Release"
      shouldUseLaunchSchemeArgsEnv = "YES"
      savedToolIdentifier = ""
      useCustomWorkingDirectory = "NO"
      debugDocumentVersioning = "YES">
      <BuildableProductRunnable
         runnableDebuggingMode = "0">
         <BuildableReference
            BuildableIdentifier = "primary"
            BlueprintIdentifier = "BA3CBE7419F7A93800CED4D5"
            BuildableName = "ios-simple-objc.app"
            BlueprintName = "ios-simple-objc"
            ReferencedContainer = "container:ios-simple-objc.xcodeproj">
         </BuildableReference>
      </BuildableProductRunnable>
   </ProfileAction>
   <AnalyzeAction
      buildConfiguration = "Debug">
   </AnalyzeAction>
   <ArchiveAction
      buildConfiguration = "Release"
      revealArchiveInOrganizer = "YES">
   </ArchiveAction>
</Scheme>
`
