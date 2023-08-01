package config_test

import (
	"fmt"
	"log"

	"github.com/drornir/factor3/pkg/config"
)

func init() {
	log.SetPrefix(">>> ")
	log.SetFlags(log.Lmsgprefix)
}

func ExampleParseString() {
	const configExample = `
	package pkg

	type Config struct {
		Username string
		Password string
		XX, YY abc.TTT
	}
	`

	conf, err := config.ParseString(configExample)
	if err != nil {
		fmt.Println(err)
		return // handle
	}

	usage := conf.Schema.Env.ShellUsage()
	fmt.Print(usage)
	// Output:
	// USERNAME=<string>
	// PASSWORD=<string>
}

func _ExampleParseString_nested() {
	const configExample = `
	package pkg

	type Config struct {
		Github struct {
			Username string
			Password string
		}
	}
	`

	conf, err := config.ParseString(configExample)
	if err != nil {
		fmt.Println(err)
		return // handle
	}

	usage := conf.Schema.Env.ShellUsage()
	fmt.Print(usage)
	// Output:
	// GITHUB_USERNAME=<string>
	// GITHUB_PASSWORD=<string>
}
