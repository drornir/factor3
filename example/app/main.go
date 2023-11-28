package main

//go:generate ../../bin/factor3 generate

import (
	factor3 "github.com/drornir/factor3/pkg/runtime"
)

//factor3:generate
type Config struct {
	//factor3:validate regex "^[0-9]+$"
	Port         string
	DBConnection string
}

func main() {
	var c Config
	if err := factor3.Load(&c); err != nil {
		panic(err)
	}
}
