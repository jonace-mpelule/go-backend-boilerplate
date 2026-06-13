package mailer

import (
	"context"

	"go.uber.org/zap"
)

type LogMailer struct {
	logger *zap.Logger
	from   string
}

func NewLog(logger *zap.Logger, from string) Mailer {
	return &LogMailer{
		logger: logger,
		from:   from,
	}
}

func (l *LogMailer) Send(ctx context.Context, to, subject, html string) error {
	l.logger.Info(
		"mail sent via log mailer",
		zap.String("from", l.from),
		zap.String("to", to),
		zap.String("subject", subject),
		zap.String("html", html),
	)
	return nil
}
