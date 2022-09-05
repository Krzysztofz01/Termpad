package main

import (
	"errors"
	"os"
)

// Structure representing the editor instance which is a warapper for text I/O
type Editor struct {
	filePath   string
	fileExists bool
	console    Console
	text       *Text
	cursor     *Cursor

	// TODO: Implement preferences
	// TODO: Changes stack
}

// Editor structure initialization funcation
func (editor *Editor) Init(filePath string, console Console) error {
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

	if console == nil {
		return errors.New("editor: invalid internal console api contract implementation")
	}

	editor.console = console
	return nil
}

// Start the editor loop
func (editor *Editor) Start() error {
	// TODO: Implementation
	return errors.New("editor: not implemented")
}

// Generate string from text structure and create or truncate target file
func (editor *Editor) SaveChanges() error {
	file, err := os.Create(editor.filePath)
	if err != nil {
		return err
	}

	textContent, err := editor.text.GetTextAsString()
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
