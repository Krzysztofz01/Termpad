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
	if err := cursor.Init(4, 0, CreateConsoleMockup(), nil); err != nil {
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
	if err := cursor.Init(90, 0, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := line.InsertBufferCharacter('d', cursor); err == nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != "Vali string" {
		t.Fail()
	}
}

func TestLineShouldRemoveCharacterHeadAtValidPositionStart(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(1, 0, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := line.RemoveBufferCharacterHead(cursor); err != nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != "alid string" {
		t.Fail()
	}
}

func TestLineShouldRemoveCharacterTailAtValidPositionStart(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(1, 0, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := line.RemoveBufferCharacterTail(cursor); err != nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != "Vlid string" {
		t.Fail()
	}
}

func TestLineShouldRemoveCharacterHeadAtValidPositionEnd(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(12, 0, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := line.RemoveBufferCharacterHead(cursor); err != nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != "Valid strin" {
		t.Fail()
	}
}

func TestLineShouldNotRemoveCharacterTailAtInvalidPositionEnd(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(12, 0, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := line.RemoveBufferCharacterTail(cursor); err == nil {
		t.Fail()
	}
}

func TestLineShouldRemoveCharacterHeadAtValidPositionMiddle(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(5, 0, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := line.RemoveBufferCharacterHead(cursor); err != nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != "Vali string" {
		t.Fail()
	}
}

func TestLineShouldRemoveCharacterTailAtValidPositionMiddle(t *testing.T) {
	line := new(Line)

	input := "Valid string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(5, 0, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := line.RemoveBufferCharacterTail(cursor); err != nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != "Validstring" {
		t.Fail()
	}
}

func TestLineShouldNotRemoveCharacterHeadAtInvalidPosition(t *testing.T) {
	line := new(Line)

	input := "Vali string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(90, 0, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := line.RemoveBufferCharacterTail(cursor); err == nil {
		t.Fail()
	}

	if *line.GetBufferAsString() != "Vali string" {
		t.Fail()
	}
}

func TestLineShouldNotRemoveCharacterTailAtInvalidPosition(t *testing.T) {
	line := new(Line)

	input := "Vali string"

	if err := line.Init(input); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(90, 0, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := line.RemoveBufferCharacterTail(cursor); err == nil {
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
	if err := cursor.Init(4, 0, CreateConsoleMockup(), nil); err != nil {
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
	if err := cursor.Init(90, 0, CreateConsoleMockup(), nil); err != nil {
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
