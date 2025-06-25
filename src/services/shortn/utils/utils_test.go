package utils

import (
	"fmt"
	"testing"
)

func TestParseUrlValid(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"https://example.com", "https://example.com"},
		{"http://example.com/path?query=1", "http://example.com/path?query=1"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.in, func(t *testing.T) {
			got, err := ParseUrl(c.in)
			if err != nil {
				t.Fatalf("ParseUrl(%q) unexpected error: %v", c.in, err)
			}
			if got != c.want {
				t.Fatalf("ParseUrl(%q)=%q want %q", c.in, got, c.want)
			}
		})
	}
}

func TestParseUrlInvalid(t *testing.T) {
	cases := []string{"", "not a url", "htp:/example.com"}
	for _, c := range cases {
		c := c
		t.Run(c, func(t *testing.T) {
			if _, err := ParseUrl(c); err == nil {
				t.Fatalf("ParseUrl(%q) expected error", c)
			}
		})
	}
}

func TestIsUrl(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cases := []string{
			"https://example.com",
			"http://example.com/path",
			"ftp://example.com",
		}
		for _, u := range cases {
			if !isUrl(u) {
				t.Errorf("isUrl(%q) = false, want true", u)
			}
		}
	})

	t.Run("invalid", func(t *testing.T) {
		cases := []string{"", "notaurl", "http://"}
		for _, u := range cases {
			if isUrl(u) {
				t.Errorf("isUrl(%q) = true, want false", u)
			}
		}
	})
}

func TestGetRandomIntInRange(t *testing.T) {
	min, max := 10, 20
	for i := 0; i < 100; i++ {
		got := GetRandomIntInRange(min, max)
		if got < min || got > max {
			t.Fatalf("GetRandomIntInRange returned %d out of range [%d,%d]", got, min, max)
		}
	}
}

func TestGetSquidShort(t *testing.T) {
	numbers := []uint64{1, 42, 123456}
	for _, n := range numbers {
		n := n
		t.Run(fmt.Sprintf("%d", n), func(t *testing.T) {
			s, err := GetSquidShort(n)
			if err != nil {
				t.Fatalf("GetSquidShort(%d) unexpected error: %v", n, err)
			}
			if s == "" {
				t.Fatalf("GetSquidShort(%d) returned empty string", n)
			}
			if len(s) < 5 {
				t.Fatalf("GetSquidShort(%d) length=%d, want >=5", n, len(s))
			}
		})
	}
}

func TestIntBytesRoundTrip(t *testing.T) {
	numbers := []int{0, 1, -1, 123456, -98765}
	for _, n := range numbers {
		b := intToBytes(n)
		got := bytesToInt(b)
		if got != n {
			t.Errorf("round trip failed: want %d got %d", n, got)
		}
	}
}
