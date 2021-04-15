package goemail

import (
	"fmt"
	"gopkg.in/validator.v2"
	"net/smtp"
)

type simpleMailer struct {
	client         *Client
	sender         string
	recipients     string
	copyRecipients string
	from           string
	to             []string
	subject        string
	message        string
}

func (m *simpleMailer) Send() error {
	message := []byte(fmt.Sprintf("From: %s\nTo: %s\nCc: %s\nSubject: %s\r\n\r\n%s\r\n", m.sender, m.recipients, m.copyRecipients, m.subject, m.message))
	return smtp.SendMail(m.client.addr, m.client.auth, m.from, m.to, message)
}

func (c *Client) NewSimpleMailer(mail *Mail) (Mailer, error) {
	if err := validator.Validate(mail); err != nil {
		return nil, err
	}

	m := simpleMailer{
		subject: mail.Subject,
		message: mail.Message,
		client:  c,
	}

	if mail.From == nil {
		m.sender = c.config.Email
		m.from = c.config.Email
	} else {
		m.sender = mail.From.String()
		m.from = mail.From.Email
	}

	for _, user := range mail.To {
		m.recipients += user.String() + ", "
		m.to = append(m.to, user.Email)
	}
	m.recipients = m.recipients[:len(m.recipients)-2]

	for _, user := range mail.Cc {
		m.copyRecipients += user.String() + ", "
		m.to = append(m.to, user.Email)
	}
	m.copyRecipients = m.copyRecipients[:len(m.copyRecipients)-2]

	for _, user := range mail.Bcc {
		m.to = append(m.to, user.Email)
	}

	return &m, nil
}
