package zrm_terminal_console

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA

type Termios syscall.Termios
