package factor3

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/drornir/factor3/pkg/log"
)

type InitArgs struct {
	Viper       *viper.Viper
	ProgramName string
	CfgFile     string
}

func InitializeViper(a InitArgs) error {
	log.GG().D(context.TODO(), "initializing viper", "programName", a.ProgramName)
	a.Viper.SetEnvPrefix(a.ProgramName)
	a.Viper.AllowEmptyEnv(true)
	a.Viper.AutomaticEnv()
	a.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if a.CfgFile != "" {
		a.Viper.SetConfigFile(a.CfgFile)
	} else {
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			configHome = filepath.Join(home, ".config", a.ProgramName)
		}

		a.Viper.AddConfigPath(configHome)
		a.Viper.SetConfigName("config")
		a.Viper.WatchConfig()
	}

	// If a config file is found, read it in.
	if err := a.Viper.ReadInConfig(); err != nil {
		if !errors.Is(err, viper.ConfigFileNotFoundError{}) {
			return fmt.Errorf("reading in config file using viper: %w", err)
		}
		// fmt.Fprintln(os.Stderr, err.Error())
	} else {
		// fmt.Fprintln(os.Stderr, "Read in config file:", a.Viper.ConfigFileUsed())
	}

	return nil
}

type viperFlagAdapter struct {
	pf        *pflag.Flag
	viperPath string
}

func (f viperFlagAdapter) HasChanged() bool    { return f.pf.Changed }
func (f viperFlagAdapter) ValueString() string { return f.pf.Value.String() }
func (f viperFlagAdapter) ValueType() string   { return f.pf.Value.Type() }
func (f viperFlagAdapter) Name() string {
	if f.viperPath == "" {
		return strings.ReplaceAll(f.pf.Name, "-", ".")
	}
	return f.viperPath
}

type viperFlagsAdapter struct {
	pfs            *pflag.FlagSet
	vipathByPFName map[string]string
}

func (fs viperFlagsAdapter) VisitAll(fn func(viper.FlagValue)) {
	fs.pfs.VisitAll(func(pf *pflag.Flag) {
		fn(viperFlagAdapter{pf: pf, viperPath: fs.vipathByPFName[pf.Name]})
	})
}
