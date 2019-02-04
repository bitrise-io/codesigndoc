package ios

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunRubyScriptForOutput(t *testing.T) {
	gemfileContent := `source 'https://rubygems.org'
gem 'json'
`

	rubyScriptContent := `require 'json'

puts "#{{ :test_key => 'test_value' }.to_json}"
`

	expectedOut := "{\"test_key\":\"test_value\"}"
	actualOut, err := runRubyScriptForOutput(rubyScriptContent, gemfileContent, "", []string{})
	require.NoError(t, err)
	require.Equal(t, expectedOut, actualOut)
}
