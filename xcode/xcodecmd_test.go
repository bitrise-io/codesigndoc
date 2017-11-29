package xcode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseSchemesFromXcodeOutput(t *testing.T) {
	xcout := `Information about project "SampleAppWithCocoapods":
    Targets:
        SampleAppWithCocoapods
        SampleAppWithCocoapodsTests

    Build Configurations:
        Debug
        Release

    If no build configuration is specified and -scheme is not passed then "Release" is used.

    Schemes:
        SampleAppWithCocoapods`
	parsedSchemes := parseSchemesFromXcodeOutput(xcout)
	require.Equal(t, []string{"SampleAppWithCocoapods"}, parsedSchemes)
}
