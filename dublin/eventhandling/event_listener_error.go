package eventhandling

import (
	"fmt"
	"github.com/matiux/dublin/dublin/domain"
)

type EventListenerError struct {
	EventListener string
	domain.Message
	OriginalError error
}

func (e EventListenerError) Error() string {
	return fmt.Sprintf(
		"Error in Event Listener `%v` with Message `%v`. Original error: %v",
		e.EventListener,
		e.Message.GetType(),
		e.OriginalError,
	)
}

