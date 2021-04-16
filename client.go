package goemail

import (
	"fmt"
	"net/smtp"
)

func NewClient(config *Config) Client {
	return &mailClient{
		config: config,
		auth:   smtp.PlainAuth("", config.Email, config.Password, config.Host),
		addr:   fmt.Sprintf("%v:%v", config.Host, config.Port),
	}
}

type mailClient struct {
	config *Config
	auth   smtp.Auth
	addr   string
}

func (c *mailClient) newMailer(subject string, body string, contentType string) Mailer {
	m := &simpleMailer{
		subject: subject,
		from:    c.config.Email,
		body:    body,
		defaultHeaders: mailHeaders{
			from:        c.config.Email,
			contentType: contentType,
		},
	}
	m.attachments = make(map[string]*Attachment)
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
