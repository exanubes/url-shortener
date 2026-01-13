package helpers

func PadStart(input []rune, target_length int, padding rune) []rune {
	length := len(input)
	if length < target_length {
		var pad []rune
		for i := length; i < target_length; i += 1 {
			pad = append(pad, padding)
		}
		return append(pad, input...)
	}

	return input
}
