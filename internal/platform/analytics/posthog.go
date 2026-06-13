package analytics

import (
	"context"
	"log"

	"github.com/posthog/posthog-go"
)

type PosthogAnalytics struct {
	client posthog.Client
}

func NewPosthog(apiKey, host string) (Analytics, error) {
	client, err := posthog.NewWithConfig(
		apiKey,
		posthog.Config{
			Endpoint: host,
		},
	)
	if err != nil {
		return nil, err
	}

	log.Println("PostHog enabled")

	return &PosthogAnalytics{client: client}, nil
}

func (p *PosthogAnalytics) Track(_ context.Context, event string, properties map[string]any) {
	err := p.client.Enqueue(posthog.Capture{
		DistinctId: "system",
		Event:      event,
		Properties: properties,
	})
	if err != nil {
		log.Println("posthog enqueue error:", err)
	}
}

func (p *PosthogAnalytics) Close() error {
	return p.client.Close()
}
