package commandhandling

import (
	"github.com/golang/mock/gomock"
	"github.com/matiux/dublin/dublin/eventdispatcher"
	"testing"
)

var eventDispatcher eventdispatcher.EventDispatcher
var command Command
var commandBus CommandBus

//var subscriber CommandHandler
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
	//subscriber = NewMockCommandHandler(ctrl)
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
