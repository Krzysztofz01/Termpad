package main

import (
	"os"
)

func main() {
	if len(os.Args) != 2 {
		// TODO: Panic - Invalid arguments notification
		os.Exit(1)
		return
	}

	targetFilePath := os.Args[1]

	console, err := CreateConsole()
	if err != nil {
		// TODO: Panic - Failed to create console API instance notification
		os.Exit(1)
		return
	}

	editor := new(Editor)
	if err := editor.Init(targetFilePath, console); err != nil {
		// TODO: Panic - Failed to create editor instance notification
		os.Exit(1)
		return
	}

	if err := editor.Start(); err != nil {
		// TODO: Handle errors while editing
		os.Exit(1)
		return
	}

	os.Exit(0)
}
