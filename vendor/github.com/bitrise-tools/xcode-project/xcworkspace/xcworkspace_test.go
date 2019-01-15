package xcworkspace

import (
	"path/filepath"
	"testing"

	"github.com/bitrise-tools/xcode-project/testhelper"
	"github.com/stretchr/testify/require"
)

func TestScheme(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-samples/xcode-project-test.git")
	workspace, err := Open(filepath.Join(dir, "XcodeProj.xcworkspace"))
	require.NoError(t, err)

	{
		scheme, container, err := workspace.Scheme("SubProjectWatchKitApp (Notification) Scheme")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(dir, "Group/SubProject/SubProject.xcodeproj"), container)
		require.Equal(t, "SubProjectWatchKitApp (Notification) Scheme", scheme.Name)
	}

	{
		scheme, container, err := workspace.Scheme("Not Exist Scheme")
		require.EqualError(t, err, "scheme Not Exist Scheme not found in XcodeProj")
		require.Equal(t, "", container)
		require.Equal(t, "", scheme.Name)
	}
}

func TestSchemes(t *testing.T) {
	dir := testhelper.GitCloneIntoTmpDir(t, "https://github.com/bitrise-samples/xcode-project-test.git")
	workspace, err := Open(filepath.Join(dir, "XcodeProj.xcworkspace"))
	require.NoError(t, err)

	schemesByContainer, err := workspace.Schemes()
	require.NoError(t, err)
	require.Equal(t, 3, len(schemesByContainer))

	{
		schemes := schemesByContainer[filepath.Join(dir, "XcodeProj.xcworkspace")]
		require.Equal(t, 1, len(schemes))
		require.Equal(t, "WorkspaceScheme", schemes[0].Name)
	}

	{
		schemes := schemesByContainer[filepath.Join(dir, "XcodeProj.xcodeproj")]
		require.Equal(t, 2, len(schemes))
		require.Equal(t, "ProjectScheme", schemes[0].Name)
		require.Equal(t, "ProjectTodayExtensionScheme", schemes[1].Name)
	}

	{
		schemes := schemesByContainer[filepath.Join(dir, "Group/SubProject/SubProject.xcodeproj")]
		require.Equal(t, 4, len(schemes))
		require.Equal(t, "SubProjectScheme", schemes[0].Name)
		require.Equal(t, "SubProjectTestScheme", schemes[1].Name)
		require.Equal(t, "SubProjectWatchKitApp (Notification) Scheme", schemes[2].Name)
		require.Equal(t, "SubProjectWatchKitAppScheme", schemes[3].Name)
	}
}

func TestWorkspaceFileLocations(t *testing.T) {
	workspaceContentsPth := testhelper.CreateTmpFile(t, "contents.xcworkspacedata", workspaceContentsContent)
	workspacePth := filepath.Dir(workspaceContentsPth)

	workspace, err := Open(workspacePth)
	require.NoError(t, err)

	workspaceDir := filepath.Dir(workspacePth)

	fileLocations, err := workspace.FileLocations()
	require.NoError(t, err)
	require.Equal(t, []string{
		filepath.Join(workspaceDir, "XcodeProj.xcodeproj"),
		filepath.Join(workspaceDir, "Group/SubProject/SubProject.xcodeproj"),
		filepath.Join(workspaceDir, "Group/SubProject/SubProjectTests/Info.plist"),
		filepath.Join(workspaceDir, "Group/SubProject/SubProject/ViewController.swift"),
		filepath.Join(workspaceDir, "Group/SubProject/SubProject/AppDelegate.swift"),
		filepath.Join(workspaceDir, "Group/SubProject/SubProject/Info.plist"),
		filepath.Join(workspaceDir, "Group/SubProject/SubProject/Assets.xcassets/Contents.json"),
		filepath.Join(workspaceDir, "Group/SubProject/SubProject/Assets.xcassets/AppIcon.appiconset/Contents.json"),
		filepath.Join(workspaceDir, "Group/SubProject/SubProject/Base.lproj/LaunchScreen.storyboard"),
		filepath.Join(workspaceDir, "Group/SubProject/SubProject/Base.lproj/Main.storyboard"),
	}, fileLocations)
}

func TestWorkspaceProjectFileLocations(t *testing.T) {
	workspaceContentsPth := testhelper.CreateTmpFile(t, "contents.xcworkspacedata", workspaceContentsContent)
	workspacePth := filepath.Dir(workspaceContentsPth)

	workspace, err := Open(workspacePth)
	require.NoError(t, err)

	workspaceDir := filepath.Dir(workspacePth)

	fileLocations, err := workspace.ProjectFileLocations()
	require.NoError(t, err)
	require.Equal(t, []string{
		filepath.Join(workspaceDir, "XcodeProj.xcodeproj"),
		filepath.Join(workspaceDir, "Group/SubProject/SubProject.xcodeproj"),
	}, fileLocations)
}

func TestOpen(t *testing.T) {
	workspaceContentsPth := testhelper.CreateTmpFile(t, "contents.xcworkspacedata", workspaceContentsContent)
	workspacePth := filepath.Dir(workspaceContentsPth)

	workspace, err := Open(workspacePth)
	require.NoError(t, err)

	require.Equal(t, filepath.Base(workspacePth), workspace.Name)
	require.Equal(t, workspacePth, workspace.Path)
	require.Equal(t, 1, len(workspace.FileRefs))
	require.Equal(t, "group:XcodeProj.xcodeproj", workspace.FileRefs[0].Location)
	require.Equal(t, 1, len(workspace.Groups))
	require.Equal(t, "group:Group", workspace.Groups[0].Location)
}

func TestIsWorkspace(t *testing.T) {
	require.True(t, IsWorkspace("./BitriseSample.xcworkspace"))
	require.False(t, IsWorkspace("./BitriseSample.xcodeproj"))
}

const workspaceContentsContent = `<?xml version="1.0" encoding="UTF-8"?>
<Workspace
   version = "1.0">
   <Group
      location = "group:Group"
      name = "Group">
      <Group
         location = "group:SubProject/SubProject"
         name = "SubProject">
         <FileRef
            location = "group:../SubProjectTests/Info.plist">
         </FileRef>
         <FileRef
            location = "group:ViewController.swift">
         </FileRef>
         <Group
            location = "group:Assets.xcassets"
            name = "Assets.xcassets">
            <Group
               location = "group:AppIcon.appiconset"
               name = "AppIcon.appiconset">
               <FileRef
                  location = "group:Contents.json">
               </FileRef>
            </Group>
            <FileRef
               location = "group:Contents.json">
            </FileRef>
         </Group>
         <Group
            location = "group:Base.lproj"
            name = "Base.lproj">
            <FileRef
               location = "group:LaunchScreen.storyboard">
            </FileRef>
            <FileRef
               location = "group:Main.storyboard">
            </FileRef>
         </Group>
         <FileRef
            location = "group:AppDelegate.swift">
         </FileRef>
         <FileRef
            location = "group:Info.plist">
         </FileRef>
      </Group>
      <FileRef
         location = "group:SubProject/SubProject.xcodeproj">
      </FileRef>
   </Group>
   <FileRef
      location = "group:XcodeProj.xcodeproj">
   </FileRef>
</Workspace>
`
