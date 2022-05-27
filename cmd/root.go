package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kontainer",
	Short: "A minimal container creator",
	Long: `Kontainer is a minimal container creator built with learning purpose.

This is NOT a production ready tool.
I wouldn't say it's even a development tool.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
