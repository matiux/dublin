package commandhandling

// CommandBus dispatches command objects to the subscribed command handlers.
type CommandBus interface {

	// Subscribe the command handler to this CommandBus.
	Subscribe(commandHandler CommandHandler)

	// Dispatch the command `command` to the proper CommandHandler.
	Dispatch(command Command) error
}
