package buildtools

// BuildTool ...
type BuildTool uint

const (
	// Msbuild ...
	Msbuild BuildTool = iota
	// Xbuild ...
	Xbuild
)
