package tools

// Runnable ...
type Runnable interface {
	PrintableCommand() string
	SetCustomOptions(options ...string)
	Run() error
}

// Printable ...
type Printable interface {
	PrintableCommand() string
}

// Editable ...
type Editable interface {
	SetCustomOptions(options ...string)
}

//
// EmptyCommand - for return type in case of failed to create a RunnableCommand
type EmptyCommand struct{}

// PrintableCommand ...
func (cmd *EmptyCommand) PrintableCommand() string { return "" }

// SetCustomOptions ...
func (cmd *EmptyCommand) SetCustomOptions(options ...string) {}

// Run ...
func (cmd *EmptyCommand) Run() error { return nil }

// ---

// PrintableSliceContains ...
func PrintableSliceContains(cmdSlice []Printable, cmd Printable) bool {
	for _, c := range cmdSlice {
		if c.PrintableCommand() == cmd.PrintableCommand() {
			return true
		}
	}
	return false
}
