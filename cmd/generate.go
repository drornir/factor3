package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/drornir/factor3/pkg/generate"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := LoadEnv()
		if err != nil {
			return fmt.Errorf("factor3: error loading env: %w", err)
		}
		app, err := generate.New(env.WorkingDirectory)
		if err != nil {
			return fmt.Errorf("factor3: error initializing app: %w", err)
		}

		if err := app.Generate(); err != nil {
			return fmt.Errorf("factor3: error generating code: %w", err)
		}

		return nil
	},
}
