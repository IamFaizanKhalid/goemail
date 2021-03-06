package goemail

import (
	"fmt"
	"net/smtp"
)

// NewClient returns mail Client using smtp credentials
func NewClient(config *Config) Client {
	return &mailClient{
		config: config,
		auth:   smtp.PlainAuth("", config.Email, config.Password, config.Host),
		addr:   fmt.Sprintf("%s:%d", config.Host, config.Port),
	}
}

type mailClient struct {
	config *Config
	auth   smtp.Auth
	addr   string
}

func (c *mailClient) newMailer(subject string, body string, contentType string) Mailer {
	m := &simpleMailer{
		client:  c,
		subject: subject,
		from:    c.config.Email,
		body:    body,
		defaultHeaders: mailHeaders{
			from:        c.config.Email,
			contentType: contentType,
		},
	}
	m.attachments = make(map[string]*attachment)
	return m
}

func (c *mailClient) NewMailer(subject string, message string) Mailer {
	return c.newMailer(subject, message, "text/plain")
}

func (c *mailClient) NewHtmlMailer(subject string, message string) Mailer {
	return c.newMailer(subject, message, "text/html")
}

func (c *mailClient) NewHtmlMailerFromTemplate(subject string, templateFile string, templateValues interface{}) (Mailer, error) {
	message, err := MessageFromHtmlTemplate(templateFile, templateValues)
	if err != nil {
		return nil, err
	}

	return c.newMailer(subject, message, "text/html"), nil
}
