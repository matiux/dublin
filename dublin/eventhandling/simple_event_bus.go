package eventhandling

import (
	"github.com/matiux/dublin/dublin/domain"
	"reflect"
)

// SimpleEventBus - synchronous publishing of events.
type SimpleEventBus struct {
	eventListeners []*EventListener
	queue          []domain.Message
	isPublishing   bool
}

func NewSimpleEventBus() *SimpleEventBus {
	return &SimpleEventBus{isPublishing: false}
}

func (eb *SimpleEventBus) Subscribe(eventListener EventListener) {
	eb.eventListeners = append(eb.eventListeners, &eventListener)
}

func (eb *SimpleEventBus) Publish(domainMessages domain.EventStream) error {

	for dm := domainMessages.GetIterator().Front(); dm != nil; dm = dm.Next() {
		eb.queue = append(eb.queue, dm.Value.(domain.Message))
	}

	if !eb.isPublishing {
		eb.isPublishing = true

		for len(eb.queue) > 0 {
			message := eb.queue[0]
			eb.queue = eb.queue[1:]

			for _, eventListener := range eb.eventListeners {

				// TODO - Perché non viene deindirizzato automaticamente con eventListener.Handle(message) ?
				if err := (*eventListener).Handle(message); err != nil {

					eb.isPublishing = false

					return EventListenerError{
						EventListener: reflect.TypeOf(eventListener).Elem().Name(),
						Message:       message,
						OriginalError: err,
					}
				}
			}
		}

		eb.isPublishing = false
	}

	return nil
}
