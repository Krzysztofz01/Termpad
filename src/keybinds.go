package main

import (
	"errors"
	"strings"
)

// TODO: Move config structure into this file
// TODO: Implement unit tests

// Structure representing the editor keyboard key-bindings for various operations
type Keybinds struct {
	save   rune
	keyMap map[rune]bool
}

// Editor keybinds structure initialization function
func (keybinds *Keybinds) Init(config *Config) error {
	if config == nil {
		return errors.New("keybinds: invalid config reference")
	}

	keybinds.keyMap = make(map[rune]bool)

	var err error = nil

	keybinds.save, err = keybinds.parseKeybindString(config.KeyBindSave)
	if err != nil {
		return err
	}

	return nil
}

// Helper funcation used to validate and extract the keybind rune from string value
func (keybind *Keybinds) parseKeybindString(keybindValue string) (rune, error) {
	if len(keybindValue) != 1 {
		return 0, errors.New("keybinds: can not parse the keybind configuration")
	}

	targetRune := rune(strings.ToLower(keybindValue)[0])

	if len(keybind.keyMap) == 0 {
		keybind.keyMap[targetRune] = true
		return targetRune, nil
	}

	if _, exist := keybind.keyMap[targetRune]; exist {
		return 0, errors.New("keybinds: ambiguous keybind configuration")
	}

	keybind.keyMap[targetRune] = true
	return targetRune, nil
}

// Return the rune (that entered with [Ctrl] key) will affect in saving the editor changes
func (keybind *Keybinds) GetSaveKeybind() rune {
	return keybind.save
}
