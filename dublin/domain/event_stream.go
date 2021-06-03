package domain

// EventStream a stream of Message(S) in sequence.
type EventStream struct {
	events []Message
}

func (es EventStream) getIterator() StreamIterator {
	return NewStreamIterator(es.events)
}

func NewEventStream(events []Message) EventStream {
	return EventStream{events}
}
