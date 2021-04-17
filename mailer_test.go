package goemail

import (
	"net/mail"
	"strings"
	"testing"
)

// TestParseMessage checks if the message is being constructed properly
func TestParseMessage(t *testing.T) {
	mailer := NewClient(&Config{
		Host:     "smtp.gmail.com",
		Port:     587,
		Email:    "user@gmail.com",
		Password: "password",
	}).NewMailer("Test", "This is a test message.")

	mailer.AddRecipients([]mail.Address{
		{
			Name:    "Random Guy",
			Address: "randomguy123@example.com",
		},
	})
	mailer.AddBlindCopyRecipients([]mail.Address{
		{
			Address: "secret01@example.com",
		},
	})
	mailer.SetSender(mail.Address{
		Name: "Faizan Khalid",
	})

	mailer.SetReplyToEmail("no-reply@example.com")

	err := mailer.AttachFile("LICENSE")
	if err != nil {
		t.Errorf(err.Error())
	}

	message := string(mailer.parseMessage())

	expectedMessage := "To: \"Random Guy\" <randomguy123@example.com>\r\n" +
		"Reply-To: no-reply@example.com\r\n" +
		"Subject: =?UTF-8?B?VGVzdA==?=\r\n" +
		"MIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=f46d043c813270fc6b04c2d223da\r\n\r\n" +
		"--f46d043c813270fc6b04c2d223da\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n\r\n" +
		"This is a test message.\r\n\r\n\r\n" +
		"--f46d043c813270fc6b04c2d223da\r\n" +
		"Content-Type: application/octet-stream\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"Content-Disposition: attachment; filename=\"=?UTF-8?B?TElDRU5TRQ==?=\"\r\n\r\n" +
		"VGhlIE1JVCBMaWNlbnNlIChNSVQpCgpDb3B5cmlnaHQgKGMpIDIwMTQgQWxleGFuZHJlIENlc2Fy\r\nbwoKUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBw\r\nZXJzb24gb2J0YWluaW5nIGEgY29weSBvZgp0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRv\r\nY3VtZW50YXRpb24gZmlsZXMgKHRoZSAiU29mdHdhcmUiKSwgdG8gZGVhbCBpbgp0aGUgU29mdHdh\r\ncmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUg\r\ncmlnaHRzIHRvCnVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwg\r\nc3VibGljZW5zZSwgYW5kL29yIHNlbGwgY29waWVzIG9mCnRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBl\r\ncm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzIGZ1cm5pc2hlZCB0byBkbyBzbywK\r\nc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgpUaGUgYWJvdmUgY29weXJpZ2h0\r\nIG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBh\r\nbGwKY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS4KClRIRSBT\r\nT0ZUV0FSRSBJUyBQUk9WSURFRCAiQVMgSVMiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5E\r\nLCBFWFBSRVNTIE9SCklNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdB\r\nUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLCBGSVRORVNTCkZPUiBBIFBBUlRJQ1VMQVIgUFVS\r\nUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SUyBP\r\nUgpDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBP\r\nVEhFUiBMSUFCSUxJVFksIFdIRVRIRVIKSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9S\r\nIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLCBPVVQgT0YgT1IgSU4KQ09OTkVDVElPTiBXSVRIIFRI\r\nRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRSBTT0ZUV0FSRS4K\r\n" +
		"\r\n--f46d043c813270fc6b04c2d223da"

	if !strings.Contains(message, expectedMessage) {
		t.Errorf("message is not being generated properly")
	}
}
