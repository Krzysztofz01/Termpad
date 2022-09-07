package main

import "errors"

// Structure representing the console display. It contains the boundries of what is/can be displayed
type Display struct {
	height    int
	width     int
	xBoundary int
	yBoundary int
	cursor    *Cursor
}

// Display structure initialization function
func (display *Display) Init(width int, height int, cursor *Cursor) error {
	display.xBoundary = 0
	display.yBoundary = 0

	if cursor == nil {
		return errors.New("display: invalid cursor struct reference")
	}

	display.cursor = cursor

	if err := display.Resize(width, height); err != nil {
		return err
	}

	return nil
}

// Function is used to change the size of target display and recalculate the offsets of ,,currently visible'' content
func (display *Display) Resize(width int, height int) error {
	if width <= 0 {
		return errors.New("display: invalid display width")
	}

	if height <= 0 {
		return errors.New("display: invalid display height")
	}

	display.width = width
	display.height = height

	xOffset := display.cursor.GetOffsetX()
	yOffset := display.cursor.GetOffsetY()

	// NOTE: Right side overflow
	for xOffset > display.xBoundary+display.width {
		display.xBoundary += 1
	}

	// NOTE: Left side overflow
	for xOffset < display.xBoundary {
		display.xBoundary -= 1
	}

	// NOTE: Top side overflow
	for yOffset > display.yBoundary+display.height {
		display.yBoundary += 1
	}

	// NOTE: Down side overflow
	for yOffset < display.yBoundary {
		display.yBoundary -= 1
	}

	return nil
}

// Return the x (horizontal) display offset
func (display *Display) GetXOffsetShift() int {
	return display.xBoundary
}

// Return the y (vertical) display offset
func (display *Display) GetYOffsetShift() int {
	return display.yBoundary
}

// Return a bool values indicating if a redraw is required according to the curent position
func (display *Display) CursorInBoundries() bool {
	xOffset := display.cursor.GetOffsetX()
	yOffset := display.cursor.GetOffsetY()

	// NOTE: Right side overflow
	if xOffset > display.xBoundary+display.width {
		return false
	}

	// NOTE: Left side overflow
	if xOffset < display.xBoundary {
		return false
	}

	// NOTE: Top side overflow
	if yOffset > display.yBoundary+display.height {
		return false
	}

	// NOTE: Down side overflow
	if yOffset < display.yBoundary {
		return false
	}

	return true
}
