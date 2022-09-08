package main

import "errors"

// FIXME: One unit padding on end of the lines
// Structure representing the current position of the editor cursor
type Cursor struct {
	xOffset int
	yOffset int

	console Console
}

// Editor cursor structure initialization function
func (cursor *Cursor) Init(xOffset int, yOffset int, console Console) error {
	if console == nil {
		return errors.New("cursor: invalid console reference")
	}

	cursor.console = console
	return cursor.SetOffsets(xOffset, yOffset)
}

// Return the X (horizontal) offset of the current editor cursor position
func (cursor *Cursor) GetOffsetX() int {
	return cursor.xOffset
}

// Return the Y (vertical) offset of the current editor cursor position
func (cursor *Cursor) GetOffsetY() int {
	return cursor.yOffset
}

// Set the X (horizontal) offset of the current editor cursor position. Also apply the offset to cursor of the underlying console API
func (cursor *Cursor) SetOffsetX(xOffset int) error {
	if xOffset < 0 {
		return errors.New("cursor: invalid x (horizontal) offset position")
	}

	cursor.xOffset = xOffset

	if err := cursor.console.SetCursorPosition(cursor.xOffset, cursor.yOffset); err != nil {
		return err
	}

	return nil
}

// Set the Y (vertical) offset of the current editor cursor position. Also apply the offset to cursor of the underlying console API
func (cursor *Cursor) SetOffsetY(yOffset int) error {
	if yOffset < 0 {
		return errors.New("cursor: invalid y (vertical) offset position")
	}

	cursor.yOffset = yOffset

	if err := cursor.console.SetCursorPosition(cursor.xOffset, cursor.yOffset); err != nil {
		return err
	}

	return nil
}

// Set the X (horizontal) and Y (vertical) offsets of the current editor cursor position. Also apply the offsets to cursor of the underlying console API
func (cursor *Cursor) SetOffsets(xOffset int, yOffset int) error {
	if err := cursor.SetOffsetX(xOffset); err != nil {
		return err
	}

	if err := cursor.SetOffsetY(yOffset); err != nil {
		return err
	}

	return nil
}
