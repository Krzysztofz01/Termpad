package main

import "errors"

// TODO: The paddingFalback is indicating if the padding is greater than the size. The logical to handle such
// situation can be implemented later (during widgets implementation)
// TODO: Add unit-test for redrawing and size change logic. This task will require to implement a more precies
// behavior to the mock console and will require mock text structures

// Structure representing the console display. It contains the calculated and fixed horizontal and vertical boundaries
type Display struct {
	height              int
	width               int
	xCalculatedBoundary int
	yCalculatedBoundary int
	paddingFallback     bool
	padding             *Padding
	cursor              *Cursor
	console             Console
}

// Display structure initialization function
func (display *Display) Init(cursor *Cursor, padding *Padding, console Console) error {
	display.xCalculatedBoundary = 0
	display.yCalculatedBoundary = 0

	pTop := 0
	pBottom := 0
	pLeft := 0
	pRight := 0

	if padding != nil {
		pTop = padding.GetTopPadding()
		pBottom = padding.GetBottomPadding()
		pLeft = padding.GetLeftPadding()
		pRight = padding.GetRightPadding()
	}

	display.padding = new(Padding)
	if err := display.padding.Init(pTop, pBottom, pLeft, pRight); err != nil {
		return err
	}

	if cursor == nil {
		return errors.New("display: invalid cursor struct reference")
	}

	display.cursor = cursor

	if console == nil {
		return errors.New("display: invalid internal console api contract implementation")
	}

	display.console = console
	width, hight := display.console.GetSize()

	if err := display.Resize(width, hight); err != nil {
		return err
	}

	return nil
}

// Function is used to recalculate the size and all boundaries/offsets of the display, according to the given width and height
func (display *Display) Resize(width int, height int) error {
	if width <= 0 {
		return errors.New("display: invalid display width value")
	}

	if height <= 0 {
		return errors.New("display: invalid display height value")
	}

	xPadding := display.GetXOffsetPadding()
	yPadding := display.GetYOffsetPadding()

	display.width = width
	display.height = height

	if display.width <= xPadding || display.height <= yPadding {
		display.paddingFallback = true
	} else {
		display.paddingFallback = false
	}

	return display.RecalculateBoundaries()
}

// Function is used to recalculate the boundaries based on the cursor position and current display size
func (display *Display) RecalculateBoundaries() error {
	xOffset := display.cursor.GetOffsetX()
	yOffset := display.cursor.GetOffsetY()

	display.xCalculatedBoundary = 0
	display.yCalculatedBoundary = 0

	// NOTE: Right side overflow
	for xOffset+1 >= display.width+display.xCalculatedBoundary {
		display.xCalculatedBoundary += 1
	}

	// NOTE: Left side overflow
	for display.xCalculatedBoundary > 0 && xOffset < display.xCalculatedBoundary {
		display.xCalculatedBoundary -= 1
	}

	// NOTE: Top side overflow
	for yOffset+1 >= display.height+display.yCalculatedBoundary {
		display.yCalculatedBoundary += 1
	}

	// NOTE: Down side overflow
	for display.yCalculatedBoundary > 0 && yOffset < display.yCalculatedBoundary {
		display.yCalculatedBoundary -= 1
	}

	return nil
}

// Return a bool value indicating whether the console size specified by the given width and height has changed (not the size of the display)
func (display *Display) HasSizeChanged(width int, height int) bool {
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
	width := display.width - display.padding.GetLeftPadding() - display.padding.GetRightPadding()
	height := display.height - display.padding.GetTopPadding() - display.padding.GetBottomPadding()

	return width, height
}

// Return the x (horizontal) display padding (left and right), specified on display initialization
func (display *Display) GetXOffsetPadding() int {
	return display.padding.GetLeftPadding() + display.padding.GetRightPadding()
}

// Return the left x (horizontal) display padding, specified on display initialization
func (display *Display) GetXLeftOffsetPadding() int {
	return display.padding.GetLeftPadding()
}

// Return the right x (horizontal) display padding, specified on display initialization
func (display *Display) GetXRightOffsetPadding() int {
	return display.padding.GetRightPadding()
}

// Return the y (vertical) display padding (top and bottom), specified on display initialization
func (display *Display) GetYOffsetPadding() int {
	return display.padding.GetTopPadding() + display.padding.GetBottomPadding()
}

// Return the y (vertical) display padding, specified on display initialization
func (display *Display) GetYTopOffsetPadding() int {
	return display.padding.GetTopPadding()
}

// Return the y (vertical) display padding, specified on display initialization
func (display *Display) GetYBottomOffsetPadding() int {
	return display.padding.GetBottomPadding()
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
	if xOffset+1 >= display.width+display.xCalculatedBoundary {
		return false
	}

	// NOTE: Left side overflow
	if display.xCalculatedBoundary > 0 && xOffset < display.xCalculatedBoundary {
		return false
	}

	// NOTE: Top side overflow
	if yOffset+1 >= display.height+display.yCalculatedBoundary {
		return false
	}

	// NOTE: Down side overflow
	if display.yCalculatedBoundary > 0 && yOffset < display.yCalculatedBoundary {
		return false
	}

	return true
}

// Request a render of all changes to the screen of the underlying console API
func (display *Display) RenderChanges() error {
	xDiff := display.xCalculatedBoundary
	yDiff := display.yCalculatedBoundary

	if err := display.cursor.CorrectUnderlyingConsolePositionDifference(xDiff, yDiff); err != nil {
		return err
	}

	return display.console.Commit()
}

// Function is rewriting text changes to the underlying console API screen, according to the display boundaries. All lines are affected
// TODO: The tHeight can be greater than dHeight - we can avoid redundand off-display rendering. We can implement this by adding range
// based functions in the display struct
func (display *Display) RedrawTextFull(text *Text) error {
	yTextLength := text.GetLineCount()

	xlPadding := display.padding.GetLeftPadding()
	xrPadding := display.padding.GetRightPadding()
	ytPadding := display.padding.GetTopPadding()
	ybPadding := display.padding.GetBottomPadding()

	for ycIndex := ytPadding; ycIndex < display.height-ybPadding; ycIndex += 1 {
		ytIndex := ycIndex + display.yCalculatedBoundary

		if ytIndex < yTextLength {
			xtLength, err := text.GetLineLengthByOffset(ytIndex)
			if err != nil {
				return err
			}

			for xcIndex := xlPadding; xcIndex < display.width-xrPadding; xcIndex += 1 {
				xtIndex := xcIndex + display.xCalculatedBoundary

				var char rune = ' '

				if xtIndex < xtLength {
					char, err = text.GetCharacterByOffsets(xtIndex, ytIndex)
					if err != nil {
						return err
					}
				}

				if err := display.console.InsertCharacter(xcIndex, ycIndex, char); err != nil {
					return err
				}
			}

			continue
		}

		for xcIndex := xlPadding; xcIndex < display.width-xrPadding; xcIndex += 1 {
			if err := display.console.InsertCharacter(xcIndex, ycIndex, ' '); err != nil {
				return err
			}
		}
	}

	return nil
}

// Function is rewriting text changes to the underlying console API screen, according to the display boundaries. Only the line specified by the cursor is affected.
func (display *Display) RedrawTextLine(text *Text, fullRedrawFallback bool) error {
	if !display.CursorInBoundries() && fullRedrawFallback {
		return display.RedrawTextFull(text)
	}

	ytOffset := display.cursor.GetOffsetY()
	ycOffset := ytOffset - display.yCalculatedBoundary

	tWidth, err := text.GetLineLengthByOffset(ytOffset)
	if err != nil {
		return err
	}

	xlPadding := display.padding.GetLeftPadding()
	xrPadding := display.padding.GetRightPadding()

	for xcIndex := xlPadding; xcIndex < display.width-xrPadding; xcIndex += 1 {
		xtIndex := xcIndex + display.xCalculatedBoundary

		var char rune = ' '
		if xtIndex < tWidth {
			char, err = text.GetCharacterByOffsets(xtIndex, ytOffset)
			if err != nil {
				return err
			}
		}

		if err := display.console.InsertCharacter(xcIndex, ycOffset, char); err != nil {
			return err
		}
	}

	return nil
}

// Function is rewriting text changes to the underlying console API screen, according to the display boundaries. All lines (including the current) below the cursor are affected.
// TODO: The ycIndex < ytLength condition prevents the overwrting of previous screen data (Edit: ycIndex patched to ytIndex, does this problem still exist?)
func (display *Display) RedrawTextBelow(text *Text, fullRedrawFallback bool) error {
	if !display.CursorInBoundries() && fullRedrawFallback {
		return display.RedrawTextFull(text)
	}

	ytLength := text.GetLineCount()

	xlPadding := display.padding.GetLeftPadding()
	xrPadding := display.padding.GetRightPadding()
	ybPadding := display.padding.GetBottomPadding()

	for ycIndex := display.cursor.GetOffsetY() - display.yCalculatedBoundary; ycIndex < display.height-ybPadding; ycIndex += 1 {
		ytIndex := ycIndex + display.yCalculatedBoundary

		if ytIndex < ytLength {
			xtLength, err := text.GetLineLengthByOffset(ytIndex)
			if err != nil {
				return err
			}

			for xcIndex := xlPadding; xcIndex < display.width-xrPadding; xcIndex += 1 {
				xtIndex := xcIndex + display.xCalculatedBoundary

				var char rune = ' '

				if xtIndex < xtLength {
					char, err = text.GetCharacterByOffsets(xtIndex, ytIndex)
					if err != nil {
						return err
					}
				}

				if err := display.console.InsertCharacter(xcIndex, ycIndex, char); err != nil {
					return err
				}
			}

			continue
		}

		for xcIndex := xlPadding; xcIndex < display.width-xrPadding; xcIndex += 1 {
			if err := display.console.InsertCharacter(xcIndex, ycIndex, ' '); err != nil {
				return err
			}
		}
	}

	return nil
}
