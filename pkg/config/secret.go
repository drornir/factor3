package config

import "strings"

type SecretString string

func (s SecretString) String() string {
	return strings.Repeat("*", len(s))
}

func (s SecretString) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *SecretString) UnmarshalText(text []byte) error {
	*s = SecretString(string(text))
	return nil
}
