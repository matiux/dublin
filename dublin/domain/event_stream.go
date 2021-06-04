package domain

import (
	"container/list"
)

// EventStream a stream of Message(S) in sequence.
type EventStream struct {
	messages []Message
}

func (es EventStream) getIterator() *list.List {

	l := list.New()

	for _, m := range es.messages {
		l.PushBack(m)
	}

	return l
}

func NewEventStream(messages []Message) EventStream {
	return EventStream{messages}
}
