package events

import (
	"context"
)

// EventHandler is an interface that should be implemented by event handlers.
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}
