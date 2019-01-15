package xcworkspace

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileRef(t *testing.T) {
	var fileRef FileRef
	require.NoError(t, xml.Unmarshal([]byte(fileRefContent), &fileRef))

	require.Equal(t, "group:../SubProjectTests/Info.plist", fileRef.Location)

	{
		refType, pth, err := fileRef.TypeAndPath()
		require.NoError(t, err)
		require.Equal(t, "../SubProjectTests/Info.plist", pth)
		require.Equal(t, GroupFileRefType, refType)
	}

	{
		dir := "/workspace_dir/group"
		pth, err := fileRef.AbsPath(dir)
		require.NoError(t, err)
		require.Equal(t, "/workspace_dir/SubProjectTests/Info.plist", pth)
	}

}

const fileRefContent = `<FileRef
	location = "group:../SubProjectTests/Info.plist">
</FileRef>`
