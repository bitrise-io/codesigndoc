package xcworkspace

import "github.com/bitrise-tools/xcode-project/xcodeproj"

func Example() {
	workspace, err := Open("workspace.xcworkspace")
	if err != nil {
		panic(err)
	}

	var fileRefLocations []string
	for _, fileRef := range workspace.FileRefs {
		pth, err := fileRef.AbsPath("workspace_dir")
		if err != nil {
			panic(err)
		}

		fileRefLocations = append(fileRefLocations, pth)
	}
	for _, group := range workspace.Groups {
		groupPth, err := group.AbsPath("workspace_dir")
		if err != nil {
			panic(err)
		}

		for _, fileRef := range group.FileRefs {
			pth, err := fileRef.AbsPath(groupPth)
			if err != nil {
				panic(err)
			}

			fileRefLocations = append(fileRefLocations, pth)
		}
	}

	var projects []xcodeproj.XcodeProj
	for _, fileRefLocation := range fileRefLocations {
		if !xcodeproj.IsXcodeProj(fileRefLocation) {
			continue
		}

		project, err := xcodeproj.Open(fileRefLocation)
		if err != nil {
			panic(err)
		}
		projects = append(projects, project)
	}
}
