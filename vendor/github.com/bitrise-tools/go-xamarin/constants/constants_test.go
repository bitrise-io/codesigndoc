package constants

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSDK(t *testing.T) {
	t.Log("it parses android")
	{
		projectType, err := ParseSDK("android")
		require.NoError(t, err)
		require.Equal(t, SDKAndroid, projectType)
	}

	t.Log("it parses ios")
	{
		projectType, err := ParseSDK("ios")
		require.NoError(t, err)
		require.Equal(t, SDKIOS, projectType)
	}

	t.Log("it parses tvos")
	{
		projectType, err := ParseSDK("tvos")
		require.NoError(t, err)
		require.Equal(t, SDKTvOS, projectType)
	}

	t.Log("it parses macos")
	{
		projectType, err := ParseSDK("macos")
		require.NoError(t, err)
		require.Equal(t, SDKMacOS, projectType)
	}

	t.Log("it failes for unknown type")
	{
		projectType, err := ParseSDK("go")
		require.Error(t, err)
		require.Equal(t, SDKUnknown, projectType)
	}
}

func TestParseTestFramwork(t *testing.T) {
	t.Log("it parses xamarin-uitest")
	{
		projectType, err := ParseTestFramwork("xamarin-uitest")
		require.NoError(t, err)
		require.Equal(t, TestFrameworkXamarinUITest, projectType)
	}

	t.Log("it parses nunit-test")
	{
		projectType, err := ParseTestFramwork("nunit-test")
		require.NoError(t, err)
		require.Equal(t, TestFrameworkNunitTest, projectType)
	}

	t.Log("it parses nunit-lite-test")
	{
		projectType, err := ParseTestFramwork("nunit-lite-test")
		require.NoError(t, err)
		require.Equal(t, TestFrameworkNunitLiteTest, projectType)
	}

	t.Log("it failes for unknown type")
	{
		projectType, err := ParseTestFramwork("go")
		require.Error(t, err)
		require.Equal(t, TestFrameworkUnknown, projectType)
	}
}

func TestParseProjectTypeGUID(t *testing.T) {
	t.Log("it parses XamarinAndroid GUID")
	{
		xamarinAndroidGUIDs := []string{
			"EFBA0AD7-5A72-4C68-AF49-83D382785DCF",
			"10368E6C-D01B-4462-8E8B-01FC667A7035",
		}
		for _, guid := range xamarinAndroidGUIDs {
			projectType, err := ParseProjectTypeGUID(guid)
			require.NoError(t, err)
			require.Equal(t, SDKAndroid, projectType)
		}
	}

	t.Log("it parses XamarinIOS GUID")
	{
		xamarinIOSGUIDs := []string{
			"E613F3A2-FE9C-494F-B74E-F63BCB86FEA6",
			"6BC8ED88-2882-458C-8E55-DFD12B67127B",
			"F5B4F3BC-B597-4E2B-B552-EF5D8A32436F",
			"FEACFBD2-3405-455C-9665-78FE426C6842",
			"8FFB629D-F513-41CE-95D2-7ECE97B6EEEC",
			"EE2C853D-36AF-4FDB-B1AD-8E90477E2198",
		}
		for _, guid := range xamarinIOSGUIDs {
			projectType, err := ParseProjectTypeGUID(guid)
			require.NoError(t, err)
			require.Equal(t, SDKIOS, projectType)
		}
	}

	t.Log("it parses XamarinTvOS GUID")
	{
		xamarinTvOSGUIDs := []string{
			"06FA79CB-D6CD-4721-BB4B-1BD202089C55",
		}
		for _, guid := range xamarinTvOSGUIDs {
			projectType, err := ParseProjectTypeGUID(guid)
			require.NoError(t, err)
			require.Equal(t, SDKTvOS, projectType)
		}
	}

	t.Log("it parses MonoMac & XamarinMac GUID")
	{
		monoMacGUIDs := []string{
			"1C533B1C-72DD-4CB1-9F6B-BF11D93BCFBE",
			"948B3504-5B70-4649-8FE4-BDE1FB46EC69",
		}
		xamarinMacGUIDs := []string{
			"42C0BBD9-55CE-4FC1-8D90-A7348ABAFB23",
			"A3F8F2AB-B479-4A4A-A458-A89E7DC349F1",
		}
		macOSGUIDs := append(monoMacGUIDs, xamarinMacGUIDs...)
		for _, guid := range macOSGUIDs {
			projectType, err := ParseProjectTypeGUID(guid)
			require.NoError(t, err)
			require.Equal(t, SDKMacOS, projectType)
		}
	}

	t.Log("it failes for unkown GUID")
	{
		xamarinTvOSGUIDs := []string{
			"06FA79CB5",
		}
		for _, guid := range xamarinTvOSGUIDs {
			projectType, err := ParseProjectTypeGUID(guid)
			require.Error(t, err)
			require.Equal(t, SDKUnknown, projectType)
		}
	}
}

func TestParseOutputType(t *testing.T) {
	t.Log("it parses apk")
	{
		outputType, err := ParseOutputType("apk")
		require.NoError(t, err)
		require.Equal(t, OutputTypeAPK, outputType)
	}

	t.Log("it parses xcarchive")
	{
		outputType, err := ParseOutputType("xcarchive")
		require.NoError(t, err)
		require.Equal(t, OutputTypeXCArchive, outputType)
	}

	t.Log("it parses ipa")
	{
		outputType, err := ParseOutputType("ipa")
		require.NoError(t, err)
		require.Equal(t, OutputTypeIPA, outputType)
	}

	t.Log("it parses dsym")
	{
		outputType, err := ParseOutputType("dsym")
		require.NoError(t, err)
		require.Equal(t, OutputTypeDSYM, outputType)
	}

	t.Log("it parses pkg")
	{
		outputType, err := ParseOutputType("pkg")
		require.NoError(t, err)
		require.Equal(t, OutputTypePKG, outputType)
	}

	t.Log("it parses app")
	{
		outputType, err := ParseOutputType("app")
		require.NoError(t, err)
		require.Equal(t, OutputTypeAPP, outputType)
	}

	t.Log("it parses dll")
	{
		outputType, err := ParseOutputType("dll")
		require.NoError(t, err)
		require.Equal(t, OutputTypeDLL, outputType)
	}

	t.Log("it failes for unknown type")
	{
		outputType, err := ParseOutputType("zip")
		require.Error(t, err)
		require.Equal(t, OutputTypeUnknown, outputType)
	}
}
