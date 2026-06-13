package sentryplatform

import (
	"context"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

type Client struct {
	enabled bool
}

func Init(dsn string) (*Client, error) {
	if dsn == "" {
		log.Println("Sentry disabled")
		return &Client{enabled: false}, nil
	}

	if err := sentry.Init(sentry.ClientOptions{Dsn: dsn}); err != nil {
		return nil, err
	}

	log.Println("Sentry enabled")
	return &Client{enabled: true}, nil
}

func (c *Client) Enabled() bool {
	return c != nil && c.enabled
}

func (c *Client) Flush(ctx context.Context) {
	if !c.Enabled() {
		return
	}

	timeout := 2 * time.Second
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
		if timeout < 0 {
			timeout = 0
		}
	}

	sentry.Flush(timeout)
}
