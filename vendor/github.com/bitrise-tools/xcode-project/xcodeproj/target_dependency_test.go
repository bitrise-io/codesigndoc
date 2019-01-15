package xcodeproj

import (
	"testing"

	"github.com/bitrise-tools/xcode-project/pretty"
	"github.com/bitrise-tools/xcode-project/serialized"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

func TestParseTargetDependency(t *testing.T) {
	var raw serialized.Object
	_, err := plist.Unmarshal([]byte(rawTargetDependency), &raw)
	require.NoError(t, err)

	targetDependency, err := parseTargetDependency("13E76E511F4AC94F0028096E", raw)
	require.NoError(t, err)
	// fmt.Printf("targetDependency:\n%s\n", pretty.Object(targetDependency))
	require.Equal(t, expectedTargetDependency, pretty.Object(targetDependency))
}

const rawTargetDependency = `
{
	13E76E511F4AC94F0028096E /* PBXTargetDependency */ = {
		isa = PBXTargetDependency;
		target = 13E76E461F4AC94F0028096E /* share-extension */;
		targetProxy = 13E76E501F4AC94F0028096E /* PBXContainerItemProxy */;
	};

	13E76E461F4AC94F0028096E /* share-extension */ = {
		isa = PBXNativeTarget;
		buildConfigurationList = 13E76E3A1F4AC90A0028096E /* Build configuration list for PBXNativeTarget "share-extension" */;
		buildPhases = (
			13E76E431F4AC94F0028096E /* Sources */,
			13E76E441F4AC94F0028096E /* Frameworks */,
			13E76E451F4AC94F0028096E /* Resources */,
		);
		buildRules = (
		);
		dependencies = (
		);
		name = "share-extension";
		productName = "share-extension";
		productReference = 13E76E471F4AC94F0028096E /* share-extension.appex */;
		productType = "com.apple.product-type.app-extension";
	};

	13E76E3A1F4AC90A0028096E /* Build configuration list for PBXNativeTarget "code-sign-test" */ = {
		isa = XCConfigurationList;
		buildConfigurations = (
		);
		defaultConfigurationIsVisible = 0;
		defaultConfigurationName = Release;
	};

	13E76E0E1F4AC90A0028096E /* code-sign-test.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = "code-sign-test.app"; sourceTree = BUILT_PRODUCTS_DIR; };
	13E76E471F4AC94F0028096E /* share-extension.appex */ = {isa = PBXFileReference; explicitFileType = "wrapper.app-extension"; includeInIndex = 0; path = "share-extension.appex"; sourceTree = BUILT_PRODUCTS_DIR; };
}`

const expectedTargetDependency = `{
	"ID": "13E76E511F4AC94F0028096E",
	"Target": {
		"Type": "PBXNativeTarget",
		"ID": "13E76E461F4AC94F0028096E",
		"Name": "share-extension",
		"BuildConfigurationList": {
			"ID": "13E76E3A1F4AC90A0028096E",
			"DefaultConfigurationName": "Release",
			"BuildConfigurations": null
		},
		"Dependencies": null,
		"ProductReference": {
			"Path": "share-extension.appex"
		},
		"ProductType": "com.apple.product-type.app-extension"
	}
}`
