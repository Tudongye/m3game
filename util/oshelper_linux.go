package util

import "syscall"

func SysCallSIGUSR1() syscall.Signal {
	return syscall.SIGUSR1
}
func SysCallSIGUSR2() syscall.Signal {
	return syscall.SIGUSR2
}
