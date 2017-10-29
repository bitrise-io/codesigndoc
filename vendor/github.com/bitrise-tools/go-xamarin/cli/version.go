package cli

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xamarin/version"
	"github.com/urfave/cli"
)

// VersionOutputModel ...
type VersionOutputModel struct {
	Version     string `json:"version"`
	BuildNumber string `json:"build_number"`
	Commit      string `json:"commit"`
}

const (
	// FormatRaw ...
	FormatRaw = "raw"
	// FormatJSON ...
	FormatJSON = "json"
	// FormatYML ...
	FormatYML = "yml"
)

// Print ...
func print(versionOutput VersionOutputModel, format string) {
	switch format {
	case FormatJSON:
		serBytes, err := json.Marshal(versionOutput)
		if err != nil {
			log.Errorf("failed to print output, error: %s", err)
			return
		}
		fmt.Printf("%s\n", serBytes)
	case FormatYML:
		serBytes, err := yaml.Marshal(versionOutput)
		if err != nil {
			log.Errorf("failed to print output, error: %s", err)
			return
		}
		fmt.Printf("%s\n", serBytes)
	default:
		fmt.Printf("version: %s\nbuild number: %s\ncommit: %s\n", versionOutput.Version, versionOutput.BuildNumber, versionOutput.Commit)
	}
}

func versionCmd(c *cli.Context) error {
	format := c.String("format")

	versionOutput := VersionOutputModel{
		Version:     version.VERSION,
		BuildNumber: version.BuildNumber,
		Commit:      version.Commit,
	}

	print(versionOutput, format)

	return nil
}
