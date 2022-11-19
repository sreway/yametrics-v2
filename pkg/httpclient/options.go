package httpclient

import (
	"crypto/tls"
	"crypto/x509"
)

type Option func(client *Client)

func WithBaseURL(addr string) Option {
	return func(s *Client) {
		s.SetBaseURL(addr)
	}
}

func WithCerts(certs ...*x509.Certificate) Option {
	tlsConfig := tls.Config{}
	tlsConfig.RootCAs = x509.NewCertPool()
	for _, cert := range certs {
		tlsConfig.RootCAs.AddCert(cert)
	}
	return func(c *Client) {
		c.SetTLSClientConfig(&tlsConfig)
	}
}
