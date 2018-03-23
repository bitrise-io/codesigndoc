package cmd

import (
	"reflect"
	"testing"
	"time"

	"github.com/bitrise-tools/codesigndoc/bitriseclient"
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

func TestGetAppFromUserSelection(t *testing.T) {
	tests := []struct {
		name            string
		selectedApp     string
		appList         []bitriseclient.Application
		wantSeledtedApp bitriseclient.Application
		wantErr         bool
	}{
		{
			name:        "Success",
			selectedApp: "bitrise-xcodearchivetest (git@bitbucket.org:Birmachera/bitrise-xcodearchivetest.git)",
			appList: []bitriseclient.Application{
				bitriseclient.Application{
					Slug:        "683dec47f6f7db20",
					Title:       "bitrise-xcodearchivetest",
					ProjectType: "ios",
					Provider:    "bitbucket",
					RepoOwner:   "Birmachera",
					RepoURL:     "git@bitbucket.org:Birmachera/bitrise-xcodearchivetest.git",
					RepoSlug:    "bitrise-xcodearchivetest",
					IsDisabled:  false,
					Status:      1,
					IsPublic:    false,
					Owner: bitriseclient.Owner{
						AccountType: "user",
						Name:        "BirmacherAkos",
						Slug:        "f88644b20a74fb29",
					},
				},
			},
			wantSeledtedApp: bitriseclient.Application{
				Slug:        "683dec47f6f7db20",
				Title:       "bitrise-xcodearchivetest",
				ProjectType: "ios",
				Provider:    "bitbucket",
				RepoOwner:   "Birmachera",
				RepoURL:     "git@bitbucket.org:Birmachera/bitrise-xcodearchivetest.git",
				RepoSlug:    "bitrise-xcodearchivetest",
				IsDisabled:  false,
				Status:      1,
				IsPublic:    false,
				Owner: bitriseclient.Owner{
					AccountType: "user",
					Name:        "BirmacherAkos",
					Slug:        "f88644b20a74fb29",
				},
			},
			wantErr: false,
		},
		{
			name:        "Second",
			selectedApp: "Second (git@bitbucket.org:Birmachera/second.git)",
			appList: []bitriseclient.Application{
				bitriseclient.Application{
					Slug:        "683dec47f6f7db20",
					Title:       "bitrise-xcodearchivetest",
					ProjectType: "ios",
					Provider:    "bitbucket",
					RepoOwner:   "Birmachera",
					RepoURL:     "git@bitbucket.org:Birmachera/bitrise-xcodearchivetest.git",
					RepoSlug:    "bitrise-xcodearchivetest",
					IsDisabled:  false,
					Status:      1,
					IsPublic:    false,
					Owner: bitriseclient.Owner{
						AccountType: "user",
						Name:        "BirmacherAkos",
						Slug:        "f88644b20a74fb29",
					},
				},
				bitriseclient.Application{
					Slug:        "68sdsdsdsd3dec47f6f7db20",
					Title:       "Second",
					ProjectType: "ios",
					Provider:    "bitbucket",
					RepoOwner:   "Birmachera",
					RepoURL:     "git@bitbucket.org:Birmachera/second.git",
					RepoSlug:    "bitrise-xcodearchivetest",
					IsDisabled:  false,
					Status:      1,
					IsPublic:    false,
					Owner: bitriseclient.Owner{
						AccountType: "user",
						Name:        "BirmacherAkos",
						Slug:        "f88644b20a74fb29",
					},
				},
			},
			wantSeledtedApp: bitriseclient.Application{
				Slug:        "68sdsdsdsd3dec47f6f7db20",
				Title:       "Second",
				ProjectType: "ios",
				Provider:    "bitbucket",
				RepoOwner:   "Birmachera",
				RepoURL:     "git@bitbucket.org:Birmachera/second.git",
				RepoSlug:    "bitrise-xcodearchivetest",
				IsDisabled:  false,
				Status:      1,
				IsPublic:    false,
				Owner: bitriseclient.Owner{
					AccountType: "user",
					Name:        "BirmacherAkos",
					Slug:        "f88644b20a74fb29",
				},
			},
			wantErr: false,
		},
		{
			name:        "Failed",
			selectedApp: " (git@bitbucket.)",
			appList: []bitriseclient.Application{
				bitriseclient.Application{
					Slug:        "683dec47f6f7db20",
					Title:       "bitrise-xcodearchivetest",
					ProjectType: "ios",
					Provider:    "bitbucket",
					RepoOwner:   "Birmachera",
					RepoURL:     "git@bitbucket.org:Birmachera/bitrise-xcodearchivetest.git",
					RepoSlug:    "bitrise-xcodearchivetest",
					IsDisabled:  false,
					Status:      1,
					IsPublic:    false,
					Owner: bitriseclient.Owner{
						AccountType: "user",
						Name:        "BirmacherAkos",
						Slug:        "f88644b20a74fb29",
					},
				},
				bitriseclient.Application{
					Slug:        "68sdsdsdsd3dec47f6f7db20",
					Title:       "Second",
					ProjectType: "ios",
					Provider:    "bitbucket",
					RepoOwner:   "Birmachera",
					RepoURL:     "git@bitbucket.org:Birmachera/second.git",
					RepoSlug:    "bitrise-xcodearchivetest",
					IsDisabled:  false,
					Status:      1,
					IsPublic:    false,
					Owner: bitriseclient.Owner{
						AccountType: "user",
						Name:        "BirmacherAkos",
						Slug:        "f88644b20a74fb29",
					},
				},
			},
			wantSeledtedApp: bitriseclient.Application{},
			wantErr:         true,
		},
		{
			name:        "No repo url",
			selectedApp: "bitrise-xcodearchivetest ()",
			appList: []bitriseclient.Application{
				bitriseclient.Application{
					Slug:        "683dec47f6f7db20",
					Title:       "bitrise-xcodearchivetest",
					ProjectType: "ios",
					Provider:    "bitbucket",
					RepoOwner:   "Birmachera",
					RepoURL:     "git@bitbucket.org:Birmachera/bitrise-xcodearchivetest.git",
					RepoSlug:    "bitrise-xcodearchivetest",
					IsDisabled:  false,
					Status:      1,
					IsPublic:    false,
					Owner: bitriseclient.Owner{
						AccountType: "user",
						Name:        "BirmacherAkos",
						Slug:        "f88644b20a74fb29",
					},
				},
				bitriseclient.Application{
					Slug:        "68sdsdsdsd3dec47f6f7db20",
					Title:       "Second",
					ProjectType: "ios",
					Provider:    "bitbucket",
					RepoOwner:   "Birmachera",
					RepoURL:     "git@bitbucket.org:Birmachera/second.git",
					RepoSlug:    "bitrise-xcodearchivetest",
					IsDisabled:  false,
					Status:      1,
					IsPublic:    false,
					Owner: bitriseclient.Owner{
						AccountType: "user",
						Name:        "BirmacherAkos",
						Slug:        "f88644b20a74fb29",
					},
				},
			},
			wantSeledtedApp: bitriseclient.Application{},
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSeledtedApp, err := getAppFromUserSelection(tt.selectedApp, tt.appList)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAppFromUserSelection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSeledtedApp, tt.wantSeledtedApp) {
				t.Errorf("getAppFromUserSelection() = %v, want %v", gotSeledtedApp, tt.wantSeledtedApp)
			}
		})
	}
}
