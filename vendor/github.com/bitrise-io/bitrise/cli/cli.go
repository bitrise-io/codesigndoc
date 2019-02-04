package cli

import (
	"errors"
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise/bitrise"
	"github.com/bitrise-io/bitrise/configs"
	"github.com/bitrise-io/bitrise/plugins"
	"github.com/bitrise-io/bitrise/version"
	"github.com/urfave/cli"
)

func initLogFormatter() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "15:04:05",
	})
}

func before(c *cli.Context) error {
	/*
		return err will print app's help also,
		use log.Fatal to avoid print help.
	*/

	initLogFormatter()
	initHelpAndVersionFlags()

	// Debug mode?
	if c.Bool(DebugModeKey) {
		// set for other tools, as an ENV
		if err := os.Setenv(configs.DebugModeEnvKey, "true"); err != nil {
			log.Fatalf("Failed to set DEBUG env, error: %s", err)
		}
		configs.IsDebugMode = true
		log.Warn("=> Started in DEBUG mode")
	}

	// Log level
	// If log level defined - use it
	logLevelStr := c.String(LogLevelKey)
	if logLevelStr == "" && configs.IsDebugMode {
		// if no Log Level defined and we're in Debug Mode - set loglevel to debug
		logLevelStr = "debug"
		log.Warn("=> LogLevel set to debug")
	}
	if logLevelStr == "" {
		// if still empty: set the default
		logLevelStr = "info"
	}

	level, err := log.ParseLevel(logLevelStr)
	if err != nil {
		log.Fatalf("Failed parse log level, error: %s", err)
	}

	if err := os.Setenv(configs.LogLevelEnvKey, level.String()); err != nil {
		log.Fatalf("Failed to set LOGLEVEL env, error: %s", err)
	}
	log.SetLevel(level)

	// CI Mode check
	if c.Bool(CIKey) {
		// if CI mode indicated make sure we set the related env
		//  so all other tools we use will also get it
		if err := os.Setenv(configs.CIModeEnvKey, "true"); err != nil {
			log.Fatalf("Failed to set CI env, error: %s", err)
		}
		configs.IsCIMode = true
	}

	if err := configs.InitPaths(); err != nil {
		log.Fatalf("Failed to initialize required paths, error: %s", err)
	}

	// Pull Request Mode check
	if c.Bool(PRKey) {
		// if PR mode indicated make sure we set the related env
		//  so all other tools we use will also get it
		if err := os.Setenv(configs.PRModeEnvKey, "true"); err != nil {
			log.Fatalf("Failed to set PR env, error: %s", err)
		}
		configs.IsPullRequestMode = true
	}

	pullReqID := os.Getenv(configs.PullRequestIDEnvKey)
	if pullReqID != "" {
		configs.IsPullRequestMode = true
	}

	IsPR := os.Getenv(configs.PRModeEnvKey)
	if IsPR == "true" {
		configs.IsPullRequestMode = true
	}

	return nil
}

func printVersion(c *cli.Context) {
	fmt.Fprintf(c.App.Writer, "%v\n", c.App.Version)
}

// Run ...
func Run() {
	if err := plugins.InitPaths(); err != nil {
		log.Fatalf("Failed to initialize plugin path, error: %s", err)
	}

	initAppHelpTemplate()

	// Parse cl
	cli.VersionPrinter = printVersion

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "Bitrise Automations Workflow Runner"
	app.Version = version.VERSION

	app.Author = ""
	app.Email = ""

	app.Before = before

	app.Flags = flags
	app.Commands = commands

	app.Action = func(c *cli.Context) error {
		pluginName, pluginArgs, isPlugin := plugins.ParseArgs(c.Args())
		if isPlugin {
			plugin, found, err := plugins.LoadPlugin(pluginName)
			if err != nil {
				return fmt.Errorf("Failed to get plugin (%s), error: %s", pluginName, err)
			}
			if !found {
				return fmt.Errorf("Plugin (%s) not installed", pluginName)
			}

			if err := bitrise.RunSetupIfNeeded(version.VERSION, false); err != nil {
				log.Fatalf("Setup failed, error: %s", err)
			}

			if err := plugins.RunPluginByCommand(plugin, pluginArgs); err != nil {
				return fmt.Errorf("Failed to run plugin (%s), error: %s", pluginName, err)
			}
		} else {
			if err := cli.ShowAppHelp(c); err != nil {
				return fmt.Errorf("Failed to show help, error: %s", err)
			}
			return errors.New("")
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
