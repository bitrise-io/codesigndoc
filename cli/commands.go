package cli

import "github.com/urfave/cli"

const (
	// LogLevelEnvKey ...
	LogLevelEnvKey = "LOGLEVEL"
	// LogLevelKey ...
	LogLevelKey      = "loglevel"
	logLevelKeyShort = "l"

	// HelpKey ...
	HelpKey      = "help"
	helpKeyShort = "h"

	// VersionKey ...
	VersionKey      = "version"
	versionKeyShort = "v"

	// FileParamKey ...
	FileParamKey = "file"
	// SchemeParamKey ...
	SchemeParamKey = "scheme"

	// AllowExportParamKey ...
	AllowExportParamKey = "allow-export"

	// AskPassParamKey ...
	AskPassParamKey = "ask-pass"
)

var (
	commands = []cli.Command{
		{
			Name:   "scan",
			Usage:  "Scan",
			Action: scan,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  FileParamKey,
					Value: "",
					Usage: "Xcode Project/Workspace file",
				},
				cli.StringFlag{
					Name:  SchemeParamKey,
					Value: "",
					Usage: "Xcode Scheme",
				},
				cli.BoolFlag{
					Name:  AllowExportParamKey,
					Usage: "Automatically allow export of discovered files",
				},
				cli.BoolFlag{
					Name:  AskPassParamKey,
					Usage: "Ask for .p12 password, instead of using an empty password",
				},
			},
		},
	}

	appFlags = []cli.Flag{
		cli.StringFlag{
			Name:   LogLevelKey + ", " + logLevelKeyShort,
			Value:  "info",
			Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
			EnvVar: LogLevelEnvKey,
		},
	}
)

func init() {
	// Override default help and version flags
	cli.HelpFlag = cli.BoolFlag{
		Name:  HelpKey + ", " + helpKeyShort,
		Usage: "Show help.",
	}

	cli.VersionFlag = cli.BoolFlag{
		Name:  VersionKey + ", " + versionKeyShort,
		Usage: "Print the version.",
	}
}
