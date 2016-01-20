package cli

import "github.com/codegangsta/cli"

const (
	// WorkdirEnvKey ...
	WorkdirEnvKey = "BITRISE_MACHINE_WORKDIR"
	// WorkdirKey ...
	WorkdirKey = "workdir"

	// LogLevelEnvKey ...
	LogLevelEnvKey = "LOGLEVEL"
	// LogLevelKey ...
	LogLevelKey      = "loglevel"
	logLevelKeyShort = "l"

	// EnvironmentParamKey ...
	EnvironmentParamKey      = "environment"
	environmentParamKeyShort = "e"

	// HelpKey ...
	HelpKey      = "help"
	helpKeyShort = "h"

	// VersionKey ...
	VersionKey      = "version"
	versionKeyShort = "v"

	// --- Command flags

	// TimeoutFlagKey ...
	TimeoutFlagKey = "timeout"
	// AbortCheckURLFlagKey ...
	AbortCheckURLFlagKey = "abort-check-url"
	// LogFormatFlagKey ...
	LogFormatFlagKey = "logformat"
	// ForceFlagKey ...
	ForceFlagKey = "force"
)

var (
	commands = []cli.Command{
		{
			Name:   "scan",
			Usage:  "Scan",
			Action: scan,
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
