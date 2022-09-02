package main

import (
	"errors"

	"github.com/gdamore/tcell"
)

// Structure implementing the console contract based on the console API exposed by Tcell library
type ConsoleTcell struct {
	screen tcell.Screen
}

// Create a new instance of the Tcell based console
func CreateConsole() (Console, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		// TODO: Logging
		return nil, err
	}

	if err := screen.Init(); err != nil {
		// TODO: Logging
		return nil, err
	}

	// TODO: Implement color related preferences etc.
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(style)

	screen.DisableMouse()

	return &ConsoleTcell{
		screen: screen,
	}, nil
}

func (console *ConsoleTcell) InsertCharacter(xIndex int, yIndex int, char rune) error {
	if xIndex < 0 {
		return errors.New("console: invalid x (horizontal) out of bound index requested to insert")
	}

	if yIndex < 0 {
		return errors.New("console: invalid y (vertical) out of bound index requested to insert")
	}

	console.screen.SetContent(xIndex, yIndex, char, nil, tcell.StyleDefault)
	return nil
}

func (console *ConsoleTcell) RemoveCharacter(xIndex int, yIndex int) error {
	if xIndex < 0 {
		return errors.New("console: invalid x (horizontal) out of bound index requested to remove")
	}

	if yIndex < 0 {
		return errors.New("console: invalid y (vertical) out of bound index requested to remove")
	}

	console.screen.SetContent(xIndex, yIndex, 0, nil, tcell.StyleDefault)
	return nil
}

func (console *ConsoleTcell) Commit() error {
	console.screen.Show()
	return nil
}

func (console *ConsoleTcell) Clear() error {
	console.screen.Clear()
	return nil
}

func (console *ConsoleTcell) WatchConsoleEvent() interface{} {
	for {
		switch event := console.screen.PollEvent().(type) {
		case *tcell.EventKey:
			{
				return ConsoleEventKeyPress{
					Char:     event.Rune(),
					Key:      event.Key(),
					Modifier: event.Modifiers(),
				}
			}

		case *tcell.EventResize:
			{
				width, height := event.Size()
				return ConsoleEventResize{
					Width:  width,
					Height: height,
				}
			}
		}
	}
}

func (console *ConsoleTcell) GetWidth() int {
	width, _ := console.screen.Size()

	return width
}

func (console *ConsoleTcell) GetHeight() int {
	_, height := console.screen.Size()

	return height
}

func (console *ConsoleTcell) GetSize() (int, int) {
	width, height := console.screen.Size()

	return width, height
}

func (console *ConsoleTcell) Dispose() error {
	console.screen.Fini()
	return nil
}
