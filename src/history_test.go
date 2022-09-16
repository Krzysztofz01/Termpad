package main

import "testing"

func TestHistoryShouldInit(t *testing.T) {
	history := new(History)
	if err := history.Init(); err != nil {
		t.Fail()
	}
}

func TestHistoryShouldPushText(t *testing.T) {
	history := new(History)
	if err := history.Init(); err != nil {
		t.Fail()
	}

	text1 := new(Text)
	if err := text1.Init("Hello World!", false); err != nil {
		t.Fail()
	}

	if history.CanPop() {
		t.Fail()
	}

	if err := history.Push(*text1); err != nil {
		t.Fail()
	}

	if !history.CanPop() {
		t.Fail()
	}

	text2 := new(Text)
	if err := text2.Init("Hello again!", false); err != nil {
		t.Fail()
	}

	if err := history.Push(*text2); err != nil {
		t.Fail()
	}
}

func TestHistoryShouldPopText(t *testing.T) {
	history := new(History)
	if err := history.Init(); err != nil {
		t.Fail()
	}

	text1 := new(Text)
	text1String := "Hello World!"
	if err := text1.Init(text1String, false); err != nil {
		t.Fail()
	}

	if err := history.Push(*text1); err != nil {
		t.Fail()
	}

	text2 := new(Text)
	text2String := "Hello again!"
	if err := text2.Init(text2String, false); err != nil {
		t.Fail()
	}

	if err := history.Push(*text2); err != nil {
		t.Fail()
	}

	var err error
	var targetText *Text
	var expectedText *string

	targetText, err = history.Pop()
	if err != nil {
		t.Fail()
	}

	expectedText, err = targetText.GetTextAsString(false)
	if err != nil {
		t.Fail()
	}

	if *expectedText != text2String {
		t.Fail()
	}

	targetText, err = history.Pop()
	if err != nil {
		t.Fail()
	}

	expectedText, err = targetText.GetTextAsString(false)
	if err != nil {
		t.Fail()
	}

	if *expectedText != text1String {
		t.Fail()
	}
}
