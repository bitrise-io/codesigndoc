package xcworkspace

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
)

// FileRef ...
type FileRef struct {
	Location string `xml:"location,attr"`
}

// FileRefType ...
type FileRefType string

// Known FileRefTypes
const (
	AbsoluteFileRefType  FileRefType = "absolute"
	GroupFileRefType     FileRefType = "group"
	ContainerFileRefType FileRefType = "container"
)

// TypeAndPath ...
func (f FileRef) TypeAndPath() (FileRefType, string, error) {
	s := strings.Split(f.Location, ":")
	if len(s) != 2 {
		return "", "", fmt.Errorf("unknown file reference location (%s)", f.Location)
	}

	switch s[0] {
	case "absolute":
		return AbsoluteFileRefType, s[1], nil
	case "group":
		return GroupFileRefType, s[1], nil
	case "container":
		return ContainerFileRefType, s[1], nil
	default:
		return "", "", fmt.Errorf("unknown file reference type: %s", s[0])
	}
}

// AbsPath ...
func (f FileRef) AbsPath(dir string) (string, error) {
	t, pth, err := f.TypeAndPath()
	if err != nil {
		return "", err
	}

	var absPth string
	switch t {
	case AbsoluteFileRefType:
		absPth = pth
	case GroupFileRefType, ContainerFileRefType:
		absPth = filepath.Join(dir, pth)
	}

	return pathutil.AbsPath(absPth)
}
