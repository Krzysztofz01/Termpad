package main

import "errors"

// Helper structure representing the inner margin (e.g. content padding for display struct)
type Padding struct {
	top    int
	bottom int
	left   int
	right  int
}

// Padding structure initialization function
func (padding *Padding) Init(top int, bottom int, left int, right int) error {
	if top < 0 {
		return errors.New("padding: invalid top padding value")
	}

	if bottom < 0 {
		return errors.New("padding: invalid bottom padding value")
	}

	if left < 0 {
		return errors.New("padding: invalid left padding value")
	}

	if right < 0 {
		return errors.New("padding: invalid right padding value")
	}

	padding.top = top
	padding.bottom = bottom
	padding.left = left
	padding.right = right

	return nil
}

// Return the top padding value
func (padding *Padding) GetTopPadding() int {
	return padding.top
}

// Return the bottom padding value
func (padding *Padding) GetBottomPadding() int {
	return padding.bottom
}

// Return the left padding value
func (padding *Padding) GetLeftPadding() int {
	return padding.left
}

// Return the right padding value
func (padding *Padding) GetRightPadding() int {
	return padding.right
}
