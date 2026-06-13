package mailer

import "context"

type NoopMailer struct{}

func NewNoop() Mailer {
	return &NoopMailer{}
}

func (n *NoopMailer) Send(context.Context, string, string, string) error {
	return nil
}
