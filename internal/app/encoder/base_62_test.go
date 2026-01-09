package encoder

import "testing"

func TestBase62Encoding(t *testing.T) {
	encoding := Base62Encoding{}

	result := encoding.Encode(11_157)

	if result != "2TX" {
		t.Fatalf("Nope! Expected '2TX', received: %s", result)
	}
}

func TestBase62Decoding(t *testing.T) {
	encoding := Base62Encoding{}

	result, _ := encoding.Decode("2TX")

	if result != uint(11_157) {
		t.Fatalf("Nope! Expected 11_157, received: %d", result)
	}
}
