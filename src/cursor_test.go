package main

import "testing"

func TestCursorShouldInitialzeForValidDefaultOffsets(t *testing.T) {
	cursor := new(Cursor)

	if err := cursor.Init(0, 0, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}
}

func TestCursorShouldInitialzeForValidNotDefaultOffsets(t *testing.T) {
	cursor := new(Cursor)

	if err := cursor.Init(4, 7, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}
}

func TestCursorShouldNotInitializeForInvalidXOffset(t *testing.T) {
	cursor := new(Cursor)

	if err := cursor.Init(-2, 4, CreateConsoleMockup(), nil); err == nil {
		t.Fail()
	}
}

func TestCursorShouldNotInitializeForInvalidYOffset(t *testing.T) {
	cursor := new(Cursor)

	if err := cursor.Init(4, -5, CreateConsoleMockup(), nil); err == nil {
		t.Fail()
	}
}

func TestCursorShouldReturnCorrectXOffset(t *testing.T) {
	cursor := new(Cursor)

	expectedValue := 2

	if err := cursor.Init(expectedValue, 4, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	actualValue := cursor.GetOffsetX()

	if expectedValue != actualValue {
		t.Fail()
	}
}

func TestCursorShouldReturnCorrectYOffset(t *testing.T) {
	cursor := new(Cursor)

	expectedValue := 4

	if err := cursor.Init(2, expectedValue, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	actualValue := cursor.GetOffsetY()

	if expectedValue != actualValue {
		t.Fail()
	}
}
