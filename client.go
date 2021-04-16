package goemail

import (
	"fmt"
	"net/smtp"
)

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		auth:   smtp.PlainAuth("", config.Email, config.Password, config.Host),
		addr:   fmt.Sprintf("%v:%v", config.Host, config.Port),
	}
}

type Client struct {
	config *Config
	auth   smtp.Auth
	addr   string
}
