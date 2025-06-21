package config

import "testing"

func TestFillEnvStruct(t *testing.T) {
	type testConfig struct {
		Foo string `env:"FOO"`
		Bar int    `env:"BAR"`
	}

	t.Setenv("FOO", "hello")
	t.Setenv("BAR", "123")

	var cfg testConfig
	if err := FillEnvStruct(&cfg); err != nil {
		t.Fatalf("FillEnvStruct returned error: %v", err)
	}

	if cfg.Foo != "hello" || cfg.Bar != 123 {
		t.Errorf("got %+v", cfg)
	}
}
