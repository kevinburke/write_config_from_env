package yamlconfig

import (
	"testing"

	yaml "gopkg.in/yaml.v2"
)

var boolTests = []struct {
	in       string
	expected bool
}{
	{`value: "true"`, true},
	{`value: "TRUE"`, true},
	{`value: "yes"`, true},
	{`value: "t"`, true},
	{`value: true`, true},
	{`value: "1"`, true},
	{`value: "false"`, false},
	{`value: "f"`, false},
	{`value: false`, false},
}

type V struct {
	Value Bool `yaml:"value"`
}

func TestBool(t *testing.T) {
	for _, tt := range boolTests {
		v := new(V)
		if err := yaml.Unmarshal([]byte(tt.in), v); err != nil {
			t.Fatal(err)
		}
		if bool(v.Value) != tt.expected {
			t.Errorf("Unmarshal(%q): want %t got %t", tt.in, tt.expected, v.Value)
		}
	}
}

func TestBoolError(t *testing.T) {
	v := new(V)
	err := yaml.Unmarshal([]byte(`value: "unknown"`), v)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
	want := `could not convert string to bool: "unknown"`
	if err.Error() != want {
		t.Errorf("bad error: got %v want %q", err, want)
	}

	err = yaml.Unmarshal([]byte(`value: ""`), v)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
	want = `cannot unmarshal empty string into type bool`
	if err.Error() != want {
		t.Errorf("bad error: got %v want %q", err, want)
	}
}

var intTests = []struct {
	in       string
	expected int
}{
	{`value: "7"`, 7},
	{`value: "-3"`, -3},
	{`value: "0"`, 0},
	{`value: 7`, 7},
	{`value: -3`, -3},
	{`value: 0`, 0},
}

type IV struct {
	Value Int `yaml:"value"`
}

func TestInt(t *testing.T) {
	for _, tt := range intTests {
		v := new(IV)
		if err := yaml.Unmarshal([]byte(tt.in), v); err != nil {
			t.Fatal(err)
		}
		if int(v.Value) != tt.expected {
			t.Errorf("Unmarshal(%q): want %d got %d", tt.in, tt.expected, v.Value)
		}
	}
}

func TestIntError(t *testing.T) {
	v := new(IV)
	err := yaml.Unmarshal([]byte(`value: "unknown"`), v)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
	want := `strconv.Atoi: parsing "unknown": invalid syntax`
	if err.Error() != want {
		t.Errorf("bad error: got %v want %q", err, want)
	}

	err = yaml.Unmarshal([]byte(`value: "5.01"`), v)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
	want = `strconv.Atoi: parsing "5.01": invalid syntax`
	if err.Error() != want {
		t.Errorf("bad error: got %v want %q", err, want)
	}
}
