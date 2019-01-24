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

package do

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/owensengoku/pixie/internal/httpsvc"
	"github.com/owensengoku/pixie/internal/tcping"
	"github.com/sirupsen/logrus"
)

type Action struct {
	Command  string          `json:"command"`
	Argument json.RawMessage `json:"argument"`
}

type Response struct {
	Code    int
	Message string
	Result  string
}

// Service serves duck & healthcheck endpoints
type Service struct {
	httpsvc.Service
}

// Start fulfills the g.Start contract.
// When stop is closed the http server will shutdown.
func (svc *Service) Start(stop <-chan struct{}) error {
	registerDo(&svc.ServeMux, svc.FieldLogger)

	return svc.Service.Start(stop)
}

func registerDo(mux *http.ServeMux, logger logrus.FieldLogger) {
	mux.HandleFunc("/do", func(w http.ResponseWriter, r *http.Request) {
		var a Action
		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			panic(err)
		}

		switch a.Command {
		case "ping":
			arg := &tcping.Arg{
				Protocol: "tcp",
				Host:     "tw.yahoo.com",
				Port:     80,
				Counter:  1,
				Timeout:  time.Second * 1,
				Interval: time.Second * 1,
			}
			err := json.Unmarshal(a.Argument, arg)
			if err != nil {
				panic(err)
			}
			logger.Debugf("%#v", arg)
			ret := tcping.Run(logger, arg)
			resp := Response{
				Code:    0,
				Message: "OK",
				Result:  ret,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

	})
}
