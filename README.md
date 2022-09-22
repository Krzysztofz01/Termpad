# Termpad
**Termpad** is a lightweight, minimalist terminal text editor. It is ideal for making small changes to text or configuration files quickly, but also for editing large texts and programming. The main motivation behind this project was to create a *"true"* multi-platform text editor and to practice developing applications using the **Go programming language**. The application has been written in such a way that you can easily swap the implementation of the *,,console API''* used to interact with the console, which makes the project independent from external libraries and allows you to easily test the application. Currently, Termpad uses the **[TCell](https://github.com/gdamore/tcell)** library to generate the console UI.

## Installation
Regardless of the operating system, the following programs are required:
- **git** *(or other operation-system-specific alternative implementation)*
- **go 19.0+** *(or other operation-system-specific alternative implementation)*

```sh
# The installation script will be implemented in the future

# Clone the source code from the repository
git clone https://github.com/Krzysztofz01/Termpad.git

# Go to the source code directory
cd Termpad/src

# Build the project
go build

# The first program run will generate a fresh configuration file
```

## Configuration
The properties of the configuration file may differ depending on the version

```json
{
 "history-configuration": {
  "history-stack-size": 256 // The size of stack containing the changes to which we can revert
 },
 "keybinds-configuration": {
  "keybind-save": "s", // Keybind used for saving the changes
  "keybind-exit": "x" // Keybind used for closing the program
 },
 "cursor-configuration": {
  "cursor-style": "block", // Style of the cursor. Available options [block, line, bar]
  "use-animations": false // Enable/disable cursor animations
 },
 "text-configuration": {
  "use-platform-specific-eol-sequence": true // Usege of operation system specific EOL. Example: CRLF of Windows and LF for GNU/Linux distros
 }
}
```