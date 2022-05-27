package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vitorfhc/kontainer/container"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a process in a containerized environment",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cntr := &container.Container{
			Entrypoint: []string{"/bin/bash", "-c"},
			Cmd:        []string{"echo hello world"},
			Env:        os.Environ(),
			Namespaces: container.UserNamespace,
		}
		cntr.Init()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
