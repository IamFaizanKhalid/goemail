package goemail

import (
	"fmt"
	"gopkg.in/validator.v2"
	"net/smtp"
)

func NewClient(config *Config) (*Client, error) {
	return &Client{
		config: config,
		auth:   smtp.PlainAuth("", config.Email, config.Password, config.Host),
		addr:   fmt.Sprintf("%v:%v", config.Host, config.Port),
	}, validator.Validate(config)
}

type Client struct {
	config *Config
	auth   smtp.Auth
	addr   string
}
