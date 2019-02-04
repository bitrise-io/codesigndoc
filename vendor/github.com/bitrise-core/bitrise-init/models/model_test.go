package models

import (
	"fmt"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/require"
)

func TestNewOption(t *testing.T) {
	actual := NewOption("Project (or Workspace) path", "BITRISE_PROJECT_PATH")
	expected := &OptionNode{
		Title:          "Project (or Workspace) path",
		EnvKey:         "BITRISE_PROJECT_PATH",
		ChildOptionMap: map[string]*OptionNode{},
		Components:     []string{},
	}

	require.Equal(t, expected, actual)
}

func TestGetValues(t *testing.T) {
	option := OptionNode{
		ChildOptionMap: map[string]*OptionNode{},
	}
	option.ChildOptionMap["assembleAndroidTest"] = &OptionNode{}
	option.ChildOptionMap["assembleDebug"] = &OptionNode{}
	option.ChildOptionMap["assembleRelease"] = &OptionNode{}

	values := option.GetValues()

	expectedMap := map[string]bool{
		"assembleAndroidTest": false,
		"assembleDebug":       false,
		"assembleRelease":     false,
	}

	for _, value := range values {
		delete(expectedMap, value)
	}

	require.Equal(t, 0, len(expectedMap))
}

func TestLastOptions(t *testing.T) {
	// 1. level
	opt0 := NewOption("OPT0", "OPT0_KEY")

	// 2. level
	opt01 := NewOption("OPT01", "OPT01_KEY") // has no child
	opt01.AddConfig("test", nil)
	opt0.AddOption("value1", opt01)

	opt02 := NewOption("OPT02", "OPT02_KEY")
	opt0.AddOption("value2", opt02)

	// 3. level
	opt021 := NewOption("OPT021", "OPT021_KEY")
	opt02.AddOption("value1", opt021)

	// 4. level
	opt0211 := NewOption("OPT0211", "OPT0211_KEY") // has no child
	opt021.AddOption("value1", opt0211)

	opt0212 := NewOption("OPT0212", "OPT0212_KEY")
	opt021.AddOption("value2", opt0212)

	// 5. level
	opt02121 := NewOption("OPT02121", "OPT02121_KEY") // has no child
	opt0212.AddOption("value1", opt02121)

	lastOptions := opt0.LastChilds()
	require.Equal(t, true, len(lastOptions) == 3, fmt.Sprintf("%d", len(lastOptions)))

	optionsMap := map[string]bool{}
	for _, opt := range lastOptions {
		optionsMap[opt.Title] = true
	}

	require.Equal(t, true, optionsMap["OPT01"])
	require.Equal(t, true, optionsMap["OPT0211"])
	require.Equal(t, true, optionsMap["OPT02121"])

	{
		optionJSON := `{
	"title": "Project (or Workspace) path",
	"env_key": "BITRISE_PROJECT_PATH",
	"value_map": {
		"BitriseTest.xcodeproj": {
			"title": "Scheme name",
			"env_key": "BITRISE_SCHEME",
			"value_map": {
				"BitriseTest": {
					"config": "ios-test-config"
				},
				"BitriseTest-tvOS": {
					"config": "ios-test-config"
				}
			}
		}
	}
}`

		var option OptionNode
		require.NoError(t, json.Unmarshal([]byte(optionJSON), &option))

		lastOptions := option.LastChilds()
		optionsMap := map[string]bool{}
		for _, opt := range lastOptions {
			optionsMap[opt.String()] = true
		}

		require.Equal(t, true, optionsMap[`{
	"title": "Scheme name",
	"env_key": "BITRISE_SCHEME",
	"value_map": {
		"BitriseTest": {
			"config": "ios-test-config"
		},
		"BitriseTest-tvOS": {
			"config": "ios-test-config"
		}
	}
}`])
	}

	{
		optionJSON := `{
	"title": "Gradlew file path",
	"env_key": "GRADLEW_PATH",
	"value_map": {
		"$HOME/Develop/react/AwesomeProject/android/gradlew": {
			"title": "Path to the gradle file to use",
			"env_key": "GRADLE_BUILD_FILE_PATH",
			"value_map": {
				"$HOME/Develop/react/AwesomeProject/android/build.gradle": {
					"config": "android-config"
				}
			}
		}
	}
}`

		var option OptionNode
		require.NoError(t, json.Unmarshal([]byte(optionJSON), &option))

		lastOptions := option.LastChilds()
		optionsMap := map[string]bool{}
		for _, opt := range lastOptions {
			optionsMap[opt.String()] = true
		}

		require.Equal(t, true, optionsMap[`{
	"title": "Path to the gradle file to use",
	"env_key": "GRADLE_BUILD_FILE_PATH",
	"value_map": {
		"$HOME/Develop/react/AwesomeProject/android/build.gradle": {
			"config": "android-config"
		}
	}
}`])
	}

	{
		optionJSON := `{
	"title": "project_dir",
	"env_key": "PROJECT_DIR",
	"value_map": {
		"$HOME/Develop/react/AwesomeProject": {
			"title": "Gradlew file path",
			"env_key": "GRADLEW_PATH",
			"value_map": {
				"$HOME/Develop/react/AwesomeProject/android/gradlew": {
					"title": "Path to the gradle file to use",
					"env_key": "GRADLE_BUILD_FILE_PATH",
					"value_map": {
						"$HOME/Develop/react/AwesomeProject/android/build.gradle": {}
					}
				}
			}
		}
	}
}`

		var option OptionNode
		require.NoError(t, json.Unmarshal([]byte(optionJSON), &option))

		lastOptions := option.LastChilds()
		optionsMap := map[string]bool{}
		for _, opt := range lastOptions {
			optionsMap[opt.String()] = true
		}

		require.Equal(t, true, optionsMap[`{
	"title": "Path to the gradle file to use",
	"env_key": "GRADLE_BUILD_FILE_PATH",
	"value_map": {
		"$HOME/Develop/react/AwesomeProject/android/build.gradle": {}
	}
}`])
	}
}

func TestCopy(t *testing.T) {
	// 1. level
	opt0 := NewOption("OPT0", "OPT0_KEY")

	// 2. level
	opt01 := NewOption("OPT01", "OPT01_KEY")
	opt01.AddOption("value01", nil)

	opt0.AddOption("value1", opt01)

	opt02 := NewConfigOption("name")
	opt0.AddConfig("value2", opt02)

	// make a copy
	opt0Copy := opt0.Copy()

	// Ensure copy is the same
	require.Equal(t, opt0.Title, opt0Copy.Title)
	require.Equal(t, opt0.EnvKey, opt0Copy.EnvKey)

	opt01Copy := opt0Copy.ChildOptionMap["value1"]
	require.Equal(t, opt01.Title, opt01Copy.Title)
	require.Equal(t, opt01.EnvKey, opt01Copy.EnvKey)
	require.Equal(t, 1, len(opt01Copy.ChildOptionMap))
	_, ok := opt01Copy.ChildOptionMap["value01"]
	require.Equal(t, true, ok)
	require.Equal(t, "", opt01Copy.Config)

	opt02Copy := opt0Copy.ChildOptionMap["value2"]
	require.Equal(t, opt02.Title, opt02Copy.Title)
	require.Equal(t, opt02.EnvKey, opt02Copy.EnvKey)
	require.Equal(t, 0, len(opt02Copy.ChildOptionMap))
	require.Equal(t, "name", opt02Copy.Config)

	// Ensure copy is a new object
	opt0Copy.Title = "OPT0_COPY"
	require.Equal(t, "OPT0", opt0.Title)

	opt01Copy.Title = "OPT01_COPY"
	require.Equal(t, "OPT01", opt01.Title)

	opt02Copy.Config = "name_copy"
	require.Equal(t, "name", opt02.Config)
}

func TestComponents(t *testing.T) {
	// 1. level
	opt0 := NewOption("OPT0", "OPT0_KEY")

	// 2. level
	opt01 := NewOption("OPT01", "OPT01_KEY") // has no child
	opt0.AddOption("value1", opt01)

	opt02 := NewOption("OPT02", "OPT02_KEY")
	opt0.AddOption("value2", opt02)

	// 3. level
	opt021 := NewOption("OPT021", "OPT021_KEY")
	opt02.AddOption("value1", opt021)

	// 4. level
	opt0211 := NewOption("OPT0211", "OPT0211_KEY") // has no child
	opt021.AddOption("value1", opt0211)

	opt0212 := NewOption("OPT0212", "OPT0212_KEY")
	opt021.AddOption("value2", opt0212)

	// 5. level
	opt02121 := NewOption("OPT02121", "OPT02121_KEY") // has no child
	opt0212.AddOption("value1", opt02121)

	require.Equal(t, []string{}, opt0.Components)
	require.Equal(t, []string{"value1"}, opt01.Components)
	require.Equal(t, []string{"value2"}, opt02.Components)
	require.Equal(t, []string{"value2", "value1"}, opt021.Components)
	require.Equal(t, []string{"value2", "value1", "value1"}, opt0211.Components)
	require.Equal(t, []string{"value2", "value1", "value2"}, opt0212.Components)
	require.Equal(t, []string{"value2", "value1", "value2", "value1"}, opt02121.Components)
}

func TestHead(t *testing.T) {
	// 1. level
	opt0 := NewOption("OPT0", "OPT0_KEY")

	// 2. level
	opt01 := NewOption("OPT01", "OPT01_KEY") // has no child
	opt0.AddOption("value1", opt01)

	opt02 := NewOption("OPT02", "OPT02_KEY")
	opt0.AddOption("value2", opt02)

	// 3. level
	opt021 := NewOption("OPT021", "OPT021_KEY")
	opt02.AddOption("value1", opt021)

	require.Equal(t, (*OptionNode)(nil), opt0.Head)
	require.Equal(t, opt0, opt01.Head)
	require.Equal(t, opt0, opt02.Head)
	require.Equal(t, opt0, opt021.Head)
}

func TestParent(t *testing.T) {
	// 1. level
	opt0 := NewOption("OPT0", "OPT0_KEY")

	// 2. level
	opt01 := NewOption("OPT01", "OPT01_KEY") // has no child
	opt0.AddOption("value1", opt01)

	opt02 := NewOption("OPT02", "OPT02_KEY")
	opt0.AddOption("value2", opt02)

	// 3. level
	opt021 := NewOption("OPT021", "OPT021_KEY")
	opt02.AddOption("value1", opt021)

	{
		parent, underKey, ok := opt0.Parent()
		require.Equal(t, (*OptionNode)(nil), parent)
		require.Equal(t, "", underKey)
		require.Equal(t, false, ok)
	}

	{
		parent, underKey, ok := opt01.Parent()
		require.Equal(t, opt0, parent)
		require.Equal(t, "value1", underKey)
		require.Equal(t, true, ok)
	}

	{
		parent, underKey, ok := opt02.Parent()
		require.Equal(t, opt0, parent)
		require.Equal(t, "value2", underKey)
		require.Equal(t, true, ok)
	}

	{
		parent, underKey, ok := opt021.Parent()
		require.Equal(t, opt02, parent)
		require.Equal(t, "value1", underKey)
		require.Equal(t, true, ok)
	}
}

func TestRemoveConfigs(t *testing.T) {
	optionJSON := `{
	"title": "Project (or Workspace) path",
	"env_key": "BITRISE_PROJECT_PATH",
	"value_map": {
		"BitriseTest.xcodeproj": {
			"title": "Scheme name",
			"env_key": "BITRISE_SCHEME",
			"value_map": {
				"BitriseTest": {
					"config": "ios-test-config"
				},
				"BitriseTest-tvOS": {
					"config": "ios-test-config"
				}
			}
		}
	}
}`

	var option OptionNode
	require.NoError(t, json.Unmarshal([]byte(optionJSON), &option))

	option.RemoveConfigs()

	require.Equal(t, `{
	"title": "Project (or Workspace) path",
	"env_key": "BITRISE_PROJECT_PATH",
	"value_map": {
		"BitriseTest.xcodeproj": {
			"title": "Scheme name",
			"env_key": "BITRISE_SCHEME",
			"value_map": {
				"BitriseTest": {},
				"BitriseTest-tvOS": {}
			}
		}
	}
}`, option.String())
}
