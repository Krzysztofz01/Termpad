package main

import "github.com/gdamore/tcell"

// Contract abstraction for the underlying console API
type Console interface {
	// Set a given character at given console position
	InsertCharacter(xIndex int, yIndex int, char rune) error

	// Remove a given character at given console position
	RemoveCharacter(xIndex int, yIndex int) error

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

	// Finalize the screen and release resources
	Dispose() error
}

// Structure representing the key press console event
type ConsoleEventKeyPress struct {
	Char rune

	// TODO: Implement internal key structure to avoid contract layer third-party dependencies
	Key      tcell.Key
	Modifier tcell.ModMask
}

type ConsoleEventResize struct {
	Width  int
	Height int
}
