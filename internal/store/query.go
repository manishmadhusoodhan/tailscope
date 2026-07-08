package store

type Query struct {
	TraceID string

	CorrelationID string

	Level string

	Pod string

	Namespace string

	Text string

	Limit int
}
