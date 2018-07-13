package profileutil

import "testing"
import "github.com/stretchr/testify/require"

func TestIsXcodeManaged(t *testing.T) {
	xcodeManagedNames := []string{
		"XC iOS: custom.bundle.id",
		"XC tvOS: custom.bundle.id",
		"iOS Team Provisioning Profile: another.custom.bundle.id",
		"tvOS Team Provisioning Profile: another.custom.bundle.id",
		"iOS Team Store Provisioning Profile: my.bundle.id",
		"tvOS Team Store Provisioning Profile: my.bundle.id",
		"Mac Team Provisioning Profile: my.bundle.id",
		"Mac Team Store Provisioning Profile: my.bundle.id",
	}
	nonXcodeManagedNames := []string{
		"Test Profile Name",
		"iOS Distribution Profile: test.bundle.id",
		"iOS Dev",
		"tvOS Distribution Profile: test.bundle.id",
		"tvOS Dev",
		"Mac Distribution Profile: test.bundle.id",
		"Mac Dev",
	}

	for _, profileName := range xcodeManagedNames {
		require.Equal(t, true, IsXcodeManaged(profileName))
	}

	for _, profileName := range nonXcodeManagedNames {
		require.Equal(t, false, IsXcodeManaged(profileName))
	}
}
