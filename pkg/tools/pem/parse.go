package pem

import (
	"crypto/tls"
	"encoding/pem"
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

	for {
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, block.Bytes)
		}
	}

	return &cert, nil
}
