package codesign

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

func TestFilterLatestProfiles(t *testing.T) {
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

	filtered := FilterLatestProfiles(profiles)
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
