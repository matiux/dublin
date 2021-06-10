package commandhandling

import "fmt"

type CommandHandlerError struct {
	CommandHandler string
	Command
	OriginalError error
}

func (e CommandHandlerError) Error() string {
	return fmt.Sprintf(
		"Error in Command Handler `%v` with Command `%v`. Original error: %v",
		e.CommandHandler,
		e.Command.GetName(),
		e.OriginalError,
	)
}
