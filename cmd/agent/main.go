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
	flag.Parse()
}

func main() {
	a, err := agent.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	a.Run()
}
