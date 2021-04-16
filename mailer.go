package goemail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"net/mail"
	"net/smtp"
	"os"
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

func (m *simpleMailer) AddRecipients(emails []mail.Address) {
	if len(emails) == 0 {
		return
	}

	for _, email := range emails {
		m.recipients += email.String() + ", "
		m.to = append(m.to, email.Address)
	}
	m.recipients = removeLastComma(m.recipients)
}

func (m *simpleMailer) AddCopyRecipients(emails []mail.Address) {
	if len(emails) == 0 {
		return
	}

	for _, email := range emails {
		m.defaultHeaders.cc += email.String() + ", "
		m.to = append(m.to, email.Address)
	}
	m.defaultHeaders.cc = removeLastComma(m.defaultHeaders.cc)
}

func (m *simpleMailer) AddBlindCopyRecipients(emails []mail.Address) {
	for _, email := range emails {
		m.to = append(m.to, email.Address)
	}
}

func (m *simpleMailer) SetSender(u mail.Address) {
	m.from = u.Address
	m.defaultHeaders.from = u.String()
}

func (m *simpleMailer) SetReplyToEmail(email string) {
	m.defaultHeaders.replyTo = email
}

func (m *simpleMailer) SetSubject(subject string) {
	m.subject = subject
}

func (m *simpleMailer) AddHeader(key string, value string) {
	m.customHeaders = append(m.customHeaders, Header{Key: key, Value: value})
}

func (m *simpleMailer) attach(fileName string, binary []byte, inline bool) {
	m.attachments[fileName] = &Attachment{
		Filename: fileName,
		Data:     binary,
		Inline:   inline,
	}
}

func (m *simpleMailer) readAndAttach(file *os.File, inline bool) error {
	var binary []byte
	_, err := file.Read(binary)
	if err != nil {
		return err
	}

	m.attach(file.Name(), binary, inline)
	return nil
}

func (m *simpleMailer) openReadAndAttach(filePath string, inline bool) error {
	binary, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(filePath)

	m.attach(fileName, binary, inline)
	return nil
}

func (m *simpleMailer) AddInlineFile(filePath string) error {
	return m.openReadAndAttach(filePath, true)
}

func (m *simpleMailer) AttachFile(filePath string) error {
	return m.openReadAndAttach(filePath, false)
}

func (m *simpleMailer) AttachOpenedFile(file *os.File) error {
	return m.readAndAttach(file, false)
}

func (m *simpleMailer) AttachFileBytes(fileName string, binary []byte) {
	m.attach(fileName, binary, false)
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

func (m *simpleMailer) Send() error {
	return smtp.SendMail(m.client.addr, m.client.auth, m.from, m.to, m.parseMessage())
}
