package main

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/drornir/factor3/pkg/factor3"
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
	rootConfig Config
	rootCmd    = &cobra.Command{
		Use: "myprogram",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("# config = %#v\n", rootConfig)
		},
	}
)

func init() {
	viperInstance := viper.New()
	err := factor3.InitializeViper(factor3.InitArgs{
		Viper:       viperInstance,
		ProgramName: "myprogram",
		CfgFile:     "tests/configs/config.json",
	})
	cobra.CheckErr(err)

	pflags := rootCmd.Flags()
	loader, err := factor3.Bind(&rootConfig, viperInstance, pflags)
	cobra.CheckErr(err)

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
