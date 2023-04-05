package encrypt

import (
	"testing"
)

var (
	tests = map[string]string{
		"123.56": "6Li1QWyBl4TMDfNG3PXdqA",
		//"1000":           "B9NLmxrXbZ9kG3YCDy6eqg",
		//"15610.01010321": "yr9ldYGGll50YZ5AjV03QQ",
		//"834":            "Rsw3Q6ujA9wxS7LbtFzJDQ",
	}
	key = "0123456789acdefh"
)

func TestDecode(t *testing.T) {
	for price, encoded := range tests {
		value, err := AesEBCDecrypt(encoded, key)
		if err != nil {
			t.Fatal(err.Error())
		}
		if price != value {
			t.Fatalf("decode error, expected %v but got %v", price, value)
		}
	}
}

func TestEncode(t *testing.T) {
	for price, encoded := range tests {
		value, err := AesEBCEncrypt(price, key)
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log(value)
		if value != encoded {
			t.Fatalf("encode error, expected %v but got %v", encoded, value)
		}
	}
}
