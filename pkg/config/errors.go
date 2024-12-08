package config

import (
	"errors"
	"fmt"
)

type ParseError struct {
	Err error
	// Value is the input pointer to the struct, unless it's invalid
	Value any
}

func (e ParseError) Error() string {
	return fmt.Sprintf("config parse error on %T: %s", e.Value, e.Err.Error())
}
func (e ParseError) Unwrap() error { return e.Err }

type LoadError struct {
	Errs []error
}

func (e LoadError) Error() string {
	return fmt.Sprintf("config load errors: %s", errors.Join(e.Errs...).Error())
}
func (e LoadError) Unwrap() []error { return e.Errs }
