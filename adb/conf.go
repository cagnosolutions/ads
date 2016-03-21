package adb

import "syscall"

var (
	SYS_PAGE = syscall.Getpagesize()
)
