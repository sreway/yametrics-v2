package main

import (
	"flag"
	"fmt"

	"github.com/sreway/yametrics-v2/services/agent/agent"
	"github.com/sreway/yametrics-v2/services/agent/config"

	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func buildInfo() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildVersion, buildDate, buildCommit)
}

func init() {
	buildInfo()
	flag.StringVar(&config.DefaultServerAddress, "a", config.DefaultServerAddress,
		"server address: host:port")
	flag.DurationVar(&config.DefaultReportInterval, "r", config.DefaultReportInterval, "report interval")
	flag.DurationVar(&config.DefaultPollInterval, "p", config.DefaultPollInterval, "poll interval")
	flag.StringVar(&config.DefaultSecretKey, "k", config.DefaultSecretKey, "encrypt key")
	flag.StringVar(&config.DefaultServerPublicKey, "crypto-key", config.DefaultServerPublicKey,
		"server public key")
	flag.StringVar(&config.DefaultRealIP, "ip", config.DefaultRealIP, "agent IP address")
	flag.StringVar(&config.DefaultConfigFile, "config", config.DefaultConfigFile,
		"json configuration file")
	flag.StringVar(&config.DefaultConfigFile, "c", config.DefaultConfigFile, "json configuration file")
	flag.BoolVar(&config.DefaultUseGRPC, "grpc", config.DefaultUseGRPC, "use grpc")
	flag.Parse()
}

func main() {
	a, err := agent.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	a.Run()
}
