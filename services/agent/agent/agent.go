package agent

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sreway/yametrics-v2/services/agent/internal/usecases/sender/grpc"
	"github.com/sreway/yametrics-v2/services/agent/internal/usecases/sender/http"

	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
	"github.com/sreway/yametrics-v2/services/agent/config"
	"github.com/sreway/yametrics-v2/services/agent/internal/usecases"
	"github.com/sreway/yametrics-v2/services/agent/internal/usecases/collector"
)

type Agent struct {
	config    *config.Config
	collector usecases.Collector
	sender    usecases.Sender
}

func (a *Agent) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	exitch := make(chan int)
	wg := new(sync.WaitGroup)
	wg.Add(4)

	go a.Collect(ctx, wg)
	go a.Send(ctx, wg)

	go func() {
		for {
			systemSignal := <-signals
			switch systemSignal {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Info("signal triggered")
				exitch <- 0
			default:
				log.Info("unknown signal")
				exitch <- 1
			}
		}
	}()

	exitCode := <-exitch
	cancel()
	wg.Wait()
	err := a.sender.Close()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(exitCode)
}

func (a *Agent) Collect(ctx context.Context, wg *sync.WaitGroup) {
	a.collector.Collect(ctx, wg, a.config.PollInterval)
}

func (a *Agent) Send(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	tick := time.NewTicker(a.config.ReportInterval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			metrics := a.collector.Expose()
			err := a.sender.Send(ctx, metrics)
			if err != nil {
				log.Error(err.Error())
			} else {
				log.Info("Success send metrics")
				a.collector.ResetCounter()
			}

		case <-ctx.Done():
			return
		}
	}
}

func New(opts ...config.OptionAgent) (*Agent, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		err = opt(cfg)
		if err != nil {
			return nil, err
		}
	}

	a := new(Agent)
	a.config = cfg
	a.collector = collector.New(cfg.Key)
	if cfg.UseGRPC {
		a.sender, err = grpc.New(cfg)
		if err != nil {
			return nil, err
		}
		return a, nil
	}
	a.sender, err = http.New(cfg)
	if err != nil {
		return nil, err
	}
	return a, nil
}
