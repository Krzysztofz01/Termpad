package main

import (
	"errors"

	"github.com/gdamore/tcell/v2"
)

// Structure implementing the console contract based on the console API exposed by Tcell library
type ConsoleTcell struct {
	screen tcell.Screen
}

// Create a new instance of the Tcell based console
func CreateConsole() (Console, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	screen.DisableMouse()
	screen.ShowCursor(0, 0)

	console := ConsoleTcell{
		screen: screen,
	}

	if err := console.SetCursorStyle(BarCursorStatic); err != nil {
		return nil, err
	}

	return &console, nil
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

func (console *ConsoleTcell) SetCursorPosition(xIndex int, yIndex int) error {
	console.screen.ShowCursor(xIndex, yIndex)
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
					Key:      console.translateNamedKey(event.Key()),
					Modifier: console.translateModifierKey(event.Modifiers()),
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

func (console *ConsoleTcell) SetCursorStyle(cursorStyle CursorStyle) error {
	if cursorStyle == NoCursor {
		console.screen.HideCursor()
		return nil
	}

	var selectedCursorStyle tcell.CursorStyle

	switch cursorStyle {
	case BarCursorStatic:
		selectedCursorStyle = tcell.CursorStyleSteadyBar
	case BarCursorDynamic:
		selectedCursorStyle = tcell.CursorStyleBlinkingBar
	case BlockCursorStatic:
		selectedCursorStyle = tcell.CursorStyleSteadyBlock
	case BlockCursorDynamic:
		selectedCursorStyle = tcell.CursorStyleBlinkingBlock
	case LineCursorStatic:
		selectedCursorStyle = tcell.CursorStyleSteadyUnderline
	case LineCursorDynamic:
		selectedCursorStyle = tcell.CursorStyleBlinkingUnderline
	default:
		return errors.New("console: invalid internal cursor style")
	}

	console.screen.SetCursorStyle(selectedCursorStyle)
	return nil
}

func (console *ConsoleTcell) Dispose() error {
	console.screen.Fini()
	return nil
}

// Helper funcation used for converting implementation specific to contract specific named key representation
func (console *ConsoleTcell) translateNamedKey(key tcell.Key) NamedKey {
	switch key {
	case tcell.KeyRune:
		return KeyPrintable
	case tcell.KeyUp:
		return KeyUp
	case tcell.KeyDown:
		return KeyDown
	case tcell.KeyRight:
		return KeyRight
	case tcell.KeyLeft:
		return KeyLeft
	case tcell.KeyPgUp:
		return KeyPgUp
	case tcell.KeyPgDn:
		return KeyPgDn
	case tcell.KeyHome:
		return KeyHome
	case tcell.KeyEnd:
		return KeyEnd
	case tcell.KeyInsert:
		return KeyInsert
	case tcell.KeyDEL:
		return KeyDelete
	case tcell.KeyPause:
		return KeyPause
	case tcell.KeyBacktab:
		return KeyBacktab
	case tcell.KeyCR:
		return KeyEnter
	case tcell.KeyTAB:
		return KeyTab
	case tcell.KeyESC:
		return KeyEscape
	case tcell.KeyBS:
		return KeyBackspace
	case tcell.KeyF1:
		return KeyF1
	case tcell.KeyF2:
		return KeyF2
	case tcell.KeyF3:
		return KeyF3
	case tcell.KeyF4:
		return KeyF4
	case tcell.KeyF5:
		return KeyF5
	case tcell.KeyF6:
		return KeyF6
	case tcell.KeyF7:
		return KeyF7
	case tcell.KeyF8:
		return KeyF8
	case tcell.KeyF9:
		return KeyF9
	case tcell.KeyF10:
		return KeyF10
	case tcell.KeyF11:
		return KeyF11
	case tcell.KeyF12:
		return KeyF12
	default:
		return NamedKey(-1)
	}
}

// Helper funcation used for converting implementation specific to contract specific modifier key representation
func (console *ConsoleTcell) translateModifierKey(mod tcell.ModMask) ModifierKey {
	switch mod {
	case tcell.ModShift:
		return ModifierShift
	case tcell.ModCtrl:
		return ModifierCtrl
	case tcell.ModAlt:
		return ModifierAlt
	default:
		return ModifierNone
	}
}
