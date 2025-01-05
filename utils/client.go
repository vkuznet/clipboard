package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ServerUrl provides server url based on server configuration
func ServerUrl(config *Config) string {
	rurl := fmt.Sprintf("http://localhost:%d", config.Port)
	if config.ServerCert != "" {
		rurl = strings.Replace(rurl, "http://", "https://", -1)
	}
	return rurl
}

// Client returns proper HTTP client with pool of root CAs based on server configuration
func Client(config *Config) *http.Client {
	caCert, err := ioutil.ReadFile(config.ServerCert) // Load the self-signed cert
	if err != nil {
		fmt.Println("Error reading certificate:", err)
		return nil
	}

	// Create a CertPool and add the certificate
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a custom HTTPS client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	return client
}
