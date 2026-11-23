package frame

// MenuEntry describes a command's trigger words and how it appears in a menu.
type MenuEntry struct {
	// Triggers lists the words that invoke this entry.
	Triggers []string `json:"triggers"`
	// ArgumentHint describes the arguments shown after the trigger in help text.
	ArgumentHint string `json:"argument_hint,omitempty"`
	// Usage describes the command in the menu.
	Usage string `json:"usage,omitempty"`
}

// GameMenuEntry binds a menu command to an in-game chat handler.
type GameMenuEntry struct {
	MenuEntry
	// OnTrigger receives the chat that triggered the entry, with the command arguments in Msg.
	OnTrigger func(chat *GameChat)
}

// BackendMenuEntry binds a menu command to a terminal handler.
type BackendMenuEntry struct {
	MenuEntry
	// OnTrigger receives the arguments following the terminal trigger.
	OnTrigger func(cmds []string)
}

// GameMenuSetter registers commands in the in-game menu.
type GameMenuSetter interface {
	// AddGameMenuEntry registers an in-game command and its chat handler.
	AddGameMenuEntry(*GameMenuEntry)
}

// Printer writes to a terminal, log file or configured combination of the two.
type Printer interface {
	// Write writes a line of text.
	Write(string)
	// Print writes the arguments as formatted by fmt.Sprint.
	Print(...any)
	// Printf writes a formatted message without adding a newline.
	Printf(format string, a ...any)
	// Println writes the arguments as a line, separated by spaces.
	Println(...any)
	// Printfln writes a formatted message followed by a newline.
	Printfln(format string, a ...any)
}

// MultiOutDst groups the output destinations for a component. Use Terminal or Log for a specific
// destination, or a severity such as Warning for configured routing.
type MultiOutDst struct {
	Printer
	// Terminal writes to the terminal only.
	Terminal Printer
	// Log writes to the log only.
	Log Printer
	// TerminalAndLog writes to both the terminal and log.
	TerminalAndLog Printer
	// Debug writes debug-level output to configured destinations.
	Debug Printer
	// Info writes informational output to configured destinations.
	Info Printer
	// Success writes success-level output to configured destinations.
	Success Printer
	// Warning writes warning-level output to configured destinations.
	Warning Printer
	// Error writes error-level output to configured destinations.
	Error Printer
}

// BackendIO provides the terminal menu, input and component output.
type BackendIO interface {
	// AddBackendMenuEntry registers a terminal command and its argument handler.
	AddBackendMenuEntry(*BackendMenuEntry)
	// GetTerminalInput returns an operation that waits for the next terminal input line.
	GetTerminalInput() AsyncResult[string]
	// Out returns the component's output destinations.
	Out() *MultiOutDst
	// ForkOut returns output destinations with the given prefix and log file.
	ForkOut(prefix string, logFile string) *MultiOutDst
}
