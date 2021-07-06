package commandhandling

import "github.com/matiux/dublin/dublin/eventdispatcher"

const event_command_success = "dublin.command_handling.command_success"
const event_command_failure = "dublin.command_handling.command_failure"

// EventDispatchingCommandBus dispatches events signalling whether a command was executed successfully or if it failed.
// Is a Command bus decorator that dispatches events.
type EventDispatchingCommandBus struct {
	CommandBus
	eventdispatcher.EventDispatcher
}

func (cb *EventDispatchingCommandBus) Subscribe(commandHandler CommandHandler) {
	cb.CommandBus.Subscribe(commandHandler)
}

func (cb *EventDispatchingCommandBus) Dispatch(command Command) error {

	_ = cb.CommandBus.Dispatch(command)

	arguments := map[string]interface{}{
		"command": command,
	}

	cb.EventDispatcher.Dispatch(event_command_success, arguments)

	return nil
}
