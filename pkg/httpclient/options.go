package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"net"
)

type Option func(client *Client)

func WithBaseURL(addr string) Option {
	return func(c *Client) {
		c.SetBaseURL(addr)
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

func WithRealIP(ip net.IP) Option {
	return func(c *Client) {
		c.SetHeader("X-Real-IP", ip.String())
	}
}
