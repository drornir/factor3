package main

//go:generate ../../bin/factor3 generate

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	factor3 "github.com/drornir/factor3/pkg/runtime"
)

//factor3:generate --filename ./example/app/config.yaml --env-prefix EX
type Config struct {
	DBConnection string
	//factor3:pflag port
	Port string
	//factor3:pflag some-number n
	SomeNumber int
	//factor3:pflag some-flag
	SomeFlag bool
}

func main() {
	var c Config
	if err := factor3.Load(&c, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if b, err := json.MarshalIndent(c, "", "\t"); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(b))
	}
}
