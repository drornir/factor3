package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/drornir/factor3/pkg/config"
	"github.com/drornir/factor3/pkg/example"
	"github.com/drornir/factor3/pkg/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const ProgramName = "waterboy"

var (
	flagConfigFile string
	flagLogFormat  string
	flagLogLevel   string

	viperInstance = viper.New()

	globalConfig       example.Config
	globalConfigLoader *config.Loader
	globalConfigLock   sync.RWMutex
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s", ProgramName),
	Short: "Productivity tool for developers that work with a Product Team",
	Long: `As developers, we all need to connect our coding to "The Business".
This means we need to open a ticket, a task or a card, however you name it, and
connect it to our Pull Request, or at least the branch name.

With the rise of popularity of SOC2, it's usually unavoidable.

waterboy is a cli tool meant to make use of this knowledge graph in order to streamline
the process of maintaining this graph, but also query it in the rare case that 'git blame' is not enough
	`,
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
	l, err := config.Bind(&globalConfig, viperInstance, RootCmd.Flags())
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
	if err := config.Initialize(config.InitArgs{
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
