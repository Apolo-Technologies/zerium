// Copyright 2017 The zerium Authors
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
	"fmt"
	"sort"

	"github.com/abt/zerium/log"
)

// deployZrmstats queries the user for various input on deploying an zrmstats
// monitoring server, after which it executes it.
func (w *wizard) deployZrmstats() {
	// Select the server to interact with
	server := w.selectServer()
	if server == "" {
		return
	}
	client := w.servers[server]

	// Retrieve any active zrmstats configurations from the server
	infos, err := checkZrmstats(client, w.network)
	if err != nil {
		infos = &zrmstatsInfos{
			port:   80,
			host:   client.server,
			secret: "",
		}
	}
	// Figure out which port to listen on
	fmt.Println()
	fmt.Printf("Which port should zrmstats listen on? (default = %d)\n", infos.port)
	infos.port = w.readDefaultInt(infos.port)

	// Figure which virtual-host to deploy zrmstats on
	if infos.host, err = w.ensureVirtualHost(client, infos.port, infos.host); err != nil {
		log.Error("Failed to decide on zrmstats host", "err", err)
		return
	}
	// Port and proxy settings retrieved, figure out the secret and boot zrmstats
	fmt.Println()
	if infos.secret == "" {
		fmt.Printf("What should be the secret password for the API? (must not be empty)\n")
		infos.secret = w.readString()
	} else {
		fmt.Printf("What should be the secret password for the API? (default = %s)\n", infos.secret)
		infos.secret = w.readDefaultString(infos.secret)
	}
	// Gather any blacklists to ban from reporting
	fmt.Println()
	fmt.Printf("Keep existing IP %v blacklist (y/n)? (default = yes)\n", infos.banned)
	if w.readDefaultString("y") != "y" {
		// The user might want to clear the entire list, although generally probably not
		fmt.Println()
		fmt.Printf("Clear out blacklist and start over (y/n)? (default = no)\n")
		if w.readDefaultString("n") != "n" {
			infos.banned = nil
		}
		// Offer the user to explicitly add/remove certain IP addresses
		fmt.Println()
		fmt.Println("Which additional IP addresses should be blacklisted?")
		for {
			if ip := w.readIPAddress(); ip != "" {
				infos.banned = append(infos.banned, ip)
				continue
			}
			break
		}
		fmt.Println()
		fmt.Println("Which IP addresses should not be blacklisted?")
		for {
			if ip := w.readIPAddress(); ip != "" {
				for i, addr := range infos.banned {
					if ip == addr {
						infos.banned = append(infos.banned[:i], infos.banned[i+1:]...)
						break
					}
				}
				continue
			}
			break
		}
		sort.Strings(infos.banned)
	}
	// Try to deploy the zrmstats server on the host
	trusted := make([]string, 0, len(w.servers))
	for _, client := range w.servers {
		if client != nil {
			trusted = append(trusted, client.address)
		}
	}
	if out, err := deployZrmstats(client, w.network, infos.port, infos.secret, infos.host, trusted, infos.banned); err != nil {
		log.Error("Failed to deploy zrmstats container", "err", err)
		if len(out) > 0 {
			fmt.Printf("%s\n", out)
		}
		return
	}
	// All ok, run a network scan to pick any changes up
	w.networkStats(false)
}
