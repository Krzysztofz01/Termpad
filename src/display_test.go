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
	if err := display.Init(cursor, nil, console); err != nil {
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
	if err := display.Init(cursor, nil, console); err != nil {
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
	if err := display.Init(cursor, nil, console); err != nil {
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
	if err := display.Init(cursor, nil, console); err != nil {
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
	if err := display.Init(cursor, nil, console); err != nil {
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
	if err := display.Init(cursor, nil, console); err != nil {
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
	if err := display.Init(cursor, nil, console); err != nil {
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
	if err := display.Init(cursor, nil, console); err != nil {
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
	if err := display.Init(cursor, nil, console); err != nil {
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
	if err := display.Init(cursor, nil, console); err != nil {
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
	if err := display.Init(cursor, nil, console); err != nil {
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
	if err := display.Init(cursor, nil, console); err != nil {
		t.Fail()
	}

	if err := cursor.SetOffsets(25, 13); err != nil {
		t.Fail()
	}

	if display.CursorInBoundries() {
		t.Fail()
	}
}

func TestDisplayShouldReturnCorrectSpecifiedPadding(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	padding := new(Padding)
	if err := padding.Init(0, 4, 0, 8); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, padding, console); err != nil {
		t.Fail()
	}

	if display.GetXOffsetPadding() != 8 {
		t.Fail()
	}

	if display.GetYOffsetPadding() != 4 {
		t.Fail()
	}
}

func TestDisplayShouldReturnCorrectUnspecifiedPadding(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, nil, console); err != nil {
		t.Fail()
	}

	if display.GetXOffsetPadding() != 0 {
		t.Fail()
	}

	if display.GetYOffsetPadding() != 0 {
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
	if err := display.Init(cursor, nil, console); err != nil {
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

func TestDisplayShouldReturnCorrectTextSizeForInitializedPadding(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	padding := new(Padding)
	if err := padding.Init(0, 4, 0, 8); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, padding, console); err != nil {
		t.Fail()
	}

	width, height := display.GetTextDisplaySize()

	if width != 2 {
		t.Fail()
	}

	if height != 6 {
		t.Fail()
	}
}

func TestDisplayShouldReturnCorrectTextSizeForUninitializedPadding(t *testing.T) {
	console := CreateConsoleMockup()

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, console, nil); err != nil {
		t.Fail()
	}

	display := new(Display)
	if err := display.Init(cursor, nil, console); err != nil {
		t.Fail()
	}

	width, height := display.GetTextDisplaySize()

	if width != 10 {
		t.Fail()
	}

	if height != 10 {
		t.Fail()
	}
}
