package encoder

import "testing"

func TestBase62Encoding(t *testing.T) {
	encoding := Base62Codec{}

	result := encoding.Encode(11_157)

	if result != "2TX" {
		t.Fatalf("Nope! Expected '2TX', received: %s", result)
	}
}

func TestBase62Decoding(t *testing.T) {
	encoding := Base62Codec{}

	result, _ := encoding.Decode("2TX")

	if result != uint64(11_157) {
		t.Fatalf("Nope! Expected 11_157, received: %d", result)
	}
}
