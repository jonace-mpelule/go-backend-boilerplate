package analytics

import "context"

type Analytics interface {
	Track(ctx context.Context, event string, properties map[string]any)
	Close() error
}
