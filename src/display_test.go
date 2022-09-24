package main

import (
	"testing"
)

func TestDisplayShouldInitialzieForValidParams(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}
}

func TestDisplayShouldNotInitialzieForInvalidParams(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, nil, nil); err == nil {
		t.Fail()
	}
}

func TestDisplayShouldResizeForValidUpdatedSize(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if err := display.Resize(20, 30); err != nil {
		t.Fail()
	}

	width, height := display.GetFullDisplaySize()

	if width != 20 {
		t.Fail()
	}

	if height != 30 {
		t.Fail()
	}
}

func TestDisplayShouldNotResizeForInvalidUpdatedSize(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if err := display.Resize(-1, 0); err == nil {
		t.Fail()
	}
}

func TestDisplayHasSizeChangedShouldCorrectlyIndicateChangedSize(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if !display.HasSizeChanged(20, 30) {
		t.Fail()
	}
}

func TestDisplayHasSizeChangedShouldCorrectlyIndicateUnchangedSize(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if display.HasSizeChanged(10, 10) {
		t.Fail()
	}
}

func TestDisplayShouldCalculateCorrectHeightBoundaries(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(5, 12, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if display.GetXOffsetShift() != 0 {
		t.Fail()
	}

	if display.GetYOffsetShift() != 4 {
		t.Fail()
	}
}

func TestDisplayShouldCalculateCorrectWidthBoundaries(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(15, 7, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if display.GetXOffsetShift() != 7 {
		t.Fail()
	}

	if display.GetYOffsetShift() != 0 {
		t.Fail()
	}
}

func TestDisplayShouldCalculateCorrectHeightAndWidthBoundaries(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(21, 17, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if display.GetXOffsetShift() != 13 {
		t.Fail()
	}

	if display.GetYOffsetShift() != 9 {
		t.Fail()
	}
}

func TestDisplayShouldIndicateInBoundries(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if err := cursor.SetOffsets(6, 3); err != nil {
		t.Fail()
	}

	if !display.CursorInBoundries() {
		t.Fail()
	}
}

func TestDisplayShouldIndicateOutOfBoundriesForHeight(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if err := cursor.SetOffsets(6, 14); err != nil {
		t.Fail()
	}

	if display.CursorInBoundries() {
		t.Fail()
	}
}

func TestDisplayShouldIndicateOutOfBoundriesForWidth(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if err := cursor.SetOffsets(16, 4); err != nil {
		t.Fail()
	}

	if display.CursorInBoundries() {
		t.Fail()
	}
}

func TestDisplayShouldIndicateOutOfBoundriesForHeightAndWidth(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if err := cursor.SetOffsets(25, 13); err != nil {
		t.Fail()
	}

	if display.CursorInBoundries() {
		t.Fail()
	}
}

func TestDisplayShouldReturnFullSize(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	width, height := display.GetFullDisplaySize()

	if width != 10 {
		t.Fail()
	}

	if height != 10 {
		t.Fail()
	}
}

func TestDisplayShouldReturnCorrectTextSizeWhenPaddingIsApplied(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, &DisplayConfig{LineNumerationEnabled: true}); err != nil {
		t.Fail()
	}

	width, height := display.GetTextDisplaySize()

	if width != 10-NumerationWidth {
		t.Fail()
	}

	if height != 10-MenuHeight {
		t.Fail()
	}
}

func TestDisplayShouldReturnCorrectTextSizeWhenNoPaddingApplied(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	width, height := display.GetTextDisplaySize()

	if width != 10 {
		t.Fail()
	}

	if height != 10-MenuHeight {
		t.Fail()
	}
}

func TestDisplayShouldReturnCorrectPaddingsWhenPaddingIsApplied(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, &DisplayConfig{LineNumerationEnabled: true}); err != nil {
		t.Fail()
	}

	if display.GetXLeftOffsetPadding() != NumerationWidth {
		t.Fail()
	}

	if display.GetXRightOffsetPadding() != 0 {
		t.Fail()
	}

	if display.GetYTopOffsetPadding() != 0 {
		t.Fail()
	}

	if display.GetYBottomOffsetPadding() != MenuHeight {
		t.Fail()
	}
}

func TestDisplayShouldReturnCorrectPaddingsWhenPaddingIsNotApplied(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	if display.GetXLeftOffsetPadding() != 0 {
		t.Fail()
	}

	if display.GetXRightOffsetPadding() != 0 {
		t.Fail()
	}

	if display.GetYTopOffsetPadding() != 0 {
		t.Fail()
	}

	if display.GetYBottomOffsetPadding() != MenuHeight {
		t.Fail()
	}
}

func TestDisplayShouldReturnCorrectTextHorizontalRangeWhenPaddingApplied(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, &DisplayConfig{LineNumerationEnabled: true}); err != nil {
		t.Fail()
	}

	xMinIndex, xCount := display.GetTextDisplayHorizontalRange()

	if xMinIndex != NumerationWidth {
		t.Fail()
	}

	if xCount != 10-NumerationWidth {
		t.Fail()
	}
}

func TestDisplayShouldReturnCorrectTextHorizontalRangeWhenNoPaddingApplied(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	xMinIndex, xCount := display.GetTextDisplayHorizontalRange()

	if xMinIndex != 0 {
		t.Fail()
	}

	if xCount != 10 {
		t.Fail()
	}
}

func TestDisplayShouldReturnCorrectTextVerticalRangeWhenPaddingApplied(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, &DisplayConfig{LineNumerationEnabled: true}); err != nil {
		t.Fail()
	}

	yMinIndex, yCount := display.GetTextDisplayVerticalRange()

	if yMinIndex != 0 {
		t.Fail()
	}

	if yCount != 10-MenuHeight {
		t.Fail()
	}
}

func TestDisplayShouldReturnCorrectTextVerticalRangeWhenNoPaddingApplied(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, console, nil); err != nil {
		t.Fail()
	}

	yMinIndex, yCount := display.GetTextDisplayVerticalRange()

	if yMinIndex != 0 {
		t.Fail()
	}

	if yCount != 10-MenuHeight {
		t.Fail()
	}
}
