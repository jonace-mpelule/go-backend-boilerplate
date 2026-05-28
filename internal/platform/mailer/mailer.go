package mailer

import "context"

type Mailer interface {
	Send(
		ctx context.Context,
		to string,
		subject string,
		html string,
	) error
}
