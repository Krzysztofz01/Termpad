package main

import (
	"errors"
	"os"
	"path/filepath"
)

// TODO: Move key handler to helper struct

// Structure representing the editor instance which is a warapper for text I/O
type Editor struct {
	filePath   string
	fileName   string
	fileExists bool
	console    Console
	display    *Display
	text       *Text
	cursor     *Cursor
	history    *History
	config     *Config
	keybinds   *Keybinds
	menu       *Menu
}

// Editor structure initialization funcation
func (editor *Editor) Init(filePath string, console Console, config *Config) error {
	if len(filePath) <= 0 {
		return errors.New("editor: invalid path passed to editor")
	}

	editor.filePath = filePath
	editor.fileName = filepath.Base(editor.filePath)

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

	if config == nil {
		return errors.New("editor: invalid config reference")
	}

	editor.config = config

	editor.text = new(Text)
	if err := editor.text.Init(fileTextContent, !editor.fileExists, &editor.config.TextConfiguration); err != nil {
		return err
	}

	editor.cursor = new(Cursor)
	if err := editor.cursor.Init(0, 0, console, &editor.config.CursorConfiguration); err != nil {
		return err
	}

	editor.history = new(History)
	if err := editor.history.Init(&editor.config.HistoryConfiguration); err != nil {
		return err
	}

	editor.keybinds = new(Keybinds)
	if err := editor.keybinds.Init(&editor.config.KeybindsConfiguration); err != nil {
		return err
	}

	editor.menu = new(Menu)
	if err := editor.menu.Init(editor.fileName, editor.text.GetEndOfLineSequenceName()); err != nil {
		return err
	}

	if err := editor.menuUpdateInformation(); err != nil {
		return err
	}

	editorPadding := new(Padding)
	if err := editorPadding.Init(0, MenuHeight, 0, 0); err != nil {
		return err
	}

	editor.display = new(Display)
	if err := editor.display.Init(editor.cursor, editorPadding, editor.console); err != nil {
		return err
	}

	if err := editor.display.RedrawTextFull(editor.text); err != nil {
		return err
	}

	if err := editor.display.RedrawMenu(editor.menu); err != nil {
		return err
	}

	if err := editor.display.RenderChanges(); err != nil {
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
	var breakEditorLoop bool = false
	var err error = nil

	// NOTE: Reseting the content of the menu notification
	if err := editor.menu.SetNotificationText(""); err != nil {
		return false, err
	}

	// NOTE: The [Ctrl] key modifier was applied
	if event.Modifier == ModifierCtrl {
		if event.Key == KeyPrintable {
			switch event.Char {
			case editor.keybinds.GetSaveKeybind():
				err = editor.handleKeybindSave()
			case editor.keybinds.GetExitKeybind():
				breakEditorLoop, err = editor.handleKeybindExit()
			default:
				err = errors.New("editor: can not handle given input")
			}
		} else {
			switch event.Key {
			case KeyLeft:
				err = editor.handleKeysCtrlArrowLeft()
			case KeyRight:
				err = editor.handleKeysCtrlArrowRight()
			default:
				err = errors.New("editor: can not handle given input")
			}
		}

	}

	// NOTE: The [Alt] key modifier was applied
	if event.Modifier == ModifierAlt {
		switch event.Char {
		default:
			err = errors.New("editor: can not handle given input")
		}
	}

	// NOTE: The [Shift] or none key modifier applied
	if event.Modifier == ModifierNone || event.Modifier == ModifierShift {
		switch event.Key {
		case KeyPrintable:
			err = editor.handleKeyPrintableCharacter(event.Char)
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
	}

	if err != nil {
		return false, err
	}

	if !editor.display.CursorInBoundries() {
		if err := editor.display.RecalculateBoundaries(); err != nil {
			return false, err
		}

		if err := editor.display.RedrawTextFull(editor.text); err != nil {
			return false, err
		}
	}

	if err := editor.menuUpdateInformation(); err != nil {
		return false, err
	}

	if err := editor.display.RedrawMenu(editor.menu); err != nil {
		return false, err
	}

	return breakEditorLoop, editor.display.RenderChanges()
}

// Handling function for the ConsoleEventResize console event. The funcation returns a bool value indicating if the editor loop should be broken
func (editor *Editor) handleConsoleEventResize(event ConsoleEventResize) (bool, error) {
	if !editor.display.HasSizeChanged(event.Width, event.Height) {
		return false, nil
	}

	if err := editor.display.Resize(event.Width, event.Height); err != nil {
		return false, err
	}

	if err := editor.display.RedrawTextFull(editor.text); err != nil {
		return false, err
	}

	if err := editor.menuUpdateInformation(); err != nil {
		return false, err
	}

	if err := editor.display.RedrawMenu(editor.menu); err != nil {
		return false, err
	}

	return false, editor.display.RenderChanges()
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

// Helper function used to update the cursor position and file modification informations displayed on the menu widget
func (editor *Editor) menuUpdateInformation() error {
	if err := editor.menu.SetCursorPositionText(*editor.cursor); err != nil {
		return err
	}

	if err := editor.menu.SetFileModificationState(editor.text.IsModified()); err != nil {
		return err
	}

	return nil
}

// NOTE: This section contains all the key-specific handler functions

// [<] Handle left arrow key. Handling the movement of the cursor to the left, considering both x and y axis
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

	if err := editor.display.RedrawTextBelow(editor.text, true); err != nil {
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
		if err := editor.display.RedrawTextFull(editor.text); err != nil {
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

	if err := editor.display.RedrawTextLine(editor.text, true); err != nil {
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

		if err := editor.display.RedrawTextBelow(editor.text, true); err != nil {
			return err
		}

		return nil
	}

	// NOTE: The default case when we are removing a character
	if err := editor.text.RemoveCharacterTail(editor.cursor); err != nil {
		return err
	}

	if err := editor.display.RedrawTextLine(editor.text, true); err != nil {
		return err
	}

	return nil
}

// [ASCII 0x20 - 0x7E] Handle printable character insertion.
func (editor *Editor) handleKeyPrintableCharacter(char rune) error {
	if err := editor.text.InsertCharacter(char, editor.cursor); err != nil {
		return err
	}

	if err := editor.display.RedrawTextLine(editor.text, true); err != nil {
		return err
	}

	xOffset := editor.cursor.GetOffsetX()
	if err := editor.cursor.SetOffsetX(xOffset + 1); err != nil {
		return err
	}

	return nil
}

// [Ctrl] + [<] Handle multi-key left jump to next word
func (editor *Editor) handleKeysCtrlArrowLeft() error {
	xOffset := editor.cursor.GetOffsetX()
	yOffset := editor.cursor.GetOffsetY()

	if xOffset == 0 && yOffset == 0 {
		return nil
	}

	if xOffset == 0 && yOffset > 0 {
		yOffset -= 1
		targetXLength, err := editor.text.GetLineLengthByOffset(yOffset)
		if err != nil {
			return err
		}

		return editor.cursor.SetOffsets(targetXLength, yOffset)
	}

	xOffset -= 1
	targetXIndex := -1

	for xIndex := xOffset; xIndex > 0; xIndex -= 1 {
		char, err := editor.text.GetCharacterByOffsets(xIndex, yOffset)
		if err != nil {
			return err
		}

		if char == ' ' {
			if targetXIndex != -1 {
				return editor.cursor.SetOffsetX(targetXIndex)
			}
		} else {
			targetXIndex = xIndex
		}
	}

	return editor.cursor.SetOffsetX(0)
}

// [Ctrl] + [>] Handle multi-key right jump to next word
func (editor *Editor) handleKeysCtrlArrowRight() error {
	xOffset := editor.cursor.GetOffsetX()
	yOffset := editor.cursor.GetOffsetY()
	yOffsetMax := editor.text.GetLineCount() - 1

	currentXLength, err := editor.text.GetLineLengthByOffset(yOffset)
	if err != nil {
		return err
	}

	if xOffset == currentXLength && yOffset == yOffsetMax {
		return nil
	}

	// NOTE: Other text editors jump to next word after switching to line below.
	// The current implementation always jumps to the start of the line, just like
	// the left-jump always jumps to end on switching to the line above.
	if xOffset == currentXLength && yOffset < yOffsetMax {
		yOffset += 1
		return editor.cursor.SetOffsets(0, yOffset)
	}

	targetSpacePassed := false

	for xIndex := xOffset; xIndex < currentXLength; xIndex += 1 {
		char, err := editor.text.GetCharacterByOffsets(xIndex, yOffset)
		if err != nil {
			return err
		}

		if char == ' ' {
			targetSpacePassed = true
		} else {
			if targetSpacePassed {
				return editor.cursor.SetOffsetX(xIndex)
			}
		}
	}

	return editor.cursor.SetOffsetX(currentXLength)
}

// [Ctrl] + [ASCII 0x20 - 0x7E (defined by configuration)] Handle file save keybind
func (editor *Editor) handleKeybindSave() error {
	if err := editor.SaveChanges(); err != nil {
		return err
	}

	if err := editor.menu.SetNotificationText("Changes saved successful."); err != nil {
		return err
	}

	if err := editor.text.ResetModificationState(); err != nil {
		return err
	}

	return nil
}

// [Ctrl] + [ASCII 0x20 - 0x7E (defined by configuration)] Handle program exit keybind. The funcation
// is returning a bool value that idicates if the program loop should be broken.
// TODO: Exit confirmation on un-saved changes
// TODO: Implement ,,alternate screen” in order to restore previous console content after program exit
func (editor *Editor) handleKeybindExit() (bool, error) {
	return true, nil
}
