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
		if err := file.Close(); err != nil {
			return err
		}
		return err
	}

	if _, err := file.WriteString(*textContent); err != nil {
		if err := file.Close(); err != nil {
			return err
		}
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	if !editor.fileExists {
		editor.fileExists = true
	}

	return nil
}
