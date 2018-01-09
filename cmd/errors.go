package cmd

import "github.com/bitrise-io/go-utils/colorstring"

// Tool ...
type Tool string

const (
	toolXcode   Tool = "Xcode"
	toolXamarin Tool = "Visual Studio"
)

// ArchiveError ...
type ArchiveError struct {
	tool Tool
	msg  string
}

// Error ...
func (e ArchiveError) Error() string {
	return `
------------------------------` + `
First of all ` + colorstring.Red("please make sure that you can Archive your app from "+e.tool+".") + `
codesigndoc only works if you can archive your app from ` + string(e.tool) + `.
If you can, and you get a valid IPA file if you export from ` + string(e.tool) + `,
` + colorstring.Red("please create an issue") + ` on GitHub at: https://github.com/bitrise-tools/codesigndoc/issues
with as many details & logs as you can share!
------------------------------

` + colorstring.Redf("Error: %s", e.msg)
}

// MissingCodesigningFilesError ...
type MissingCodesigningFilesError struct {
	msg string
}

const collectCodesigningFilesInfo = `To collect available code sign files, we search for installed Provisioning Profiles:"
- which has installed Codesign Identity in your Keychain"
- which can provision your application target's bundle ids"
- which has the project defined Capabilities set"
- which matches to the selected ipa export method"
`

func (e MissingCodesigningFilesError) Error() string {
	return e.msg + "\n" + collectCodesigningFilesInfo
}
