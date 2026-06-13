package mailer

import (
	"context"
	"fmt"
	"net/smtp"
)

type SMTPMailer struct {
	addr string
	auth smtp.Auth
	from string
}

func NewSMTP(host string, port int, username, password, from string) Mailer {
	var auth smtp.Auth
	if username != "" || password != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	return &SMTPMailer{
		addr: fmt.Sprintf("%s:%d", host, port),
		auth: auth,
		from: from,
	}
}

func (s *SMTPMailer) Send(context.Context, string, string, string) error {
	return smtp.SendMail(
		s.addr,
		s.auth,
		s.from,
		nil,
		nil,
	)
}
