package project

import (
	"github.com/bitrise-io/xcode-project/xcodeproj"
	"github.com/bitrise-io/xcode-project/xcscheme"
	"github.com/bitrise-io/xcode-project/xcworkspace"
)

// HasScheme represents a struct that implements Scheme.
type HasScheme interface {
	Scheme(string) (*xcscheme.Scheme, string, error)
}

// Scheme returns the project or workspace scheme by name.
func Scheme(pth string, name string) (*xcscheme.Scheme, string, error) {
	var p HasScheme
	var err error
	if xcodeproj.IsXcodeProj(pth) {
		p, err = xcodeproj.Open(pth)
	} else {
		p, err = xcworkspace.Open(pth)
	}
	if err != nil {
		return nil, "", err
	}
	return p.Scheme(name)
}
