package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/drornir/factor3"
)

type Config struct {
	Port         string
	DBConnection string
}

func main() {
	c, err := factor3.Load[Config]()
	if err != nil {
		log.Panicf("error loading conf: %s", err)
	}

	fmt.Println("here it is")
	j, _ := json.MarshalIndent(c, "", "  ")
	fmt.Printf("%s\n", j)
}
