package simulator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLatestSimulatorInfoFromSimctlOut(t *testing.T) {
	t.Log("iOS iPhone 6s Plus")
	{
		info, version, err := getLatestSimulatorInfoFromSimctlOut(simctlListOut, "iOS", "iPhone 6s Plus")
		require.NoError(t, err)
		require.Equal(t, "iOS 10.2", version)
		require.Equal(t, "iPhone 6s Plus", info.Name)
		require.Equal(t, "3924585A-F638-46BC-9F93-99354C058D19", info.ID)
		require.Equal(t, "Shutdown", info.Status)
		require.Equal(t, "", info.StatusOther)
	}
}

func TestGetSimulatorInfoFromSimctlOut(t *testing.T) {
	t.Log("iOS 10.1 iPhone 6s Plus")
	{
		info, err := getSimulatorInfoFromSimctlOut(simctlListOut, "iOS 10.1", "iPhone 6s Plus")
		require.NoError(t, err)
		require.Equal(t, "iPhone 6s Plus", info.Name)
		require.Equal(t, "8255740B-6952-4802-93FF-52B1F0D34BBF", info.ID)
		require.Equal(t, "Shutdown", info.Status)
		require.Equal(t, "", info.StatusOther)
	}
}

func TestIs64BitArchitecture(t *testing.T) {
	t.Log("64-bit simulators")
	{
		iPhoneSimulators := []string{
			"iPhone 5S",
			"iPhone 6", "iPhone 6 Plus", "iPhone 6S", "iPhone 6S Plus",
			"iPhone SE",
			"iPhone 7", "iPhone 7 Plus",
		}

		for _, iPhoneSimulator := range iPhoneSimulators {
			is64bit, err := Is64BitArchitecture(iPhoneSimulator)
			require.NoError(t, err, iPhoneSimulator)
			require.Equal(t, true, is64bit, iPhoneSimulator)
		}

		iPadSimulators := []string{
			"iPad Mini 2", "iPad Mini 3", "iPad Mini 4",
			"iPad Air", "iPad Air 2",
			"iPad Pro (12.9 inch)", "iPad Pro (9.7 inch)",
		}

		for _, iPadSimulator := range iPadSimulators {
			is64bit, err := Is64BitArchitecture(iPadSimulator)
			require.NoError(t, err, iPadSimulator)
			require.Equal(t, true, is64bit, iPadSimulator)
		}
	}

	t.Log("not 64-bit simulators")
	{
		iPhoneSimulators := []string{
			"iPhone 5", "iPhone 5C",
			"iPhone 4", "iPhone 4s",
			"iPhone 3G", "iPhone 3GS",
		}

		for _, iPhoneSimulator := range iPhoneSimulators {
			is64bit, err := Is64BitArchitecture(iPhoneSimulator)
			require.NoError(t, err, iPhoneSimulator)
			require.Equal(t, false, is64bit, iPhoneSimulator)
		}

		iPadSimulators := []string{
			"iPad", "iPad 2", "iPad 3", "iPad 4",
			"iPad Mini",
		}

		for _, iPadSimulator := range iPadSimulators {
			is64bit, err := Is64BitArchitecture(iPadSimulator)
			require.NoError(t, err, iPadSimulator)
			require.Equal(t, false, is64bit, iPadSimulator)
		}
	}

}

const simctlListOut = `== Device Types ==
iPhone 4s (com.apple.CoreSimulator.SimDeviceType.iPhone-4s)
iPhone 5 (com.apple.CoreSimulator.SimDeviceType.iPhone-5)
iPhone 5s (com.apple.CoreSimulator.SimDeviceType.iPhone-5s)
iPhone 6 (com.apple.CoreSimulator.SimDeviceType.iPhone-6)
iPhone 6 Plus (com.apple.CoreSimulator.SimDeviceType.iPhone-6-Plus)
iPhone 6s (com.apple.CoreSimulator.SimDeviceType.iPhone-6s)
iPhone 6s Plus (com.apple.CoreSimulator.SimDeviceType.iPhone-6s-Plus)
iPhone 7 (com.apple.CoreSimulator.SimDeviceType.iPhone-7)
iPhone 7 Plus (com.apple.CoreSimulator.SimDeviceType.iPhone-7-Plus)
iPhone SE (com.apple.CoreSimulator.SimDeviceType.iPhone-SE)
iPad 2 (com.apple.CoreSimulator.SimDeviceType.iPad-2)
iPad Retina (com.apple.CoreSimulator.SimDeviceType.iPad-Retina)
iPad Air (com.apple.CoreSimulator.SimDeviceType.iPad-Air)
iPad Air 2 (com.apple.CoreSimulator.SimDeviceType.iPad-Air-2)
iPad Pro (9.7-inch) (com.apple.CoreSimulator.SimDeviceType.iPad-Pro--9-7-inch-)
iPad Pro (12.9-inch) (com.apple.CoreSimulator.SimDeviceType.iPad-Pro)
Apple TV 1080p (com.apple.CoreSimulator.SimDeviceType.Apple-TV-1080p)
Apple Watch - 38mm (com.apple.CoreSimulator.SimDeviceType.Apple-Watch-38mm)
Apple Watch - 42mm (com.apple.CoreSimulator.SimDeviceType.Apple-Watch-42mm)
Apple Watch Series 2 - 38mm (com.apple.CoreSimulator.SimDeviceType.Apple-Watch-Series-2-38mm)
Apple Watch Series 2 - 42mm (com.apple.CoreSimulator.SimDeviceType.Apple-Watch-Series-2-42mm)
== Runtimes ==
iOS 8.1 (8.1 - 12B411) (com.apple.CoreSimulator.SimRuntime.iOS-8-1)
iOS 9.3 (9.3 - 13E233) (com.apple.CoreSimulator.SimRuntime.iOS-9-3)
iOS 10.1 (10.1 - 14B72) (com.apple.CoreSimulator.SimRuntime.iOS-10-1)
iOS 10.2 (10.2 - 14C89) (com.apple.CoreSimulator.SimRuntime.iOS-10-2)
tvOS 10.1 (10.1 - 14U591) (com.apple.CoreSimulator.SimRuntime.tvOS-10-1)
watchOS 3.1 (3.1 - 14S471a) (com.apple.CoreSimulator.SimRuntime.watchOS-3-1)
== Devices ==
-- iOS 8.1 --
-- iOS 9.3 --
-- iOS 10.1 --
    iPhone 5 (BF890074-6F3B-4947-863B-719389EE130A) (Shutdown)
    iPhone 5s (6578E45A-DBE6-4132-BB9D-EFB9E46F3DB1) (Shutdown)
    iPhone 6 (FE7403F6-2DC4-437C-8F22-36ED833F43D5) (Shutdown)
    iPhone 6 Plus (8F988216-8CC9-4C34-B134-F91A0D243107) (Shutdown)
    iPhone 6s (1AE8A4B3-72D0-4CC7-92AC-DA499C035E2D) (Shutdown)
    iPhone 6s Plus (8255740B-6952-4802-93FF-52B1F0D34BBF) (Shutdown)
    iPhone 7 (C3615E82-9864-4C62-9071-4D6F3A950158) (Creating)
    iPhone 7 Plus (92F985EF-AFA7-4A6A-AE59-898DAA787214) (Creating)
    iPhone SE (C78E0A01-B3EC-43B2-B219-C76BFA730142) (Shutdown)
    iPad Retina (47997B81-9B9B-4EF5-A02C-5A3146965952) (Shutdown)
    iPad Air (ED4281C1-DC83-449A-B740-24DBC3FAE7F9) (Shutdown)
    iPad Air 2 (466D9824-416A-414F-8F17-CFDEB124AB1F) (Shutdown)
    iPad Pro (9.7 inch) (A2ED8637-F295-4322-9EAD-6F55ACD209AB) (Shutdown)
    iPad Pro (12.9 inch) (8423BCC0-203A-4076-B1D8-BFAA2601A722) (Shutdown)
-- iOS 10.2 --
    iPhone 5 (F9CCC75E-1CF8-4FD6-8973-0696134858EA) (Shutdown)
    iPhone 5s (4E15C24C-DFC9-4C4A-AA4D-2132AF6B42D7) (Shutdown)
    iPhone 6 (F428080E-B6DA-493B-8101-F372FA537E2C) (Shutdown)
    iPhone 6 Plus (3D6E672E-5D89-4B8E-A83C-7D9D30831AD2) (Shutdown)
    iPhone 6s (5656804A-BB52-4453-AAC9-7B439D744B52) (Shutdown)
    iPhone 6s Plus (3924585A-F638-46BC-9F93-99354C058D19) (Shutdown)
    iPhone 7 (DBAEDE5B-34F6-431C-898D-CECA63658E63) (Shutdown)
    iPhone 7 Plus (A0385845-47BC-4424-8D02-D223810F85DA) (Shutdown)
    iPhone SE (E1CD6D89-2C33-4534-A340-47E0E9168970) (Shutdown)
    iPad Retina (F7D8A427-7594-4A24-9CE9-09BE38090C23) (Shutdown)
    iPad Air (7C37463A-5587-4BC5-BBA6-F32B52252BAB) (Shutdown)
    iPad Air 2 (87D6D52A-A522-4D93-B13F-EFED1BBDD9DF) (Shutdown)
    iPad Pro (9.7 inch) (B38E2FB8-6375-4020-85C2-49FF4C6DC9A6) (Shutdown)
    iPad Pro (12.9 inch) (49288CB3-EA62-4472-8BF6-2F81C4E2F62E) (Shutdown)
-- tvOS 10.1 --
    Apple TV 1080p (FBB62097-8171-4A99-95B2-91A90A039BEC) (Shutdown)
-- watchOS 3.1 --
    Apple Watch - 38mm (DD9BAB53-23EB-4CA4-B042-DC2EF2C84B7D) (Shutdown)
    Apple Watch - 42mm (6DC5251E-C2F3-4D89-AD33-DB189AFDEDD7) (Shutdown)
    Apple Watch Series 2 - 38mm (D7D06B47-C8D6-4072-9F3C-BB5058826286) (Shutdown)
    Apple Watch Series 2 - 42mm (F47234F6-D237-4BFC-A337-4D88B63ADE77) (Shutdown)
-- Unavailable: com.apple.CoreSimulator.SimRuntime.tvOS-10-0 --
    Apple TV 1080p (C246D04C-3C24-4FFC-AB37-34F471A6265A) (Shutdown) (unavailable, runtime profile not found)
== Device Pairs ==
604413C4-EE63-4A88-B7E4-91018674FE8A (active, disconnected)
    Watch: Apple Watch Series 2 - 38mm (D7D06B47-C8D6-4072-9F3C-BB5058826286) (Shutdown)
    Phone: iPhone 7 (DBAEDE5B-34F6-431C-898D-CECA63658E63) (Shutdown)
0CE457A5-3458-4545-9503-F2A1584169C6 (active, disconnected)
    Watch: Apple Watch Series 2 - 42mm (F47234F6-D237-4BFC-A337-4D88B63ADE77) (Shutdown)
    Phone: iPhone 7 Plus (A0385845-47BC-4424-8D02-D223810F85DA) (Shutdown)`
