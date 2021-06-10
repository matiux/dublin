package commandhandling

import (
	"reflect"
)

// SimpleCommandBus is a simple synchronous dispatching of commands.
type SimpleCommandBus struct {
	commandHandlers []*CommandHandler
	queue           []Command
	isDispatching   bool
}

func NewSimpleCommandBus() *SimpleCommandBus {
	return &SimpleCommandBus{isDispatching: false}
}

func (scb *SimpleCommandBus) Subscribe(commandHandler CommandHandler) {
	scb.commandHandlers = append(scb.commandHandlers, &commandHandler)
}

func (scb *SimpleCommandBus) Dispatch(command Command) error {

	scb.queue = append(scb.queue, command)

	if !scb.isDispatching {
		scb.isDispatching = true

		for len(scb.queue) > 0 {
			command := scb.queue[0]
			scb.queue = scb.queue[1:]

			for _, commandHandler := range scb.commandHandlers {

				// TODO - Perché non viene deindirizzato automaticamente con commandHandler.Handle(command) ?
				if err := (*commandHandler).Handle(command); err != nil {

					scb.isDispatching = false

					return CommandHandlerError{
						CommandHandler: reflect.TypeOf(commandHandler).Elem().Name(),
						Command:        command,
						OriginalError:  err,
					}
				}
			}
		}

		scb.isDispatching = false
	}

	return nil
}
