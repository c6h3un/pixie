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
	"github.com/owensengoku/pixie/internal/duck"
	"github.com/owensengoku/pixie/internal/httpsvc"
	"github.com/sirupsen/logrus"
)

func main() {
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

	args := os.Args[1:]
	switch kingpin.MustParse(app.Parse(args)) {

	case serve.FullCommand():
		log.Infof("args: %v", args)
		var g workgroup.Group

		flag.Parse()

		g.Add(ducksvc.Start)
		g.Run()
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
