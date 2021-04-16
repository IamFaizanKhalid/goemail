package goemail

import (
	"net/mail"
	"os"
)

type Client interface {
	NewMailer(subject string, body string) Mailer
	NewHtmlMailer(subject string, body string) Mailer
	NewHtmlMailerFromTemplate(subject string, templateFile string, templateValues interface{}) (Mailer, error)
}

type Mailer interface {
	AddRecipients(emails []mail.Address)
	AddCopyRecipients(emails []mail.Address)
	AddBlindCopyRecipients(emails []mail.Address)
	AddHeader(key string, value string)
	AddInlineFile(filePath string) error
	SetSender(u mail.Address)
	SetReplyToEmail(email string)
	SetSubject(subject string)
	AttachFile(filePath string) error
	AttachOpenedFile(file *os.File) error
	AttachFileBytes(fileName string, binary []byte)
	Send() error
}

type Config struct {
	Host     string
	Port     int
	Email    string
	Password string
}
