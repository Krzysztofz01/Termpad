package main

import (
	"errors"
	"fmt"
	"strings"
)

// TODO: Implement text notifications that are showing only for few seconds. This concurrent operation can not interfere with
// ,,main'' redrawing (usage of semaphores/mutexes?)

// Structure representing the menu widget that is rendered on the bottom of the editor display
type Menu struct {
	height             int
	notificationText   string
	cursorPositionText string
	fileNameText       string
	eolSequenceText    string
	fileModified       bool
}

// Menu widget structure initialization funcation
func (menu *Menu) Init(menuHeight int, fileName string, eolSequenceName string) error {
	// TODO: Currently there is not support for height different than ,,1''
	if menuHeight != 1 {
		return errors.New("menu: invalid menu height specified")
	}

	if len(fileName) <= 0 {
		return errors.New("menu: invalid file name specified")
	}

	if len(eolSequenceName) <= 0 {
		return errors.New("menu: invalid end-of-line sequence name specified")
	}

	menu.height = menuHeight
	menu.fileNameText = fileName
	menu.eolSequenceText = eolSequenceName

	menu.notificationText = ""
	menu.cursorPositionText = ""
	menu.fileModified = false

	return nil
}

// Function used to update the menu notification text
func (menu *Menu) SetNotificationText(notification string) error {
	menu.notificationText = notification
	return nil
}

// Function used to update the menu cursor position text
func (menu *Menu) SetCursorPositionText(cursor Cursor) error {
	xOffset := cursor.GetOffsetX()
	yOffset := cursor.GetOffsetY()

	menu.cursorPositionText = fmt.Sprintf("[%d;%d]", xOffset, yOffset)
	return nil
}

// Function used to update the file modifiation indication variable
func (menu *Menu) SetFileModificationState(modified bool) error {
	menu.fileModified = modified
	return nil
}

// Return a buffer containg the content of the menu, ready for rendering
func (menu *Menu) GenerateOutputBuffer(width int) ([]rune, error) {
	if width <= 0 {
		return nil, errors.New("menu: invalid width specified to generate output buffer")
	}

	mlNotification := width / 2
	mlInfo := width - mlNotification

	var notificationPart string
	if len(menu.notificationText) >= mlNotification {
		notificationPart = fmt.Sprintf("%s...", menu.notificationText[:mlNotification-3])
	} else {
		notificationBuilder := strings.Builder{}
		notificationBuilder.WriteString(menu.notificationText)

		for notificationBuilder.Len() < mlNotification {
			notificationBuilder.WriteRune(' ')
		}

		notificationPart = notificationBuilder.String()
	}

	const separator = " | "

	informationContentBuilder := strings.Builder{}
	informationContentBuilder.WriteString(menu.fileNameText)
	informationContentBuilder.WriteString(separator)
	informationContentBuilder.WriteString(menu.eolSequenceText)
	informationContentBuilder.WriteString(separator)
	informationContentBuilder.WriteString(menu.cursorPositionText)

	informationContent := informationContentBuilder.String()
	if menu.fileModified {
		informationContent = fmt.Sprintf("*%s", informationContent)
	}

	var informationPart string = ""
	if len(informationContent) <= mlInfo {
		informationBuilder := strings.Builder{}

		for informationBuilder.Len()+len(informationContent) < mlInfo {
			informationBuilder.WriteRune(' ')
		}

		informationBuilder.WriteString(informationContent)
		informationPart = informationBuilder.String()
	}

	outputBuffer := append([]rune(notificationPart), []rune(informationPart)...)
	return outputBuffer, nil
}
