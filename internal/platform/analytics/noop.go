package analytics

import "context"

type NoopAnalytics struct{}

func NewNoop() Analytics {
	return &NoopAnalytics{}
}

func (n *NoopAnalytics) Track(context.Context, string, map[string]any) {}

func (n *NoopAnalytics) Close() error {
	return nil
}
