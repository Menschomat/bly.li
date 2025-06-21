package tracking

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Menschomat/bly.li/shared/model"
	"github.com/redis/go-redis/v9"
)

func TestParseClickEventValid(t *testing.T) {
	click := model.ShortClick{
		Short:     "abc",
		Timestamp: time.Now().UTC(),
		Ip:        "1.1.1.1",
		UsrAgent:  "test-agent",
	}
	b, err := json.Marshal(click)
	if err != nil {
		t.Fatalf("failed to marshal click: %v", err)
	}
	msg := redis.XMessage{
		ID:     "1",
		Values: map[string]interface{}{"data": string(b)},
	}
	got, ok := parseClickEvent(msg)
	if !ok {
		t.Fatalf("expected ok=true, got false")
	}
	if got != click {
		t.Errorf("unexpected decoded click: %+v, want %+v", got, click)
	}
}

func TestParseClickEventInvalid(t *testing.T) {
	cases := []struct {
		name string
		msg  redis.XMessage
	}{
		{"missing data", redis.XMessage{ID: "missing", Values: map[string]interface{}{}}},
		{"bad json", redis.XMessage{ID: "badjson", Values: map[string]interface{}{"data": "notjson"}}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := parseClickEvent(tc.msg); ok {
				t.Errorf("parseClickEvent(%s) expected false", tc.msg.ID)
			}
		})
	}
}
