package utils

import "testing"

func TestParseUrlValid(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"https://example.com", "https://example.com"},
		{"http://example.com/path?query=1", "http://example.com/path?query=1"},
	}
	for _, c := range cases {
		got, err := ParseUrl(c.in)
		if err != nil {
			t.Errorf("ParseUrl(%q) unexpected error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseUrl(%q)=%q want %q", c.in, got, c.want)
		}
	}
}

func TestParseUrlInvalid(t *testing.T) {
	cases := []string{"", "not a url", "htp:/example.com"}
	for _, c := range cases {
		if _, err := ParseUrl(c); err == nil {
			t.Errorf("ParseUrl(%q) expected error", c)
		}
	}
}
