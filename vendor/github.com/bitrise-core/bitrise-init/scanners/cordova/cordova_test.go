package cordova

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConfigXMLContent(t *testing.T) {
	widget, err := parseConfigXMLContent(testConfigXMLContent)
	require.NoError(t, err)
	require.Equal(t, "http://cordova.apache.org/ns/1.0", widget.XMLNSCDV)
}

const testConfigXMLContent = `<?xml version='1.0' encoding='utf-8'?>
<widget id="com.bitrise.cordovasample" version="0.9.0" xmlns="http://www.w3.org/ns/widgets" xmlns:cdv="http://cordova.apache.org/ns/1.0">
    <name>CordovaOnBitrise</name>
    <description>A sample Apache Cordova application that builds on Bitrise.</description>
    <content src="index.html" />
    <access origin="*" />
    <plugin name="cordova-plugin-whitelist" spec="1" />
    <allow-intent href="http://*/*" />
    <allow-intent href="https://*/*" />
    <allow-intent href="tel:*" />
    <allow-intent href="sms:*" />
    <allow-intent href="mailto:*" />
    <allow-intent href="geo:*" />
    <engine name="ios" spec="~4.3.1" />
    <platform name="android">
        <allow-intent href="market:*" />
    </platform>
    <platform name="ios">
        <allow-intent href="itms:*" />
        <allow-intent href="itms-apps:*" />
    </platform>
    <engine name="android" spec="~6.1.2" />
</widget>`
