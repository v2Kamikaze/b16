package mailer

import "net/smtp"

type Mailer interface {
	Send(subject string, body string, to ...string) error
}

type SMTPClient interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}
