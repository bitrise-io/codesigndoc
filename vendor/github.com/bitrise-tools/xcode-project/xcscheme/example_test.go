package xcscheme

import (
	"fmt"
)

func Example() {
	scheme, err := Open("scheme.xcscheme")
	if err != nil {
		panic(err)
	}

	fmt.Printf("archive action's default configuration name: %s\n", scheme.ArchiveAction.BuildConfiguration)

	for _, entry := range scheme.BuildAction.BuildActionEntries {
		if entry.BuildForArchiving == "YES" {

		}

		if entry.BuildForTesting == "YES" {

		}

		targetContainerProjectPth, err := entry.BuildableReference.ReferencedContainerAbsPath("scheme_container_project.xcodeproj")
		if err != nil {
			panic(err)
		}

		fmt.Printf("target with id: %s can be found in project: %s\n", entry.BuildableReference.BlueprintIdentifier, targetContainerProjectPth)
	}
}
