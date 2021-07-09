package eventhandling

import "github.com/matiux/dublin/dublin/domain"

// EventListener handles dispatched events.
type EventListener interface {
	Handle(domainMessage domain.Message) error
}
