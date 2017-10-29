package cli

import (
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xamarin/version"
	"github.com/urfave/cli"
)

// Run ...
func Run() {
	app := cli.NewApp()
	app.Name = "xamarin-builder"
	app.Usage = "Build xamarin projects"
	app.Version = version.VERSION

	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Errorf("Finished with error: %s", err)
		os.Exit(1)
	}
}
