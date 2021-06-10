package commandhandling

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_it_dispatches_commands_to_subscribed_handlers(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := TestCommand{Name: "matiux"}

	commandHandlerMock1 := NewMockCommandHandler(ctrl)
	commandHandlerMock2 := NewMockCommandHandler(ctrl)

	// Assert that Handle() is invoked.
	// Asserts that the first and only call to Handle() is passed cmd. Anything else will fail.
	commandHandlerMock1.EXPECT().Handle(cmd).Times(1)
	commandHandlerMock2.EXPECT().Handle(cmd).Times(1)

	commandBus := NewSimpleCommandBus()
	commandBus.Subscribe(commandHandlerMock1)
	commandBus.Subscribe(commandHandlerMock2)
	commandBus.Dispatch(cmd)
}

func Test_it_does_not_handle_new_commands_before_all_commandhandlers_have_run(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd1 := TestCommand{Name: "mat"}
	cmd2 := TestCommand{Name: "teo"}

	commandHandlerMock := NewMockCommandHandler(ctrl)

	gomock.InOrder(
		commandHandlerMock.EXPECT().Handle(cmd1).Times(1),
		commandHandlerMock.EXPECT().Handle(cmd2).Times(1),
	)

	commandBus := NewSimpleCommandBus()

	simpleCommandTestHandler := TestCommandHandler{CommandBus: commandBus, handled: false, dispatchableCommand: cmd2}

	commandBus.Subscribe(&simpleCommandTestHandler)
	commandBus.Subscribe(commandHandlerMock)
	commandBus.Dispatch(cmd1)

	assert.Equal(t, simpleCommandTestHandler.fetchedValue, "mat")
}

func Test_it_should_still_handle_commands_after_exception(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd1 := TestCommand{Name: "mat"}
	cmd2 := TestCommand{Name: "teo"}

	commandHandler := NewMockCommandHandler(ctrl)
	simpleHandler := NewMockCommandHandler(ctrl)

	gomock.InOrder(
		commandHandler.EXPECT().Handle(cmd1).Times(1).Return(fmt.Errorf("an error")),
		commandHandler.EXPECT().Handle(cmd2).Times(1),
		simpleHandler.EXPECT().Handle(cmd2).Times(1),
	)

	commandBus := NewSimpleCommandBus()
	commandBus.Subscribe(commandHandler)
	commandBus.Subscribe(simpleHandler)

	if handlerError := commandBus.Dispatch(cmd1); handlerError != nil {
		assert.Equal(t, "Error in Command Handler `CommandHandler` with Command `TestCommand`. Original error: an error", handlerError.Error())
	}

	_ = commandBus.Dispatch(cmd2)
}

type TestCommand struct {
	Name string
}

func (c TestCommand) GetName() string {
	return reflect.TypeOf(c).Name()
}

type TestCommandHandler struct {
	CommandHandlerUtil
	CommandBus
	handled             bool
	dispatchableCommand Command
	fetchedValue        string
}

func (h *TestCommandHandler) Handle(cmd Command) error {

	method := h.CommandHandlerUtil.getMethodName(cmd)

	if _, ok := reflect.TypeOf(h).MethodByName(method); ok {
		reflect.ValueOf(h).MethodByName(method).Call([]reflect.Value{reflect.ValueOf(cmd)})
	}

	return nil
	//i, ok := interface{}(h).(interface{HandleSimpleCommandTest(cmd SimpleCommandTest)})
	//fmt.Println(i, ok)
}

func (h *TestCommandHandler) HandleTestCommand(cmd TestCommand) {
	if !h.handled {
		_ = h.CommandBus.Dispatch(h.dispatchableCommand)
		h.handled = true
		h.fetchedValue = cmd.Name
	}
}
