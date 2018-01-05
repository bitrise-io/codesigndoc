package cmd

import "github.com/bitrise-io/go-utils/colorstring"

// XcodeArchiveError ...
type XcodeArchiveError struct {
	msg string
}

func (e XcodeArchiveError) Error() string {
	return ArchiveError{"Xcode", e.msg}.Error()
}

// XamarinArchiveError ...
type XamarinArchiveError struct {
	msg string
}

func (e XamarinArchiveError) Error() string {
	return ArchiveError{"Visual Studio", e.msg}.Error()
}

// ArchiveError ...
type ArchiveError struct {
	tool string
	msg  string
}

// Error ...
func (e ArchiveError) Error() string {
	return `
------------------------------` + `
First of all ` + colorstring.Red("please make sure that you can Archive your app from "+e.tool+".") + `
codesigndoc only works if you can archive your app from ` + e.tool + `.
If you can, and you get a valid IPA file if you export from ` + e.tool + `,
` + colorstring.Red("please create an issue") + ` on GitHub at: https://github.com/bitrise-tools/codesigndoc/issues
with as many details & logs as you can share!
------------------------------

` + colorstring.Redf("Error: %s", e.msg)
}
