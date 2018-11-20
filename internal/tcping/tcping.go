package tcping

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cloverstd/tcping/ping"
	"github.com/sirupsen/logrus"
)

// Work for doing something with tcping
type Worker struct {
	*ping.Target
	logrus.FieldLogger
	result chan string
}

type Arg struct {
	Protocol string        `json:"pingProto"`
	Host     string        `json:"pingAddr"`
	Port     int           `json:"pingPort"`
	Counter  int           `json:"pingCounter"`
	Timeout  time.Duration `json:"pingTimeout"`
	Interval time.Duration `json:"pingInterval"`
}

func Run(logger logrus.FieldLogger, arg *Arg) string {
	worker, err := NewWorker(arg.Protocol, arg.Host, arg.Port, arg.Counter, arg.Timeout, arg.Interval)
	if err != nil {
		panic(err)
	}
	worker.FieldLogger = logger
	var stop <-chan struct{}
	worker.Start(stop)
	return worker.Result()
}

func NewWorker(protocol, host string, port, counter int, timeout, intervel time.Duration) (*Worker, error) {
	p, err := ping.NewProtocol(protocol)
	if err != nil {
		return nil, err
	}

	return &Worker{
		Target: &ping.Target{
			Timeout:  timeout,
			Interval: intervel,
			Host:     host,
			Port:     port,
			Counter:  counter,
			Protocol: p,
		},
		result: make(chan string, 1),
	}, nil
}

func (w *Worker) Start(stop <-chan struct{}) (err error) {
	defer func() {
		if err != nil {
			w.WithError(err).Error("terminated with error")
		} else {
			w.Info("stopped")
		}
	}()

	var pinger ping.Pinger
	switch w.Target.Protocol {
	case ping.TCP:
		pinger = ping.NewTCPing()
	case ping.HTTP, ping.HTTPS:
		// Temporarily use GET as default
		pinger = ping.NewHTTPing(http.MethodGet)
	default:
		return fmt.Errorf("protocol: %s not support for pinger\n", w.Target.Protocol.String())
	}

	pinger.SetTarget(w.Target)
	pingerDone := pinger.Start()
	select {
	case <-pingerDone:
		break
	case <-stop:
		pinger.Stop()
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		break
	}
	r := pinger.Result().String()
	fmt.Println(r)
	w.result <- r
	return nil
}

func (w *Worker) Result() string {
	return <-w.result
}
