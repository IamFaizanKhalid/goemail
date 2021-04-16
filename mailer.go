package goemail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"net/smtp"
	"path/filepath"
	"time"
)

type simpleMailer struct {
	client         *mailClient
	defaultHeaders mailHeaders
	customHeaders  []Header
	to             []string
	from           string
	recipients     string
	subject        string
	body           string
	attachments    map[string]*Attachment
}

type mailHeaders struct {
	from        string
	to          string
	cc          string
	replyTo     string
	contentType string
}

func (m *simpleMailer) AddRecipients(emails []User) {
	if len(emails) == 0 {
		return
	}

	for _, email := range emails {
		m.recipients += email.String() + ", "
		m.to = append(m.to, email.Email)
	}
	m.recipients = removeLastComma(m.recipients)
}

func (m *simpleMailer) AddCopyRecipients(emails []User) {
	if len(emails) == 0 {
		return
	}

	for _, email := range emails {
		m.defaultHeaders.cc += email.String() + ", "
		m.to = append(m.to, email.Email)
	}
	m.defaultHeaders.cc = removeLastComma(m.defaultHeaders.cc)
}

func (m *simpleMailer) AddBlindCopyRecipients(emails []User) {
	for _, email := range emails {
		m.to = append(m.to, email.Email)
	}
}

func (m *simpleMailer) AddSender(u User) {
	m.from = u.Email
	m.defaultHeaders.from = u.String()
}

func (m *simpleMailer) AddReplyToMail(email string) {
	m.defaultHeaders.replyTo = email
}

func (m *simpleMailer) AddSubject(subject string) {
	m.subject = subject
}

func (m *simpleMailer) AddHeader(key string, value string) {
	m.customHeaders = append(m.customHeaders, Header{Key: key, Value: value})
}

func (m *simpleMailer) attach(file string, inline bool) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	_, filename := filepath.Split(file)

	m.attachments[filename] = &Attachment{
		Filename: filename,
		Data:     data,
		Inline:   inline,
	}

	return nil
}

func (m *simpleMailer) AttachFile(filePath string) error {
	return m.attach(filePath, false)
}

func (m *simpleMailer) InsertFile(filePath string) error {
	return m.attach(filePath, true)
}

func removeLastComma(s string) string {
	return s[:len(s)-2]
}

func (m *simpleMailer) parseMessage() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("From: " + m.defaultHeaders.from + "\r\n")
	buf.WriteString("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n")
	buf.WriteString("To: " + m.defaultHeaders.to + "\r\n")
	if m.defaultHeaders.cc != "" {
		buf.WriteString("Cc: " + m.defaultHeaders.cc + "\r\n")
	}
	if m.defaultHeaders.replyTo != "" {
		buf.WriteString("Reply-To: " + m.defaultHeaders.replyTo + "\r\n")
	}

	// fix  Encode
	var coder = base64.StdEncoding
	var subject = "=?UTF-8?B?" + coder.EncodeToString([]byte(m.subject)) + "?="
	buf.WriteString("Subject: " + subject + "\r\n")

	buf.WriteString("MIME-Version: 1.0\r\n")

	// Add custom headers
	for _, header := range m.customHeaders {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", header.Key, header.Value))
	}

	boundary := "f46d043c813270fc6b04c2d223da"

	if len(m.attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("\r\n--" + boundary + "\r\n")
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", m.defaultHeaders.contentType))
	buf.WriteString(m.body)
	buf.WriteString("\r\n")

	if len(m.attachments) > 0 {
		for _, attachment := range m.attachments {
			buf.WriteString("\r\n\r\n--" + boundary + "\r\n")

			if attachment.Inline {
				buf.WriteString("Content-Type: message/rfc822\r\n")
				buf.WriteString("Content-Disposition: inline; filename=\"" + attachment.Filename + "\"\r\n\r\n")
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")

				buf.Write(attachment.Data)
			} else {
				ext := filepath.Ext(attachment.Filename)
				mimetype := mime.TypeByExtension(ext)
				if mimetype != "" {
					buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", mimetype))
				} else {
					buf.WriteString("Content-Type: application/octet-stream\r\n")
				}
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")

				buf.WriteString("Content-Disposition: attachment; filename=\"=?UTF-8?B?")
				buf.WriteString(coder.EncodeToString([]byte(attachment.Filename)))
				buf.WriteString("?=\"\r\n\r\n")

				b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
				base64.StdEncoding.Encode(b, attachment.Data)

				// write base64 content in lines of up to 76 chars
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i+1)%76 == 0 {
						buf.WriteString("\r\n")
					}
				}
			}

			buf.WriteString("\r\n--" + boundary)
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

func (m *simpleMailer) SendEmail() error {
	return smtp.SendMail(m.client.addr, m.client.auth, m.from, m.to, m.parseMessage())
}
