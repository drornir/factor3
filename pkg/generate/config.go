package generate

import (
	"strings"

	"github.com/spf13/pflag"
)

func parseConfig(args string) (Config, error) {
	var c Config
	fset := pflag.NewFlagSet("config", pflag.ContinueOnError)

	fset.StringVarP(&c.ConfigFileName, "filename", "f", "", "")

	err := fset.Parse(strings.Split(args, " "))
	return c, err
}
