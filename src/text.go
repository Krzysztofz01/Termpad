package main

import (
	"errors"
	"runtime"
	"strings"
)

// TODO: Implement function that will insert multiple characters as array. This function
// would be very usefull for implementing clipboard pasting.

// A structure representing the text, which is a container for the List structures
type Text struct {
	lines             []*Line
	modified          bool
	endOfLineSequence string
	config            *TextConfig
}

// Text structure initialization funcation
func (text *Text) Init(textString string, newFile bool, textConfig *TextConfig) error {
	if textConfig == nil {
		defaultConfig := CreateDefaultTextConfig()
		text.config = &defaultConfig
	} else {
		text.config = textConfig
	}

	if strings.Contains(textString, "\r\n") {
		text.endOfLineSequence = "CRLF"
	} else {
		text.endOfLineSequence = "LF"
	}

	// NOTE: Removing the 0x0D CR (Carriage Return)
	textString = strings.Replace(textString, "\r", "", -1)

	// NOTE: Spliting the text by 0x0A LF (Line Feed)
	textStringLines := strings.Split(textString, "\n")

	text.modified = newFile

	if len(textStringLines) == 0 {
		line := new(Line)
		if err := line.Init(""); err != nil {
			return err
		}

		text.lines = make([]*Line, 1)
		text.lines[0] = line
		return nil
	}

	text.lines = make([]*Line, len(textStringLines))
	for textStringLineIndex, textStringLineValue := range textStringLines {
		line := new(Line)
		if err := line.Init(textStringLineValue); err != nil {
			return err
		}

		text.lines[textStringLineIndex] = line
	}

	return nil
}

// Return the count of lines (height)
func (text *Text) GetLineCount() int {
	return len(text.lines)
}

// Return the length of the line based on given y (vertical) offset
func (text *Text) GetLineLengthByOffset(yOffset int) (int, error) {
	if yOffset < 0 {
		return 0, errors.New("text: invalid y (vertical) negative offset requested to get")
	}

	if yOffset >= len(text.lines) {
		return 0, errors.New("text: invalid y (vertical) out of bound offset requested to get")
	}

	targetLine := text.lines[yOffset]
	return targetLine.GetBufferLength(), nil
}

// Return the length of the line based on given cursor position
func (text *Text) GetLineLengthByCursor(cursor *Cursor) (int, error) {
	return text.GetLineLengthByOffset(cursor.GetOffsetY())
}

// Place a given character inside specific line at specific offset given by the cursor position
func (text *Text) InsertCharacter(char rune, cursor *Cursor) error {
	yOffset := cursor.GetOffsetY()

	if yOffset < 0 {
		return errors.New("text: invalid y (vertical) negative offset requested to insert")
	}

	if yOffset > len(text.lines) {
		return errors.New("text: invalid y (vertical) out of bound offset requested to insert")
	}

	text.modified = true

	targetLine := text.lines[yOffset]
	return targetLine.InsertBufferCharacter(char, cursor)
}

// Remove a character at specific line at specific position before the position given by the offset of the given cursor
func (text *Text) RemoveCharacterHead(cursor *Cursor) error {
	yOffset := cursor.GetOffsetY()

	if yOffset < 0 {
		return errors.New("text: invalid y (vertical) negative offset requested to remove")
	}

	if yOffset > len(text.lines) {
		return errors.New("text: invalid y (vertical) out of bound offset requested to remove")
	}

	text.modified = true

	targetLine := text.lines[yOffset]
	return targetLine.RemoveBufferCharacterHead(cursor)
}

// Remove a character at specific line at specific position behind the position given by the offset of the given cursor
func (text *Text) RemoveCharacterTail(cursor *Cursor) error {
	yOffset := cursor.GetOffsetY()

	if yOffset < 0 {
		return errors.New("text: invalid y (vertical) negative offset requested to remove")
	}

	if yOffset > len(text.lines) {
		return errors.New("text: invalid y (vertical) out of bound offset requested to remove")
	}

	text.modified = true

	targetLine := text.lines[yOffset]
	return targetLine.RemoveBufferCharacterTail(cursor)
}

// Handle line inserting and line breaking
func (text *Text) InsertLine(cursor *Cursor) error {
	yOffset := cursor.GetOffsetY()
	xOffset := cursor.GetOffsetX()

	if yOffset < 0 {
		return errors.New("text: invalid y (vertical) negative offset requested to split")
	}

	if yOffset > len(text.lines) {
		return errors.New("text: invalid y (vertical) out of bound offset requested to split")
	}

	text.modified = true

	targetLine := text.lines[yOffset]

	// NOTE: Breaking the line at the start of the line
	if xOffset == 0 {
		line := new(Line)
		if err := line.Init(""); err != nil {
			return err
		}

		if err := text.appendLinesAtIndex(line, yOffset); err != nil {
			return err
		}

		return nil
	}

	// NOTE: B the line at the end of the line
	if xOffset == targetLine.GetBufferLength() {
		line := new(Line)
		if err := line.Init(""); err != nil {
			return err
		}

		targetLine := text.lines[yOffset]

		if err := text.appendLinesAtIndex(targetLine, yOffset); err != nil {
			return err
		}

		text.lines[yOffset+1] = line
		return nil
	}

	// NOTE: Breaking the line in middle of the line
	targetLineBufferSlice := targetLine.GetBufferAsSlice()

	targetLineHead, err := text.bufferToLine(targetLineBufferSlice[:xOffset])
	if err != nil {
		return err
	}

	targetLineTail, err := text.bufferToLine(targetLineBufferSlice[xOffset:])
	if err != nil {
		return err
	}

	if err := text.appendLinesAtIndex(targetLineHead, yOffset); err != nil {
		return err
	}

	text.lines[yOffset+1] = targetLineTail
	return nil
}

// Handle combining two lines into one. The current line specified by the cursors yOffset will be appended
// to the end of the line above. The param lineStepDown inverts this logic by appending the line bellow
// to the current line specified by the cursors yOffset
func (text *Text) CombineLine(cursor *Cursor, lineStepDown bool) error {
	yOffset := cursor.GetOffsetY()

	if lineStepDown {
		yOffset += 1
	}

	if yOffset < 1 {
		return errors.New("text: invalid y (vertical) negative or out of bound offset requested to combine")
	}

	if yOffset > len(text.lines) {
		return errors.New("text: invalid y (vertical) out of bound offset requested to combine")
	}

	text.modified = true

	currentLineBuffer := text.lines[yOffset].GetBufferAsSlice()
	targetLineBuffer := text.lines[yOffset-1].GetBufferAsSlice()
	combinedLineBuffer := append(targetLineBuffer, currentLineBuffer...)

	combinedLine, err := text.bufferToLine(combinedLineBuffer)
	if err != nil {
		return err
	}

	text.lines[yOffset-1] = combinedLine

	linesHead := text.lines[:yOffset]
	linesTail := text.lines[yOffset+1:]
	text.lines = append(linesHead, linesTail...)

	return nil
}

// Helper function to for creating line structures from line buffers
func (text *Text) bufferToLine(lineBuffer []rune) (*Line, error) {
	builder := strings.Builder{}

	for _, char := range lineBuffer {
		if _, err := builder.WriteRune(char); err != nil {
			return nil, err
		}
	}

	line := new(Line)
	if err := line.Init(builder.String()); err != nil {
		return nil, err
	}

	return line, nil
}

// Helper funcation to insert lines to line container at given index
func (text *Text) appendLinesAtIndex(line *Line, index int) error {
	if index < 0 || index > len(text.lines) {
		return errors.New("text: invalid index to append lines container")
	}

	if index == len(text.lines) {
		text.lines = append(text.lines, line)
		return nil
	}

	linesHead := text.lines[:index+1]
	linesTail := text.lines[index:]

	text.lines = append(linesHead, linesTail...)
	text.lines[index] = line
	return nil
}

// Return a character based on the given x (horizontal) and y (vertical) offsets
func (text *Text) GetCharacterByOffsets(xOffset int, yOffset int) (rune, error) {
	if yOffset < 0 {
		return 0, errors.New("text: invalid y (vertical) negative offset requested to get")
	}

	if yOffset > len(text.lines) {
		return 0, errors.New("text: invalid y (vertical) out of bound offset requested to get")
	}

	targetLine := text.lines[yOffset]

	char, err := targetLine.GetBufferCharacterByOffset(xOffset)
	if err != nil {
		return 0, err
	}

	return char, nil
}

// Return a character based on the x (horizontal) and y (vertical) specified by the given cursor
func (text *Text) GetCharacterByCursor(cursor *Cursor) (rune, error) {
	return text.GetCharacterByOffsets(cursor.GetOffsetX(), cursor.GetOffsetY())
}

// Return the text in form of single string
func (text *Text) GetTextAsString() (*string, error) {
	builder := strings.Builder{}

	lineSeparator := "\n"
	if text.config.UsePlatformSpecificEndOfLineSequence {
		os := runtime.GOOS
		switch os {
		case "windows":
			lineSeparator = "\r\n"
		case "darwin":
			lineSeparator = "\n"
		case "linux":
			lineSeparator = "\n"
		default:
			lineSeparator = "\n"
		}
	}

	for index, line := range text.lines {
		if _, err := builder.WriteString(*line.GetBufferAsString()); err != nil {
			return nil, err
		}

		if index+1 < len(text.lines) {
			if _, err := builder.WriteString(lineSeparator); err != nil {
				return nil, err
			}
		}
	}

	builderText := builder.String()
	return &builderText, nil
}

// Return a bool value indicating if the current text differs from the persistent text
func (text *Text) IsModified() bool {
	return text.modified
}

// Reset the modification state to indicate that the current and persistent text are the same
func (text *Text) ResetModificationState() error {
	text.modified = false
	return nil
}

// Return the end-of-line sequence name (CRLF/LF)
func (text *Text) GetEndOfLineSequenceName() string {
	return text.endOfLineSequence
}

// A structure containing the configuration for the text structure
type TextConfig struct {
	UsePlatformSpecificEndOfLineSequence bool `json:"use-platform-specific-eol-sequence"`
}

// Return a new isntance of the text configuration with default values
func CreateDefaultTextConfig() TextConfig {
	return TextConfig{
		UsePlatformSpecificEndOfLineSequence: true,
	}
}
