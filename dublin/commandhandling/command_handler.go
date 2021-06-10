package commandhandling

// CommandHandler handles dispatched commands.
type CommandHandler interface {
	Handle(cmd Command) error
}
