package helpers

import (
	"fmt"
	"strings"
)

func PadStart(input string, target_length int, padding string) string {
	length := len(input)
	if length < target_length {
		pad := strings.Repeat(padding, target_length-length)
		return fmt.Sprintf("%s%s", pad, input)
	}

	return input
}
