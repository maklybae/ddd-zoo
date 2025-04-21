package events

import (
	"context"
	"fmt"
	"sync"
)

type Dispatcher interface {
	RegisterHandler(eventType string, handler EventHandler)
	Dispatch(ctx context.Context, event Event)
}

type EventDispatcher struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

func (d *EventDispatcher) RegisterHandler(eventType string, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.handlers[eventType]; !exists {
		d.handlers[eventType] = []EventHandler{}
	}

	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *EventDispatcher) Dispatch(ctx context.Context, event Event) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if handlers, ok := d.handlers[event.Name()]; ok {
		for _, handler := range handlers {
			if err := handler.Handle(ctx, event); err != nil {
				fmt.Println("Error handling event:", err)
			}
		}
	}
}
