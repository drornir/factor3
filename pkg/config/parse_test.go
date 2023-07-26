package config_test

import (
	"fmt"

	"github.com/drornir/factor3/pkg/config"
)

func ExampleParseString() {
	const configExample = `
	type Config struct {
		Username string
		Password string
	}
	`

	conf, err := config.ParseString(configExample)
	if err != nil {
		return // handle
	}

	usage := conf.Schema.Env.Usage.Shell()
	fmt.Print(usage)
	// Output:
	// USERNAME=<string>
	// PASSWORD=<string>
}

func ExampleParseString_nested() {
	const configExample = `
	type Config struct {
		Github struct {
			Username string
			Password string
		}
	}
	`

	conf, err := config.ParseString(configExample)
	if err != nil {
		return // handle
	}

	usage := conf.Schema.Env.Usage.Shell()
	fmt.Print(usage)
	// Output:
	// GITHUB_USERNAME=<string>
	// GITHUB_PASSWORD=<string>
}
