package sentryplatform

import (
	"log"

	"github.com/getsentry/sentry-go"
)

func Init(dsn string) error {
	if dsn == "" {
		log.Println("Sentry disabled")
		return nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
	})

	if err != nil {
		return err
	}

	log.Println("Sentry enabled")

	return nil
}
