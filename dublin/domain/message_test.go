package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_it_has_getters(t *testing.T) {

	id := "message_id"
	payload := SomeEvent{}
	playhead := 15
	metadata := Metadata{MetadataValuesT{"foo": "bar"}}
	messageType := "domain.SomeEvent"

	domainMessage := RecordMessageNow(id, playhead, metadata, payload)

	assert.Equal(t, id, domainMessage.id)
	assert.Equal(t, payload, domainMessage.payload)
	assert.Equal(t, playhead, domainMessage.playhead)
	assert.Equal(t, metadata, domainMessage.metadata)
	assert.Equal(t, messageType, domainMessage.GetType())
}

func Test_it_returns_a_new_instance_with_more_metadata_on_and_metadata(t *testing.T) {

	domainMessage := RecordMessageNow("message_id", 42, NewMetadata(MetadataValuesT{}), "payload")
	newDomainMessage := domainMessage.andMetadata(NewMetadataKV("key", "value"))

	assert.NotEqual(t, domainMessage, newDomainMessage)

	assert.Len(t, domainMessage.metadata.values, 0)
	assert.Len(t, newDomainMessage.metadata.values, 1)
}

func Test_it_keeps_all_data_the_same_expect_metadata_on_and_metadata(t *testing.T) {

	domainMessage := RecordMessageNow("message_id", 42, NewMetadata(MetadataValuesT{}), "payload")
	newDomainMessage := domainMessage.andMetadata(NewMetadataKV("key", "value"))

	assert.Equal(t, domainMessage.id, newDomainMessage.id)
	assert.Equal(t, domainMessage.playhead, newDomainMessage.playhead)
	assert.Equal(t, domainMessage.payload, newDomainMessage.payload)
	assert.Equal(t, domainMessage.recordedOn, newDomainMessage.recordedOn)

	assert.NotEqual(t, domainMessage.metadata, newDomainMessage.metadata)

}

func Test_it_merges_the_metadata_instances_on_and_metadata(t *testing.T) {

	domainMessage := RecordMessageNow("message_id", 42, NewMetadataKV("key", "value"), "payload").
		andMetadata(NewMetadataKV("foo", "bar"))

	expected := Metadata{MetadataValuesT{"key": "value", "foo": "bar"}}

	assert.Equal(t, expected, domainMessage.metadata)
}

type SomeEvent struct {
}
