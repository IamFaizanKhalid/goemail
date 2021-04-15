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

func (c *Client) NewMailer(elems *Mail) (*Mailer, error) {
	if err := validator.Validate(elems); err != nil {
		return nil, err
	}

	m := Mailer{
		subject: elems.Subject,
		message: elems.Message,
		client:  c,
	}

	if elems.From == nil {
		m.sender = c.config.Email
		m.from = c.config.Email
	} else {
		m.sender = elems.From.String()
		m.from = elems.From.Email
	}

	for _, user := range elems.To {
		m.recipients += user.String() + ", "
		m.to = append(m.to, user.Email)
	}
	m.recipients = m.recipients[:len(m.recipients)-2]

	for _, user := range elems.Cc {
		m.copyRecipients += user.String() + ", "
		m.to = append(m.to, user.Email)
	}
	m.copyRecipients = m.copyRecipients[:len(m.copyRecipients)-2]

	for _, user := range elems.Bcc {
		m.to = append(m.to, user.Email)
	}

	return &m, nil
}
