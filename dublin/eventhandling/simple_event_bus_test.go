package eventhandling

import (
	"github.com/golang/mock/gomock"
	"github.com/matiux/dublin/dublin/domain"
	"testing"
)

var eventBus *SimpleEventBus

func setup(t *testing.T) {

	eventBus = NewSimpleEventBus()
}

func Test_it_subscribes_an_event_listener(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	setup(t)

	domainMessage := createDomainMessage(map[string]string{"foo": "bar"})
	eventListener := NewMockEventListener(ctrl)

	eventListener.EXPECT().Handle(domainMessage).Times(1)

	eventBus.Subscribe(eventListener)

	_ = eventBus.Publish(domain.NewEventStream([]domain.Message{domainMessage}))

}

func Test_it_publishes_events_to_subscribed_event_listeners(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	setup(t)

	domainMessage1 := createDomainMessage([]string{})
	domainMessage2 := createDomainMessage([]string{})

	domainEventStream := domain.NewEventStream([]domain.Message{domainMessage1, domainMessage2})

	eventListener1 := NewMockEventListener(ctrl)
	gomock.InOrder(
		eventListener1.EXPECT().Handle(domainMessage1),
		eventListener1.EXPECT().Handle(domainMessage2),
	)

	eventListener2 := NewMockEventListener(ctrl)
	gomock.InOrder(
		eventListener2.EXPECT().Handle(domainMessage1),
		eventListener2.EXPECT().Handle(domainMessage2),
	)

	eventBus.Subscribe(eventListener1)
	eventBus.Subscribe(eventListener2)
	_ = eventBus.Publish(domainEventStream)
}

func Test_it_does_not_dispatch_new_events_before_all_listeners_have_run(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	setup(t)

	domainMessage1 := createDomainMessage(map[string]string{"foo": "bar"})
	domainMessage2 := createDomainMessage(map[string]string{"foo": "bas"})

	domainEventStream := domain.NewEventStream([]domain.Message{domainMessage1})

	eventListener1 := SimpleEventBusTestListener{
		eventBus,
		false,
		domain.NewEventStream([]domain.Message{domainMessage2}),
	}

	eventListener2 := NewMockEventListener(ctrl)
	gomock.InOrder(
		eventListener2.EXPECT().Handle(domainMessage1).Times(1),
		eventListener2.EXPECT().Handle(domainMessage2).Times(1),
	)

	eventBus.Subscribe(&eventListener1)
	eventBus.Subscribe(eventListener2)
	_ = eventBus.Publish(domainEventStream)
}

func Test_it_should_still_publish_events_after_exception(t *testing.T) {


}

type SimpleEventBusTestEvent interface{}

func createDomainMessage(payload SimpleEventBusTestEvent) domain.Message {

	return domain.RecordMessageNow("1", 1, domain.NewMetadata(domain.MetadataValuesT{}), payload)
}

type SimpleEventBusTestListener struct {
	EventBus
	handled           bool
	publishableStream domain.EventStream
}

func (eb *SimpleEventBusTestListener) Handle(domainMessage domain.Message) error {

	if !eb.handled {
		_ = eb.EventBus.Publish(eb.publishableStream)
		eb.handled = true
	}

	return nil
}
