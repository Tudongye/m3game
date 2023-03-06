package util

import "syscall"

func SysCallSIGUSR1() syscall.Signal {
	return syscall.SIGINT
}
func SysCallSIGUSR2() syscall.Signal {
	return syscall.SIGINT
}
