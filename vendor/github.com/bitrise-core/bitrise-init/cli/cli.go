package cli

import (
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-core/bitrise-init/version"
	"github.com/urfave/cli"
)

// Run ...
func Run() {
	// Parse cl
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Version)
	}

	app := cli.NewApp()

	app.Name = path.Base(os.Args[0])
	app.Usage = "Bitrise Init Tool"
	app.Version = version.VERSION
	app.Author = ""
	app.Email = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "loglevel, l",
			Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
			EnvVar: "LOGLEVEL",
		},
		cli.BoolFlag{
			Name:   "ci",
			Usage:  "If true it indicates that we're used by another tool so don't require any user input!",
			EnvVar: "CI",
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			ForceColors:     true,
			TimestampFormat: "15:04:05",
		})

		// Log level
		logLevelStr := c.String("loglevel")
		if logLevelStr == "" {
			logLevelStr = "info"
		}

		level, err := log.ParseLevel(logLevelStr)
		if err != nil {
			return err
		}
		log.SetLevel(level)

		return nil
	}

	app.Commands = []cli.Command{
		versionCommand,
		configCommand,
		manualConfigCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
