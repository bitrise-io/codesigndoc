package xcodebuild

import "github.com/bitrise-io/go-utils/command"

const (
	toolName = "xcodebuild"
)

// CommandModel ...
type CommandModel interface {
	PrintableCmd() string
	Command() *command.Model
}
