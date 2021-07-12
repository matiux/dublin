package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_it_returns_all_events_when_traversing(t *testing.T) {

	expected := []Message{RecordMessageNow("message_id", 42, NewMetadata(MetadataValuesT{}), "payload")}
	domainMessage := RecordMessageNow("message_id", 42, NewMetadata(MetadataValuesT{}), "payload")

	iterator := NewEventStream([]Message{domainMessage}).GetIterator()

	var events []Message

	for e := iterator.Front(); e != nil; e = e.Next() {
		events = append(events, e.Value.(Message))
	}

	assert.Len(t, events, 1)
	assert.Equal(t, expected[0].id, events[0].id)
}
