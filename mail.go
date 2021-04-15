package goemail

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	client         *Client
	sender         string
	recipients     string
	copyRecipients string
	from           string
	to             []string
	subject        string
	message        string
}

func (m *Mailer) Send() error {
	message := []byte(fmt.Sprintf("From: %s\nTo: %s\nCc: %s\nSubject: %s\r\n\r\n%s\r\n", m.sender, m.recipients, m.copyRecipients, m.subject, m.message))
	return smtp.SendMail(m.client.addr, m.client.auth, m.from, m.to, message)
}
