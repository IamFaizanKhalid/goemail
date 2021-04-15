package mail

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	client     *Client
	sender     string
	recipients string
	from       string
	to         []string
	subject    string
	message    string
}

func (m *Mailer) Send() error {
	message := []byte(fmt.Sprintf("To: %s\nFrom: %s\nSubject: %s\r\n\r\n%s\r\n", m.recipients, m.sender, m.subject, m.message))
	return smtp.SendMail(m.client.addr, m.client.auth, m.from, m.to, message)
}
