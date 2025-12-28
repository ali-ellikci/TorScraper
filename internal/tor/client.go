package tor

import (
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

func NewTorClient() (*http.Client, error) {
	dialer, err := proxy.SOCKS5(
		"tcp",
		"127.0.0.1:9050",
		nil,
		proxy.Direct,
	)

	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}

	return client, nil

}
