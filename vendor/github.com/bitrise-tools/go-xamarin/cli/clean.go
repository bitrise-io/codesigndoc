package cli

import (
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xamarin/analyzers/project"
	"github.com/bitrise-tools/go-xamarin/builder"
	"github.com/bitrise-tools/go-xamarin/tools/buildtools"
	"github.com/urfave/cli"
)

func cleanCmd(c *cli.Context) error {
	solutionPth := c.String(solutionFilePathKey)

	fmt.Println("")
	log.Infof("Config:")
	log.Printf("- solution: %s", solutionPth)
	fmt.Println("")

	if solutionPth == "" {
		return fmt.Errorf("missing required input: %s", solutionFilePathKey)
	}

	buildTool := buildtools.Xbuild

	builder, err := builder.New(solutionPth, nil, buildTool)
	if err != nil {
		return err
	}

	callback := func(project project.Model, dir string) {
		log.Printf("  cleaning project: %s (removing: %s)", project.Name, dir)
	}

	log.Infof("Cleaning solution: %s", solutionPth)
	if err := builder.CleanAll(callback); err != nil {
		return err
	}

	return nil
}
