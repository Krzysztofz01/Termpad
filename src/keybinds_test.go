package main

import "testing"

func TestKeybindsShouldInitializeForDefaultConfig(t *testing.T) {
	keybinds := new(Keybinds)
	if err := keybinds.Init(nil); err != nil {
		t.Fail()
	}
}

func TestKeybindsShouldNotInitializeForInvalidConfigAmbiguousKeys(t *testing.T) {
	config := KeybindsConfig{
		SaveKeybind: "s",
		ExitKeybind: "s",
	}

	keybinds := new(Keybinds)
	if err := keybinds.Init(&config); err == nil {
		t.Fail()
	}
}

func TestKeybindsShouldNotInitializeForInvalidConfigParsingFailed(t *testing.T) {
	config := KeybindsConfig{
		SaveKeybind: "hello",
		ExitKeybind: "world",
	}

	keybinds := new(Keybinds)
	if err := keybinds.Init(&config); err == nil {
		t.Fail()
	}
}

func TestKeybindsGettersShouldReturnCorrectValue(t *testing.T) {
	config := KeybindsConfig{
		SaveKeybind: "s",
		ExitKeybind: "x",
	}

	keybinds := new(Keybinds)
	if err := keybinds.Init(&config); err != nil {
		t.Fail()
	}

	keybind := keybinds.GetSaveKeybind()
	if keybind != 's' {
		t.Fail()
	}
}
