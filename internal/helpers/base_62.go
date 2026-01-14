package helpers

var base62Set = func() [256]bool {
	set := [256]bool{}
	for _, c := range "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" {
		set[byte(c)] = true
	}
	return set
}()

func IsBase62(input string) bool {
	for index := 0; index < len(input); index++ {
		if !base62Set[input[index]] {
			return false
		}
	}
	return len(input) > 0
}
