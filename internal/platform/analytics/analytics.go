package analytics

type Analtics interface {
	Track(
		event string,
		properties map[string]any,
	)
}
