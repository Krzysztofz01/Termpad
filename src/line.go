package main

import (
	"errors"
	"strings"
)

// Structure representing a single line of text in the editor
type Line struct {
	buffer []rune
}

// Line structure initialization funcation
func (line *Line) Init(stringLine string) error {
	line.buffer = make([]rune, len(stringLine))

	for index, char := range stringLine {
		line.buffer[index] = char
	}

	return nil
}

// Return the line (buffer) length
func (line *Line) GetBufferLength() int {
	return len(line.buffer)
}

// Return the line in string representation
func (line *Line) GetBufferAsString() *string {
	builder := strings.Builder{}

	for _, char := range line.buffer {
		builder.WriteRune(char)
	}

	lineString := builder.String()
	return &lineString
}

// Return the line in rune slice representation
func (line *Line) GetBufferAsSlice() []rune {
	return line.buffer
}

// Insert a given rune at the position specified by the given cursor
func (line *Line) InsertBufferCharacter(char rune, cursor *Cursor) error {
	xOffset := cursor.GetOffsetX()
	if xOffset < 0 {
		return errors.New("line: invalid x (horizontal) negative offset requested to insert")
	}

	if xOffset > len(line.buffer) {
		return errors.New("line: invalid x (horizontal) out of bound offset requested to insert")
	}

	if len(line.buffer) == xOffset {
		line.buffer = append(line.buffer, char)
		return nil
	}

	bufferHead := line.buffer[:xOffset+1]
	bufferTail := line.buffer[xOffset:]

	line.buffer = append(bufferHead, bufferTail...)
	line.buffer[xOffset] = char

	return nil
}

// Remove a rune at the position specified by the given cursor
func (line *Line) RemoveBufferCharacter(cursor *Cursor) error {
	xOffset := cursor.GetOffsetX()
	if xOffset < 0 {
		return errors.New("line: invalid x (horizontal) negative offset requested to remove")
	}

	if xOffset > len(line.buffer) {
		return errors.New("line: invalid x (horizontal) out of bound offset requested to remove")
	}

	bufferHead := line.buffer[:xOffset]
	bufferTail := line.buffer[xOffset+1:]

	line.buffer = append(bufferHead, bufferTail...)

	return nil
}

// Return the rune at the position specified by the given offset
func (line *Line) GetBufferCharacterByOffset(xOffset int) (rune, error) {
	if xOffset < 0 {
		return 0, errors.New("line: invalid x (horizontal) negative offset requested to get")
	}

	if xOffset > len(line.buffer) {
		return 0, errors.New("line: invalid x (horizontal) out of bound offset requested to get")
	}

	targetChar := line.buffer[xOffset]
	return targetChar, nil
}

// Return the rune at the position specified by the given cursor
func (line *Line) GetBufferCharacterByCursor(cursor *Cursor) (rune, error) {
	return line.GetBufferCharacterByOffset(cursor.GetOffsetX())
}
