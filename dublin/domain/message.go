package domain

import (
	"reflect"
	"time"
)

// Message represents an important change in the domain.
type Message struct {
	playhead   int
	metadata   Metadata
	payload    interface{}
	id         string
	recordedOn time.Time
}

func (dm Message) getType() string {
	return reflect.TypeOf(dm.payload).String()
}

func (dm Message) andMetadata(metadata Metadata) Message {
	newMetadata := dm.metadata.merge(metadata)

	return Message{dm.playhead, newMetadata, dm.payload, dm.id, dm.recordedOn}
}

func (dm Message) recordNow(id string, playhead int, metadata Metadata, payload interface{}) Message {
	return Message{playhead, metadata, payload, id, time.Now()}
}

//func NewDomainMessage(id string, playhead int, metadata Metadata, payload interface{}) DomainMessage {
//	return DomainMessage{playhead, metadata, payload, id, time.Now()}
//}
