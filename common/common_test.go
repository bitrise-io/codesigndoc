package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBundleIDFromAppID(t *testing.T) {
	for anAppID, expectedBundleID := range map[string]string{
		"01TeaM02ID.com.company.app": "com.company.app",
		"A.com.company.app":          "com.company.app",
		"1-TEAM02ID.com.company.app": "", // invalid char in team ID, should be only letters and numbers
		"*.com.company.app":          "", // invalid char in team ID, should be only letters and numbers
		"01TEAM02ID.*":               "",
		"01TEAM02ID":                 "",
		"":                           "",
	} {
		t.Log(" * anAppID:", anAppID)
		require.Equal(t, expectedBundleID, BundleIDFromAppID(anAppID))
	}
}
