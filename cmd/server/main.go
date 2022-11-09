package main

import (
	"flag"
	"log"

	"github.com/sreway/yametrics-v2/services/server/config"
	"github.com/sreway/yametrics-v2/services/server/server"
)

func main() {
	flag.StringVar(&config.DefaultAddress, "a", config.DefaultAddress, "address: host:port")
	flag.DurationVar(&config.DefaultStoreInterval, "i", config.DefaultStoreInterval, "store interval")
	flag.BoolVar(&config.DefaultRestore, "r", config.DefaultRestore, "restoring metrics at startup")
	flag.StringVar(&config.DefaultStoreFile, "f", config.DefaultStoreFile, "store file")
	flag.StringVar(&config.DefaultKey, "k", config.DefaultKey, "encrypt key")
	flag.StringVar(&config.DefaultDSN, "d", config.DefaultDSN, "PosgreSQL data source name")
	flag.Parse()

	srv, err := server.New()
	if err != nil {
		log.Fatalln(err)
	}
	srv.Run()
}
