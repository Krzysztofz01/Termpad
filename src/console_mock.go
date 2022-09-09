package main

// Structure implementing the console contract, used as a mockup for testing purposes
type ConsoleMock struct {
}

// Create a new instance of the Tcell based console
func CreateConsoleMockup() Console {
	return &ConsoleMock{}
}

func (console *ConsoleMock) InsertCharacter(xIndex int, yIndex int, char rune) error {
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
	return 0
}

func (console *ConsoleMock) GetHeight() int {
	return 0
}

func (console *ConsoleMock) GetSize() (int, int) {
	return 0, 0
}

func (console *ConsoleMock) Dispose() error {
	return nil
}
