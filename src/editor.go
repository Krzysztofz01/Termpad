package main

import (
	"errors"
	"os"
	"runtime"
)

// TODO: Verify if the cursor can be out of display now, when it is wrapping the console API
// TODO: Move key handler to helper struct
// TODO: Better wrapper approach for keeping sync during operation on both internal and console API components

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
	if err := editor.cursor.Init(0, 0, console); err != nil {
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

	// TODO: This will change after implementation of all redrawing and rendering functions
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
			{
				editorBreak, err := editor.handleConsoleEventKeyPress(event)
				if err != nil || editorBreak {
					return err
				}
			}

		case ConsoleEventResize:
			{
				editorBreak, err := editor.handleConsoleEventResize(event)
				if err != nil || editorBreak {
					return err
				}
			}
		}
	}
}

// Handling function for the ConsoleEventKeyPress console event. The funcation returns a bool value indicating if the editor loop should be broken
func (editor *Editor) handleConsoleEventKeyPress(event ConsoleEventKeyPress) (bool, error) {
	if event.Key == KeyPrintable {
		if err := editor.handleKeyPrintableCharacter(event.Char); err != nil {
			return false, err
		}

		xOffset := editor.cursor.GetOffsetX() + 1
		if err := editor.cursor.SetOffsetX(xOffset); err != nil {
			return false, err
		}

		return false, editor.renderChanges()
	}

	var breakEditorLoop bool = false
	var err error = nil

	switch event.Key {
	case KeyEnter:
		err = editor.handleKeyEnter()
	case KeyLeft:
		err = editor.handleKeyLeftArrow()
	case KeyRight:
		err = editor.handleKeyRightArrow()
	case KeyUp:
		err = editor.handleKeyUpArrow()
	case KeyDown:
		err = editor.handleKeyDownArrow()
	default:
		err = errors.New("editor: can not handle given input")
	}

	if err != nil {
		return false, err
	}

	return breakEditorLoop, editor.renderChanges()
}

// Handling function for the ConsoleEventResize console event. The funcation returns a bool value indicating if the editor loop should be broken
func (editor *Editor) handleConsoleEventResize(event ConsoleEventResize) (bool, error) {
	// TODO: Implement some checks to verify if the redraw is even required
	return false, editor.renderChanges()
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

// Handle the underlying console API render. If the cursor is out of display boundary the whole screen will be rewriten
//
// TODO: Currently every change is redrawing the whole screen, a better appraochs is required to render only the line that
// has changes (or lines, in case of line break/insert)
// FIXME: This funcation should not call any drawing funcation. It should only pass changes to underling API
func (editor *Editor) renderChanges() error {
	// if !editor.display.CursorInBoundries() {
	// 	if err := editor.redrawText(); err != nil {
	// 		return err
	// 	}
	// }

	// if err := editor.redrawFull(); err != nil {
	// 	return err
	// }

	return editor.console.Commit()
}

// Function is rewriting text changes to the underlying console API screen, according to the display boundaries. All lines are affected
//
// TODO: The text can be longer than the screen. The editor will require a functionality
// to move the current visible content. The current implementation is ,,naive” and
// does not handle screen overflow.
//
// The implementation of the display may solve the problem. But further tests are required
// to verfiy all posible offset shift posibilities
func (editor *Editor) redrawFull() error {
	if err := editor.console.Clear(); err != nil {
		return err
	}

	tHeight := editor.text.GetLineCount()
	xShift := editor.display.GetXOffsetShift()
	yShift := editor.display.GetYOffsetShift()

	for ytIndex := yShift; ytIndex < tHeight; ytIndex += 1 {
		tWidth, err := editor.text.GetLineLengthByOffset(ytIndex)
		if err != nil {
			return err
		}

		for xtIndex := xShift; xtIndex < tWidth; xtIndex += 1 {
			char, err := editor.text.GetCharacterByOffsets(xtIndex, ytIndex)
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

// Function is rewriting text changes to the underlying console API screen, according to the display boundaries. Only the line specified by the cursor is affected.
// TODO: Verify if the y (vertical) offset is correclty calculated
func (editor *Editor) redrawLine(fullRedrawFallback bool) error {
	if editor.display.CursorInBoundries() && fullRedrawFallback {
		return editor.redrawFull()
	}

	ytOffset := editor.cursor.GetOffsetY()
	ycOffset := ytOffset - editor.display.GetYOffsetShift()

	tWidth, err := editor.text.GetLineLengthByOffset(ytOffset)
	if err != nil {
		return err
	}

	// TODO: Decide where to access this information (console/display)
	cWidth := editor.console.GetWidth()

	xShfit := editor.display.GetXOffsetShift()

	for xcIndex := 0; xcIndex < cWidth; xcIndex += 1 {
		xtIndex := xcIndex + xShfit

		var char rune = ' '
		if xtIndex < tWidth {
			char, err = editor.text.GetCharacterByOffsets(xtIndex, ytOffset)
			if err != nil {
				return err
			}
		}

		if err := editor.console.InsertCharacter(xcIndex, ycOffset, char); err != nil {
			return err
		}
	}

	return nil
}

// Function is rewriting text changes to the underlying console API screen, according to the display boundaries. All lines (including the current) below the cursor are affected.
// TODO: Implement full redraw fallback after resolving TODO:9
func (editor *Editor) redrawBelow(fullRedrawFallback bool) error {
	// NOTE: Perform loop starting from the current y index
	// NOTE: Dont loop x in range of line length, but in display width
	// NOTE: Fill with rune until line end, than fill with spaces (this is how clear works)
	return errors.New("editor: not implemeneted")
}

// NOTE: This section contains all the key-specific handler functions

// [<] Handle left arrow key. Handling the movement of the cursor to the left, considering both x and y axis
//
// TODO: Partial cursor position change. May cause invalid cursor state.
// Cursor position handling function callers should backup prev position
// in order to restore it if the operation returns an error
func (editor *Editor) handleKeyLeftArrow() error {
	xOffset := editor.cursor.GetOffsetX()
	if xOffset > 0 {
		xOffset -= 1
		if err := editor.cursor.SetOffsetX(xOffset); err != nil {
			return err
		}

		return nil
	}

	yOffset := editor.cursor.GetOffsetY()
	if yOffset > 0 {
		yOffset -= 1
		if err := editor.cursor.SetOffsetY(yOffset); err != nil {
			return err
		}

		xLength, err := editor.text.GetLineLengthByCursor(editor.cursor)
		if err != nil {
			return err
		}

		xOffset = xLength
		if err := editor.cursor.SetOffsetX(xOffset); err != nil {
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
func (editor *Editor) handleKeyRightArrow() error {
	xOffset := editor.cursor.GetOffsetX()

	lineLength, err := editor.text.GetLineLengthByCursor(editor.cursor)
	if err != nil {
		return err
	}

	if xOffset < lineLength {
		xOffset += 1
		if err := editor.cursor.SetOffsetX(xOffset); err != nil {
			return err
		}

		return nil
	}

	yOffset := editor.cursor.GetOffsetY()
	if yOffset < editor.text.GetLineCount()-1 {
		yOffset += 1
		xOffset = 0
		if err := editor.cursor.SetOffsets(xOffset, yOffset); err != nil {
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
func (editor *Editor) handleKeyUpArrow() error {
	yOffset := editor.cursor.GetOffsetY()
	if yOffset == 0 {
		return nil
	}

	yOffset -= 1
	if err := editor.cursor.SetOffsetY(yOffset); err != nil {
		return err
	}

	targetXLength, err := editor.text.GetLineLengthByCursor(editor.cursor)
	if err != nil {
		return err
	}

	xOffset := editor.cursor.GetOffsetX()
	if xOffset >= targetXLength {
		if err := editor.cursor.SetOffsetX(targetXLength); err != nil {
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
func (editor *Editor) handleKeyDownArrow() error {
	yOffset := editor.cursor.GetOffsetY()
	if yOffset == editor.text.GetLineCount()-1 {
		return nil
	}

	yOffset += 1
	if err := editor.cursor.SetOffsetY(yOffset); err != nil {
		return err
	}

	targetXLength, err := editor.text.GetLineLengthByCursor(editor.cursor)
	if err != nil {
		return err
	}

	xOffset := editor.cursor.GetOffsetX()
	if xOffset >= targetXLength {
		if err := editor.cursor.SetOffsetX(targetXLength); err != nil {
			return err
		}

		return nil
	}

	return nil
}

// [Enter] Handle line breaking via the enter key.
func (editor *Editor) handleKeyEnter() error {
	if err := editor.text.InsertLine(editor.cursor); err != nil {
		return err
	}

	// FIXME: Here is a posiblity for implementing a more efficient render. But this is working for now. Use redrawLine()
	if err := editor.redrawFull(); err != nil {
		return err
	}

	targetXOffset := 0
	targetYOffset := editor.cursor.GetOffsetY() + 1
	if err := editor.cursor.SetOffsets(targetXOffset, targetYOffset); err != nil {
		return err
	}

	return nil
}

// [ASCII 0x20 - 0x7E] Handle printable character insertion.
func (editor *Editor) handleKeyPrintableCharacter(char rune) error {
	if err := editor.text.InsertCharacter(char, editor.cursor); err != nil {
		return err
	}

	if err := editor.redrawLine(true); err != nil {
		return err
	}

	return nil
}
