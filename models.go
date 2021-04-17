package goemail

import (
	"net/mail"
	"os"
)

// Client is the mail client which holds the connection details to the smtp server
type Client interface {
	// NewMailer returns a Mailer object to send an email with the given subject and body
	NewMailer(subject string, body string) Mailer

	// NewHtmlMailer returns a Mailer object to send an email with the given subject and html body
	NewHtmlMailer(subject string, body string) Mailer

	// NewHtmlMailerFromTemplate returns a Mailer object to send an email with the given subject
	// and by using the given html template as the body.
	NewHtmlMailerFromTemplate(subject string, templateFile string, templateValues interface{}) (Mailer, error)
}

// Mailer can be used to build and send an email
type Mailer interface {
	// AddRecipients adds `To` to the email
	AddRecipients(emails []mail.Address)

	// AddCopyRecipients adds `Cc` to the email
	AddCopyRecipients(emails []mail.Address)

	// AddBlindCopyRecipients adds `Bcc` to the email
	AddBlindCopyRecipients(emails []mail.Address)

	// AddHeader adds a custom header to the email
	AddHeader(key string, value string)

	// AddInlineFile embeds a file in the email body
	AddInlineFile(filePath string) error

	// SetSender sets `From` in the email
	SetSender(u mail.Address)

	// SetReplyToEmail sets `Reply-To` in the email
	SetReplyToEmail(email string)

	// UpdateSubject update the subject of the email
	UpdateSubject(subject string)

	// AttachFile opens the given file and adds it as an attachment to the email
	AttachFile(filePath string) error

	// AttachFile adds the given file as an attachment to the email
	AttachOpenedFile(file *os.File) error

	// AttachFileBytes puts the given bytes as an attached file to the email
	AttachFileBytes(fileName string, binary []byte)

	// Send sends the email
	Send() error
}

// Config to get Client using smtp credentials
type Config struct {
	// Host represents the smtp server's host
	Host string
	// Port represents the smtp server's port
	Port int
	// Email is the email used to sign in to the smtp server
	Email string
	// Password is the password used to sign in to the smtp server
	Password string
}
