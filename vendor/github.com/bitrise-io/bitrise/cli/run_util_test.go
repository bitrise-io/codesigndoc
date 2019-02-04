package cli

import (
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/bitrise-io/bitrise/bitrise"
	"github.com/bitrise-io/bitrise/configs"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pointers"
	"github.com/stretchr/testify/require"
)

func TestIsPRMode(t *testing.T) {
	prModeEnv := os.Getenv(configs.PRModeEnvKey)
	prIDEnv := os.Getenv(configs.PullRequestIDEnvKey)

	// cleanup Envs after these tests
	defer func() {
		require.NoError(t, os.Setenv(configs.PRModeEnvKey, prModeEnv))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, prIDEnv))
	}()

	t.Log("Should be false for: prGlobalFlag: nil, prModeEnv: '', prIDEnv: ''")
	{
		require.NoError(t, os.Setenv(configs.PRModeEnvKey, ""))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, ""))

		pr, err := isPRMode(nil, []envmanModels.EnvironmentItemModel{})
		require.NoError(t, err)
		require.Equal(t, false, pr)
	}

	t.Log("Should be false for: prGlobalFlag: nil, prModeEnv: '', prIDEnv: '', secrets: false")
	{
		inventoryStr := `
envs:
- PR: "false"
- PULL_REQUEST_ID: ""
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.PRModeEnvKey, ""))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, ""))

		pr, err := isPRMode(nil, inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, false, pr)
	}

	t.Log("Should be false for: prGlobalFlag: nil, prModeEnv: 'false', prIDEnv: '', secrets: ''")
	{
		inventoryStr := `
envs:
- PR: ""
- PULL_REQUEST_ID: ""
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.PRModeEnvKey, "false"))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, ""))

		pr, err := isPRMode(nil, inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, false, pr)
	}

	t.Log("Should be false for: prGlobalFlag: false, prModeEnv: 'true', prIDEnv: 'ID', secrets: 'true'")
	{
		inventoryStr := `
envs:
- PR: "true"
- PULL_REQUEST_ID: "ID"
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.PRModeEnvKey, "true"))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, "ID"))

		pr, err := isPRMode(pointers.NewBoolPtr(false), inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, false, pr)
	}

	t.Log("Should be true for: prGlobalFlag: true, prModeEnv: '', prIDEnv: ''")
	{
		require.NoError(t, os.Setenv(configs.PRModeEnvKey, ""))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, ""))

		pr, err := isPRMode(pointers.NewBoolPtr(true), []envmanModels.EnvironmentItemModel{})
		require.NoError(t, err)
		require.Equal(t, true, pr)
	}

	t.Log("Should be true for: prGlobalFlag: true, prModeEnv: '', prIDEnv: '', secrets: false")
	{
		inventoryStr := `
envs:
- PR: "false"
- PULL_REQUEST_ID: ""
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.PRModeEnvKey, ""))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, ""))

		pr, err := isPRMode(pointers.NewBoolPtr(true), inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, true, pr)
	}

	t.Log("Should be true for: prGlobalFlag: nil, prModeEnv: 'true', prIDEnv: '', secrets: false")
	{
		inventoryStr := `
envs:
- PR: "false"
- PULL_REQUEST_ID: ""
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.PRModeEnvKey, "true"))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, ""))

		pr, err := isPRMode(nil, inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, true, pr)
	}

	t.Log("Should be true for: prGlobalFlag: nil, prModeEnv: 'false', prIDEnv: 'some', secrets: false")
	{
		inventoryStr := `
envs:
- PR: "false"
- PULL_REQUEST_ID: ""
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.PRModeEnvKey, "false"))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, "some"))

		pr, err := isPRMode(nil, inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, true, pr)
	}

	t.Log("Should be true for: prGlobalFlag: nil, prModeEnv: '', prIDEnv: '', secrets: true")
	{
		inventoryStr := `
envs:
- PR: "true"
- PULL_REQUEST_ID: ""
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.PRModeEnvKey, ""))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, ""))

		pr, err := isPRMode(nil, inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, true, pr)
	}

	t.Log("Should be true for: prGlobalFlag: nil, prModeEnv: 'false', prIDEnv: '', secrets: true")
	{
		inventoryStr := `
envs:
- PR: ""
- PULL_REQUEST_ID: "some"
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.PRModeEnvKey, "false"))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, ""))

		pr, err := isPRMode(nil, inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, true, pr)
	}

	t.Log("Should be true for: prGlobalFlag: true, prModeEnv: 'false', prIDEnv: '', secrets: false")
	{
		inventoryStr := `
envs:
- PR: "false"
- PULL_REQUEST_ID: ""
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.PRModeEnvKey, "false"))
		require.NoError(t, os.Setenv(configs.PullRequestIDEnvKey, ""))

		pr, err := isPRMode(pointers.NewBoolPtr(true), inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, true, pr)
	}
}

func TestIsCIMode(t *testing.T) {
	ciModeEnv := os.Getenv(configs.CIModeEnvKey)

	defer func() {
		require.NoError(t, os.Setenv(configs.CIModeEnvKey, ciModeEnv))
	}()

	t.Log("Should be false for: ciGlobalFlag: nil, ciModeEnv: 'false'")
	{
		require.NoError(t, os.Setenv(configs.CIModeEnvKey, "false"))

		ci, err := isCIMode(nil, []envmanModels.EnvironmentItemModel{})
		require.NoError(t, err)
		require.Equal(t, false, ci)
	}

	t.Log("Should be false for: ciGlobalFlag: false, ciModeEnv: 'false' secrets: false")
	{
		inventoryStr := `
envs:
- CI: "false"
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.CIModeEnvKey, "false"))

		ci, err := isCIMode(pointers.NewBoolPtr(false), inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, false, ci)
	}

	t.Log("Should be true for: ciGlobalFlag: true, ciModeEnv: 'false'")
	{
		require.NoError(t, os.Setenv(configs.CIModeEnvKey, ""))

		ci, err := isCIMode(pointers.NewBoolPtr(true), []envmanModels.EnvironmentItemModel{})
		require.NoError(t, err)
		require.Equal(t, true, ci)
	}

	t.Log("Should be true for: ciGlobalFlag: true, ciModeEnv: '' secrets: false")
	{
		inventoryStr := `
envs:
- CI: "false"
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.CIModeEnvKey, ""))

		ci, err := isCIMode(pointers.NewBoolPtr(true), inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, true, ci)
	}

	t.Log("Should be true for: ciGlobalFlag: nil, ciModeEnv: 'true' secrets: false")
	{
		inventoryStr := `
envs:
- CI: ""
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.CIModeEnvKey, "true"))

		ci, err := isCIMode(nil, inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, true, ci)
	}

	t.Log("Should be true for: ciGlobalFlag: nil, ciModeEnv: '' secrets: true")
	{
		inventoryStr := `
envs:
- CI: "true"
`
		inventory, err := bitrise.InventoryModelFromYAMLBytes([]byte(inventoryStr))
		require.NoError(t, err)

		require.NoError(t, os.Setenv(configs.CIModeEnvKey, ""))

		ci, err := isCIMode(nil, inventory.Envs)
		require.NoError(t, err)
		require.Equal(t, true, ci)
	}
}

func TestExpandEnvs(t *testing.T) {
	configStr := `
format_version: 1.3.0
default_step_lib_source: "https://github.com/bitrise-io/bitrise-steplib.git"

workflows:
  test:
    envs:
    - ENV0: "Hello"
    - ENV1: "$ENV0 world"
    steps:
    - script:
        inputs:
        - content: |
            #!/bin/bash
            envman add --key ENV2 --value "$ENV1 !"
    - script:
        inputs:
        - content: |
            #!/bin/bash
            echo "ENV2: $ENV2"
            if [ "$ENV2" != "Hello world !" ] ; then
              echo "Actual ($ENV2), excpected (Hello world !)"
              exit 1
            fi
`

	require.NoError(t, configs.InitPaths())

	config, warnings, err := bitrise.ConfigModelFromYAMLBytes([]byte(configStr))
	require.NoError(t, err)
	require.Equal(t, 0, len(warnings))

	buildRunResults, err := runWorkflowWithConfiguration(time.Now(), "test", config, []envmanModels.EnvironmentItemModel{})
	require.NoError(t, err)
	require.Equal(t, 2, len(buildRunResults.SuccessSteps))
	require.Equal(t, 0, len(buildRunResults.FailedSteps))
	require.Equal(t, 0, len(buildRunResults.FailedSkippableSteps))
	require.Equal(t, 0, len(buildRunResults.SkippedSteps))
}

func TestEvaluateInputs(t *testing.T) {
	configStr := `
format_version: 1.3.0
default_step_lib_source: "https://github.com/bitrise-io/bitrise-steplib.git"

workflows:
  test:
    envs:
    - TEST_KEY: "test value"
    steps:
    - script:
        title: "Template test"
        inputs:
        - content: |
            #!/bin/bash
            set -v
            {{if .IsCI}}
            exit 1
            {{else}}
            exit 0
            {{end}}
          opts:
            is_template: true
    - script:
        title: "Template test"
        inputs:
        - content: |
            #!/bin/bash
            set -v
            {{if enveq "TEST_KEY" "test value"}}
            exit 0
            {{else}}
            exit 1
            {{end}}
          opts:
            is_template: true
`

	require.NoError(t, configs.InitPaths())

	config, warnings, err := bitrise.ConfigModelFromYAMLBytes([]byte(configStr))
	require.NoError(t, err)
	require.Equal(t, 0, len(warnings))

	buildRunResults, err := runWorkflowWithConfiguration(time.Now(), "test", config, []envmanModels.EnvironmentItemModel{})
	require.Equal(t, nil, err)
	require.Equal(t, 0, len(buildRunResults.SkippedSteps))
	require.Equal(t, 2, len(buildRunResults.SuccessSteps))
	require.Equal(t, 0, len(buildRunResults.FailedSteps))
	require.Equal(t, 0, len(buildRunResults.FailedSkippableSteps))
}

func TestGetBitriseConfigFromBase64Data(t *testing.T) {
	configStr := `
format_version: 0.9.10
default_step_lib_source: "https://github.com/bitrise-io/bitrise-steplib.git"

workflows:
  target:
    title: target
`
	configBytes := []byte(configStr)
	configBase64Str := base64.StdEncoding.EncodeToString(configBytes)

	config, warnings, err := GetBitriseConfigFromBase64Data(configBase64Str)
	require.NoError(t, err)
	require.Equal(t, 0, len(warnings))

	require.Equal(t, "0.9.10", config.FormatVersion)
	require.Equal(t, "https://github.com/bitrise-io/bitrise-steplib.git", config.DefaultStepLibSource)

	workflow, found := config.Workflows["target"]
	require.Equal(t, true, found)
	require.Equal(t, "target", workflow.Title)
}

func TestGetInventoryFromBase64Data(t *testing.T) {
	inventoryStr := `
envs:
  - MY_HOME: $HOME
    opts:
      is_expand: true
`
	inventoryBytes := []byte(inventoryStr)
	inventoryBase64Str := base64.StdEncoding.EncodeToString(inventoryBytes)

	inventory, err := GetInventoryFromBase64Data(inventoryBase64Str)
	require.NoError(t, err)

	env := inventory[0]

	key, value, err := env.GetKeyValuePair()
	require.NoError(t, err)
	require.Equal(t, "MY_HOME", key)
	require.Equal(t, "$HOME", value)

	opts, err := env.GetOptions()
	require.NoError(t, err)
	require.Equal(t, true, *opts.IsExpand)
}

func TestInvalidStepID(t *testing.T) {
	configStr := `
format_version: 1.3.0
default_step_lib_source: "https://github.com/bitrise-io/bitrise-steplib.git"

workflows:
  target:
    title: Invalid step id
    steps:
    - invalid-step:
    - invalid-step:
    - invalid-step:
`

	require.NoError(t, configs.InitPaths())

	config, warnings, err := bitrise.ConfigModelFromYAMLBytes([]byte(configStr))
	require.NoError(t, err)
	require.Equal(t, 0, len(warnings))

	results, err := runWorkflowWithConfiguration(time.Now(), "target", config, []envmanModels.EnvironmentItemModel{})
	require.Equal(t, 1, len(results.StepmanUpdates))
}
