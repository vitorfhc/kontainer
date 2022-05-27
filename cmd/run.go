package cmd

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vitorfhc/kontainer/container"
)

var (
	cntr       *container.Container
	entrypoint string
	command    string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a process in a containerized environment",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if strings.Trim(entrypoint, " ") != "" {
			cntr.Entrypoint = strings.Split(entrypoint, " ")
		}

		if strings.Trim(command, " ") != "" {
			cntr.Cmd = strings.Split(command, " ")
		}

		err := cntr.Run()
		if err != nil {
			kErr := cntr.Kill()
			if kErr != nil {
				logrus.Warn("Failed to kill container: ", kErr)
			}
			logrus.Fatalf("Failed to run container: %s", err)
		}
	},
}

func init() {
	cntr = container.NewWithDefaults()

	runCmd.Flags().StringVarP(&entrypoint, "entrypoint", "e", "", "Entrypoint of the container")
	runCmd.Flags().StringVarP(&command, "cmd", "c", "", "Command to execute")

	rootCmd.AddCommand(runCmd)
}
