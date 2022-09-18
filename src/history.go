package main

import "errors"

// Structure representing the edistors text history stack (LIFO)
type History struct {
	nodes []Text
	count int

	config *HistoryConfig
}

// History structure initialization function
func (history *History) Init(historyConfig *HistoryConfig) error {
	if historyConfig == nil {
		defaultConfig := CreateDefaultHistoryConfig()
		history.config = &defaultConfig
	} else {
		history.config = historyConfig
	}

	if history.config.HistoryStackSize <= 0 {
		return errors.New("history: invalid stack size specified in the configuration")
	}

	history.nodes = make([]Text, history.config.HistoryStackSize)
	history.count = 0

	return nil
}

// Add the given state of text to the history stack
func (history *History) Push(text Text) error {
	if history.count < history.config.HistoryStackSize {
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

// A structure containing the configuration for the history structure
type HistoryConfig struct {
	HistoryStackSize int `json:"history-stack-size"`
}

// Return a new isntance of the text configuration with default values
func CreateDefaultHistoryConfig() HistoryConfig {
	return HistoryConfig{
		HistoryStackSize: 256,
	}
}
