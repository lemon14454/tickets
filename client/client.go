package main

import (
	"fmt"
	"net/http"
	"ticket/client/request"
	"time"
)

type Client struct {
	serverIP   string
	serverPort int
	Request    *request.Client
}

func NewClient(ip string, port int) *Client {
	url := fmt.Sprintf("http://%s:%d", ip, port)
	httpClient := &http.Client{Timeout: 10 * time.Second}

	return &Client{
		serverIP:   ip,
		serverPort: port,
		Request: &request.Client{
			Url:        url,
			HttpClient: *httpClient,
		},
	}
}
