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

	if _, err := os.Stat(editor.filePath); err != nil {
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

	editor.text = new(Text)
	if err := editor.text.Init(fileTextContent, !editor.fileExists); err != nil {
		return err
	}

	editor.cursor = new(Cursor)
	if err := editor.cursor.Init(0, 0); err != nil {
		return err
	}

	editor.history = new(History)
	if err := editor.history.Init(); err != nil {
		return err
	}

	if console == nil {
		return errors.New("editor: invalid internal console api contract implementation")
	}

	editor.console = console

	if config == nil {
		return errors.New("editor: invalid config reference")
	}

	editor.config = config

	if err := editor.redrawText(); err != nil {
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
				if event.Key == KeyPrintable {
					if err := editor.handleNamedPrintableKey(event); err != nil {
						// TODO: Handle error
					}
				}

				if err := editor.handleNamedKey(event); err != nil {
					// TODO: Handle error
				}
			}
			break

		case ConsoleEventResize:
			{
				// TODO: Check if cursor is on screen to prevent redundant full content redraws
				if err := editor.redrawText(); err != nil {
					// TODO: Handle error
				}
			}
			break
		}
	}

	return errors.New("editor: not implemented")
}

func (editor *Editor) handleNamedPrintableKey(event ConsoleEventKeyPress) error {
	// TODO: Implement
	return errors.New("editor: not implemented")
}

func (editor *Editor) handleNamedKey(event ConsoleEventKeyPress) error {
	// TODO: Implement
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

// Function is handling the movement of cursor to the left in x and y axis
//
// TODO: Partial cursor position change. May cause invalid cursor state.
// Cursor position handling function callers should backup prev position
// in order to restore it if the operation returns an error
func (editor *Editor) moveCursorLeft() error {
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
		if err := editor.cursor.SetOffsetY(xOffset); err != nil {
			return err
		}

		xLength, err := editor.text.GetLineLength(editor.cursor)
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

// Function is handling the movement of cursor to the right in x and y axis
//
// TODO: Partial cursor position change. May cause invalid cursor state.
// Cursor position handling function callers should backup prev position
// in order to restore it if the operation returns an error
func (editor *Editor) moveCursorRight() error {
	xOffset := editor.cursor.GetOffsetX()

	lineLength, err := editor.text.GetLineLength(editor.cursor)
	if err != nil {
		return err
	}

	if xOffset < lineLength {
		xOffset += 1
		if err := editor.cursor.SetOffsetY(xOffset); err != nil {
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

// Function is handling the movement of cursor to the line above in x and y axis
//
// TODO: Partial cursor position change. May cause invalid cursor state.
// Cursor position handling function callers should backup prev position
// in order to restore it if the operation returns an error
func (editor *Editor) moveCursorUp() error {
	yOffset := editor.cursor.GetOffsetY()
	if yOffset == 0 {
		return nil
	}

	yOffset -= 1
	currentXLength, err := editor.text.GetLineLength(editor.cursor)
	if err != nil {
		return err
	}

	if err := editor.cursor.SetOffsetY(yOffset); err != nil {
		return err
	}

	targetXLength, err := editor.text.GetLineLength(editor.cursor)
	if err != nil {
		return err
	}

	if targetXLength < currentXLength {
		if err := editor.cursor.SetOffsetX(targetXLength); err != nil {
			return err
		}

		return nil
	}

	return nil
}

// Function is handling the movement of cursor to the line below in x and y axis
//
// TODO: Partial cursor position change. May cause invalid cursor state.
// Cursor position handling function callers should backup prev position
// in order to restore it if the operation returns an error
func (editor *Editor) moveCursorDown() error {
	yOffset := editor.cursor.GetOffsetY()
	if yOffset == editor.text.GetLineCount()-1 {
		return nil
	}

	yOffset += 1
	currentXLength, err := editor.text.GetLineLength(editor.cursor)
	if err != nil {
		return err
	}

	if err := editor.cursor.SetOffsetY(yOffset); err != nil {
		return err
	}

	targetXLength, err := editor.text.GetLineLength(editor.cursor)
	if err != nil {
		return err
	}

	if targetXLength < currentXLength {
		if err := editor.cursor.SetOffsetX(targetXLength); err != nil {
			return err
		}

		return nil
	}

	return nil
}

// Funcation is rendering all the text to the screen
//
// TODO: The text can be longer than the screen. The editor will require a functionality
// to move the current visible content. The current implementation is ,,naive” and
// does not handle screen overflow.
func (editor *Editor) redrawText() error {
	if err := editor.console.Clear(); err != nil {
		return err
	}

	initCursor := new(Cursor)
	if err := initCursor.Init(0, 0); err != nil {
		return err
	}

	height := editor.text.GetLineCount()
	for yIndex := 0; yIndex < height; yIndex += 1 {
		if err := initCursor.SetOffsets(0, yIndex); err != nil {
			return err
		}

		width, err := editor.text.GetLineLength(initCursor)
		if err != nil {
			return err
		}

		for xIndex := 0; xIndex < width; xIndex += 1 {
			if err := initCursor.SetOffsets(xIndex, yIndex); err != nil {
				return err
			}

			char, err := editor.text.GetCharacter(initCursor)
			if err != nil {
				return err
			}

			if err := editor.console.InsertCharacter(xIndex, yIndex, char); err != nil {
				return err
			}
		}
	}

	return editor.console.Commit()
}
