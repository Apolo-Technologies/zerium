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
	"os"
	"os/signal"
	"strings"

	"github.com/apolo-technologies/zerium/cmd/utils"
	"github.com/apolo-technologies/zerium/abtconsole"
	"github.com/apolo-technologies/zerium/node"
	"github.com/apolo-technologies/zerium/rpc"
	"gopkg.in/urfave/cli.v1"
)

var (
	abtconsoleFlags = []cli.Flag{utils.JSpathFlag, utils.ExecFlag, utils.PreloadJSFlag}

	abtconsoleCommand = cli.Command{
		Action:   utils.MigrateFlags(localConsole),
		Name:     "abtconsole",
		Usage:    "Start an interactive JavaScript environment",
		Flags:    append(append(append(nodeFlags, rpcFlags...), abtconsoleFlags...), whisperFlags...),
		Category: "CONSOLE COMMANDS",
		Description: `
The Gabt abtconsole is an interactive shell for the JavaScript runtime environment
which exposes a node admin interface as well as the Ðapp JavaScript API.
See https://github.com/apolo-technologies/zerium/wiki/Javascipt-Console.`,
	}

	attachCommand = cli.Command{
		Action:    utils.MigrateFlags(remoteConsole),
		Name:      "attach",
		Usage:     "Start an interactive JavaScript environment (connect to node)",
		ArgsUsage: "[endpoint]",
		Flags:     append(abtconsoleFlags, utils.DataDirFlag),
		Category:  "CONSOLE COMMANDS",
		Description: `
The Gabt abtconsole is an interactive shell for the JavaScript runtime environment
which exposes a node admin interface as well as the Ðapp JavaScript API.
See https://github.com/apolo-technologies/zerium/wiki/Javascipt-Console.
This command allows to open a abtconsole on a running gabt node.`,
	}

	javascriptCommand = cli.Command{
		Action:    utils.MigrateFlags(ephemeralConsole),
		Name:      "js",
		Usage:     "Execute the specified JavaScript files",
		ArgsUsage: "<jsfile> [jsfile...]",
		Flags:     append(nodeFlags, abtconsoleFlags...),
		Category:  "CONSOLE COMMANDS",
		Description: `
The JavaScript VM exposes a node admin interface as well as the Ðapp
JavaScript API. See https://github.com/apolo-technologies/zerium/wiki/Javascipt-Console`,
	}
)

// localConsole starts a new gabt node, attaching a JavaScript abtconsole to it at the
// same time.
func localConsole(ctx *cli.Context) error {
	// Create and start the node based on the CLI flags
	node := makeFullNode(ctx)
	startNode(ctx, node)
	defer node.Stop()

	// Attach to the newly started node and start the JavaScript abtconsole
	client, err := node.Attach()
	if err != nil {
		utils.Fatalf("Failed to attach to the inproc gabt: %v", err)
	}
	config := abtconsole.Config{
		DataDir: utils.MakeDataDir(ctx),
		DocRoot: ctx.GlobalString(utils.JSpathFlag.Name),
		Client:  client,
		Preload: utils.MakeConsolePreloads(ctx),
	}

	abtconsole, err := abtconsole.New(config)
	if err != nil {
		utils.Fatalf("Failed to start the JavaScript abtconsole: %v", err)
	}
	defer abtconsole.Stop(false)

	// If only a short execution was requested, evaluate and return
	if script := ctx.GlobalString(utils.ExecFlag.Name); script != "" {
		abtconsole.Evaluate(script)
		return nil
	}
	// Otherwise print the welcome screen and enter interactive mode
	abtconsole.Welcome()
	abtconsole.Interactive()

	return nil
}

// remoteConsole will connect to a remote gabt instance, attaching a JavaScript
// abtconsole to it.
func remoteConsole(ctx *cli.Context) error {
	// Attach to a remotely running gabt instance and start the JavaScript abtconsole
	client, err := dialRPC(ctx.Args().First())
	if err != nil {
		utils.Fatalf("Unable to attach to remote gabt: %v", err)
	}
	config := abtconsole.Config{
		DataDir: utils.MakeDataDir(ctx),
		DocRoot: ctx.GlobalString(utils.JSpathFlag.Name),
		Client:  client,
		Preload: utils.MakeConsolePreloads(ctx),
	}

	abtconsole, err := abtconsole.New(config)
	if err != nil {
		utils.Fatalf("Failed to start the JavaScript abtconsole: %v", err)
	}
	defer abtconsole.Stop(false)

	if script := ctx.GlobalString(utils.ExecFlag.Name); script != "" {
		abtconsole.Evaluate(script)
		return nil
	}

	// Otherwise print the welcome screen and enter interactive mode
	abtconsole.Welcome()
	abtconsole.Interactive()

	return nil
}

// dialRPC returns a RPC client which connects to the given endpoint.
// The check for empty endpoint implements the defaulting logic
// for "gabt attach" and "gabt monitor" with no argument.
func dialRPC(endpoint string) (*rpc.Client, error) {
	if endpoint == "" {
		endpoint = node.DefaultIPCEndpoint(clientIdentifier)
	} else if strings.HasPrefix(endpoint, "rpc:") || strings.HasPrefix(endpoint, "ipc:") {
		// Backwards compatibility with gabt < 1.5 which required
		// these prefixes.
		endpoint = endpoint[4:]
	}
	return rpc.Dial(endpoint)
}

// ephemeralConsole starts a new gabt node, attaches an ephemeral JavaScript
// abtconsole to it, executes each of the files specified as arguments and tears
// everything down.
func ephemeralConsole(ctx *cli.Context) error {
	// Create and start the node based on the CLI flags
	node := makeFullNode(ctx)
	startNode(ctx, node)
	defer node.Stop()

	// Attach to the newly started node and start the JavaScript abtconsole
	client, err := node.Attach()
	if err != nil {
		utils.Fatalf("Failed to attach to the inproc gabt: %v", err)
	}
	config := abtconsole.Config{
		DataDir: utils.MakeDataDir(ctx),
		DocRoot: ctx.GlobalString(utils.JSpathFlag.Name),
		Client:  client,
		Preload: utils.MakeConsolePreloads(ctx),
	}

	abtconsole, err := abtconsole.New(config)
	if err != nil {
		utils.Fatalf("Failed to start the JavaScript abtconsole: %v", err)
	}
	defer abtconsole.Stop(false)

	// Evaluate each of the specified JavaScript files
	for _, file := range ctx.Args() {
		if err = abtconsole.Execute(file); err != nil {
			utils.Fatalf("Failed to execute %s: %v", file, err)
		}
	}
	// Wait for pending callbacks, but stop for Ctrl-C.
	abort := make(chan os.Signal, 1)
	signal.Notify(abort, os.Interrupt)

	go func() {
		<-abort
		os.Exit(0)
	}()
	abtconsole.Stop(true)

	return nil
}
