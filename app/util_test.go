package app

import "testing"

var integerTestCases = map[string]bool{
	"123":    true,
	"0":      true,
	"-456":   true,
	"123.45": false,
	"abc":    false,
	"":       false,
}

var uuidTestCases = map[string]bool{
	"01234567-89ab-cdef-0123-456789abcdef":  true,
	"":                                      false,
	"0123456789abcdef0123456789abcdef1234":  false,
	"01234567-89ab-cdef-0123-456789abcde":   false,
	"01234567-89ab-cdef-0123-456789abcdefg": false,
}

func TestIsInteger(t *testing.T) {
	for str, expected := range integerTestCases {
		result := IsInteger(str)
		if result != expected {
			t.Errorf("IsInteger('%s') returned %v, expected %v", str, result, expected)
		}
	}
}

func TestIsUUID(t *testing.T) {
	for str, expected := range uuidTestCases {
		result := IsUUID(str)
		if result != expected {
			t.Errorf("IsUUID('%s') returned %v, expected %v", str, result, expected)
		}
	}
}

func BenchmarkIsInteger(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for str := range integerTestCases {
			IsInteger(str)
		}
	}
}

func BenchmarkIsUUID(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for str := range uuidTestCases {
			IsUUID(str)
		}
	}
}
