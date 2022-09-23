package main

import (
	"fmt"
	"strings"
)

const (
	// TODO: Extract this to configurable Numeration struct property
	NumerationWidth = 3
)

// TODO: Implement unit tests

// Structure representing the numeration widget that is rendered on the left side of the editor display
type Numeration struct {
	useRelativeNumeration bool
}

// Numeration widget structure initialization funcation
func (numeration *Numeration) Init() error {
	numeration.useRelativeNumeration = false
	return nil
}

// Return a value indicating if the numeration is relative to the cursor position
func (numeration *Numeration) IsNumerationRelative() bool {
	return numeration.useRelativeNumeration
}

// Function is toggling the numeration between absolute and relative
func (numeration *Numeration) ToggleRelativeNumeration(toggle bool) error {
	numeration.useRelativeNumeration = toggle
	return nil
}

// Return a buffer containg the content of the numeration, ready for rendering
func (numeration *Numeration) GenerateOutputBuffer(yOffset int, yShift int, height int) ([]string, error) {
	numerationBuffer := make([]string, height)
	numerationBuilder := strings.Builder{}

	for yIndex := 0; yIndex < height; yIndex += 1 {
		numerationBuilder.Reset()

		var number string
		if numeration.useRelativeNumeration {
			number = fmt.Sprint(abs(yOffset - yIndex))
		} else {
			number = fmt.Sprint(yIndex + yShift)
		}

		if len(number) > NumerationWidth {
			numerationBuffer[yIndex] = fmt.Sprintf("*%s", number[len(number)-NumerationWidth-1:])
		} else {
			for xFill := 0; xFill < NumerationWidth-len(number); xFill += 1 {
				numerationBuilder.WriteRune(' ')
			}
			numerationBuilder.WriteString(number)

			numerationBuffer[yIndex] = numerationBuilder.String()
		}

	}

	return numerationBuffer, nil
}

// Math absolute value helper funcation. Golang provides math utilites only for floats
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
