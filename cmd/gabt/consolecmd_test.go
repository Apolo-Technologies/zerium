// Copyright 2016 The zerium Authors
// This file is part of zerium.
//
// zerium is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// zerium is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with zerium. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apolo-technologies/zerium/params"
)

const (
	ipcAPIs  = "admin:1.0 debug:1.0 zrm:1.0 miner:1.0 net:1.0 personal:1.0 rpc:1.0 shh:1.0 txpool:1.0 zae:1.0"
	httpAPIs = "zrm:1.0 net:1.0 rpc:1.0 zae:1.0"
)

// Tests that a node embedded within a zaeconsole can be started up properly and
// then terminated by closing the input stream.
func TestConsoleWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"

	// Start a zaed zaeconsole, make sure it's cleaned up and terminate the zaeconsole
	zaed := runzaed(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--zeriumbase", coinbase, "--shh",
		"zaeconsole")

	// Gather all the infos the welcome message needs to contain
	zaed.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	zaed.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	zaed.SetTemplateFunc("gover", runtime.Version)
	zaed.SetTemplateFunc("zaedver", func() string { return params.Version })
	zaed.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	zaed.SetTemplateFunc("apis", func() string { return ipcAPIs })

	// Verify the actual welcome message to the required template
	zaed.Expect(`
Welcome to the zaed JavaScript zaeconsole!

instance: zaed/v{{zaedver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{.Zeriumbase}}
at block: 0 ({{niltime}})
 datadir: {{.Datadir}}
 modules: {{apis}}

> {{.InputLine "exit"}}
`)
	zaed.ExpectExit()
}

// Tests that a zaeconsole can be attached to a running node via various means.
func TestIPCAttachWelcome(t *testing.T) {
	// Configure the instance for IPC attachement
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	var ipc string
	if runtime.GOOS == "windows" {
		ipc = `\\.\pipe\zaed` + strconv.Itoa(trulyRandInt(100000, 999999))
	} else {
		ws := tmpdir(t)
		defer os.RemoveAll(ws)
		ipc = filepath.Join(ws, "zaed.ipc")
	}
	// Note: we need --shh because testAttachWelcome checks for default
	// list of ipc modules and shh is included there.
	zaed := runzaed(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--zeriumbase", coinbase, "--shh", "--ipcpath", ipc)

	time.Sleep(2 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, zaed, "ipc:"+ipc, ipcAPIs)

	zaed.Interrupt()
	zaed.ExpectExit()
}

func TestHTTPAttachWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P
	zaed := runzaed(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--zeriumbase", coinbase, "--rpc", "--rpcport", port)

	time.Sleep(2 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, zaed, "http://localhost:"+port, httpAPIs)

	zaed.Interrupt()
	zaed.ExpectExit()
}

func TestWSAttachWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P

	zaed := runzaed(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--zeriumbase", coinbase, "--ws", "--wsport", port)

	time.Sleep(2 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, zaed, "ws://localhost:"+port, httpAPIs)

	zaed.Interrupt()
	zaed.ExpectExit()
}

func testAttachWelcome(t *testing.T, zaed *testzaed, endpoint, apis string) {
	// Attach to a running zaed note and terminate immediately
	attach := runzaed(t, "attach", endpoint)
	defer attach.ExpectExit()
	attach.CloseStdin()

	// Gather all the infos the welcome message needs to contain
	attach.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	attach.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	attach.SetTemplateFunc("gover", runtime.Version)
	attach.SetTemplateFunc("zaedver", func() string { return params.Version })
	attach.SetTemplateFunc("zeriumbase", func() string { return zaed.Zeriumbase })
	attach.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	attach.SetTemplateFunc("ipc", func() bool { return strings.HasPrefix(endpoint, "ipc") })
	attach.SetTemplateFunc("datadir", func() string { return zaed.Datadir })
	attach.SetTemplateFunc("apis", func() string { return apis })

	// Verify the actual welcome message to the required template
	attach.Expect(`
Welcome to the zaed JavaScript zaeconsole!

instance: zaed/v{{zaedver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{zeriumbase}}
at block: 0 ({{niltime}}){{if ipc}}
 datadir: {{datadir}}{{end}}
 modules: {{apis}}

> {{.InputLine "exit" }}
`)
	attach.ExpectExit()
}

// trulyRandInt generates a crypto random integer used by the zaeconsole tests to
// not clash network ports with other tests running cocurrently.
func trulyRandInt(lo, hi int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(hi-lo)))
	return int(num.Int64()) + lo
}
