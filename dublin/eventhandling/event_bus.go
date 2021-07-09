package eventhandling

import "github.com/matiux/dublin/dublin/domain"

// EventBus publish events to the subscribed event listeners.
type EventBus interface {

	// Subscribe subscribes the event listener to the event bus.
	Subscribe(eventListener EventListener)

	// Publish publishes the events from the domain event stream to the listeners.
	Publish(domainMessages domain.EventStream) error
}
