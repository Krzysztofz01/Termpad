package main

import (
	"errors"
	"strings"
)

// Structure representing the editor keyboard key-bindings for various operations
type Keybinds struct {
	save   rune
	exit   rune
	keyMap map[rune]bool
	config *KeybindsConfig
}

// Editor keybinds structure initialization function
func (keybinds *Keybinds) Init(keybindsConfig *KeybindsConfig) error {
	if keybindsConfig == nil {
		defaultConfig := CreateDefaultKeybindsConfig()
		keybinds.config = &defaultConfig
	} else {
		keybinds.config = keybindsConfig
	}

	keybinds.keyMap = make(map[rune]bool)
	var err error = nil

	keybinds.save, err = keybinds.parseKeybindString(keybinds.config.SaveKeybind)
	if err != nil {
		return err
	}

	keybinds.exit, err = keybinds.parseKeybindString(keybinds.config.ExitKeybind)
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

// Return the rune (that entered with [Ctrl] key) will affect in exiting the program
func (keybind *Keybinds) GetExitKeybind() rune {
	return keybind.exit
}

// A structure containing the configuration for the keybinds structure
type KeybindsConfig struct {
	SaveKeybind string `json:"keybind-save"`
	ExitKeybind string `json:"keybind-exit"`
}

// Return a new isntance of the keybinds configuration with default values
func CreateDefaultKeybindsConfig() KeybindsConfig {
	return KeybindsConfig{
		SaveKeybind: "s",
		ExitKeybind: "x",
	}
}
