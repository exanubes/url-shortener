package domain

type alphabet [256]bool

func create_valid_alphabet_set() alphabet {
	set := [256]bool{}
	for _, c := range "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" {
		set[byte(c)] = true
	}
	return set
}

func (alphabet alphabet) contains(input string) bool {

	for index := 0; index < len(input); index++ {
		if !alphabet[input[index]] {
			return false
		}
	}
	return len(input) > 0
}

var allowed_alphabet = create_valid_alphabet_set()
