package generate

import (
	"fmt"

	"github.com/mattn/go-shellwords"
	"github.com/spf13/pflag"
)

func parseConfig(args string) (Config, error) {
	var c Config
	fset := pflag.NewFlagSet("config", pflag.ContinueOnError)

	fset.StringVarP(&c.ConfigFileName, "filename", "f", "", "")
	fset.StringVarP(&c.EnvPrefix, "env-prefix", "e", "", "")

	argv, err := shellwords.Parse(args)
	if err != nil {
		return c, fmt.Errorf("parsing %q in shell syntax: %w", args, err)
	}

	if err := fset.Parse(argv); err != nil {
		return c, fmt.Errorf("parsing flags from %q: %w", args, err)
	}

	return c, nil
}
