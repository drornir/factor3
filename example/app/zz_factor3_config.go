package main

import (
	"fmt"
)

func (s *Config) Factor3Load() error {
	fmt.Printf("type name: %s\n", "Config")
	fmt.Printf("annotations: [//factor3:generate]\n")
	fmt.Printf("fields:\n")
	fmt.Printf("\t-Port (string)\n")
	fmt.Printf("\t-DBConnection (string)\n")

	return nil
}
