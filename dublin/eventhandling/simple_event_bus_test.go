package eventhandling

import "testing"

var eventBus *SimpleEventBus

func setup(t *testing.T) {

	eventBus = NewSimpleEventBus()
}

func Test_it_subscribes_an_event_listener(t *testing.T) {

	setup(t)
}

type TestEvent struct {
	Name string
}