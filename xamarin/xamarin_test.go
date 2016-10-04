package xamarin

import (
	"testing"

	"github.com/bitrise-io/go-utils/testutil"
	"github.com/bitrise-tools/codesigndoc/common"
	"github.com/bitrise-tools/codesigndoc/provprofile"
	"github.com/stretchr/testify/require"
)

func Test_parseCodeSigningSettingsFromOutput(t *testing.T) {
	t.Log("Multiple Code Signing settings")
	{
		logOutput := `...

   	Target DeployOutputFiles:
   		Copying file from '/Users/USER/develop/tmp/xamarin-test/XamTest/iOS/obj/iPhone/Release/XamTest.iOS.exe.mdb' to
   '/Users/USER/develop/tmp/xamarin-test/XamTest/iOS/bin/iPhone/Release/XamTest.iOS.exe.mdb'
   		Copying file from '/Users/USER/develop/tmp/xamarin-test/XamTest/iOS/obj/iPhone/Release/XamTest.iOS.exe' to
   '/Users/USER/develop/tmp/xamarin-test/XamTest/iOS/bin/iPhone/Release/XamTest.iOS.exe'

   	Target _DetectSigningIdentity:
   		DetectSigningIdentity Task
   		  AppBundleName: XamTest.iOS
   		  AppManifest: Info.plist
   		  Keychain: <null>
   		  ProvisioningProfile: <null>
   		  RequireCodesigning: True
   		  SdkPlatform: iPhoneOS
   		  SdkIsSimulator: False
   		  SigningKey: iPhone Developer
   		Detected signing identity:
   		  Code Signing Key: "iPhone Developer: First Last (F72Z82XD37)" (CBAB0B7E123AE7AD790EE801EDFB45035360DB3F)
   		  Provisioning Profile: "iOS Team Provisioning Profile: *" (87af6d83-cb65-4dbe-aee7-f97a87d6fec1)
   		  Bundle Id: io.bitrise.xamtest
   		  App Id: 01SA2B3CDL.io.bitrise.xamtest

   	Target _ComputeBundleResourceOutputPaths:
   		ComputeBundleResourceOutputPaths Task
   		  AppBundleDir: bin/iPhone/Release/XamTest.iOS.app
   		  BundleIdentifier: io.bitrise.xamtest
   		  BundleResources:
   		    obj/iPhone/Release/ibtool-link/LaunchScreen.storyboardc/01J-lp-oVM-view-Ze5-6b-2t3.nib
   		    obj/iPhone/Release/ibtool-link/LaunchScreen.storyboardc/Info.plist
   		    obj/iPhone/Release/ibtool-link/LaunchScreen.storyboardc/UIViewController-01J-lp-oVM.nib
   		    obj/iPhone/Release/ibtool-link/Main.storyboardc/BYZ-38-t0r-view-8bC-Xf-vdC.nib
   		    obj/iPhone/Release/ibtool-link/Main.storyboardc/Info.plist
   		    obj/iPhone/Release/ibtool-link/Main.storyboardc/UIViewController-BYZ-38-t0r.nib
   		  IntermediateOutputPath: obj/iPhone/Release/
   		  OutputPath: bin/iPhone/Release/

...

   			Target _DetectSdkLocations:
   				DetectSdkLocations Task
   				  TargetFrameworkIdentifier: Xamarin.iOS
   				  TargetArchitectures: ARMv7, ARM64
   				  SdkVersion: 10.0
   				  XamarinSdkRoot: /Library/Frameworks/Xamarin.iOS.framework/Versions/Current
   				  SdkRoot: /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS10.0.sdk
   				  SdkDevPath: /Applications/Xcode.app/Contents/Developer
   				  SdkUsrPath: /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/usr
   				  SdkPlatform: iPhoneOS
   				  SdkIsSimulator: False

   			Target _DetectSigningIdentity:
   				DetectSigningIdentity Task
   				  AppBundleName: TodayExt
   				  AppManifest: Info.plist
   				  Keychain: <null>
   				  ProvisioningProfile: <null>
   				  RequireCodesigning: True
   				  SdkPlatform: iPhoneOS
   				  SdkIsSimulator: False
   				  SigningKey: iPhone Developer
   				Detected signing identity:
   				  Code Signing Key: "iPhone Distribution: BFirst BLast (B72Z82XD37)" (CBAB0B7E123AE7AD790EE801EDFB45035360DB3F)
   				  Provisioning Profile: "Prov Profile 43" (97af6d83-cb65-4dbe-aee7-f97a87d6fec1)
   				  Bundle Id: io.bitrise.xamtest.ios.todayext
   				  App Id: 01SA2B3CDL.io.bitrise.xamtest.ios.todayext
   		Done building project "/Users/USER/develop/tmp/xamarin-test/XamTest/TodayExt/TodayExt.csproj".


   	Target _CompileToNative:
   		MTouch Task
   		  AppBundleDir: bin/iPhone/Release/XamTest.iOS.app
   		  AppExtensionReferences:
   		    /Users/USER/develop/tmp/xamarin-test/XamTest/TodayExt/bin/iPhone/Release/TodayExt.appex
   		  AppManifest: bin/iPhone/Release/XamTest.iOS.app/Info.plist
   		  Architectures: ARMv7, ARM64
   		  ArchiveSymbols: <null>
   		  BitcodeEnabled: False
   		  CompiledEntitlements: obj/iPhone/Release/Entitlements.xcent
...


`

		parsedCodeSigningSettings, err := parseCodeSigningSettingsFromOutput(logOutput)
		require.NoError(t, err)
		testutil.EqualSlicesWithoutOrder(t, []common.CodeSigningIdentityInfo{
			common.CodeSigningIdentityInfo{Title: "iPhone Developer: First Last (F72Z82XD37)"},
			common.CodeSigningIdentityInfo{Title: "iPhone Distribution: BFirst BLast (B72Z82XD37)"},
		}, parsedCodeSigningSettings.Identities)
		testutil.EqualSlicesWithoutOrder(t, []provprofile.ProvisioningProfileInfo{
			provprofile.ProvisioningProfileInfo{Title: "iOS Team Provisioning Profile: *", UUID: "87af6d83-cb65-4dbe-aee7-f97a87d6fec1"},
			provprofile.ProvisioningProfileInfo{Title: "Prov Profile 43", UUID: "97af6d83-cb65-4dbe-aee7-f97a87d6fec1"},
		}, parsedCodeSigningSettings.ProvProfiles)
		require.Equal(t, []string{"01SA2B3CDL"},
			parsedCodeSigningSettings.TeamIDs)
		testutil.EqualSlicesWithoutOrder(t, []string{
			"01SA2B3CDL.io.bitrise.xamtest",
			"01SA2B3CDL.io.bitrise.xamtest.ios.todayext",
		},
			parsedCodeSigningSettings.AppBundleIDs)
	}
}
