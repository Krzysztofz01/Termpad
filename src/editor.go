package main

import (
	"errors"
	"os"
	"runtime"
)

// Structure representing the editor instance which is a warapper for text I/O
type Editor struct {
	filePath   string
	fileExists bool
	console    Console
	display    *Display
	text       *Text
	cursor     *Cursor
	history    *History
	config     *Config
}

// Editor structure initialization funcation
func (editor *Editor) Init(filePath string, console Console, config *Config) error {
	if len(filePath) <= 0 {
		return errors.New("editor: invalid path passed to editor")
	}

	editor.filePath = filePath

	if _, err := os.Stat(editor.filePath); err == nil {
		editor.fileExists = true
	} else if errors.Is(err, os.ErrNotExist) {
		editor.fileExists = false
	} else {
		return errors.New("editor: can not determine if the file is accesable")
	}

	fileTextContent := ""
	if editor.fileExists {
		fileData, err := os.ReadFile(editor.filePath)
		if err != nil {
			return err
		}

		fileTextContent = string(fileData)
	}

	if console == nil {
		return errors.New("editor: invalid internal console api contract implementation")
	}

	editor.console = console

	editor.text = new(Text)
	if err := editor.text.Init(fileTextContent, !editor.fileExists); err != nil {
		return err
	}

	editor.cursor = new(Cursor)
	if err := editor.cursor.Init(0, 0); err != nil {
		return err
	}

	if err := editor.setCursorPosition(0, 0); err != nil {
		return err
	}

	editor.history = new(History)
	if err := editor.history.Init(); err != nil {
		return err
	}

	editor.display = new(Display)
	width, height := editor.console.GetSize()
	if err := editor.display.Init(width, height, editor.cursor); err != nil {
		return err
	}

	if config == nil {
		return errors.New("editor: invalid config reference")
	}

	editor.config = config

	if err := editor.renderChanges(); err != nil {
		return err
	}

	return nil
}

// Start the editor loop
func (editor *Editor) Start() error {
	for {
		ev := editor.console.WatchConsoleEvent()
		switch event := ev.(type) {
		case ConsoleEventKeyPress:
			// TODO: Handle returned error. Different bahaviour for critical and non-fatal errors
			editor.handleConsoleEventKeyPress(event)
			break

		case ConsoleEventResize:
			// TODO: Handle returned error. Different bahaviour for critical and non-fatal errors
			editor.handleConsoleEventResize(event)
			break
		}
	}

	return nil
}

// Handling function for the ConsoleEventKeyPress console event
func (editor *Editor) handleConsoleEventKeyPress(event ConsoleEventKeyPress) error {
	if event.Key == KeyPrintable {
		if err := editor.insertCharacter(event.Char); err != nil {
			return err
		}

		targetXOffset := editor.cursor.GetOffsetX() + 1
		targetYOffset := editor.cursor.GetOffsetY()

		if err := editor.setCursorPosition(targetXOffset, targetYOffset); err != nil {
			return err
		}

		return editor.renderChanges()
	}

	var err error = nil

	switch event.Key {
	case KeyEnter:
		err = editor.handleEnterKey()
		break
	case KeyLeft:
		err = editor.handleLeftArrowKey()
		break
	case KeyRight:
		err = editor.handleRightArrowKey()
		break
	case KeyUp:
		err = editor.handleUpArrowKey()
		break
	case KeyDown:
		err = editor.handleDownArrowKey()
	}

	if err != nil {
		return err
	}

	return editor.renderChanges()
}

// Handling function for the ConsoleEventResize console event
func (editor *Editor) handleConsoleEventResize(event ConsoleEventResize) error {
	return errors.New("editor: not implemented")
}

// Generate string from text structure and create or truncate target file
func (editor *Editor) SaveChanges() error {
	file, err := os.Create(editor.filePath)
	if err != nil {
		return err
	}

	useCarriageReturn := false
	if editor.config.UsePlatformSpecificEndOfLineSequence {
		if runtime.GOOS == "windows" {
			useCarriageReturn = true
		}
	}

	textContent, err := editor.text.GetTextAsString(useCarriageReturn)
	if err != nil {
		if fileErr := file.Close(); fileErr != nil {
			return fileErr
		}
		return err
	}

	if _, err := file.WriteString(*textContent); err != nil {
		if fileErr := file.Close(); fileErr != nil {
			return fileErr
		}
		return err
	}

	if fileErr := file.Close(); fileErr != nil {
		return fileErr
	}

	if !editor.fileExists {
		editor.fileExists = true
	}

	return nil
}

// [<] Handle left arrow key. Handling the movement of the cursor to the left, considering both x and y axis
//
// TODO: Partial cursor position change. May cause invalid cursor state.
// Cursor position handling function callers should backup prev position
// in order to restore it if the operation returns an error
func (editor *Editor) handleLeftArrowKey() error {
	xOffset := editor.cursor.GetOffsetX()
	yOffset := editor.cursor.GetOffsetY()

	if xOffset > 0 {
		xOffset -= 1
		if err := editor.setCursorPosition(xOffset, yOffset); err != nil {
			return err
		}

		return nil
	}

	yOffset = editor.cursor.GetOffsetY()
	if yOffset > 0 {
		yOffset -= 1
		if err := editor.setCursorPosition(xOffset, yOffset); err != nil {
			return err
		}

		xLength, err := editor.text.GetLineLength(editor.cursor)
		if err != nil {
			return err
		}

		xOffset = xLength
		if err := editor.setCursorPosition(xOffset, yOffset); err != nil {
			return err
		}

		return nil
	}

	return nil
}

// [>] Handle right arrow key. Handling the movement of the cursor to the right, considering both x and y axis
//
// TODO: Partial cursor position change. May cause invalid cursor state.
// Cursor position handling function callers should backup prev position
// in order to restore it if the operation returns an error
func (editor *Editor) handleRightArrowKey() error {
	xOffset := editor.cursor.GetOffsetX()
	yOffset := editor.cursor.GetOffsetY()

	lineLength, err := editor.text.GetLineLength(editor.cursor)
	if err != nil {
		return err
	}

	if xOffset < lineLength {
		xOffset += 1
		if err := editor.setCursorPosition(xOffset, yOffset); err != nil {
			return err
		}

		return nil
	}

	yOffset = editor.cursor.GetOffsetY()
	if yOffset < editor.text.GetLineCount()-1 {
		yOffset += 1
		xOffset = 0
		if err := editor.setCursorPosition(xOffset, yOffset); err != nil {
			return err
		}

		return nil
	}

	return nil
}

// [/\] Handle up arrow key. Handling the movement of the cursor to the line above, considering both y and x axis
//
// TODO: Partial cursor position change. May cause invalid cursor state.
// Cursor position handling function callers should backup prev position
// in order to restore it if the operation returns an error
//
// FIXME: Sometimes the cursor is moving to the end of the line after beeing on a shorter one
func (editor *Editor) handleUpArrowKey() error {
	xOffset := editor.cursor.GetOffsetX()
	yOffset := editor.cursor.GetOffsetY()

	if yOffset == 0 {
		return nil
	}

	yOffset -= 1
	currentXLength, err := editor.text.GetLineLength(editor.cursor)
	if err != nil {
		return err
	}

	if err := editor.setCursorPosition(xOffset, yOffset); err != nil {
		return err
	}

	targetXLength, err := editor.text.GetLineLength(editor.cursor)
	if err != nil {
		return err
	}

	if targetXLength < currentXLength {
		if err := editor.setCursorPosition(targetXLength, yOffset); err != nil {
			return err
		}

		return nil
	}

	return nil
}

// [\/] Handle down arrow key. Handling the movement of the cursor to the line below, considering both y and x axis
//
// TODO: Partial cursor position change. May cause invalid cursor state.
// Cursor position handling function callers should backup prev position
// in order to restore it if the operation returns an error
//
// FIXME: Sometimes the cursor is moving to the end of the line after beeing on a shorter one
func (editor *Editor) handleDownArrowKey() error {
	xOffset := editor.cursor.GetOffsetX()
	yOffset := editor.cursor.GetOffsetY()

	if yOffset == editor.text.GetLineCount()-1 {
		return nil
	}

	yOffset += 1
	currentXLength, err := editor.text.GetLineLength(editor.cursor)
	if err != nil {
		return err
	}

	if err := editor.setCursorPosition(xOffset, yOffset); err != nil {
		return err
	}

	targetXLength, err := editor.text.GetLineLength(editor.cursor)
	if err != nil {
		return err
	}

	if targetXLength < currentXLength {
		if err := editor.setCursorPosition(targetXLength, yOffset); err != nil {
			return err
		}

		return nil
	}

	return nil
}

// [Enter] Handle line breaking via the enter key.
func (editor *Editor) handleEnterKey() error {
	if err := editor.text.InsertLine(editor.cursor); err != nil {
		return err
	}

	// FIXME: Here is a posiblity for implementing a more efficient render. But this is working for now
	if err := editor.redrawText(); err != nil {
		return err
	}

	targetXOffset := 0
	targetYOffset := editor.cursor.GetOffsetY() + 1
	if err := editor.setCursorPosition(targetXOffset, targetYOffset); err != nil {
		return err
	}

	return nil
}

// Handle printable character insertion. This function is handling both text structure and underlying console API
func (editor *Editor) insertCharacter(char rune) error {
	if err := editor.text.InsertCharacter(char, editor.cursor); err != nil {
		return err
	}

	xcIndex := editor.cursor.GetOffsetX() - editor.display.GetXOffsetShift()
	ycIndex := editor.cursor.GetOffsetY() - editor.display.GetYOffsetShift()

	if err := editor.console.InsertCharacter(xcIndex, ycIndex, char); err != nil {
		return err
	}

	return nil
}

// Handle cursor position change. This funcation is handling both the cursor struct and the underlying console API
func (editor *Editor) setCursorPosition(xOffset int, yOffset int) error {
	if err := editor.cursor.SetOffsets(xOffset, yOffset); err != nil {
		return err
	}

	if err := editor.console.SetCursorPosition(xOffset, yOffset); err != nil {
		return err
	}

	return nil
}

// Handle the underlying console API render. If the cursor is out of display boundary the whole screen will be rewriten
//
// TODO: Currently every change is redrawing the whole screen, a better appraochs is required to render only the line that
// has changes (or lines, in case of line break/insert)
func (editor *Editor) renderChanges() error {
	// if !editor.display.CursorInBoundries() {
	// 	if err := editor.redrawText(); err != nil {
	// 		return err
	// 	}
	// }

	if err := editor.redrawText(); err != nil {
		return err
	}

	return editor.console.Commit()
}

// Function is clearing and rewriting changes to the underlying console API screen, according to the display boundaries
//
// TODO: The text can be longer than the screen. The editor will require a functionality
// to move the current visible content. The current implementation is ,,naive” and
// does not handle screen overflow.
//
// The implementation of the display may solve the problem. But further tests are required
// to verfiy all posible offset shift posibilities
func (editor *Editor) redrawText() error {
	if err := editor.console.Clear(); err != nil {
		return err
	}

	redrawCursor := new(Cursor)
	if err := redrawCursor.Init(0, 0); err != nil {
		return err
	}

	tHeight := editor.text.GetLineCount()
	xShift := editor.display.GetXOffsetShift()
	yShift := editor.display.GetYOffsetShift()

	for ytIndex := yShift; ytIndex < tHeight; ytIndex += 1 {
		if err := redrawCursor.SetOffsetY(ytIndex); err != nil {
			return err
		}

		tWidth, err := editor.text.GetLineLength(redrawCursor)
		if err != nil {
			return err
		}

		for xtIndex := xShift; xtIndex < tWidth; xtIndex += 1 {
			if err := redrawCursor.SetOffsetX(xtIndex); err != nil {
				return err
			}

			char, err := editor.text.GetCharacter(redrawCursor)
			if err != nil {
				return err
			}

			xcIndex := xtIndex - xShift
			ycIndex := ytIndex - yShift

			if err := editor.console.InsertCharacter(xcIndex, ycIndex, char); err != nil {
				return err
			}
		}
	}

	return nil
}
