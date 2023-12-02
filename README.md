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
  // ... c is ready
}
```

You need to include `//go:generate factor3 generate` once, somewhere in the package,
in order to trigger code generation.

After that, you need annotate every struct you want with `//factor3:generate` in
order to generate an implementation of `Factor3Load()`, which loads the data into
the struct.

For this example, the corresponding yaml is the following (notice it's snakecase)

```yaml
port: "3001"
db_connection: root:pass@localhost:3306
some_number: 4
some_flag: true
```

And the env vars are 

```sh
export EX_PORT=9090
export EX_DB_CONNECTION=root:pass@localhost:3307
export EX_SOME_NUMBER=6
export EX_SOME_FLAG=false
```

Setting flags is optional (like `DBConnection` in the example), but you can
explicitly set them by using the `//factor3:pflag` annotation.
It accepts between one and two string arguments. For example,
`//factor3:pflag some-number n` will declare a flag called `some-number`, with 
the short form `n`.

## Install

Using go install

```sh
go install github.com/drornir/factor3
```

## Annotations API

Adding configuration to the code generation is done using one line comments
which start with `//factor3:`. Notice there is no space after the `//`. We call
these comments "annotations".

"Top level" annotations means comments that come directly above the struct type.

"In Struct" annotations means comments that come directly above a field in the struct.

### Top Level

#### `factor3:generate`

`//factor3:generate [--filename <FILE_PATH>] [--env-prefix <PREFIX>]`

This is **required** to trigger the code generation for a certain struct type.

- `filename` is a path to where you place the json or yaml file with values for
  this struct
- `env-prefix` is a string that will be prepended to all env vars lookups. ⚠️ Note
  that an underscore will also be prepended between your prefix and the rest of 
  the env var name (for `--env-prefix=PRE` and env var `VAR`, the result will be `PRE_VAR`).

### In Struct

#### `//factor3:pflag`

`//factor3:pflag <LONG_FLAG_NAME> [SHORT_FLAG_NAME]`

Will add flag support for this field. At least one argument is required,
which is the long form of the flag(`LONG_FLAG_NAME`). Optionally, you can set 
`SHORT_FLAG_NAME`, which must be a one character string representing the short flag letter.

The underlying package in use is not the standard library `flag`, but 
`github.com/spf13/pflag`, the POSIX compliant alternative.
