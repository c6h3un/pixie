// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/heptio/workgroup"
	"github.com/owensengoku/pixie/internal/do"
	"github.com/owensengoku/pixie/internal/duck"
	"github.com/owensengoku/pixie/internal/httpsvc"
	"github.com/owensengoku/pixie/internal/tcping"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	log := logrus.StandardLogger()
	app := kingpin.New("pixie", "A simple application for test connection ability.")

	serve := app.Command("serve", "Serve for test connection")

	ducksvc := duck.Service{
		Service: httpsvc.Service{
			FieldLogger: log.WithField("context", "ducksvc"),
		},
	}

	serve.Flag("duck-address", "address the serve http endpoint will bind").Default("0.0.0.0").StringVar(&ducksvc.Addr)
	serve.Flag("duck-port", "port the serve http endpoint will bind").Default("8000").IntVar(&ducksvc.Port)

	dosvc := do.Service{
		Service: httpsvc.Service{
			FieldLogger: log.WithField("context", "dosvc"),
		},
	}

	serve.Flag("do-address", "address the serve http endpoint will bind").Default("0.0.0.0").StringVar(&dosvc.Addr)
	serve.Flag("do-port", "port the serve http endpoint will bind").Default("8001").IntVar(&dosvc.Port)

	ping := app.Command("ping", "ping for tcp or http(s)")

	pingProto := ping.Flag("ping-protocol", "protocol for ping").Default("http").String()
	pingAddr := ping.Flag("ping-address", "address of ping target").Default("127.0.0.1").String()
	pingPort := ping.Flag("ping-port", "port of ping target").Default("8000").Int()
	pingCounter := ping.Flag("ping-coutner", "repeat times for ping").Default("4").Int()
	pingTimeout := ping.Flag("ping-timeout", "timeout of ping").Default("1s").Duration()
	pingInterval := ping.Flag("ping-interval", "port of ping target").Default("1s").Duration()

	args := os.Args[1:]
	log.Infof("args: %v", args)
	switch kingpin.MustParse(app.Parse(args)) {

	case serve.FullCommand():

		var g workgroup.Group

		flag.Parse()

		g.Add(ducksvc.Start)
		g.Add(dosvc.Start)
		g.Run()
	case ping.FullCommand():
		pingWorker, err := tcping.NewWorker(*pingProto, *pingAddr, *pingPort, *pingCounter, *pingTimeout, *pingInterval)
		check(err)
		pingWorker.FieldLogger = log.WithField("context", "tcping")

		var stop <-chan struct{}
		pingWorker.Start(stop)
	default:
		app.Usage(args)
		os.Exit(2)
	}
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
