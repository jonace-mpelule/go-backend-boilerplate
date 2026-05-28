package analytics

type NoopAnalytics struct {
}

func NewNoop() *NoopAnalytics {
	return &NoopAnalytics{}
}

func (n *NoopAnalytics) Track(
	event string,
	properties map[string]any,
) {

}
