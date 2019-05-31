package tools

import "io"

// Runnable ...
type Runnable interface {
	String() string
	SetCustomOptions(options ...string)
	Run(outWriter, errWriter io.Writer) error
}

// Printable ...
type Printable interface {
	String() string
}

// Editable ...
type Editable interface {
	SetCustomOptions(options ...string)
}

//
// EmptyCommand - for return type in case of failed to create a RunnableCommand
type EmptyCommand struct{}

// String ...
func (cmd *EmptyCommand) String() string { return "" }

// SetCustomOptions ...
func (cmd *EmptyCommand) SetCustomOptions(options ...string) {}

// Run ...
func (cmd *EmptyCommand) Run() error { return nil }

// ---

// PrintableSliceContains ...
func PrintableSliceContains(cmdSlice []Printable, cmd Printable) bool {
	for _, c := range cmdSlice {
		if c.String() == cmd.String() {
			return true
		}
	}
	return false
}
