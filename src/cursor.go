package main

import (
	"errors"
	"strings"
)

// Structure representing the current position of the editor cursor
type Cursor struct {
	xOffset int
	yOffset int

	console Console
	config  *CursorConfig
}

// Editor cursor structure initialization function
func (cursor *Cursor) Init(xOffset int, yOffset int, console Console, cursorConfig *CursorConfig) error {
	if cursorConfig == nil {
		defaultConfig := CreateDefaultCursorConfig()
		cursor.config = &defaultConfig
	} else {
		cursor.config = cursorConfig
	}

	if console == nil {
		return errors.New("cursor: invalid console reference")
	}

	cursor.console = console

	if err := cursor.applyCursorStyle(); err != nil {
		return err
	}

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

// Helper function used to apply cursor style to underlying console API
func (cursor *Cursor) applyCursorStyle() error {
	style := strings.ToLower(cursor.config.CursorStyle)

	if cursor.config.UseAnimations {
		switch style {
		case "bar":
			return cursor.console.SetCursorStyle(BarCursorDynamic)
		case "block":
			return cursor.console.SetCursorStyle(BlockCursorDynamic)
		case "line":
			return cursor.console.SetCursorStyle(LineCursorDynamic)
		default:
			return errors.New("cursor: invalid cursor config style name")
		}
	}

	switch style {
	case "bar":
		return cursor.console.SetCursorStyle(BarCursorStatic)
	case "block":
		return cursor.console.SetCursorStyle(BlockCursorStatic)
	case "line":
		return cursor.console.SetCursorStyle(LineCursorStatic)
	default:
		return errors.New("cursor: invalid cursor config style name")
	}
}

// A structure containing the configuration for the cursor structure
type CursorConfig struct {
	// NOTE: Available options: "bar", "block", "line"
	CursorStyle   string `json:"cursor-style"`
	UseAnimations bool   `json:"use-animations"`
}

// Return a new isntance of the cursor configuration with default values
func CreateDefaultCursorConfig() CursorConfig {
	return CursorConfig{
		CursorStyle:   "bar",
		UseAnimations: false,
	}
}
