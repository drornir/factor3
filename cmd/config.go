package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/drornir/factor3/pkg/example"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

// ConfigCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the current configuration",
	Long:  `Print the current configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("all keys", viperInstance.AllKeys())
		fmt.Println("all settings reconstructed")
		fmt.Println("---")
		yenc := yaml.NewEncoder(os.Stdout)
		yenc.SetIndent(2)
		yenc.Encode(viperInstance.AllSettings())
		// fmt.Println("---")
		// viperInstance.Debug()
		fmt.Println("---")
		fmt.Println("my own:")
		conf := example.Global()
		yenc.Encode(conf)
		fmt.Println("---")
		jenc := json.NewEncoder(os.Stdout)
		jenc.SetIndent("", "  ")
		jenc.Encode(conf)
		fmt.Println("---")
	},
}

func init() {
	RootCmd.AddCommand(ConfigCmd)
}
