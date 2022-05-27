package container

import (
	"fmt"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Container is the type that holds all the information about a container
type Container struct {
	Entrypoint []string
	Cmd        []string
	Env        []string
	Namespaces uintptr
}

// Init initializes the container by forking the process and setting the namespaces
func (c *Container) Init() error {
	forkAttrs := &syscall.ProcAttr{
		Env:   c.Env,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys: &syscall.SysProcAttr{
			Unshareflags: c.Namespaces,
		},
	}
	forkArgs := []string{os.Args[0], "run", "--child"}
	pid, err := syscall.ForkExec("/proc/self/exe", forkArgs, forkAttrs)
	if err != nil {
		errWrapped := fmt.Errorf("Failed to intialize container with ForkExec: %w", err)
		return errWrapped
	}
	logrus.Infof("Initialized container with PID %d", pid)

	_, err = syscall.Wait4(pid, nil, 0, nil)
	if err != nil {
		errWrapped := fmt.Errorf("Failed to wait for container with Wait4: %w", err)
		return errWrapped
	}

	return nil
}
