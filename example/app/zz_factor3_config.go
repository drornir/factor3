// GENERATED FILE DO NOT EDIT
package main

import (
	"fmt"
)

func (s *Config) Factor3Load() error {
	fmt.Printf("type name: Config\n")
	fmt.Printf("annotations: //factor3:generate,\n")
	fmt.Printf("fields:\n")

	fmt.Printf("\t-name=Port, type=string\n")
	fmt.Printf("\t-annotations: //factor3:validate opts...,\n")

	fmt.Printf("\t-name=DBConnection, type=string\n")
	fmt.Printf("\t-annotations: //factor3:validate opts...,\n")

	return nil
}
