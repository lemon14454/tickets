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

	tr := &http.Transport{
		ResponseHeaderTimeout: time.Hour,
		MaxConnsPerHost:   99999,
		DisableKeepAlives: true,
	}

	httpClient := &http.Client{
		Transport: tr,
	}

	return &Client{
		serverIP:   ip,
		serverPort: port,
		Request: &request.Client{
			Url:        url,
			HttpClient: *httpClient,
		},
	}
}

func (client *Client) Login(username, password string) error {
	res, err := client.Request.LoginUser(username, password)

	if err != nil {
		return err
	}

	client.Request.AccessToken = &res.AccessToken
	client.Request.RefreshToken = &res.RefreshToken

	return nil
}

func (client *Client) Register(username, password string) error {

	mail := fmt.Sprintf("%s@mail.com", username)
	_, err := client.Request.RegisterUser(username, mail, password)

	if err != nil {
		return err
	}

	return nil
}
