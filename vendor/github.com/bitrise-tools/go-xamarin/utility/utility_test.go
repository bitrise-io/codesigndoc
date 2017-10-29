package utility

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToConfig(t *testing.T) {
	t.Log("creates config from configuration and platform")
	{
		configuration := "Release"
		platform := "iPhone"
		config := ToConfig(configuration, platform)
		require.Equal(t, "Release|iPhone", config)
	}

	t.Log("creates config from configuration")
	{
		configuration := "Release"
		platform := ""
		config := ToConfig(configuration, platform)
		require.Equal(t, "Release|", config)
	}

	t.Log("creates config from platform")
	{
		configuration := ""
		platform := "iPhone"
		config := ToConfig(configuration, platform)
		require.Equal(t, "|iPhone", config)
	}

	t.Log("creates empty config")
	{
		configuration := ""
		platform := ""
		config := ToConfig(configuration, platform)
		require.Equal(t, "|", config)
	}
}

func TestFixWindowsPath(t *testing.T) {
	t.Log("fixes absolute windows path")
	{
		pth := `\bin\iPhoneSimulator\Debug`
		require.Equal(t, "/bin/iPhoneSimulator/Debug", FixWindowsPath(pth))
	}

	t.Log("fixes relative windows path")
	{
		pth := `bin\iPhoneSimulator\Debug`
		require.Equal(t, "bin/iPhoneSimulator/Debug", FixWindowsPath(pth))
	}

	t.Log("fixes relative windows path")
	{
		pth := `..\CreditCardValidator\CreditCardValidator.csproj`
		require.Equal(t, "../CreditCardValidator/CreditCardValidator.csproj", FixWindowsPath(pth))
	}

	t.Log("do not modify absolute unix path")
	{
		pth := "/bin/iPhoneSimulator/Debug"
		require.Equal(t, "/bin/iPhoneSimulator/Debug", FixWindowsPath(pth))
	}

	t.Log("do not modify relative unix path")
	{
		pth := "bin/iPhoneSimulator/Debug"
		require.Equal(t, "bin/iPhoneSimulator/Debug", FixWindowsPath(pth))
	}

	t.Log("do not modify relative unix path")
	{
		pth := "../CreditCardValidator/CreditCardValidator.csproj"
		require.Equal(t, "../CreditCardValidator/CreditCardValidator.csproj", FixWindowsPath(pth))
	}
}

func TestSplitAndStripList(t *testing.T) {
	t.Log("splits string list")
	{
		list := "{FEACFBD2-3405-455C-9665-78FE426C6842};{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}"
		split := SplitAndStripList(list, ";")
		require.Equal(t, 2, len(split))
		require.Equal(t, "{FEACFBD2-3405-455C-9665-78FE426C6842}", split[0])
		require.Equal(t, "{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}", split[1])
	}

	t.Log("splits string list")
	{
		list := "ARMv7, ARM64"
		split := SplitAndStripList(list, ",")
		require.Equal(t, 2, len(split))
		require.Equal(t, "ARMv7", split[0])
		require.Equal(t, "ARM64", split[1])
	}

	t.Log("do not split unless proper separator")
	{
		list := "ARMv7, ARM64"
		split := SplitAndStripList(list, ";")
		require.Equal(t, 1, len(split))
		require.Equal(t, "ARMv7, ARM64", split[0])
	}
}
