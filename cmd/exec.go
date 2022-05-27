package cmd

import (
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Executes the command given as argument",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := syscall.ForkExec(args[0], args, &syscall.ProcAttr{
			Env: os.Environ(),
			Files: []uintptr{
				os.Stdin.Fd(),
				os.Stdout.Fd(),
				os.Stderr.Fd(),
			},
		})
		if err != nil {
			logrus.Fatalf("Failed to execute command: %s", err)
		}
		_, err = syscall.Wait4(pid, nil, 0, nil)
		if err != nil {
			logrus.Fatalf("Failed to wait for command with Wait4: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
