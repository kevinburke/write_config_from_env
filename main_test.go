package main

import (
	"bytes"
	"testing"
)

var escapeTests = []struct {
	in   string
	want string
}{
	{"foo", "foo"},
	{"foo-bar", "'foo-bar'"},
}

func TestEscape(t *testing.T) {
	for _, tt := range escapeTests {
		got := escape(tt.in)
		if got != tt.want {
			t.Errorf("escape(%q): got %q, want %q", tt.in, got, tt.want)
		}
	}
}

var expectedConfig = `port: 56789
public_host: 'localhost:7'
twilio_account_sid: AC123
timezones:
  - America/Los_Angeles
  - America/New_York


`

func TestWriteConfig(t *testing.T) {
	t.Parallel()
	env := []string{
		"PORT=56789",
		"PUBLIC_HOST=localhost:7",
		"TWILIO_ACCOUNT_SID=AC123",
		"TIMEZONES=America/Los_Angeles,America/New_York",
	}
	buf := new(bytes.Buffer)
	writeConfig(buf, env, nil)
	if s := buf.String(); s != expectedConfig {
		t.Errorf("expected config to be %q, got %q", expectedConfig, s)
	}
}