package commandhandling

import (
	"github.com/matiux/dublin/dublin/eventdispatcher"
)

const EventCommandSuccess = "dublin.command_handling.command_success"
const EventCommandFailure = "dublin.command_handling.command_failure"

// EventDispatchingCommandBus dispatches events signal<l>ing whether a command was executed successfully or if it failed.
// Is a Command bus decorator that dispatches events.
type EventDispatchingCommandBus struct {
	CommandBus
	eventdispatcher.EventDispatcher
}

func (cb *EventDispatchingCommandBus) Subscribe(commandHandler CommandHandler) {
	cb.CommandBus.Subscribe(commandHandler)
}

func (cb *EventDispatchingCommandBus) Dispatch(command Command) error {

	defer func() {
		if r := recover(); r != nil {
			_ = cb.EventDispatcher.Dispatch(EventCommandFailure, map[string]interface{}{
				"command":   command,
				"exception": r,
			})
		}
	}()

	cb.commandBusOrFail(command)
	cb.eventDispatchingCommandBusOrFail(command)

	return nil
}

func (cb *EventDispatchingCommandBus) commandBusOrFail(command Command) {

	if err := cb.CommandBus.Dispatch(command); err != nil {
		panic(err)
	}
}

func (cb *EventDispatchingCommandBus) eventDispatchingCommandBusOrFail(command Command) {

	err := cb.EventDispatcher.Dispatch(EventCommandSuccess, map[string]interface{}{
		"command": command,
	})

	if err != nil {
		panic(err)
	}
}

func NewEventDispatchingCommandBus(commandBus CommandBus, eventDispatcher eventdispatcher.EventDispatcher) *EventDispatchingCommandBus {
	return &EventDispatchingCommandBus{commandBus, eventDispatcher}
}
