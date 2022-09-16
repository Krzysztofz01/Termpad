package main

import "errors"

// TODO: Extract to preference
const (
	historyStackSize = 512
)

// Structure representing the edistors text history stack (LIFO)
type History struct {
	nodes []Text
	count int
}

// History structure initialization function
func (history *History) Init() error {
	history.nodes = make([]Text, historyStackSize)
	history.count = 0

	return nil
}

// Add the given state of text to the history stack
func (history *History) Push(text Text) error {
	if history.count < historyStackSize {
		history.nodes[history.count] = text
		history.count += 1

		return nil
	}

	history.nodes = history.nodes[1:]
	history.nodes[history.count-1] = text

	return nil
}

// Return the next (LIFO) text state from the history stack
func (history *History) Pop() (*Text, error) {
	if history.count == 0 {
		return nil, errors.New("history: can not retrieve text from empty history stack")
	}

	text := history.nodes[history.count-1]
	history.count -= 1

	return &text, nil
}

// Return a bool value indicating if there are any text structs on the history stack
func (history *History) CanPop() bool {
	return history.count > 0
}
