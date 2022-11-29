package pem

import (
	"crypto/tls"
	"encoding/pem"
	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
	"os"
)

func ParsePEM(path string) (*tls.Certificate, error) {
	var (
		cert  tls.Certificate
		block *pem.Block
	)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, data = pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		log.Fatal("failed to decode PEM block containing certificate key")
	}

	cert.Certificate = append(cert.Certificate, block.Bytes)

	return &cert, nil
}
