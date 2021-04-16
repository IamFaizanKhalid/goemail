package goemail

import "fmt"

type Mailer interface {
	Send() error
}

type Config struct {
	Host     string
	Port     string
	Email    string
	Password string
}

type Mail struct {
	To      []User
	Cc      []User
	Bcc     []User
	From    *User
	Subject string
	Message string
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
