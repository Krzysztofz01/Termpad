package main

import "errors"

// TODO: Validate if padding is not greate than the overall size. This can lead to crashed on terminal resize.

// Structure representing the console display. It contains the calculated and fixed horizontal and vertical boundaries
type Display struct {
	height              int
	width               int
	xCalculatedBoundary int
	yCalculatedBoundary int
	xFixedBoundary      int
	yFixedBoundary      int
	cursor              *Cursor
	console             Console
}

// Display structure initialization function
func (display *Display) Init(cursor *Cursor, padding *Padding, console Console) error {
	display.xCalculatedBoundary = 0
	display.yCalculatedBoundary = 0
	display.xFixedBoundary = 0
	display.yFixedBoundary = 0

	// NOTE: There is currently no support for ,,none-default-console'' dimensions like top, left
	if padding != nil {
		display.xFixedBoundary = padding.GetRightPadding()
		display.yFixedBoundary = padding.GetBottomPadding()
	}

	if cursor == nil {
		return errors.New("display: invalid cursor struct reference")
	}

	display.cursor = cursor

	if console == nil {
		return errors.New("display: invalid internal console api contract implementation")
	}

	display.console = console

	if err := display.Resize(); err != nil {
		return err
	}

	return nil
}

// Function is used to recalculate the size and all boundaries/offsets of the display (currently visible content)
func (display *Display) Resize() error {
	width, height := display.console.GetSize()

	if width <= 0 {
		return errors.New("display: invalid display width value")
	}

	if height <= 0 {
		return errors.New("display: invalid display height value")
	}

	display.width = width
	display.height = height

	xOffset := display.cursor.GetOffsetX()
	yOffset := display.cursor.GetOffsetY()

	// NOTE: Right side overflow
	for xOffset > display.xCalculatedBoundary+display.width {
		display.xCalculatedBoundary += 1
	}

	// NOTE: Left side overflow
	for xOffset < display.xCalculatedBoundary {
		display.xCalculatedBoundary -= 1
	}

	// NOTE: Top side overflow
	for yOffset > display.yCalculatedBoundary+display.height {
		display.yCalculatedBoundary += 1
	}

	// NOTE: Down side overflow
	for yOffset < display.yCalculatedBoundary {
		display.yCalculatedBoundary -= 1
	}

	return nil
}

// Return a bool value indicating whether the console size specified by the underlying console API has changed (not the size of the display)
func (display *Display) HasSizeChanged() bool {
	width, height := display.console.GetSize()

	if display.width != width {
		return true
	}

	if display.height != height {
		return true
	}

	return false
}

// Return the full width and height of the display, which is the raw size deriving from the underlying console API
func (display *Display) GetFullDisplaySize() (int, int) {
	return display.width, display.height
}

// Return the width and height provided for the text. The sizes are affected by the specified display padding
func (display *Display) GetTextDisplaySize() (int, int) {
	return display.width - display.xFixedBoundary, display.height - display.yFixedBoundary
}

// Return the x (horizontal) display padding, specified on display initialization
func (display *Display) GetXOffsetPadding() int {
	return display.xFixedBoundary
}

// Return the y (vertical) display padding, specified on display initialization
func (display *Display) GetYOffsetPadding() int {
	return display.yFixedBoundary
}

// Return the x (horizontal) display offset shift, calculated from the display size and cursor position
func (display *Display) GetXOffsetShift() int {
	return display.xCalculatedBoundary
}

// Return the y (vertical) display offset shift, calculated frmo the display size and cursor position
func (display *Display) GetYOffsetShift() int {
	return display.yCalculatedBoundary
}

// Return a bool value indicating whether the cursor is currenlty ,,visible'' according to the offsets (boundaries)
func (display *Display) CursorInBoundries() bool {
	xOffset := display.cursor.GetOffsetX()
	yOffset := display.cursor.GetOffsetY()

	// NOTE: Right side overflow
	if xOffset > display.xCalculatedBoundary+display.width {
		return false
	}

	// NOTE: Left side overflow
	if display.xCalculatedBoundary > 0 && xOffset < display.xCalculatedBoundary {
		return false
	}

	// NOTE: Top side overflow
	if yOffset > display.yCalculatedBoundary+display.height {
		return false
	}

	// NOTE: Down side overflow
	if display.yCalculatedBoundary > 0 && yOffset < display.yCalculatedBoundary {
		return false
	}

	return true
}
