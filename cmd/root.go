package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/drornir/factor3/pkg/example"
	"github.com/drornir/factor3/pkg/factor3"
	"github.com/drornir/factor3/pkg/log"
)

const ProgramName = "example"

var (
	flagConfigFile string
	flagLogFormat  string
	flagLogLevel   string

	viperInstance = viper.New()

	globalConfig       example.Config
	globalConfigLoader *factor3.Loader
	globalConfigLock   sync.RWMutex
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s", ProgramName),
	Short: "example",
	Long:  "",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Main() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// setup core global flags before reading config
	RootCmd.PersistentFlags().StringVarP(&flagConfigFile, "config", "c", "", "config file (default is $XDG_CONFIG_HOME/"+ProgramName+"/config[.yaml])")
	RootCmd.PersistentFlags().StringVarP(&flagLogFormat, "log-format", "", "logfmt", "either 'logfmt' or 'json'")
	RootCmd.PersistentFlags().StringVarP(&flagLogLevel, "log-level", "l", "info", "'trace', 'debug', 'info', 'warn[ing]', 'error'")

	// setup global logger before reading config using viper
	cobra.OnInitialize(initLogger)

	// setup reading config file and
	l, err := factor3.Bind(&globalConfig, viperInstance, RootCmd.Flags())
	if err != nil {
		cobra.CheckErr(fmt.Errorf("config.Bind: %w", err))
	}
	globalConfigLoader = l
	cobra.OnInitialize(initViper)
}

func initInitViper() func() {
	return func() { initViper() }
}

func initViper() {
	if err := factor3.InitializeViper(factor3.InitArgs{
		Viper:       viperInstance,
		ProgramName: ProgramName,
		CfgFile:     flagConfigFile,
	}); err != nil {
		cobra.CheckErr(fmt.Errorf("config.Initialize: %w", err))
	}

	globalConfigLock.Lock()
	defer globalConfigLock.Unlock()
	if err := globalConfigLoader.Load(); err != nil {
		err = fmt.Errorf("config.Load: error loading config: %w", err)
		log.GG().E(context.TODO(), "loading config", "error", err)
		cobra.CheckErr(err)
	}
	example.SetGlobal(globalConfig)
	viperInstance.OnConfigChange(func(in fsnotify.Event) {
		globalConfigLock.Lock()
		defer globalConfigLock.Unlock()
		if err := globalConfigLoader.Load(); err != nil {
			log.GG().E(context.TODO(), "error loading config file",
				"file_name", in.Name, "error", err)
			return
		}
		example.SetGlobal(globalConfig)
	})
}

func initLogger() {
	logOut := os.Stdout // TODO  configurable

	var sloggerHandler slog.Handler
	sloggerOpts := &slog.HandlerOptions{
		// AddSource: log.ParseLevel(flagLogLevel) <= slog.LevelDebug,
		AddSource:   false,
		Level:       log.PflagLeveler{Flag: &flagLogLevel},
		ReplaceAttr: log.SlogReplacerMinimal(),
	}
	if flagLogFormat == "logfmt" || flagLogFormat == "" {
		sloggerHandler = slog.NewTextHandler(logOut, sloggerOpts)
	} else {
		sloggerHandler = slog.NewJSONHandler(logOut, sloggerOpts)
	}

	slogger := slog.New(sloggerHandler)
	log.SetGlobal(log.WrapSlogger(slogger))
}
