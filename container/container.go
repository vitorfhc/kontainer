package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
)

// Container is the type that holds all the information about a container.
type Container struct {
	PID          int
	Namespaces   uintptr
	Entrypoint   []string
	Cmd          []string
	Env          []string
	UserMapping  []string
	GroupMapping []string
}

// NewWithDefaults returns a container with default values.
//
// It inherits all environment variables, uses user and mount namespace,
// and runs te command "id" with the entrypoint "/bin/bash -c".
// It also maps the inside root user and group to outside user and group 1000.
func NewWithDefaults() *Container {
	return &Container{
		Env:        os.Environ(),
		Namespaces: UserNamespace | MountNamespace,
		Entrypoint: []string{"/bin/bash", "-c"},
		Cmd:        []string{"id"},
		UserMapping: []string{
			"0 1000 1",
		},
		GroupMapping: []string{
			"0 1000 1",
		},
	}
}

// Run initializes the container by forking the process and setting the namespaces.
func (c *Container) Run() error {
	forkAttrs := &syscall.ProcAttr{
		Env:   c.Env,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys: &syscall.SysProcAttr{
			Unshareflags: c.Namespaces | UserNamespace | MountNamespace,
		},
	}

	execCmd := c.buildExecCommand()
	pid, err := syscall.ForkExec("/proc/self/exe", execCmd, forkAttrs)
	if err != nil {
		errWrapped := fmt.Errorf("Failed to intialize container with ForkExec: %w", err)
		return errWrapped
	}
	c.PID = pid

	err = c.writeUserMapping()
	if err != nil {
		errWrapped := fmt.Errorf("Failed to write user mapping: %w", err)
		return errWrapped
	}

	err = c.writeGroupMapping()
	if err != nil {
		errWrapped := fmt.Errorf("Failed to write group mapping: %w", err)
		return errWrapped
	}

	_, err = syscall.Wait4(pid, nil, 0, nil)
	if err != nil {
		errWrapped := fmt.Errorf("Failed to wait for container with Wait4: %w", err)
		return errWrapped
	}

	return nil
}

func (c *Container) writeUserMapping() error {
	uidMapPath := fmt.Sprintf("/proc/%d/uid_map", c.PID)
	toWrite := ""
	for _, userMapping := range c.UserMapping {
		toWrite += userMapping + "\n"
	}
	err := ioutil.WriteFile(uidMapPath, []byte(toWrite), 0644)
	return err
}

func (c *Container) writeGroupMapping() error {
	// https://unix.stackexchange.com/questions/692177/echo-to-gid-map-fails-but-uid-map-success
	setGroupsPath := fmt.Sprintf("/proc/%d/setgroups", c.PID)
	ioutil.WriteFile(setGroupsPath, []byte("deny"), 0644)

	gidMapPath := fmt.Sprintf("/proc/%d/gid_map", c.PID)
	toWrite := ""
	for _, groupMapping := range c.GroupMapping {
		toWrite += groupMapping + "\n"
	}
	err := ioutil.WriteFile(gidMapPath, []byte(toWrite), 0644)
	return err
}

func (c *Container) buildExecCommand() []string {
	execCommand := []string{os.Args[0], "exec", "--"}

	if len(c.Entrypoint) > 0 {
		execCommand = append(execCommand, c.Entrypoint...)
	}

	if len(c.Cmd) > 0 {
		execCommand = append(execCommand, c.Cmd...)
	}

	return execCommand
}

// Kill sends a SIGKILL to the container main process.
func (c *Container) Kill() error {
	err := syscall.Kill(c.PID, syscall.SIGKILL)
	return err
}
