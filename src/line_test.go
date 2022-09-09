package main

import "testing"

func TestLineShouldInitializeForNonEmptyString(t *testing.T) {
	line := new(Line)

	if err := line.Init("Valid string value"); err != nil {
		t.Fail()
	}
}

func TestLineShouldInitializeForEmptyString(t *testing.T) {
	line := new(Line)

	if err := line.Init(""); err != nil {
		t.Fail()
	}
}

func TestLineBufferShouldHaveCorrectLength(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	if line.GetBufferLength() != len(input) {
		t.Fail()
	}
}

func TestLineShouldProduceSameOutputAsString(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != input {
		t.Fail()
	}
}

func TestLineShouldProduceSameOutputAsRuneSlice(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	runeSlice := line.GetBufferAsSlice()

	for index, char := range runeSlice {
		if rune(input[index]) != char {
			t.Fail()
		}
	}
}

func TestLineShouldInsertCharacterAtValidPosition(t *testing.T) {
	line := new(Line)

	input := "Vali string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(4, 0, CreateConsoleMockup()); err != nil {
		t.Fail()
	}

	if err := line.InsertBufferCharacter('d', cursor); err != nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != "Valid string" {
		t.Fail()
	}
}

func TestLineShouldNotInsertCharacterAtInvalidPosition(t *testing.T) {
	line := new(Line)

	input := "Vali string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(90, 0, CreateConsoleMockup()); err != nil {
		t.Fail()
	}

	if err := line.InsertBufferCharacter('d', cursor); err == nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != "Vali string" {
		t.Fail()
	}
}

func TestLineShouldRemoveCharacterAtValidPosition(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(4, 0, CreateConsoleMockup()); err != nil {
		t.Fail()
	}

	if err := line.RemoveBufferCharacter(cursor); err != nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != "Vali string" {
		t.Fail()
	}
}

func TestLineShouldNotRemoveCharacterAtInvalidPosition(t *testing.T) {
	line := new(Line)

	input := "Vali string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(90, 0, CreateConsoleMockup()); err != nil {
		t.Fail()
	}

	if err := line.RemoveBufferCharacter(cursor); err == nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != "Vali string" {
		t.Fail()
	}
}

func TestLineShouldReturCharacterAtValidPosition(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(4, 0, CreateConsoleMockup()); err != nil {
		t.Fail()
	}

	char, err := line.GetBufferCharacterByCursor(cursor)
	if err != nil {
		t.Fail()
	}

	if char != 'd' {
		t.Fail()
	}
}

func TestLineShouldNotReturCharacterAtInvalidPosition(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(90, 0, CreateConsoleMockup()); err != nil {
		t.Fail()
	}

	char, err := line.GetBufferCharacterByCursor(cursor)
	if err == nil {
		t.Fail()
	}

	if char != 0 {
		t.Fail()
	}
}
