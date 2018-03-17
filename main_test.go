package main

import (
	"bytes"
	"testing"
)

var expectedConfig = `flag_new_router: "true"
port: "56789"
public_host: localhost:7
timezones:
- America/Los_Angeles
- America/New_York
twilio_account_sid: AC123
`

func TestWriteConfig(t *testing.T) {
	t.Parallel()
	env := []string{
		"FLAG_NEW_ROUTER=true",
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
