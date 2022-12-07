package server

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sreway/yametrics-v2/services/server/internal/delivery/grpc"
	"github.com/sreway/yametrics-v2/services/server/internal/delivery/http"

	"github.com/sreway/yametrics-v2/services/server/config"
	"github.com/sreway/yametrics-v2/services/server/internal/usecases/adapters/delivery"
	"github.com/sreway/yametrics-v2/services/server/internal/usecases/adapters/storage"
	metricService "github.com/sreway/yametrics-v2/services/server/internal/usecases/metric"

	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
)

type Server struct {
	config   *config.Config
	delivery delivery.Delivery
	storage  storage.Storage
}

func (s *Server) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	exitch := make(chan int)
	wg := new(sync.WaitGroup)
	wg.Add(1)

	err := s.InitStorage(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	service := metricService.New(s.storage, s.config.SecretKey)

	if s.config.UseGRPC {
		if s.delivery, err = grpc.New(service, &s.config.Delivery); err != nil {
			log.Fatal(err.Error())
		}
	} else {
		if s.delivery, err = http.New(service, &s.config.Delivery); err != nil {
			log.Fatal(err.Error())
		}
	}

	go func() {
		defer wg.Done()
		err = s.delivery.Run(ctx, &s.config.Delivery)
		if err != nil {
			log.Error(err.Error())
			signals <- syscall.SIGSTOP
		}
	}()

	go func() {
		for {
			systemSignal := <-signals
			switch systemSignal {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Info("Server: signal triggered")
				if store, ok := s.storage.(storage.MemoryStorage); ok {
					if s.config.MemoryStorage.StoreFile != "" {
						err = store.Store()
						if err != nil {
							log.Info(err.Error())
						}
					}
				}
				exitch <- 0
			default:
				log.Error("Server: signal triggered")
				exitch <- 1
			}
		}
	}()

	exitCode := <-exitch
	cancel()
	wg.Wait()

	err = s.storage.Close()

	if err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(exitCode)
}

func New(opts ...config.OptionServer) (*Server, error) {
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

	s := new(Server)
	s.config = cfg
	return s, nil
}
