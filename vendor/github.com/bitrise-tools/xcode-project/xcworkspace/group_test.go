package xcworkspace

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroup(t *testing.T) {
	var group Group
	require.NoError(t, xml.Unmarshal([]byte(groupContent), &group))

	containerDir := "/workspace_dir"

	require.Equal(t, "group:Group", group.Location)
	require.Equal(t, 1, len(group.FileRefs))
	require.Equal(t, 1, len(group.Groups))

	pth, err := group.AbsPath(containerDir)
	require.NoError(t, err)
	require.Equal(t, "/workspace_dir/Group", pth)

	{
		fileRef := group.FileRefs[0]
		require.Equal(t, "group:SubProject/SubProject.xcodeproj", fileRef.Location)

		pth, err := fileRef.AbsPath(pth)
		require.NoError(t, err)
		require.Equal(t, "/workspace_dir/Group/SubProject/SubProject.xcodeproj", pth)
	}

	{
		subGroup := group.Groups[0]

		require.Equal(t, "group:SubProject/SubProject", subGroup.Location)
		require.Equal(t, 3, len(subGroup.FileRefs))
		require.Equal(t, 2, len(subGroup.Groups))

		pth, err := subGroup.AbsPath(pth)
		require.NoError(t, err)
		require.Equal(t, "/workspace_dir/Group/SubProject/SubProject", pth)

		group := subGroup

		{
			fileRef := group.FileRefs[0]
			require.Equal(t, "group:ViewController.swift", fileRef.Location)

			pth, err := fileRef.AbsPath(pth)
			require.NoError(t, err)
			require.Equal(t, "/workspace_dir/Group/SubProject/SubProject/ViewController.swift", pth)
		}

		{
			subGroup := group.Groups[0]

			require.Equal(t, "group:Assets.xcassets", subGroup.Location)
			require.Equal(t, 1, len(subGroup.FileRefs))
			require.Equal(t, 1, len(subGroup.Groups))

			pth, err := subGroup.AbsPath(pth)
			require.NoError(t, err)
			require.Equal(t, "/workspace_dir/Group/SubProject/SubProject/Assets.xcassets", pth)
		}
	}
}

func TestFileLocations(t *testing.T) {
	var group Group
	require.NoError(t, xml.Unmarshal([]byte(groupContent), &group))

	containerDir := "/workspace_dir"

	fileLocations, err := group.FileLocations(containerDir)
	require.NoError(t, err)
	require.Equal(t, 8, len(fileLocations))
	require.Equal(t, []string{
		"/workspace_dir/Group/SubProject/SubProject.xcodeproj",
		"/workspace_dir/Group/SubProject/SubProject/ViewController.swift",
		"/workspace_dir/Group/SubProject/SubProject/AppDelegate.swift",
		"/workspace_dir/Group/SubProject/SubProject/Info.plist",
		"/workspace_dir/Group/SubProject/SubProject/Assets.xcassets/Contents.json",
		"/workspace_dir/Group/SubProject/SubProject/Assets.xcassets/AppIcon.appiconset/Contents.json",
		"/workspace_dir/Group/SubProject/SubProject/Base.lproj/LaunchScreen.storyboard",
		"/workspace_dir/Group/SubProject/SubProject/Base.lproj/Main.storyboard",
	}, fileLocations)
}

const groupContent = `<Group
location = "group:Group"
name = "Group">
<Group
   location = "group:SubProject/SubProject"
   name = "SubProject">
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
</Group>`
