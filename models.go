package goemail

import (
	"fmt"
	"os"
)

type Client interface {
	NewMailer(subject string, body string) Mailer
	NewHtmlMailer(subject string, body string) Mailer
	NewHtmlMailerFromTemplate(subject string, templateFile string, templateValues interface{}) (Mailer, error)
}

type Mailer interface {
	AddRecipients(emails []User)
	AddCopyRecipients(emails []User)
	AddBlindCopyRecipients(emails []User)
	AddHeader(key string, value string)
	AddInlineFile(file *os.File) error
	SetSender(u User)
	SetReplyToEmail(email string)
	SetSubject(subject string)
	AttachFile(file *os.File) error
	AttachFileBytes(fileName string, binary []byte)
	Send() error
}

type Config struct {
	Host     string
	Port     string
	Email    string
	Password string
}

type Attachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

type Header struct {
	Key   string
	Value string
}

type User struct {
	Name  string
	Email string
}

func (u *User) String() string {
	if u.Name == "" {
		return u.Email
	}
	return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}
