// GENERATED FILE DO NOT EDIT
package main

import (
	"encoding/json"
)

type zz_factor3_Config struct {
	Filename  string
	EnvPrefix string
}

type zz_factor3_JSONer[T any] struct{ t *T }

func (s *zz_factor3_JSONer[T]) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, s.t)
}
func (s *zz_factor3_JSONer[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.t)
}
