package mailer

import (
	"errors"
	"fmt"
	"net/smtp"
)

var ErrFailedToSendMail = errors.New("failed to send mail")

type DefaultMailer struct {
	Host     string
	Port     int
	From     string
	Username string
	Password string
}

type MailerParams struct {
	Host     string
	Port     int
	From     string
	Username string
	Password string
}

func NewMailer(params MailerParams) Mailer {
	return &DefaultMailer{
		Host:     params.Host,
		Port:     params.Port,
		From:     params.From,
		Username: params.Username,
		Password: params.Password,
	}
}

func (m *DefaultMailer) SendMail(subject, body string) error {
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	headers := ""
	headers += fmt.Sprintf("From: %s\r\n", m.From)
	headers += fmt.Sprintf("To: %s\r\n", m.From)
	headers += fmt.Sprintf("Subject: %s\r\n", subject)
	headers += "MIME-Version: 1.0\r\n"
	headers += "Content-Type: text/html; charset=UTF-8\r\n"
	headers += "\r\n"

	msg := []byte(headers + body)

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)

	if err := smtp.SendMail(addr, auth, m.From, []string{m.From}, msg); err != nil {
		return errors.Join(err, ErrFailedToSendMail)
	}

	return nil
}
