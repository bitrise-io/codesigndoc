package xcworkspace

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
)

// Group ...
type Group struct {
	Location string    `xml:"location,attr"`
	FileRefs []FileRef `xml:"FileRef"`
	Groups   []Group   `xml:"Group"`
}

// AbsPath ...
func (g Group) AbsPath(dir string) (string, error) {
	s := strings.Split(g.Location, ":")
	if len(s) != 2 {
		return "", fmt.Errorf("unknown group location (%s)", g.Location)
	}
	pth := filepath.Join(dir, s[1])
	return pathutil.AbsPath(pth)
}

// FileLocations ...
func (g Group) FileLocations(dir string) ([]string, error) {
	var fileLocations []string

	groupPth, err := g.AbsPath(dir)
	if err != nil {
		return nil, err
	}

	for _, fileRef := range g.FileRefs {
		fileLocation, err := fileRef.AbsPath(groupPth)
		if err != nil {
			return nil, err
		}
		fileLocations = append(fileLocations, fileLocation)
	}

	for _, group := range g.Groups {
		groupFileLocations, err := group.FileLocations(groupPth)
		if err != nil {
			return nil, err
		}
		fileLocations = append(fileLocations, groupFileLocations...)
	}

	return fileLocations, nil
}
