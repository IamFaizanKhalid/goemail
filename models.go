package goemail

import "fmt"

type Mailer interface {
	Send() error
}

type Config struct {
	Host     string `validate:"nonzero"`
	Port     string `validate:"regexp=^\d{2,5}$"`
	Email    string `validate:"regexp=^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$"`
	Password string `validate:"min=5,max=256"`
}

type Mail struct {
	To      []User `validate:"nonzero"`
	Cc      []User
	Bcc     []User
	From    *User
	Subject string `validate:"nonzero"`
	Message string `validate:"nonzero"`
}

type User struct {
	Name  string
	Email string `validate:"regexp=^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$"`
}

func (u *User) String() string {
	if u.Name == "" {
		return u.Email
	}
	return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}
