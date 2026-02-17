package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/v2code/b16/internal/logger"
)

const LOG_EMAIL_PREFIX = "EMAIL SEND"

var ErrFailedToSendMail = errors.New("failed to send mail")

type DefaultMailer struct {
	host     string
	port     int
	from     string
	username string
	password string
	client   SMTPClient
}

type MailerParams struct {
	Host     string
	Port     int
	From     string
	Username string
	Password string
}

func NewDefaultMailer(params MailerParams, client SMTPClient) Mailer {
	return &DefaultMailer{
		host:     params.Host,
		port:     params.Port,
		from:     params.From,
		username: params.Username,
		password: params.Password,
		client:   client,
	}
}

func (m *DefaultMailer) Send(subject string, body string, to ...string) error {
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	buffer := bytes.Buffer{}

	buffer.WriteString("From: " + m.from + "\r\n")
	buffer.WriteString("To: " + m.from + "\r\n")
	buffer.WriteString("Subject: " + subject + "\r\n")
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buffer.WriteString("\r\n")
	buffer.WriteString(body)

	logger.Debug(LOG_EMAIL_PREFIX, "Sending email to", to)

	if err := m.client.SendMail(m.BuildAddr(), auth, m.from, to, buffer.Bytes()); err != nil {
		logger.Error(LOG_EMAIL_PREFIX, "Error sending email", err.Error())
		return ErrFailedToSendMail
	}

	logger.Debug(LOG_EMAIL_PREFIX, "Email sent successfully to", to)

	return nil
}

func (m *DefaultMailer) BuildAddr() string {
	return fmt.Sprintf("%s:%d", m.host, m.port)
}
