package cmd

import (
	"testing"
	"time"

	"github.com/bitrise-tools/go-xcode/profileutil"
	"github.com/stretchr/testify/require"
)

func createTime(t *testing.T, timeStr string) time.Time {
	date, err := time.Parse("2006.01.02", timeStr)
	require.NoError(t, err)
	return date
}

func timeToString(date time.Time) string {
	return date.Format("2006.01.02")
}

func Test_filterLatestProfiles(t *testing.T) {
	profiles := []profileutil.ProvisioningProfileInfoModel{
		{
			Name:           "Profile 1",
			BundleID:       "*",
			ExpirationDate: createTime(t, "2017.11.01"),
		},
		{
			Name:           "Profile 1",
			BundleID:       "*",
			ExpirationDate: createTime(t, "2017.12.01"),
		},
		{
			Name:           "Profile 2",
			BundleID:       "io.bitrise",
			ExpirationDate: createTime(t, "2017.10.01"),
		},
	}

	filtered := filterLatestProfiles(profiles)
	require.Equal(t, 2, len(filtered))

	desiredProfileExpireMap := map[string]bool{
		"2017.12.01": false,
		"2017.10.01": false,
	}
	for _, profile := range filtered {
		expire := timeToString(profile.ExpirationDate)
		if _, ok := desiredProfileExpireMap[expire]; ok {
			desiredProfileExpireMap[expire] = true
		}
	}
	for _, found := range desiredProfileExpireMap {
		require.True(t, found)
	}
}

func Test_trimProjectpath(t *testing.T) {
	type args struct {
		projpth string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No change",
			args: args{projpth: "Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj"},
			want: "Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj",
		},
		{
			name: "Quotation mark",
			args: args{projpth: "\"Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj\""},
			want: "Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj",
		},
		{
			name: "Apostrophe",
			args: args{projpth: "'Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj'"},
			want: "Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj",
		},
		{
			name: "Quotation mark With whitespace",
			args: args{projpth: "\" Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj \""},
			want: "Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj",
		},
		{
			name: "Apostrophe With whitespace",
			args: args{projpth: "' Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj '"},
			want: "Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj",
		},
		{
			name: "New line",
			args: args{projpth: "\nDevelop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj\n"},
			want: "Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj",
		},
		{
			name: "Multiple",
			args: args{projpth: "\n  \"  \nDevelop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj\n  '  ' ''''\n\n\""},
			want: "Develop/XCode/XcodeArchiveTest/XcodeArchiveTest.xcodeproj",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimProjectpath(tt.args.projpth); got != tt.want {
				t.Errorf("trimProjpth() = %v, want %v", got, tt.want)
			}
		})
	}
}
