package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		// TODO: Implement title screen
		printErrorMessage(errors.New("args: invalid program arguments"))
		os.Exit(1)
		return
	}

	targetFilePath := os.Args[1]

	config := new(Config)
	if err := config.Init(); err != nil {
		printErrorMessage(err)
		os.Exit(1)
		return
	}

	console, err := CreateConsole()
	if err != nil {
		printErrorMessage(err)
		os.Exit(1)
		return
	}

	editor := new(Editor)
	if err := editor.Init(targetFilePath, console, config); err != nil {
		printErrorMessage(err)
		console.Dispose()
		os.Exit(1)
		return
	}

	if err := editor.Start(); err != nil {
		printErrorMessage(err)
		console.Dispose()
		os.Exit(1)
		return
	}

	if err := console.Dispose(); err != nil {
		printErrorMessage(err)
		os.Exit(1)
	}

	os.Exit(0)
}

const (
	redColorCode   = "\033[31m"
	resetColorCode = "\033[0m"
)

func printErrorMessage(err error) {
	fmt.Printf("%sThe program encountered a problem! [ %s ]%s\n", redColorCode, err, resetColorCode)
}
