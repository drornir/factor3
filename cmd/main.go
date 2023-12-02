package cmd

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "factor3",
	SilenceErrors: true,
}

func init() {
}

func Main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		log.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
}
