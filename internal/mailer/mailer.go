package mailer

type Mailer interface {
	SendMail(subject, body string) error
}
