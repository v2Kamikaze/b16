package mailer

import "net/smtp"

type DefaultClient struct{}

func NewDefaultClient() SMTPClient {
	return &DefaultClient{}
}

func (c *DefaultClient) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return smtp.SendMail(addr, a, from, to, msg)
}
