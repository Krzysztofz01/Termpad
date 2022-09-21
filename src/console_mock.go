package main

const (
	MockConsoleWidth  = 10
	MockConsoleHeight = 10
)

// Structure implementing the console contract, used as a mockup for testing purposes
type ConsoleMock struct {
}

// Create a new instance of the Tcell based console. The size of the console is 10x10
func CreateConsoleMockup() Console {
	return &ConsoleMock{}
}

func (console *ConsoleMock) InsertCharacter(xIndex int, yIndex int, char rune) error {
	return nil
}

func (console *ConsoleMock) InsertCharacterWithStyle(xIndex int, yIndex int, char rune, characterStyle CharacterStyle) error {
	return nil
}

func (console *ConsoleMock) RemoveCharacter(xIndex int, yIndex int) error {
	return nil
}

func (console *ConsoleMock) SetCursorPosition(xIndex int, yIndex int) error {
	return nil
}

func (console *ConsoleMock) Commit() error {
	return nil
}

func (console *ConsoleMock) Clear() error {
	return nil
}

func (console *ConsoleMock) WatchConsoleEvent() interface{} {
	return nil
}

func (console *ConsoleMock) GetWidth() int {
	return MockConsoleWidth
}

func (console *ConsoleMock) GetHeight() int {
	return MockConsoleHeight
}

func (console *ConsoleMock) GetSize() (int, int) {
	return MockConsoleWidth, MockConsoleHeight
}

func (console *ConsoleMock) SetCursorStyle(cursorStyle CursorStyle) error {
	return nil
}

func (console *ConsoleMock) Dispose() error {
	return nil
}
