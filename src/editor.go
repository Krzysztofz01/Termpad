package main

import (
	"errors"
	"os"
	"runtime"
)

// TODO: Verify if the cursor can be out of display now, when it is wrapping the console API
// TODO: Move key handler to helper struct
// TODO: Better wrapper approach for keeping sync during operation on both internal and console API components
// TODO: Add support for cursor style preferences

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

	if err := editor.redrawFull(); err != nil {
		return err
	}

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
	case KeyBackspace:
		err = editor.handleKeyBackspace()
	case KeyDelete:
		err = editor.handleKeyDelete()
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
	applyRedraw := false
	width, height := editor.console.GetSize()

	if editor.display.HasSizeChanged(width, height) {
		if err := editor.display.Resize(width, height); err != nil {
			return false, err
		}

		applyRedraw = true
	}

	if !editor.display.CursorInBoundries() {
		applyRedraw = true
	}

	if !applyRedraw {
		return false, nil
	}

	if err := editor.redrawFull(); err != nil {
		return false, err
	}

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

// Request a render of all changes to the screen of the underlying console API
func (editor *Editor) renderChanges() error {
	return editor.console.Commit()
}

// Function is rewriting text changes to the underlying console API screen, according to the display boundaries. All lines are affected
// TODO: The tHeight can be greater than dHeight - we can avoid redundand off-display rendering
func (editor *Editor) redrawFull() error {
	ytLength := editor.text.GetLineCount()
	xcLength, ycLength := editor.console.GetSize()
	xShift := editor.display.GetXOffsetShift()
	yShift := editor.display.GetYOffsetShift()

	for ycIndex := 0; ycIndex < ycLength; ycIndex += 1 {
		ytIndex := ycIndex + yShift

		if ycIndex < ytLength {
			xtLength, err := editor.text.GetLineLengthByOffset(ytIndex)
			if err != nil {
				return err
			}

			for xcIndex := 0; xcIndex < xcLength; xcIndex += 1 {
				xtIndex := xcIndex + xShift

				var char rune = ' '

				if xtIndex < xtLength {
					char, err = editor.text.GetCharacterByOffsets(xtIndex, ytIndex)
					if err != nil {
						return err
					}
				}

				if err := editor.console.InsertCharacter(xcIndex, ycIndex, char); err != nil {
					return err
				}
			}

			continue
		}

		for xcIndex := 0; xcIndex < xcLength; xcIndex += 1 {
			if err := editor.console.InsertCharacter(xcIndex, ycIndex, ' '); err != nil {
				return err
			}
		}
	}

	return nil
}

// Function is rewriting text changes to the underlying console API screen, according to the display boundaries. Only the line specified by the cursor is affected.
func (editor *Editor) redrawLine(fullRedrawFallback bool) error {
	if !editor.display.CursorInBoundries() && fullRedrawFallback {
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
// TODO: The ycIndex < ytLength condition prevents the overwrting of previous screen data
func (editor *Editor) redrawBelow(fullRedrawFallback bool) error {
	if !editor.display.CursorInBoundries() && fullRedrawFallback {
		return editor.redrawFull()
	}

	ytLength := editor.text.GetLineCount()
	xcLength, ycLength := editor.console.GetSize()
	xShift := editor.display.GetXOffsetShift()
	yShift := editor.display.GetYOffsetShift()

	for ycIndex := editor.cursor.GetOffsetY() - yShift; ycIndex < ycLength; ycIndex += 1 {
		ytIndex := ycIndex + yShift

		if ycIndex < ytLength {
			xtLength, err := editor.text.GetLineLengthByOffset(ytIndex)
			if err != nil {
				return err
			}

			for xcIndex := 0; xcIndex < xcLength; xcIndex += 1 {
				xtIndex := xcIndex + xShift

				var char rune = ' '

				if xtIndex < xtLength {
					char, err = editor.text.GetCharacterByOffsets(xtIndex, ytIndex)
					if err != nil {
						return err
					}
				}

				if err := editor.console.InsertCharacter(xcIndex, ycIndex, char); err != nil {
					return err
				}
			}

			continue
		}

		for xcIndex := 0; xcIndex < xcLength; xcIndex += 1 {
			if err := editor.console.InsertCharacter(xcIndex, ycIndex, ' '); err != nil {
				return err
			}
		}
	}

	return nil
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

	if err := editor.redrawBelow(true); err != nil {
		return err
	}

	yOffset := editor.cursor.GetOffsetY()
	if err := editor.cursor.SetOffsets(0, yOffset+1); err != nil {
		return err
	}

	return nil
}

// [Backspace] Handle character removing via the backspace key
func (editor *Editor) handleKeyBackspace() error {
	xOffset := editor.cursor.GetOffsetX()
	yOffset := editor.cursor.GetOffsetY()

	// NOTE: The case when we are at the begining of the text
	if yOffset == 0 && xOffset == 0 {
		return nil
	}

	// NOTE: The case when we need to remove the new line
	if xOffset == 0 {
		targetLineLength, err := editor.text.GetLineLengthByOffset(yOffset - 1)
		if err != nil {
			return err
		}

		if err := editor.text.CombineLine(editor.cursor, false); err != nil {
			return err
		}

		// TODO: This look like it can be optimized, but there is currently no ,,redraw above'' function.
		// And a func with such capabilities would still redraw the content below. It can only optimize
		// endge cases when we are editing the last line in current display range. Too much hustle for
		// such negligible performance improvement.
		if err := editor.redrawFull(); err != nil {
			return err
		}

		if err := editor.cursor.SetOffsets(targetLineLength, yOffset-1); err != nil {
			return err
		}

		return nil
	}

	// NOTE: The default case when we are removing a character
	if err := editor.text.RemoveCharacterHead(editor.cursor); err != nil {
		return err
	}

	if err := editor.redrawLine(true); err != nil {
		return err
	}

	if err := editor.cursor.SetOffsetX(xOffset - 1); err != nil {
		return err
	}

	return nil
}

// [Delete] Handle character removing via the backsapce key
func (editor *Editor) handleKeyDelete() error {
	xOffset := editor.cursor.GetOffsetX()
	yOffset := editor.cursor.GetOffsetY()

	targetLineLength, err := editor.text.GetLineLengthByOffset(yOffset)
	if err != nil {
		return err
	}

	// NOTE: The case when we are at the end of the text
	if yOffset == editor.text.GetLineCount()-1 && xOffset == targetLineLength {
		return nil
	}

	// NOTE: The case when we need to remove the new line
	if xOffset == targetLineLength {
		if err := editor.text.CombineLine(editor.cursor, true); err != nil {
			return err
		}

		if err := editor.redrawBelow(true); err != nil {
			return err
		}

		return nil
	}

	// NOTE: The default case when we are removing a character
	if err := editor.text.RemoveCharacterTail(editor.cursor); err != nil {
		return err
	}

	if err := editor.redrawLine(true); err != nil {
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
