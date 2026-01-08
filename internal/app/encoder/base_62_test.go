package encoder

import "testing"

func TestBase62Encoding(t *testing.T) {
	encoding := Base62Encoding{}

	result := encoding.Encode(11157)

	if result != "2TX" {
		t.Fatalf("Nope! Expected '2TX', received: %s", result)
	}
}
