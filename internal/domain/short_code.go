package domain

import (
	"fmt"
	"strings"
)

const MAX_SHORT_CODE_SIZE = 11
const MIN_SHORT_CODE_SIZE = 0

type ShortCode struct {
	size  int
	value string
	zero  string
}

func NewShortCodeFromParam(value string) (ShortCode, error) {
	if len(value) > 11 {
		return ShortCode{}, ErrInvalidShortCode
	}
	if !allowed_alphabet.contains(value) {
		return ShortCode{}, ErrInvalidShortCode
	}
	return ShortCode{
		value: value,
		size:  len(value),
		zero:  "0",
	}, nil
}

func NewShortCode(value string, size int, zero_char string) (ShortCode, error) {
	if size <= MIN_SHORT_CODE_SIZE {
		return ShortCode{}, ErrShortCodeEmpty
	}

	if size > MAX_SHORT_CODE_SIZE {
		return ShortCode{}, ErrExceededMaxSize
	}

	if len(value) > int(size) {
		return ShortCode{}, ErrInvalidShortCode
	}
	return ShortCode{
		value: value,
		size:  size,
		zero:  zero_char,
	}, nil
}

func (code ShortCode) String() string {
	return pad_start(code.value, code.size, code.zero)
}

func pad_start(input string, target_length int, padding string) string {
	length := len(input)
	if length < target_length {
		pad := strings.Repeat(padding, target_length-length)
		return fmt.Sprintf("%s%s", pad, input)
	}

	return input
}
