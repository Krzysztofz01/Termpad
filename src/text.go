package main

import (
	"errors"
	"strings"
)

// TODO: Fix references to cursor offsets via functions
// A structure representing the text, which is a container for the List structures
type Text struct {
	lines    []*Line
	modified bool

	// TODO: Implement preferences, preferences should contain information about usage of CRLF/LF
}

// Text structure initialization funcation
func (text *Text) Init(textString string, newFile bool) error {
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

// Return the lenght of the line based on given cursor position
func (text *Text) GetLineLength(cursor *Cursor) (int, error) {
	yOffset := cursor.yOffset

	if yOffset < 0 {
		return 0, errors.New("text: invalid y (vertical) negative offset requested to get")
	}

	if yOffset > len(text.lines) {
		return 0, errors.New("text: invalid y (vertical) out of bound offset requested to get")
	}

	targetLine := text.lines[yOffset]
	return targetLine.GetBufferLength(), nil
}

// Place a given character inside specific line at specific offset given by the cursor position
func (text *Text) InsertCharacter(char rune, cursor *Cursor) error {
	yOffset := cursor.yOffset

	if yOffset < 0 {
		return errors.New("text: invalid y (vertical) negative offset requested to insert")
	}

	if yOffset > len(text.lines) {
		return errors.New("text: invalid y (vertical) out of bound offset requested to insert")
	}

	targetLine := text.lines[yOffset]
	return targetLine.InsertBufferCharacter(char, cursor)
}

// Remove a character inside specific line at specific offset given by the cursor position
func (text *Text) RemoveCharacter(cursor *Cursor) error {
	yOffset := cursor.yOffset

	if yOffset < 0 {
		return errors.New("text: invalid y (vertical) negative offset requested to remove")
	}

	if yOffset > len(text.lines) {
		return errors.New("text: invalid y (vertical) out of bound offset requested to remove")
	}

	targetLine := text.lines[yOffset]
	return targetLine.RemoveBufferCharacter(cursor)
}

// Handle line inserting and line breaking
func (text *Text) InsertLine(cursor *Cursor) error {
	yOffset := cursor.yOffset
	xOffset := cursor.xOffset

	if yOffset < 0 {
		return errors.New("text: invalid y (vertical) negative offset requested to split")
	}

	if yOffset > len(text.lines) {
		return errors.New("text: invalid y (vertical) out of bound offset requested to split")
	}

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

// Return a character based on the line and offset specified by the given cursor position
func (text *Text) GetCharacter(cursor *Cursor) (rune, error) {
	yOffset := cursor.yOffset

	if yOffset < 0 {
		return 0, errors.New("text: invalid y (vertical) negative offset requested to get")
	}

	if yOffset > len(text.lines) {
		return 0, errors.New("text: invalid y (vertical) out of bound offset requested to get")
	}

	targetLine := text.lines[yOffset]

	char, err := targetLine.GetBufferCharacter(cursor)
	if err != nil {
		return 0, err
	}

	return char, nil
}

// Return the text in form of single string
func (text *Text) GetTextAsString() (*string, error) {
	builder := strings.Builder{}

	for index, line := range text.lines {
		if _, err := builder.WriteString(*line.GetBufferAsString()); err != nil {
			return nil, err
		}

		if index+1 < len(text.lines) {
			if _, err := builder.WriteRune('\n'); err != nil {
				return nil, err
			}
		}
	}

	builderText := builder.String()
	return &builderText, nil
}
