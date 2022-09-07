package main

import (
	"testing"
)

func TestDisplayShouldInitialzieForValidSizeAndCursor(t *testing.T) {
	cursor := new(Cursor)
	if err := cursor.Init(0, 0); err != nil {
		t.Fail()
	}

	width := 10
	height := 10

	display := new(Display)
	if err := display.Init(width, height, cursor); err != nil {
		t.Fail()
	}
}

func TestDisplayShouldNotInitialzieForInvalidSize(t *testing.T) {
	cursor := new(Cursor)
	if err := cursor.Init(0, 0); err != nil {
		t.Fail()
	}

	width := 0
	height := -5

	display := new(Display)
	if err := display.Init(width, height, cursor); err == nil {
		t.Fail()
	}
}

func TestDisplayShouldCalculateCorrectHeightBoundaries(t *testing.T) {
	cursor := new(Cursor)
	if err := cursor.Init(5, 12); err != nil {
		t.Fail()
	}

	width := 10
	height := 10

	display := new(Display)
	if err := display.Init(width, height, cursor); err != nil {
		t.Fail()
	}

	if display.GetXOffsetShift() != 0 {
		t.Fail()
	}

	if display.GetYOffsetShift() != 2 {
		t.Fail()
	}
}

func TestDisplayShouldCalculateCorrectWidthBoundaries(t *testing.T) {
	cursor := new(Cursor)
	if err := cursor.Init(15, 7); err != nil {
		t.Fail()
	}

	width := 10
	height := 10

	display := new(Display)
	if err := display.Init(width, height, cursor); err != nil {
		t.Fail()
	}

	if display.GetXOffsetShift() != 5 {
		t.Fail()
	}

	if display.GetYOffsetShift() != 0 {
		t.Fail()
	}
}

func TestDisplayShouldCalculateCorrectHeightAndWidthBoundaries(t *testing.T) {
	cursor := new(Cursor)
	if err := cursor.Init(21, 17); err != nil {
		t.Fail()
	}

	width := 10
	height := 10

	display := new(Display)
	if err := display.Init(width, height, cursor); err != nil {
		t.Fail()
	}

	if display.GetXOffsetShift() != 11 {
		t.Fail()
	}

	if display.GetYOffsetShift() != 7 {
		t.Fail()
	}
}

func TestDisplayShouldIndicateInBoundries(t *testing.T) {
	cursor := new(Cursor)
	if err := cursor.Init(0, 0); err != nil {
		t.Fail()
	}

	width := 10
	height := 10

	display := new(Display)
	if err := display.Init(width, height, cursor); err != nil {
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
	cursor := new(Cursor)
	if err := cursor.Init(0, 0); err != nil {
		t.Fail()
	}

	width := 10
	height := 10

	display := new(Display)
	if err := display.Init(width, height, cursor); err != nil {
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
	cursor := new(Cursor)
	if err := cursor.Init(0, 0); err != nil {
		t.Fail()
	}

	width := 10
	height := 10

	display := new(Display)
	if err := display.Init(width, height, cursor); err != nil {
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
	cursor := new(Cursor)
	if err := cursor.Init(0, 0); err != nil {
		t.Fail()
	}

	width := 10
	height := 10

	display := new(Display)
	if err := display.Init(width, height, cursor); err != nil {
		t.Fail()
	}

	if err := cursor.SetOffsets(25, 13); err != nil {
		t.Fail()
	}

	if display.CursorInBoundries() {
		t.Fail()
	}
}
