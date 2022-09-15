package main

import "testing"

func TestPaddingShouldInitializeForValidValues(t *testing.T) {
	padding := new(Padding)
	if err := padding.Init(1, 0, 1, 0); err != nil {
		t.Fail()
	}
}

func TestPaddingShouldNotInitializeForInvalidValues(t *testing.T) {
	padding := new(Padding)
	if err := padding.Init(-1, 0, -1, 0); err == nil {
		t.Fail()
	}
}

func TestPaddingGettersShouldReturnCorrectValues(t *testing.T) {
	padding := new(Padding)
	if err := padding.Init(1, 2, 3, 4); err != nil {
		t.Fail()
	}

	if padding.GetTopPadding() != 1 {
		t.Fail()
	}

	if padding.GetBottomPadding() != 2 {
		t.Fail()
	}

	if padding.GetLeftPadding() != 3 {
		t.Fail()
	}

	if padding.GetRightPadding() != 4 {
		t.Fail()
	}
}
