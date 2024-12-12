# Factor 3

## Usage

Install using

```bash
go get github.com/drornir/factor3
```

Here's a small program to get you started with cobra and viper:

```go
package main
import (
	"fmt"
	"os"

	"github.com/drornir/factor3/pkg/factor3"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Define a struct
type (
	Config struct {
		Username string               `flag:"username" json:"username"`
		Password factor3.SecretString `flag:"password" json:"password"`
		Log      LogConfig            `flag:"log" json:"log"`
	}
	LogConfig struct {
		Level string `flag:"level" json:"level"`
	}
)

var (
	// define a variable to bind with factor3.Bind()
	rootConfig Config
	// an example cobra command that uses the config
	rootCmd    = &cobra.Command{
		Use: "myprogram",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("# config = %#v\n", rootConfig)
		},
	}
)

func init() {
	viperInstance := viper.New()

	// Setting up viper with options that fit factor3
	err := factor3.InitializeViper(factor3.InitArgs{
		Viper:       viperInstance,
		ProgramName: "myprogram", // Used as the env variables prefix
		CfgFile:     "config.json", // Optional path to config file
	})
	cobra.CheckErr(err)

	pflags := rootCmd.Flags()
	// Using Bind() we create Loader that populates the config when called
	// It also registers the flags in your pflag.FlagSet
	loader, err := factor3.Bind(&rootConfig, viperInstance, pflags)
	cobra.CheckErr(err)

	// we need to let cobra parse to commandline flags before calling Load(), so we put it in cobra.OnInitialize()
	cobra.OnInitialize(func() {
		err := loader.Load()
		cobra.CheckErr(err)
		// Advanced: You can call Load() multiple times, for example in reaction to changes to the config file.
		viperInstance.OnConfigChange(func(in fsnotify.Event) {
			if err := loader.Load(); err != nil {
				fmt.Println("error reloading config on viper.OnConfigChange")
			}
		})
	})
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```

This program just prints the configuration.
For this example, as defined above in CfgFile, create a config file json called 'config.json'.

```bash
cat <<EOF > config.json
{
  "username": "u",
  "password": "p",
  "log": {
    "level": "warn"
  }
}
EOF

## reads the file specified by `CfgFile`,
$ go run main.go
# config = main.Config{Username:"u", Password:"p", ...}

## flags that were explicitly set on the command line will override the config file
$ go run main.go --username=u_flag --password=p_flag
# config = main.Config{Username:"u_flag", Password:"p_flag", ...}

## env vars that were explicitly set will override the config file
## note that "MYPROGRAM_" comes from the ProgramName set in the `InitArgs`
$ MYPROGRAM_PASSWORD='p_env' go run main.go --username=u_flag
# config = main.Config{Username:"u_flag", Password:"p_env", ...}

## nested fields can be set using underscore ('_') for ENV, and dash ('-') for flags
$ MYPROGRAM_LOG_LEVEL=info go run main.go --log-level=debug
# config = main.Config{..., Log:main.LogConfig{Level:"debug"}}
```

## Development

### Version 0

⚠️ This project is still in development, so according to semver it is in version 0.
This means that bumping of minor versions (the `x` in `0.x.y`) signifies breaking changes.

#### Version 0.2

I started from scratch, starting from the "opposite" side now - the new implementation
has no code generation, and everything happens in runtime. The motivation for it is to
play with the API and integration with cobra and viper and then convert as much as I can to code gen.

#### Version 0.1

moved to [docs](./docs/why.md)
