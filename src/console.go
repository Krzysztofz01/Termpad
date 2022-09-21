package main

// Contract abstraction for the underlying console API
type Console interface {
	// Set a given character at given console position
	InsertCharacter(xIndex int, yIndex int, char rune) error

	// Set a given character at given console position with sepecified style attributes
	InsertCharacterWithStyle(xIndex int, yIndex int, char rune, characterStyle CharacterStyle) error

	// Remove a given character at given console position
	RemoveCharacter(xIndex int, yIndex int) error

	// Set the console cursor at given console position
	SetCursorPosition(xIndex int, yIndex int) error

	// Apply pending changes to render on the console
	Commit() error

	// Clears the whole console
	Clear() error

	// Returns an event related to behavior or interaction with the console
	WatchConsoleEvent() interface{}

	// Return the x (width) of the console
	GetWidth() int

	// Return the y (height) of the console
	GetHeight() int

	// Return the x, y (width, height) of the console
	GetSize() (int, int)

	// Set the cursor style provided by the console
	SetCursorStyle(cursorStyle CursorStyle) error

	// Finalize the screen and release resources
	Dispose() error
}

// Structure representing the key press console event
type ConsoleEventKeyPress struct {
	Char     rune
	Key      NamedKey
	Modifier ModifierKey
}

// Structure representing the display/console size change event
type ConsoleEventResize struct {
	Width  int
	Height int
}

// Structure representing the style for a given character to print on the console
// TODO: Add support for console-sepcific colors
type CharacterStyle struct {
	Bold          bool
	Italic        bool
	StrikeThrough bool
	Underline     bool
	Foreground    string
	Background    string
}

// Type representing the console cursor style, that can be provided by the console API implementation
type CursorStyle int16

const (
	NoCursor CursorStyle = iota
	BarCursorStatic
	BarCursorDynamic
	BlockCursorStatic
	BlockCursorDynamic
	LineCursorStatic
	LineCursorDynamic
)

// Type representing named keys, the first one, named printable is an universal representation for ASCII letters
type NamedKey int16

const (
	KeyPrintable NamedKey = iota
	KeyUp
	KeyDown
	KeyRight
	KeyLeft
	KeyPgUp
	KeyPgDn
	KeyHome
	KeyEnd
	KeyInsert
	KeyDelete
	KeyPause
	KeyBacktab
	KeyEnter
	KeyTab
	KeyEscape
	KeyBackspace
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
)

// Type representing modifier keys, the first one, named none indicates that no modifiers were applied
type ModifierKey int16

const (
	ModifierNone ModifierKey = iota
	ModifierShift
	ModifierCtrl
	ModifierAlt
)
