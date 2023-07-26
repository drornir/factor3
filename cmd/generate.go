package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("hello")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
