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

	fmt.Println(pinger.Result())
	return nil
}
