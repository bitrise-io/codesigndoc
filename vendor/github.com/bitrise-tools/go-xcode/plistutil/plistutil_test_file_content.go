package plistutil

const infoPlistContent = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>CFBundleName</key>
    <string>ios-simple-objc</string>
    <key>DTXcode</key>
    <string>0832</string>
    <key>DTSDKName</key>
    <string>iphoneos10.3</string>
    <key>UILaunchStoryboardName</key>
    <string>LaunchScreen</string>
    <key>DTSDKBuild</key>
    <string>14E269</string>
    <key>CFBundleDevelopmentRegion</key>
    <string>en</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>BuildMachineOSBuild</key>
    <string>16F73</string>
    <key>DTPlatformName</key>
    <string>iphoneos</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>UIMainStoryboardFile</key>
    <string>Main</string>
    <key>CFBundleSupportedPlatforms</key>
    <array>
      <string>iPhoneOS</string>
    </array>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>UIRequiredDeviceCapabilities</key>
    <array>
      <string>armv7</string>
    </array>
    <key>CFBundleExecutable</key>
    <string>ios-simple-objc</string>
    <key>DTCompiler</key>
    <string>com.apple.compilers.llvm.clang.1_0</string>
    <key>UISupportedInterfaceOrientations~ipad</key>
    <array>
      <string>UIInterfaceOrientationPortrait</string>
      <string>UIInterfaceOrientationPortraitUpsideDown</string>
      <string>UIInterfaceOrientationLandscapeLeft</string>
      <string>UIInterfaceOrientationLandscapeRight</string>
    </array>
    <key>CFBundleIdentifier</key>
    <string>Bitrise.ios-simple-objc</string>
    <key>MinimumOSVersion</key>
    <string>8.1</string>
    <key>DTXcodeBuild</key>
    <string>8E2002</string>
    <key>DTPlatformVersion</key>
    <string>10.3</string>
    <key>LSRequiresIPhoneOS</key>
    <true/>
    <key>UISupportedInterfaceOrientations</key>
    <array>
      <string>UIInterfaceOrientationPortrait</string>
      <string>UIInterfaceOrientationLandscapeLeft</string>
      <string>UIInterfaceOrientationLandscapeRight</string>
    </array>
    <key>CFBundleSignature</key>
    <string>????</string>
    <key>UIDeviceFamily</key>
    <array>
      <integer>1</integer>
      <integer>2</integer>
    </array>
    <key>DTPlatformBuild</key>
    <string>14E269</string>
  </dict>
</plist>
`

const developmentProfileContent = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AppIDName</key>
	<string>Bitrise Test</string>
	<key>ApplicationIdentifierPrefix</key>
	<array>
	<string>9NS4</string>
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
			<string>9NS4.*</string>
		</array>
		<key>get-task-allow</key>
		<true/>
		<key>application-identifier</key>
		<string>9NS4.*</string>
		<key>com.apple.developer.team-identifier</key>
		<string>9NS4</string>
	</dict>
	<key>ExpirationDate</key>
	<date>2017-09-22T11:28:46Z</date>
	<key>Name</key>
	<string>Bitrise Test Development</string>
	<key>ProvisionedDevices</key>
	<array>
		<string>b138</string>
	</array>
	<key>TeamIdentifier</key>
	<array>
		<string>9NS4</string>
	</array>
	<key>TeamName</key>
	<string>Some Dude</string>
	<key>TimeToLive</key>
	<integer>365</integer>
	<key>UUID</key>
	<string>4b617a5f</string>
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
	<string>9NS4</string>
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
			<string>9NS4.*</string>
		</array>
		<key>get-task-allow</key>
		<false/>
		<key>application-identifier</key>
		<string>9NS4.*</string>
		<key>com.apple.developer.team-identifier</key>
		<string>9NS4</string>
		<key>beta-reports-active</key>
		<true/>
	</dict>
	<key>ExpirationDate</key>
	<date>2017-09-21T13:20:06Z</date>
	<key>Name</key>
	<string>Bitrise Test App Store</string>
	<key>TeamIdentifier</key>
	<array>
		<string>9NS4</string>
	</array>
	<key>TeamName</key>
	<string>Some Dude</string>
	<key>TimeToLive</key>
	<integer>364</integer>
	<key>UUID</key>
	<string>a60668dd</string>
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
	<string>9NS4</string>
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
			<string>9NS4.*</string>
		</array>
		<key>get-task-allow</key>
		<false/>
		<key>application-identifier</key>
		<string>9NS4.*</string>
		<key>com.apple.developer.team-identifier</key>
		<string>9NS4</string>
	</dict>
	<key>ExpirationDate</key>
	<date>2017-09-21T13:20:06Z</date>
	<key>Name</key>
	<string>Bitrise Test Ad Hoc</string>
	<key>ProvisionedDevices</key>
	<array>
		<string>b138</string>
	</array>
	<key>TeamIdentifier</key>
	<array>
		<string>9NS4</string>
	</array>
	<key>TeamName</key>
	<string>Some Dude</string>
	<key>TimeToLive</key>
	<integer>364</integer>
	<key>UUID</key>
	<string>26668300</string>
	<key>Version</key>
	<integer>1</integer>
</dict>`

const enterpriseProfileContent = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AppIDName</key>
	<string>Bitrise Test</string>
	<key>ApplicationIdentifierPrefix</key>
	<array>
	<string>PF3BP78LQ8</string>
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
			<string>PF3BP78LQ8.*</string>
		</array>
		<key>get-task-allow</key>
		<false/>
		<key>application-identifier</key>
		<string>9NS4.*</string>
		<key>com.apple.developer.team-identifier</key>
		<string>9NS4</string>
	</dict>
	<key>ExpirationDate</key>
	<date>2016-10-04T13:32:46Z</date>
	<key>Name</key>
	<string>Bitrise Test Enterprise</string>
	<key>ProvisionsAllDevices</key>
	<true/>
	<key>TeamIdentifier</key>
	<array>
		<string>PF3BP78LQ8</string>
	</array>
	<key>TeamName</key>
	<string>Some Dude</string>
	<key>TimeToLive</key>
	<integer>365</integer>
	<key>UUID</key>
	<string>8d6caa15</string>
	<key>Version</key>
	<integer>1</integer>
</dict>`
