// Copyright 2019-2020 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/nats-io/nats.go"
	"gopkg.in/alecthomas/kingpin.v2"
)

type actCmd struct{}

func configureActCommand(app *kingpin.Application) {
	c := &actCmd{}
	act := app.Command("account", "Account information and status")
	act.Command("info", "Account information").Alias("nfo").Action(c.infoAction)
}

func (c *actCmd) infoAction(pc *kingpin.ParseContext) error {
	nc, mgr, err := prepareHelper("", natsOpts()...)
	kingpin.FatalIfError(err, "setup failed")

	id, _ := nc.GetClientID()
	ip, _ := nc.GetClientIP()
	rtt, _ := nc.RTT()

	fmt.Println("Connection Information:")
	fmt.Println()
	fmt.Printf("               Client ID: %v\n", id)
	fmt.Printf("               Client IP: %v\n", ip)
	fmt.Printf("                     RTT: %v\n", rtt)
	fmt.Printf("       Headers Supported: %v\n", nc.HeadersSupported())
	fmt.Printf("         Maximum Payload: %v\n", humanize.IBytes(uint64(nc.MaxPayload())))
	fmt.Printf("           Connected URL: %v\n", nc.ConnectedUrl())
	fmt.Printf("       Connected Address: %v\n", nc.ConnectedAddr())
	fmt.Printf("     Connected Server ID: %v\n", nc.ConnectedServerId())
	if nc.ConnectedServerId() != nc.ConnectedServerName() {
		fmt.Printf("   Connected Server Name: %v\n", nc.ConnectedServerName())
	}

	fmt.Println()

	info, err := mgr.JetStreamAccountInfo()
	fmt.Println("JetStream Account Information:")
	fmt.Println()
	switch err {
	case nil:
		fmt.Printf("           Memory: %s of %s\n", humanize.IBytes(info.Memory), humanize.IBytes(uint64(info.Limits.MaxMemory)))
		fmt.Printf("          Storage: %s of %s\n", humanize.IBytes(info.Store), humanize.IBytes(uint64(info.Limits.MaxStore)))

		if info.Limits.MaxStreams == -1 {
			fmt.Printf("          Streams: %s of Unlimited\n", humanize.Comma(int64(info.Streams)))
		} else {
			fmt.Printf("          Streams: %s of %s\n", humanize.Comma(int64(info.Streams)), humanize.Comma(int64(info.Limits.MaxStreams)))
		}

		if info.Limits.MaxConsumers == -1 {
			fmt.Printf("        Consumers: %s of Unlimited\n", humanize.Comma(int64(info.Consumers)))
		} else {
			fmt.Printf("        Consumers: %s of %s\n", humanize.Comma(int64(info.Consumers)), humanize.Comma(int64(info.Limits.MaxConsumers)))
		}

	case context.DeadlineExceeded:
		fmt.Printf("   No response from JetStream server")
	case nats.ErrNoResponders:
		fmt.Printf("   JetStream is not supported in this account")
	default:
		fmt.Printf("   Could not obtain account information: %s", err)
	}

	fmt.Println()

	return nil
}
