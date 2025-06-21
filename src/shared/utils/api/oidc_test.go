package api

import (
	"net/http/httptest"
	"testing"
)

func TestTokenFromHeader(t *testing.T) {
	cases := []struct {
		name   string
		header string
		want   string
	}{
		{"valid bearer", "Bearer token123", "token123"},
		{"case insensitive", "bearer token456", "token456"},
		{"malformed header", "Basic abc", ""},
		{"too short", "Bearer", ""},
	}

	for _, tc := range cases {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		req.Header.Set("Authorization", tc.header)
		if got := TokenFromHeader(req); got != tc.want {
			t.Errorf("%s: got %q want %q", tc.name, got, tc.want)
		}
	}
}
