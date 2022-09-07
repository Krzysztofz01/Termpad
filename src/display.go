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
func (display *Display) Init(height int, width int, cursor *Cursor) error {
	display.xBoundary = 0
	display.yBoundary = 0

	if err := display.Resize(height, width); err != nil {
		return err
	}

	if cursor == nil {
		return errors.New("display: invalid cursor struct reference")
	}

	display.cursor = cursor
	return nil
}

// Function is used to change the size of target display and recalculate the offsets of ,,currently visible'' content
func (display *Display) Resize(height int, width int) error {
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
	if xOffset > display.xBoundary+display.width {
		display.xBoundary += 1
	}

	// NOTE: Left side overflow
	if xOffset < display.xBoundary {
		display.xBoundary -= 1
	}

	// NOTE: Top side overflow
	if yOffset > display.yBoundary+display.height {
		display.yBoundary += 1
	}

	// NOTE: Down side overflow
	if yOffset < display.yBoundary {
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