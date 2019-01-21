package profileutil

import (
	"testing"

	"github.com/bitrise-tools/go-xcode/exportoptions"
	"github.com/bitrise-tools/go-xcode/plistutil"
	"github.com/stretchr/testify/require"
)

func TestPlistData(t *testing.T) {
	t.Log("development profile specifies development export method")
	{
		profile, err := plistutil.NewPlistDataFromContent(developmentProfileContent)
		require.NoError(t, err)
		require.Equal(t, "4b617a5f-e31e-4edc-9460-718a5abacd05", PlistData(profile).GetUUID())
		require.Equal(t, "Bitrise Test Development", PlistData(profile).GetName())
		require.Equal(t, "9NS44DLTN7.*", PlistData(profile).GetApplicationIdentifier())
		require.Equal(t, "*", PlistData(profile).GetBundleIdentifier())
		require.Equal(t, exportoptions.MethodDevelopment, PlistData(profile).GetExportMethod())
		require.Equal(t, "9NS44DLTN7", PlistData(profile).GetTeamID())
		require.Equal(t, "Some Dude", PlistData(profile).GetTeamName())
		require.Equal(t, "2016-09-22T11:28:46Z", PlistData(profile).GetCreationDate().Format("2006-01-02T15:04:05Z"))
		require.Equal(t, "2017-09-22T11:28:46Z", PlistData(profile).GetExpirationDate().Format("2006-01-02T15:04:05Z"))
		require.Equal(t, []string{"b13813075ad9b298cb9a9f28555c49573d8bc322"}, PlistData(profile).GetProvisionedDevices())
		require.Equal(t, [][]uint8{[]uint8{}}, PlistData(profile).GetDeveloperCertificates())
		require.Equal(t, false, PlistData(profile).GetProvisionsAllDevices())
	}

	t.Log("app store profile specifies app-store export method")
	{
		profile, err := plistutil.NewPlistDataFromContent(appStoreProfileContent)
		require.NoError(t, err)
		require.Equal(t, "a60668dd-191a-4770-8b1e-b453b87aa60b", PlistData(profile).GetUUID())
		require.Equal(t, "Bitrise Test App Store", PlistData(profile).GetName())
		require.Equal(t, "9NS44DLTN7.*", PlistData(profile).GetApplicationIdentifier())
		require.Equal(t, "*", PlistData(profile).GetBundleIdentifier())
		require.Equal(t, exportoptions.MethodAppStore, PlistData(profile).GetExportMethod())
		require.Equal(t, "9NS44DLTN7", PlistData(profile).GetTeamID())
		require.Equal(t, "Some Dude", PlistData(profile).GetTeamName())
		require.Equal(t, "2016-09-22T11:29:12Z", PlistData(profile).GetCreationDate().Format("2006-01-02T15:04:05Z"))
		require.Equal(t, "2017-09-21T13:20:06Z", PlistData(profile).GetExpirationDate().Format("2006-01-02T15:04:05Z"))
		require.Equal(t, []string(nil), PlistData(profile).GetProvisionedDevices())
		require.Equal(t, [][]uint8{[]uint8{}}, PlistData(profile).GetDeveloperCertificates())
		require.Equal(t, false, PlistData(profile).GetProvisionsAllDevices())
	}

	t.Log("ad hoc profile specifies ad-hoc export method")
	{
		profile, err := plistutil.NewPlistDataFromContent(adHocProfileContent)
		require.NoError(t, err)
		require.Equal(t, "26668300-5743-46a1-8e00-7023e2e35c7d", PlistData(profile).GetUUID())
		require.Equal(t, "Bitrise Test Ad Hoc", PlistData(profile).GetName())
		require.Equal(t, "9NS44DLTN7.*", PlistData(profile).GetApplicationIdentifier())
		require.Equal(t, "*", PlistData(profile).GetBundleIdentifier())
		require.Equal(t, exportoptions.MethodAdHoc, PlistData(profile).GetExportMethod())
		require.Equal(t, "9NS44DLTN7", PlistData(profile).GetTeamID())
		require.Equal(t, "Some Dude", PlistData(profile).GetTeamName())
		require.Equal(t, "2016-09-22T11:29:38Z", PlistData(profile).GetCreationDate().Format("2006-01-02T15:04:05Z"))
		require.Equal(t, "2017-09-21T13:20:06Z", PlistData(profile).GetExpirationDate().Format("2006-01-02T15:04:05Z"))
		require.Equal(t, []string{"b13813075ad9b298cb9a9f28555c49573d8bc322"}, PlistData(profile).GetProvisionedDevices())
		require.Equal(t, [][]uint8{[]uint8{}}, PlistData(profile).GetDeveloperCertificates())
		require.Equal(t, false, PlistData(profile).GetProvisionsAllDevices())
	}

	t.Log("it creates model from enterprise profile content")
	{
		profile, err := plistutil.NewPlistDataFromContent(enterpriseProfileContent)
		require.NoError(t, err)
		require.Equal(t, "8d6caa15-ac49-48f9-9bd3-ce9244add6a0", PlistData(profile).GetUUID())
		require.Equal(t, "Bitrise Test Enterprise", PlistData(profile).GetName())
		require.Equal(t, "9NS44DLTN7.com.Bitrise.Test", PlistData(profile).GetApplicationIdentifier())
		require.Equal(t, "com.Bitrise.Test", PlistData(profile).GetBundleIdentifier())
		require.Equal(t, exportoptions.MethodEnterprise, PlistData(profile).GetExportMethod())
		require.Equal(t, "9NS44DLTN7", PlistData(profile).GetTeamID())
		require.Equal(t, "Bitrise", PlistData(profile).GetTeamName())
		require.Equal(t, "2015-10-05T13:32:46Z", PlistData(profile).GetCreationDate().Format("2006-01-02T15:04:05Z"))
		require.Equal(t, "2016-10-04T13:32:46Z", PlistData(profile).GetExpirationDate().Format("2006-01-02T15:04:05Z"))
		require.Equal(t, []string(nil), PlistData(profile).GetProvisionedDevices())
		require.Equal(t, [][]uint8{[]uint8{}}, PlistData(profile).GetDeveloperCertificates())
		require.Equal(t, true, PlistData(profile).GetProvisionsAllDevices())
	}
}

func TestTVOSPlistData(t *testing.T) {
	t.Log("it creates model from tvOS appstore profile content")
	{
		profile, err := plistutil.NewPlistDataFromContent(tvOSAppStoreProfileContent)
		require.NoError(t, err)
		require.Equal(t, "dec523d5-624b-44bd-8d16-6d1d69c63276", PlistData(profile).GetUUID())
		require.Equal(t, "Bitrise app-store - (bdh.NPO-Live.bitrise.sample)", PlistData(profile).GetName())
		require.Equal(t, "72SA8V3WYL.bdh.NPO-Live.bitrise.sample", PlistData(profile).GetApplicationIdentifier())
		require.Equal(t, "bdh.NPO-Live.bitrise.sample", PlistData(profile).GetBundleIdentifier())
		require.Equal(t, exportoptions.MethodAppStore, PlistData(profile).GetExportMethod())
		require.Equal(t, "72SA8V3WYL", PlistData(profile).GetTeamID())
		require.Equal(t, "Bitrise", PlistData(profile).GetTeamName())
		require.Equal(t, "2018-10-24T11:22:30Z", PlistData(profile).GetCreationDate().Format("2006-01-02T15:04:05Z"))
		require.Equal(t, "2019-04-16T08:42:18Z", PlistData(profile).GetExpirationDate().Format("2006-01-02T15:04:05Z"))
		require.Equal(t, []string(nil), PlistData(profile).GetProvisionedDevices())
		require.Equal(t, [][]uint8{[]uint8{}}, PlistData(profile).GetDeveloperCertificates())
		require.Equal(t, false, PlistData(profile).GetProvisionsAllDevices())
	}
}

const developmentProfileContent = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AppIDName</key>
	<string>Bitrise Test</string>
	<key>ApplicationIdentifierPrefix</key>
	<array>
	<string>9NS44DLTN7</string>
	</array>
	<key>CreationDate</key>
	<date>2016-09-22T11:28:46Z</date>
	<key>Platform</key>
	<array>
		<string>iOS</string>
	</array>
	<key>DeveloperCertificates</key>
	<array>
		<data></data>
	</array>
	<key>Entitlements</key>
	<dict>
		<key>keychain-access-groups</key>
		<array>
			<string>9NS44DLTN7.*</string>
		</array>
		<key>get-task-allow</key>
		<true/>
		<key>application-identifier</key>
		<string>9NS44DLTN7.*</string>
		<key>com.apple.developer.team-identifier</key>
		<string>9NS44DLTN7</string>
	</dict>
	<key>ExpirationDate</key>
	<date>2017-09-22T11:28:46Z</date>
	<key>Name</key>
	<string>Bitrise Test Development</string>
	<key>ProvisionedDevices</key>
	<array>
		<string>b13813075ad9b298cb9a9f28555c49573d8bc322</string>
	</array>
	<key>TeamIdentifier</key>
	<array>
		<string>9NS44DLTN7</string>
	</array>
	<key>TeamName</key>
	<string>Some Dude</string>
	<key>TimeToLive</key>
	<integer>365</integer>
	<key>UUID</key>
	<string>4b617a5f-e31e-4edc-9460-718a5abacd05</string>
	<key>Version</key>
	<integer>1</integer>
</dict>`

const appStoreProfileContent = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AppIDName</key>
	<string>Bitrise Test</string>
	<key>ApplicationIdentifierPrefix</key>
	<array>
	<string>9NS44DLTN7</string>
	</array>
	<key>CreationDate</key>
	<date>2016-09-22T11:29:12Z</date>
	<key>Platform</key>
	<array>
		<string>iOS</string>
	</array>
	<key>DeveloperCertificates</key>
	<array>
		<data></data>
	</array>
	<key>Entitlements</key>
	<dict>
		<key>keychain-access-groups</key>
		<array>
			<string>9NS44DLTN7.*</string>
		</array>
		<key>get-task-allow</key>
		<false/>
		<key>application-identifier</key>
		<string>9NS44DLTN7.*</string>
		<key>com.apple.developer.team-identifier</key>
		<string>9NS44DLTN7</string>
		<key>beta-reports-active</key>
		<true/>
	</dict>
	<key>ExpirationDate</key>
	<date>2017-09-21T13:20:06Z</date>
	<key>Name</key>
	<string>Bitrise Test App Store</string>
	<key>TeamIdentifier</key>
	<array>
		<string>9NS44DLTN7</string>
	</array>
	<key>TeamName</key>
	<string>Some Dude</string>
	<key>TimeToLive</key>
	<integer>364</integer>
	<key>UUID</key>
	<string>a60668dd-191a-4770-8b1e-b453b87aa60b</string>
	<key>Version</key>
	<integer>1</integer>
</dict>`

const adHocProfileContent = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AppIDName</key>
	<string>Bitrise Test</string>
	<key>ApplicationIdentifierPrefix</key>
	<array>
	<string>9NS44DLTN7</string>
	</array>
	<key>CreationDate</key>
	<date>2016-09-22T11:29:38Z</date>
	<key>Platform</key>
	<array>
		<string>iOS</string>
	</array>
	<key>DeveloperCertificates</key>
	<array>
		<data></data>
	</array>
	<key>Entitlements</key>
	<dict>
		<key>keychain-access-groups</key>
		<array>
			<string>9NS44DLTN7.*</string>
		</array>
		<key>get-task-allow</key>
		<false/>
		<key>application-identifier</key>
		<string>9NS44DLTN7.*</string>
		<key>com.apple.developer.team-identifier</key>
		<string>9NS44DLTN7</string>
	</dict>
	<key>ExpirationDate</key>
	<date>2017-09-21T13:20:06Z</date>
	<key>Name</key>
	<string>Bitrise Test Ad Hoc</string>
	<key>ProvisionedDevices</key>
	<array>
		<string>b13813075ad9b298cb9a9f28555c49573d8bc322</string>
	</array>
	<key>TeamIdentifier</key>
	<array>
		<string>9NS44DLTN7</string>
	</array>
	<key>TeamName</key>
	<string>Some Dude</string>
	<key>TimeToLive</key>
	<integer>364</integer>
	<key>UUID</key>
	<string>26668300-5743-46a1-8e00-7023e2e35c7d</string>
	<key>Version</key>
	<integer>1</integer>
</dict>`

const enterpriseProfileContent = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AppIDName</key>
	<string>Test</string>
	<key>ApplicationIdentifierPrefix</key>
	<array>
	<string>9NS44DLTN7</string>
	</array>
	<key>CreationDate</key>
	<date>2015-10-05T13:32:46Z</date>
	<key>Platform</key>
	<array>
		<string>iOS</string>
	</array>
	<key>DeveloperCertificates</key>
	<array>
		<data></data>
	</array>
	<key>Entitlements</key>
	<dict>
		<key>keychain-access-groups</key>
		<array>
			<string>9NS44DLTN7.*</string>
		</array>
		<key>get-task-allow</key>
		<false/>
		<key>application-identifier</key>
		<string>9NS44DLTN7.com.Bitrise.Test</string>
		<key>com.apple.developer.team-identifier</key>
		<string>9NS44DLTN7</string>

	</dict>
	<key>ExpirationDate</key>
	<date>2016-10-04T13:32:46Z</date>
	<key>Name</key>
	<string>Bitrise Test Enterprise</string>
	<key>ProvisionsAllDevices</key>
	<true/>
	<key>TeamIdentifier</key>
	<array>
		<string>9NS44DLTN7</string>
	</array>
	<key>TeamName</key>
	<string>Bitrise</string>
	<key>TimeToLive</key>
	<integer>365</integer>
	<key>UUID</key>
	<string>8d6caa15-ac49-48f9-9bd3-ce9244add6a0</string>
	<key>Version</key>
	<integer>1</integer>
</dict>`

const tvOSAppStoreProfileContent = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AppIDName</key>
	<string>Bitrise  bdh NPOLive bitrise sample be2b4e3cfb0f2a967b404820aa18e09c</string>
	<key>ApplicationIdentifierPrefix</key>
	<array>
	<string>72SA8V3WYL</string>
	</array>
	<key>CreationDate</key>
	<date>2018-10-24T11:22:30Z</date>
	<key>Platform</key>
	<array>
		<string>tvOS</string>
	</array>
	<key>DeveloperCertificates</key>
	<array>
		<data></data>
	</array>
    <key>IsXcodeManaged</key>
    <false/>
	<key>Entitlements</key>
	<dict>
		<key>keychain-access-groups</key>
		<array>
			<string>72SA8V3WYL.*</string>
		</array>
		<key>get-task-allow</key>
		<false/>
		<key>application-identifier</key>
		<string>72SA8V3WYL.bdh.NPO-Live.bitrise.sample</string>
		<key>com.apple.developer.team-identifier</key>
		<string>72SA8V3WYL</string>
		<key>beta-reports-active</key>
		<true/>
	</dict>
	<key>ExpirationDate</key>
	<date>2019-04-16T08:42:18Z</date>
	<key>Name</key>
	<string>Bitrise app-store - (bdh.NPO-Live.bitrise.sample)</string>
	<key>TeamIdentifier</key>
	<array>
		<string>72SA8V3WYL</string>
	</array>
	<key>TeamName</key>
	<string>Bitrise</string>
	<key>TimeToLive</key>
	<integer>173</integer>
	<key>UUID</key>
	<string>dec523d5-624b-44bd-8d16-6d1d69c63276</string>
	<key>Version</key>
	<integer>1</integer>
</dict>`
