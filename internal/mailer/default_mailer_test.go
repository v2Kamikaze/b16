package mailer

import (
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeClient struct {
	err error
}

func (c *fakeClient) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return c.err
}

type TestMailerParams struct {
	Name      string
	Emails    []string
	ExpectErr error
	ClientErr error
}

func TestMailer_Send(t *testing.T) {

	tests := []TestMailerParams{
		{
			Name:      "successfully send email",
			Emails:    []string{"test@email.com"},
			ExpectErr: nil,
			ClientErr: nil,
		},
		{
			Name:      "failed to send email",
			Emails:    []string{"test1@email.com", "test2@email.com"},
			ExpectErr: ErrFailedToSendMail,
			ClientErr: ErrFailedToSendMail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := NewDefaultMailer(MailerParams{
				Host:     "localhost",
				Port:     25,
				From:     "b16@email.com",
				Username: "username",
				Password: "password",
			}, &fakeClient{
				err: tt.ClientErr,
			})

			body := ""
			err := m.Send("b16@email.com", body, tt.Emails...)
			assert.ErrorIs(t, err, tt.ExpectErr)
		})
	}
}
