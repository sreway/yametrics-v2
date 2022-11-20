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
	flag.StringVar(&config.DefaultDSN, "d", config.DefaultDSN, "PostgreSQL data source name")
	flag.StringVar(&config.DefaultCryptoKey, "crypto-key", config.DefaultCryptoKey,
		"x509 private key path")
	flag.StringVar(&config.DefaultCryptoCrt, "crypto-cert", config.DefaultCryptoCrt,
		"x509 certificate path")
	flag.StringVar(&config.DefaultConfigFile, "config", config.DefaultConfigFile,
		"json configuration file")
	flag.StringVar(&config.DefaultConfigFile, "c", config.DefaultConfigFile, "json configuration file")
	flag.Parse()

	srv, err := server.New()
	if err != nil {
		log.Fatalln(err)
	}
	srv.Run()
}
