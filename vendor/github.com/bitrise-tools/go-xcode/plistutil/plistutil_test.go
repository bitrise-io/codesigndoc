package plistutil

import (
	"reflect"
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

func TestPlistData_GetMapStringInterfaceArray(t *testing.T) {
	testSummariesData, err := NewPlistDataFromContent(paritalTestSummariesContent)
	if err != nil {
		t.Errorf("NewPlistDataFromContent(), got: %v, want: %v", err, nil)
	}
	const key = "Key"

	type args struct {
		forKey string
	}
	tests := []struct {
		name  string
		data  PlistData
		args  args
		want  []PlistData
		want1 bool
	}{
		{
			name: "Test ok case",
			data: PlistData{key: []interface{}{
				map[string]interface{}{"k1": "v1", "k2": "v2"},
				map[string]interface{}{"k3": "v3"},
			}},
			args: args{key},
			want: []PlistData{
				map[string]interface{}{"k1": "v1", "k2": "v2"},
				map[string]interface{}{"k3": "v3"},
			},
			want1: true,
		},
		{
			name:  "Test key not found",
			data:  PlistData{"otherKey": []PlistData{}},
			args:  args{key},
			want:  nil,
			want1: false,
		},
		{
			name: "Test failed to cast to interface{}",
			data: PlistData{key: []PlistData{
				map[string]interface{}{"k1": "v1", "k2": "v2"},
				map[string]interface{}{"k3": "v3"},
			}},
			args:  args{key},
			want:  nil,
			want1: false,
		},
		{
			name: "Failed to cast array element to map[string]interface{}",
			data: PlistData{key: []interface{}{
				map[string]string{"k1": "v1", "k2": "v2"},
				map[string]string{"k3": "v3"},
			}},
			args:  args{key},
			want:  nil,
			want1: false,
		},
		{
			name: "Intefration test with real plist data",
			data: testSummariesData,
			args: args{"Subtests"},
			want: []PlistData{
				map[string]interface{}{"TestIdentifier": "ios_simple_objcTests/testExample", "TestStatus": "Success"},
				map[string]interface{}{"TestIdentifier": "ios_simple_objcTests/testExample2", "TestStatus": "Success"},
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.data.GetMapStringInterfaceArray(tt.args.forKey)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlistData.GetMapStringInterfaceArray() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PlistData.GetMapStringInterfaceArray() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPlistData_GetFloat64(t *testing.T) {
	testSummariesData, err := NewPlistDataFromContent(paritalTestSummariesContent)
	if err != nil {
		t.Errorf("NewPlistDataFromContent(), got: %v, want: %v", err, nil)
	}
	const key = "Duration"
	const value = 0.00072991847991943359

	type args struct {
		forKey string
	}
	tests := []struct {
		name  string
		data  PlistData
		args  args
		want  float64
		want1 bool
	}{
		{
			name:  "Read float, ok",
			data:  map[string]interface{}{key: value},
			args:  args{key},
			want:  value,
			want1: true,
		},
		{
			name:  "Key not found",
			data:  map[string]interface{}{"otherKey": value},
			args:  args{key},
			want:  0,
			want1: false,
		},
		{
			name:  "Read int value, fail",
			data:  map[string]interface{}{key: 23},
			args:  args{key},
			want:  0,
			want1: false,
		},
		{
			name:  "Integration test with real plist data",
			data:  testSummariesData,
			args:  args{"Duration"},
			want:  0.34774100780487061,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.data.GetFloat64(tt.args.forKey)
			if got != tt.want {
				t.Errorf("PlistData.GetFloat64() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PlistData.GetFloat64() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
