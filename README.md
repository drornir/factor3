# Factor 3

⚠️ This package is still EXPERIMENTAL, so versioning is not used yet.

## Overview

Factor number Three of the [Twelve Factor App](https://12factor.net/config) is "Config"

They suggest to store the configuration in env vars. I completely agree.

However, it's left to us to define how the env vars are loaded into our app.

In addition to that, some defaults are better stored in a json or yaml, which
can be overridden 

This project is inspired by 

- The legendary `github.com/spf13/viper`, and
- backstage.io's `app-config.yaml` (still not implemented)

The idea is to declaratively define a big struct, and pass that struct to 
this app, and it returns your filled config.

## Example

The [example](./example/app) looks like this:

```go
package main

//go:generate ../../bin/factor3 generate

import (
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
  // ... c is ready
}
```

You need to include `//go:generate factor3 generate` once, somewhere in the package,
in order to trigger code generation.

After that, you need annotate every struct you want with `//factor3:generate` in
order to generate an implementation of `Factor3Load()`, which loads the data into
the struct.