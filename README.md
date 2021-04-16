# goemail [![Build Status](https://api.travis-ci.com/IamFaizanKhalid/goemail.svg?branch=master)](https://travis-ci.com/github/IamFaizanKhalid/goemail) [![Go Report Card](https://goreportcard.com/badge/github.com/IamFaizanKhalid/goemail)](https://goreportcard.com/report/github.com/IamFaizanKhalid/goemail) ![License](https://img.shields.io/badge/license-MIT-blue.svg)
<img align="right" src="https://mimepost.com/blog/content/images/size/w600/2021/02/Untitled_design-removebg-preview-1.png" width="250">

A wrapper over `net/smtp` to make sending email easier in go.

## Features 
- HTML email
- HTML template email
- Attach file (`[]byte`, `os.File` or by file path)

## Usage

### 1. Get Client:
Use your smtp server credentials to get client object.
```go
client := goemail.NewClient(&goemail.Config{
    Host:     "smtp.gmail.com",
    Port:     587,
    Email:    "user@example.com",
    Password: "password",
})
```

### 2. Get Mailer:
Get a new mailer object from the client for each new email to send.
```go
mailer := client.NewMailer("Test Email", "This is an email for testing.")
})
```

### 3. Build Email and Send:
Add recipients of your email.
Attach a file if you need to.
And send the email.
```go
mailer.AddRecipients([]mail.Address{
    {
        Name:  "Random Guy",
        Address: "randomguy123@example.com",
    },
})
mailer.AddBlindCopyRecipients([]mail.Address{
    {
        Address: "secret01@example.com",
    },
})
mailer.SetSender(mail.Address{
    Name:  "Faizan Khalid",
})

mailer.SetReplyToEmail("no-reply@example.com")

_ = mailer.AttachFile("../Downloads/my_file.pdf")

_ = mailer.Send()
```


## Html Template Example

```go
package main

import (
    "github.com/IamFaizanKhalid/goemail"
    "log"
    "net/mail"
)

func main() {
    client := goemail.NewClient(&goemail.Config{
        Host:     "smtp.gmail.com",
        Port:     587,
        Email:    "user@example.com",
        Password: "password",
    })
    
    templateValues := struct {
        Name  string
        Url   string
        Title string
    }{
        Name:  "John Doe",
        Url:   "http://johndoe.com",
        Title: "Welcome to my Homepage",
    }
        
    mailer, err := client.NewHtmlMailerFromTemplate("JohnDoe.com", "welcome.html", templateValues)
    if err != nil {
        log.Println(err)
        return
    }

    mailer.AddRecipients([]mail.Address{{Address: "IamFaizanKhalid@gmail.com"}})

    err = mailer.Send()
    if err != nil {
        log.Println(err)
        return
    }
}
```

## Known Issues

- `AddInlineFile()` adds gibberish in the email.
