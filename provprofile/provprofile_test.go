package provprofile

import "testing"
import "github.com/stretchr/testify/require"
import "github.com/bitrise-io/go-utils/testutil"

func Test_ProvisioningProfileModel_TeamID(t *testing.T) {
	t.Log("Empty TeamIdentifiers")
	{
		provProfile := ProvisioningProfileModel{}
		teamID, err := provProfile.TeamID()
		require.EqualError(t, err, "No TeamIdentifier specified")
		require.Equal(t, "", teamID)
	}

	t.Log("One TeamIdentifier - valid")
	{
		provProfile := ProvisioningProfileModel{
			TeamIdentifiers: []string{"TeamID1"},
		}
		teamID, err := provProfile.TeamID()
		require.NoError(t, err)
		require.Equal(t, "TeamID1", teamID)
	}

	t.Log("One empty TeamIdentifier - invalid")
	{
		provProfile := ProvisioningProfileModel{
			TeamIdentifiers: []string{""},
		}
		teamID, err := provProfile.TeamID()
		require.EqualError(t, err, "An empty item specified for TeamIdentifier")
		require.Equal(t, "", teamID)
	}

	t.Log("Two TeamIdentifiers")
	{
		provProfile := ProvisioningProfileModel{
			TeamIdentifiers: []string{
				"TeamID1", "TeamID2",
			},
		}
		teamID, err := provProfile.TeamID()
		require.EqualError(t, err, "More than one TeamIdentifier specified")
		require.Equal(t, "", teamID)
	}
}

func Test_ProvisioningProfileFileInfoModels_CollectTeamIDs(t *testing.T) {
	t.Log("Empty")
	{
		ppFileInfos := ProvisioningProfileFileInfoModels{}
		teamIDs, err := ppFileInfos.CollectTeamIDs()
		require.NoError(t, err)
		require.Equal(t, []string{}, teamIDs)
	}

	t.Log("One item")
	{
		ppFileInfos := ProvisioningProfileFileInfoModels{
			ProvisioningProfileFileInfoModel{
				ProvisioningProfileInfo: ProvisioningProfileModel{
					TeamIdentifiers: []string{"TeamID1"},
				},
			},
		}
		teamIDs, err := ppFileInfos.CollectTeamIDs()
		require.NoError(t, err)
		require.Equal(t, []string{"TeamID1"}, teamIDs)
	}

	t.Log("Two distinct items with different team IDs")
	{
		ppFileInfos := ProvisioningProfileFileInfoModels{
			{
				ProvisioningProfileInfo: ProvisioningProfileModel{
					TeamIdentifiers: []string{"TeamID1"},
				},
			},
			{
				ProvisioningProfileInfo: ProvisioningProfileModel{
					TeamIdentifiers: []string{"TeamID2"},
				},
			},
		}
		teamIDs, err := ppFileInfos.CollectTeamIDs()
		require.NoError(t, err)
		testutil.EqualSlicesWithoutOrder(t, []string{"TeamID1", "TeamID2"}, teamIDs)
	}

	t.Log("Two distinct items with same team IDs")
	{
		ppFileInfos := ProvisioningProfileFileInfoModels{
			{
				ProvisioningProfileInfo: ProvisioningProfileModel{
					TeamIdentifiers: []string{"TeamID1"},
				},
			},
			{
				ProvisioningProfileInfo: ProvisioningProfileModel{
					TeamIdentifiers: []string{"TeamID1"},
				},
			},
		}
		teamIDs, err := ppFileInfos.CollectTeamIDs()
		require.NoError(t, err)
		require.Equal(t, []string{"TeamID1"}, teamIDs)
	}

	t.Log("Two items, the second one's TeamID is invalid")
	{
		ppFileInfos := ProvisioningProfileFileInfoModels{
			{
				ProvisioningProfileInfo: ProvisioningProfileModel{
					UUID:            "UUID1",
					TeamIdentifiers: []string{"TeamID1"},
				},
			},
			{
				ProvisioningProfileInfo: ProvisioningProfileModel{
					UUID:            "UUID2",
					TeamIdentifiers: []string{""},
				},
			},
		}
		teamIDs, err := ppFileInfos.CollectTeamIDs()
		require.EqualError(t, err, "Team ID error for profile (uuid: UUID2), error: An empty item specified for TeamIdentifier")
		require.Equal(t, []string{}, teamIDs)
	}
}
