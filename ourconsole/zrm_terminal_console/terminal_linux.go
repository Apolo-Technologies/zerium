// Based on ssh/terminal:
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine

package zrm_terminal_console

import "syscall"

const ioctlReadTermios = syscall.TCGETS

type Termios syscall.Termios
