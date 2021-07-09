package commandhandling

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/matiux/dublin/dublin/eventdispatcher"
	"testing"
)

var eventDispatcher eventdispatcher.EventDispatcher
var command Command
var commandBus CommandBus
var subscriber CommandHandler
var eventDispatchingCommandBus *EventDispatchingCommandBus

//func TestMain(m *testing.M) {
//	setup()
//	code := m.Run()
//	shutdown()
//	os.Exit(code)
//}

func setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventDispatcher = eventdispatcher.NewMockEventDispatcher(ctrl)
	commandBus = NewMockCommandBus(ctrl)
	subscriber = NewMockCommandHandler(ctrl)
	command = TestCommand{Name: "matiux"}
	eventDispatchingCommandBus = NewEventDispatchingCommandBus(commandBus, eventDispatcher)
}

//func shutdown() {
//	fmt.Println("3")
//
//}
func Test_it_dispatches_the_success_event(t *testing.T) {

	setup(t)

	gomock.InOrder(
		commandBus.(*MockCommandBus).EXPECT().Dispatch(command),
		eventDispatcher.(*eventdispatcher.MockEventDispatcher).
			EXPECT().
			Dispatch(
				EventCommandSuccess,
				map[string]interface{}{
					"command": command,
				},
			),
	)

	_ = eventDispatchingCommandBus.Dispatch(command)
}

func Test_it_dispatches_the_failure_event_and_forwards_the_exception(t *testing.T) {

	setup(t)

	myError := fmt.Errorf("an myError")

	gomock.InOrder(
		commandBus.(*MockCommandBus).EXPECT().Dispatch(command).Return(myError).Times(1),
		eventDispatcher.(*eventdispatcher.MockEventDispatcher).
			EXPECT().
			Dispatch(
				EventCommandFailure,
				map[string]interface{}{"command": command, "exception": myError},
			).Times(1),
	)

	_ = eventDispatchingCommandBus.Dispatch(command)
}

//func Test_it_forwards_the_dispatched_command(t *testing.T) {
//
//	setup(t)
//	gomock.InOrder(
//		commandBus.(*MockCommandBus).EXPECT().Dispatch(command).Times(1),
//		eventDispatcher.(*eventdispatcher.MockEventDispatcher).
//			EXPECT().
//			Dispatch(
//				EventCommandSuccess,
//				map[string]interface{}{
//					"command": command,
//				},
//			),
//	)
//	_ = eventDispatchingCommandBus.Dispatch(command)
//}

func Test_it_forwards_the_subscriber(t *testing.T) {

	setup(t)

	commandBus.(*MockCommandBus).EXPECT().Subscribe(subscriber).Times(1)

	eventDispatchingCommandBus.Subscribe(subscriber)
}