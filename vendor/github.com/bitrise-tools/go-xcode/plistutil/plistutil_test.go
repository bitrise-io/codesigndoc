package plistutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAnalyzeInfoPlist(t *testing.T) {
	infoPlistData, err := NewPlistDataFromContent(infoPlistContent)
	require.NoError(t, err)

	appTitle, ok := infoPlistData.GetString("CFBundleName")
	require.Equal(t, true, ok)
	require.Equal(t, "ios-simple-objc", appTitle)

	bundleID, _ := infoPlistData.GetString("CFBundleIdentifier")
	require.Equal(t, true, ok)
	require.Equal(t, "Bitrise.ios-simple-objc", bundleID)

	version, ok := infoPlistData.GetString("CFBundleShortVersionString")
	require.Equal(t, true, ok)
	require.Equal(t, "1.0", version)

	buildNumber, ok := infoPlistData.GetString("CFBundleVersion")
	require.Equal(t, true, ok)
	require.Equal(t, "1", buildNumber)

	minOSVersion, ok := infoPlistData.GetString("MinimumOSVersion")
	require.Equal(t, true, ok)
	require.Equal(t, "8.1", minOSVersion)

	deviceFamilyList, ok := infoPlistData.GetUInt64Array("UIDeviceFamily")
	require.Equal(t, true, ok)
	require.Equal(t, 2, len(deviceFamilyList))
	require.Equal(t, uint64(1), deviceFamilyList[0])
	require.Equal(t, uint64(2), deviceFamilyList[1])
}

func TestAnalyzeEmbeddedProfile(t *testing.T) {
	profileData, err := NewPlistDataFromContent(appStoreProfileContent)
	require.NoError(t, err)

	creationDate, ok := profileData.GetTime("CreationDate")
	require.Equal(t, true, ok)
	expectedCreationDate, err := time.Parse("2006-01-02T15:04:05Z", "2016-09-22T11:29:12Z")
	require.NoError(t, err)
	require.Equal(t, true, creationDate.Equal(expectedCreationDate))

	expirationDate, ok := profileData.GetTime("ExpirationDate")
	require.Equal(t, true, ok)
	expectedExpirationDate, err := time.Parse("2006-01-02T15:04:05Z", "2017-09-21T13:20:06Z")
	require.NoError(t, err)
	require.Equal(t, true, expirationDate.Equal(expectedExpirationDate))

	deviceUDIDList, ok := profileData.GetStringArray("ProvisionedDevices")
	require.Equal(t, false, ok)
	require.Equal(t, 0, len(deviceUDIDList))

	teamName, ok := profileData.GetString("TeamName")
	require.Equal(t, true, ok)
	require.Equal(t, "Some Dude", teamName)

	profileName, ok := profileData.GetString("Name")
	require.Equal(t, true, ok)
	require.Equal(t, "Bitrise Test App Store", profileName)

	provisionsAlldevices, ok := profileData.GetBool("ProvisionsAllDevices")
	require.Equal(t, false, ok)
	require.Equal(t, false, provisionsAlldevices)
}

func TestGetBool(t *testing.T) {
	profileData, err := NewPlistDataFromContent(enterpriseProfileContent)
	require.NoError(t, err)

	allDevices, ok := profileData.GetBool("ProvisionsAllDevices")
	require.Equal(t, true, ok)
	require.Equal(t, true, allDevices)
}

func TestGetTime(t *testing.T) {
	profileData, err := NewPlistDataFromContent(developmentProfileContent)
	require.NoError(t, err)

	expire, ok := profileData.GetTime("ExpirationDate")
	require.Equal(t, true, ok)

	// 2017-09-22T11:28:46Z
	desiredExpire, err := time.Parse("2006-01-02T15:04:05Z", "2017-09-22T11:28:46Z")
	require.NoError(t, err)
	require.Equal(t, true, expire.Equal(desiredExpire))
}

func TestGetInt(t *testing.T) {
	profileData, err := NewPlistDataFromContent(developmentProfileContent)
	require.NoError(t, err)

	version, ok := profileData.GetUInt64("Version")
	require.Equal(t, true, ok)
	require.Equal(t, uint64(1), version)
}

func TestGetStringArray(t *testing.T) {
	profileData, err := NewPlistDataFromContent(developmentProfileContent)
	require.NoError(t, err)

	devices, ok := profileData.GetStringArray("ProvisionedDevices")
	require.Equal(t, true, ok)
	require.Equal(t, 1, len(devices))
	require.Equal(t, "b138", devices[0])
}

func TestGetMapStringInterface(t *testing.T) {
	profileData, err := NewPlistDataFromContent(developmentProfileContent)
	require.NoError(t, err)

	entitlements, ok := profileData.GetMapStringInterface("Entitlements")
	require.Equal(t, true, ok)

	teamID, ok := entitlements.GetString("com.apple.developer.team-identifier")
	require.Equal(t, true, ok)
	require.Equal(t, "9NS4", teamID)
}
