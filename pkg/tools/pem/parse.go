package pem

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
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

	block, _ = pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode PEM block containing certificate key")
	}

	cert.Certificate = append(cert.Certificate, block.Bytes)

	return &cert, nil
}
