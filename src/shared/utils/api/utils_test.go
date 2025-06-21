package api

import (
	"net/http/httptest"
	"testing"
)

func TestReadUserIPPrecedence(t *testing.T) {
	cases := []struct {
		name      string
		real      string
		forwarded string
		remote    string
		want      string
	}{
		{"real header preferred", "1.1.1.1", "2.2.2.2", "3.3.3.3", "1.1.1.1"},
		{"forwarded header next", "", "2.2.2.2", "3.3.3.3", "2.2.2.2"},
		{"remote addr fallback", "", "", "3.3.3.3", "3.3.3.3"},
	}

	for _, tc := range cases {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		req.RemoteAddr = tc.remote
		if tc.real != "" {
			req.Header.Set("X-Real-Ip", tc.real)
		}
		if tc.forwarded != "" {
			req.Header.Set("X-Forwarded-For", tc.forwarded)
		}
		if got := ReadUserIP(req); got != tc.want {
			t.Errorf("%s: got %q want %q", tc.name, got, tc.want)
		}
	}
}
