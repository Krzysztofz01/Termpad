package main

import "testing"

func TestTextShouldInitializeValidCrLf(t *testing.T) {
	textContent := "First line\r\nSecond line\r\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}
}

func TestTextShouldInitializeValidLf(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}
}

func TestTextShouldInitializeValidNewFile(t *testing.T) {
	text := new(Text)
	if err := text.Init("", false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}
}

func TestTextShouldReturnCorrectLineCountWithSomeLines(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	if text.GetLineCount() != 3 {
		t.Fail()
	}
}

func TestTextShouldReturnCorrectLineCountWithoutLines(t *testing.T) {
	text := new(Text)
	if err := text.Init("", false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	if text.GetLineCount() != 1 {
		t.Fail()
	}
}

func TestTextShouldReturnCorrectLineLength(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	length, err := text.GetLineLengthByOffset(1)
	if err != nil {
		t.Fail()
	}

	if length != 11 {
		t.Fail()
	}
}

func TestTextShouldNotReturnLineLengthOnInvalidPosition(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	_, err := text.GetLineLengthByOffset(4)
	if err == nil {
		t.Fail()
	}
}

func TestTextShouldInsertCharacter(t *testing.T) {
	textContent := "First line\necond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(0, 1, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.InsertCharacter('S', cursor); err != nil {
		t.Fail()
	}

	char, err := text.GetCharacterByCursor(cursor)
	if err != nil {
		t.Fail()
	}

	if char != 'S' {
		t.Fail()
	}
}

func TestTextShouldNotInsertCharacterAtInvalidPosition(t *testing.T) {
	textContent := "First line\necond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(0, 20, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.InsertCharacter('S', cursor); err == nil {
		t.Fail()
	}
}

func TestTextShouldRemoveCharactrerHead(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(1, 1, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.RemoveCharacterHead(cursor); err != nil {
		t.Fail()
	}

	if err := cursor.SetOffsetX(0); err != nil {
		t.Fail()
	}

	char, err := text.GetCharacterByCursor(cursor)
	if err != nil {
		t.Fail()
	}

	if char != 'e' {
		t.Fail()
	}
}

func TestTextShouldRemoveCharactrerTail(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(1, 1, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.RemoveCharacterTail(cursor); err != nil {
		t.Fail()
	}

	if err := cursor.SetOffsetX(2); err != nil {
		t.Fail()
	}

	char, err := text.GetCharacterByCursor(cursor)
	if err != nil {
		t.Fail()
	}

	if char != 'o' {
		t.Fail()
	}
}

func TestTextShouldNotRemoveCharactrerHeadAtInvalidPosition(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(0, 20, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.RemoveCharacterHead(cursor); err == nil {
		t.Fail()
	}
}

func TestTextShouldNotRemoveCharactrerTailAtInvalidPosition(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(0, 20, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.RemoveCharacterTail(cursor); err == nil {
		t.Fail()
	}
}

func TestTextShouldBreaklineAtLineStart(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(0, 1, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.InsertLine(cursor); err != nil {
		t.Fail()
	}

	result, err := text.GetTextAsString()
	if err != nil {
		t.Fail()
	}

	expectedTextContent := "First line\n\nSecond line\nThird line"
	if *result != expectedTextContent {
		t.Fail()
	}
}

func TestTextShouldBreaklineAtLineEnd(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(11, 1, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.InsertLine(cursor); err != nil {
		t.Fail()
	}

	result, err := text.GetTextAsString()
	if err != nil {
		t.Fail()
	}

	expectedTextContent := "First line\nSecond line\n\nThird line"
	if *result != expectedTextContent {
		t.Fail()
	}
}

func TestTextShouldBreaklineInsideLine(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(5, 1, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.InsertLine(cursor); err != nil {
		t.Fail()
	}

	result, err := text.GetTextAsString()
	if err != nil {
		t.Fail()
	}

	expectedTextContent := "First line\nSecon\nd line\nThird line"
	if *result != expectedTextContent {
		t.Fail()
	}
}

func TestTextShouldCombineLineForValidCursorPosition(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(0, 1, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.CombineLine(cursor, false); err != nil {
		t.Fail()
	}

	result, err := text.GetTextAsString()
	if err != nil {
		t.Fail()
	}

	expectedTextContent := "First lineSecond line\nThird line"
	if *result != expectedTextContent {
		t.Fail()
	}
}

func TestTextShouldCombineLineForInvalidCursorPosition(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(0, 5, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.CombineLine(cursor, false); err == nil {
		t.Fail()
	}
}

func TestTextShouldCombineLineWithStepDownForValidCursorPosition(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(0, 1, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if err := text.CombineLine(cursor, true); err != nil {
		t.Fail()
	}

	result, err := text.GetTextAsString()
	if err != nil {
		t.Fail()
	}

	expectedTextContent := "First line\nSecond lineThird line"
	if *result != expectedTextContent {
		t.Fail()
	}
}

func TestTextShouldGetCharacter(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(1, 1, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	char, err := text.GetCharacterByCursor(cursor)
	if err != nil {
		t.Fail()
	}

	if char != 'e' {
		t.Fail()
	}
}

func TestTextShouldNotGetCharacterAtInvalidPosition(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(0, 20, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	_, err := text.GetCharacterByCursor(cursor)
	if err == nil {
		t.Fail()
	}
}

func TestTextShouldConvertBackToString(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	result, err := text.GetTextAsString()
	if err != nil {
		t.Fail()
	}

	if textContent != *result {
		t.Fail()
	}
}

func TestTextShouldCorrecltyIndicateTheModificationState(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	cursor := new(Cursor)
	if err := cursor.Init(0, 0, CreateConsoleMockup(), nil); err != nil {
		t.Fail()
	}

	if text.IsModified() {
		t.Fail()
	}

	if err := text.RemoveCharacterTail(cursor); err != nil {
		t.Fail()
	}

	if !text.IsModified() {
		t.Fail()
	}

	if err := text.ResetModificationState(); err != nil {
		t.Fail()
	}

	if text.IsModified() {
		t.Fail()
	}
}

func TestTextShouldReturnCorrectEolForLFText(t *testing.T) {
	textContent := "First line\nSecond line\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	if text.GetEndOfLineSequenceName() != "LF" {
		t.Fail()
	}
}

func TestTextShouldReturnCorrectEolForCRLFText(t *testing.T) {
	textContent := "First line\r\nSecond line\r\nThird line"

	text := new(Text)
	if err := text.Init(textContent, false, GetTextTestTextConfigMockup()); err != nil {
		t.Fail()
	}

	if text.GetEndOfLineSequenceName() != "CRLF" {
		t.Fail()
	}
}

func GetTextTestTextConfigMockup() *TextConfig {
	return &TextConfig{
		UsePlatformSpecificEndOfLineSequence: false,
	}
}
