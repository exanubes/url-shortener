package encoder

import "slices"

type Base62Encoding struct{}

var dictionary = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func (_ Base62Encoding) Encode(input uint) string {
	if input == 0 {
		return "0"
	}
	var digits []rune
	current := input
	for current != 0 {
		digits = append(digits, dictionary[current%62])
		current = current / 62
	}

	slices.Reverse(digits)

	return string(digits)
}

func (_ Base62Encoding) Decode(input string) int {
	return -1
}
