package events

// Event interface that should be implemented by events.
type Event interface {
	Name() string
}
