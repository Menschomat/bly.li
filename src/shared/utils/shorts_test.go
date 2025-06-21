package utils

import "testing"

func TestIsValidShort(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"lowercase", "abc", true},
		{"uppercase", "XYZ", true},
		{"mixedCase", "GoLang", true},
		{"numeric", "123", false},
		{"mixedAlphaNumeric", "abc123", false},
		{"symbol", "short!", false},
		{"empty", "", false},
	}

	for _, tc := range cases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			if got := IsValidShort(tc.input); got != tc.want {
				t.Errorf("IsValidShort(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
