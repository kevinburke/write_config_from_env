package yamlconfig

import (
	"errors"
	"fmt"
	"strconv"
)

// Bool can unmarshal a YAML string ("true", "false") into a boolean value.
type Bool bool

func (b *Bool) UnmarshalYAML(f func(interface{}) error) error {
	var s string
	if err := f(&s); err != nil {
		var t bool
		if err := f(&t); err != nil {
			return err
		}
		*b = Bool(t)
		return nil
	}
	switch s {
	case "":
		return errors.New("cannot unmarshal empty string into type bool")
	case "true", "TRUE", "t", "1", "yes":
		*b = true
		return nil
	case "false", "FALSE", "f", "0", "no":
		*b = false
		return nil
	default:
		return fmt.Errorf("could not convert string to bool: %q", s)
	}
}

// Int can unmarshal a YAML integer ("5", "-3") into an integer value.
type Int int

func (i *Int) UnmarshalYAML(f func(interface{}) error) error {
	var it int
	if err := f(&it); err != nil {
		var s string
		if err := f(&s); err != nil {
			return nil
		}
		if val, err := strconv.Atoi(s); err != nil {
			return err
		} else {
			*i = Int(val)
			return nil
		}
	}
	*i = Int(it)
	return nil
}
