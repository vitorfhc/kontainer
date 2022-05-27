package container

import "syscall"

// Namespace flags copy from Linux flags
const (
	UserNamespace    uintptr = syscall.CLONE_NEWUSER
	PIDNamespace     uintptr = syscall.CLONE_NEWPID
	IPCNamespace     uintptr = syscall.CLONE_NEWIPC
	UTSNamespace     uintptr = syscall.CLONE_NEWUTS
	MountNamespace   uintptr = syscall.CLONE_NEWNS
	NetworkNamespace uintptr = syscall.CLONE_NEWNET
)
